package datastore

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/stretchr/testify/require"

	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var (
	ds  *Datastore
	log *zap.Logger
)

// to test unregistered type checks
type invalidVersionedEntity struct{}

func (v *invalidVersionedEntity) JSONField() string  { return "invalid" }
func (v *invalidVersionedEntity) TableName() string  { return "" }
func (v *invalidVersionedEntity) Schema() string     { return "" }
func (v *invalidVersionedEntity) Kind() string       { return "Invalid" }
func (v *invalidVersionedEntity) APIVersion() string { return "v1" }
func (v *invalidVersionedEntity) GetMeta() *v1.Meta  { return nil }

func TestMain(m *testing.M) {
	code := 0
	defer func() {
		os.Exit(code)
	}()

	// used to debug race condition in this method
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from ", r)
			fmt.Println(string(debug.Stack()))
		}
	}()

	log, _ = zap.NewProduction()

	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:12-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env:          map[string]string{"POSTGRES_PASSWORD": "password"},
		// TODO: should work, but dont, hence using the loop below to check pg is up an ready for connections.
		// WaitingFor:   wait.ForLog("database system is ready to accept connections"),
		// WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Info(err.Error())
	}
	defer func() {
		err := postgres.Terminate(ctx)
		if err != nil {
			log.Fatal(err.Error())
		}
	}()

	ip, err := postgres.Host(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	port, err := postgres.MappedPort(ctx, "5432/tcp")
	if err != nil {
		log.Fatal(err.Error())
	}

	for {
		var err error
		ds, err = NewPostgresStorage(log, ip, port.Port(), "postgres", "password", "postgres", "disable", &v1.Project{}, &v1.Tenant{})
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		err = ds.db.Ping()
		if err != nil {
			log.Error("Could not connect to postgres server", zap.Error(err))
		}
		if err == nil {
			break
		}
	}

	log.Info("connected to postgres")
	code = m.Run()
}

func TestCRUD(t *testing.T) {
	assert.NotNil(t, ds, "Datastore must not be nil")
	ctx := context.Background()
	tcr := &v1.Tenant{
		Meta:        &v1.Meta{Id: "tenant-1"},
		Name:        "A Tenant",
		Description: "A very important Tenant",
	}

	err := ds.Create(ctx, tcr)
	assert.NoError(t, err)
	assert.NotNil(t, tcr)
	// specified id is persisted
	assert.Equal(t, "tenant-1", tcr.Meta.Id)
	// initial version is set
	assert.Equal(t, int64(0), tcr.Meta.Version)
	assert.Equal(t, "A Tenant", tcr.GetName())
	assert.Equal(t, "A very important Tenant", tcr.GetDescription())

	err = ds.Create(ctx, tcr)
	assert.EqualError(t, err, "an entity of type:tenant with the id:tenant-1 already exists")

	// get existing
	var tgr v1.Tenant
	err = ds.Get(ctx, tcr.Meta.GetId(), &tgr)
	assert.NoError(t, err)
	assert.NotNil(t, &tgr)
	assert.Equal(t, "tenant-1", tgr.Meta.Id)
	assert.Equal(t, "A Tenant", tgr.GetName())
	assert.Equal(t, "A very important Tenant", tgr.GetDescription())

	// get unknown
	var tgr2 v1.Tenant
	err = ds.Get(ctx, "unknown-id", &tgr2)
	assert.Error(t, err)
	assert.EqualError(t, err, "entity of type:tenant with id:unknown-id not found")
	assert.NotNil(t, &tgr2)

	// update without meta and id
	err = ds.Update(ctx, &tgr2)
	assert.Error(t, err)
	assert.EqualError(t, err, "update of type:tenant failed, meta is nil")

	// update with unknown id
	tcr2 := &v1.Tenant{
		Meta:        &v1.Meta{Id: "tenant-2"},
		Name:        "A second Tenant",
		Description: "A not so important Tenant",
	}
	err = ds.Update(ctx, tcr2)
	assert.Error(t, err)
	assert.EqualError(t, err, "update - no entity of type:tenant with id:tenant-2 found")

	// update name
	tcr.Name = "Important Tenant"
	err = ds.Update(ctx, tcr)
	assert.NoError(t, err)
	assert.Equal(t, "Important Tenant", tcr.GetName())

	// find existing
	var tenants []v1.Tenant
	filter := make(map[string]interface{})
	// filter["tenant->>name"] = "Important Tenant"
	filter["id"] = "tenant-1"
	_, err = ds.Find(ctx, filter, nil, &tenants)
	assert.NoError(t, err)
	assert.NotNil(t, tenants)
	assert.Len(t, tenants, 1)
	assert.Equal(t, "Important Tenant", tenants[0].Name)

	// delete existing
	err = ds.Delete(ctx, tcr)
	assert.NoError(t, err)
	assert.NotNil(t, tcr)

	// delete not existing
	err = ds.Delete(ctx, tcr)
	assert.Error(t, err)
	assert.EqualError(t, err, "entity of type:tenant with id:tenant-1 not found")

}

func TestUpdateOptimisticLock(t *testing.T) {
	assert.NotNil(t, ds, "Datastore must not be nil")
	ctx := context.Background()
	tcr := &v1.Tenant{
		Meta:        &v1.Meta{Id: "tenant-2"},
		Name:        "A Tenant",
		Description: "A very important Tenant",
	}

	err := ds.Create(ctx, tcr)
	assert.NoError(t, err)
	assert.NotNil(t, tcr)
	assert.Equal(t, int64(0), tcr.Meta.Version)
	assert.Equal(t, "A Tenant", tcr.GetName())

	// get from db
	tget := &v1.Tenant{}
	err = ds.Get(ctx, tcr.Meta.Id, tget)
	require.NoError(t, err)
	require.Equal(t, tcr.Meta.Id, tget.Meta.Id)
	assert.Equal(t, int64(0), tget.Meta.Version)

	// update instance
	tcr.Name = "updated name"
	err = ds.Update(ctx, tcr)
	assert.NoError(t, err)
	assert.NotNil(t, tcr)
	// incremented version after update
	assert.Equal(t, int64(1), tcr.Meta.Version)

	// re-read from db
	var tgr v1.Tenant
	err = ds.Get(ctx, tcr.Meta.GetId(), &tgr)
	assert.NoError(t, err)
	// version is incremented
	assert.Equal(t, int64(1), tgr.Meta.Version)
	// updated data is reflected
	assert.Equal(t, "updated name", tgr.GetName())

	// try to update older version --> optimistic lock error
	tget.Name = "updated older entity"
	err = ds.Update(ctx, tget)
	require.Equal(t, err, NewOptimisticLockError(fmt.Sprintf("optimistic lock error updating tenant with id %s, existing version 1 mismatches entity version 0", tget.GetMeta().Id)))
}

func TestCreate(t *testing.T) {
	const t1 = "t1"
	assert.NotNil(t, ds, "Datastore must not be nil")
	ctx := context.Background()
	// unregistered type

	ive := &invalidVersionedEntity{}
	err := ds.Create(ctx, ive)
	assert.Error(t, err)
	assert.EqualError(t, err, "type:invalid is not registered")

	tcr1 := &v1.Tenant{
		Name:        "atenant",
		Description: "A Tenant",
	}

	// meta is nil
	err = ds.Create(ctx, tcr1)
	assert.Error(t, err)
	assert.EqualError(t, err, "create of type:tenant failed, meta is nil")

	// valid entity
	tcr1 = &v1.Tenant{
		Meta:        &v1.Meta{Id: t1},
		Name:        "atenant",
		Description: "A Tenant",
	}
	err = ds.Create(ctx, tcr1)
	assert.NoError(t, err)
	// specified id is persisted
	assert.Equal(t, t1, tcr1.Meta.Id)
	// initial version is set
	assert.Equal(t, int64(0), tcr1.Meta.Version)
	assert.Equal(t, "atenant", tcr1.GetName())
	assert.Equal(t, "A Tenant", tcr1.GetDescription())

	// create with same id
	tcr2 := &v1.Tenant{
		Meta:        &v1.Meta{Id: t1},
		Name:        "btenant",
		Description: "B Tenant",
	}
	err = ds.Create(ctx, tcr2)
	assert.Error(t, err)
	assert.EqualError(t, err, "an entity of type:tenant with the id:t1 already exists")

	// create with empty id
	tcr3 := &v1.Tenant{
		Meta:        &v1.Meta{},
		Name:        "ctenant",
		Description: "C Tenant",
	}
	err = ds.Create(ctx, tcr3)
	assert.NoError(t, err)
	assert.NotNil(t, tcr3.GetMeta().GetId())
	assert.NotEmpty(t, tcr3.GetMeta().GetId())
	assert.Contains(t, tcr3.GetMeta().GetId(), "-")
	assert.Len(t, tcr3.GetMeta().GetId(), 36)

	// create with empty kind and apiversion
	tcr4 := &v1.Tenant{
		Meta:        &v1.Meta{},
		Name:        "dtenant",
		Description: "D Tenant",
	}
	err = ds.Create(ctx, tcr4)
	assert.NoError(t, err)
	assert.NotNil(t, tcr3.GetMeta().GetApiversion())
	assert.NotEmpty(t, tcr3.GetMeta().GetApiversion())
	assert.Equal(t, tcr3.GetMeta().GetApiversion(), "v1")
	assert.NotNil(t, tcr3.GetMeta().GetKind())
	assert.NotEmpty(t, tcr3.GetMeta().GetKind())
	assert.Equal(t, tcr3.GetMeta().GetKind(), "Tenant")

	// create with wrong kind
	tcr5 := &v1.Tenant{
		Meta:        &v1.Meta{Kind: "Project"},
		Name:        "etenant",
		Description: "E Tenant",
	}
	err = ds.Create(ctx, tcr5)
	assert.Error(t, err)
	assert.EqualError(t, err, "create of type:tenant failed, kind is set to:Project but must be:Tenant")

	// create with wrong apiversion
	tcr6 := &v1.Tenant{
		Meta:        &v1.Meta{Apiversion: "v2"},
		Name:        "ftenant",
		Description: "F Tenant",
	}
	err = ds.Create(ctx, tcr6)
	assert.Error(t, err)
	assert.EqualError(t, err, "create of type:tenant failed, apiversion must be set to:v1")
}

func TestUpdate(t *testing.T) {
	const t3 = "t3"
	assert.NotNil(t, ds, "Datastore must not be nil")
	ctx := context.Background()
	// unregistered type

	ive := &invalidVersionedEntity{}
	err := ds.Update(ctx, ive)
	assert.Error(t, err)
	assert.EqualError(t, err, "type:invalid is not registered")

	// meta is nil
	tcr1 := &v1.Tenant{
		Name:        "ctenant",
		Description: "C Tenant",
	}
	err = ds.Update(ctx, tcr1)
	assert.Error(t, err)
	assert.EqualError(t, err, "update of type:tenant failed, meta is nil")

	// id is empty
	tcr1 = &v1.Tenant{
		Meta:        &v1.Meta{Id: ""},
		Name:        "ctenant",
		Description: "C Tenant",
	}
	err = ds.Update(ctx, tcr1)
	assert.Error(t, err)
	assert.EqualError(t, err, "entity of type:tenant has no id, cannot update: meta:{} name:\"ctenant\" description:\"C Tenant\"")

	// tenant with id is not found
	tcr1 = &v1.Tenant{
		Meta:        &v1.Meta{Id: t3},
		Name:        "ctenant",
		Description: "C Tenant",
	}
	err = ds.Update(ctx, tcr1)
	assert.Error(t, err)
	assert.EqualError(t, err, "update - no entity of type:tenant with id:t3 found")
	// create tenant
	err = ds.Create(ctx, tcr1)
	assert.NoError(t, err)
	assert.Equal(t, t3, tcr1.GetMeta().GetId())
	assert.Equal(t, "ctenant", tcr1.GetName())
	assert.Equal(t, "C Tenant", tcr1.GetDescription())

	tc := time.Now()
	checkHistory(ctx, t, t3, tc, "ctenant", "C Tenant")

	// now update existing
	tcr1.Description = "C Tenant 3"
	err = ds.Update(ctx, tcr1)
	assert.NoError(t, err)
	assert.Equal(t, t3, tcr1.GetMeta().GetId())
	assert.Equal(t, "ctenant", tcr1.GetName())
	assert.Equal(t, "C Tenant 3", tcr1.GetDescription())

	tu := time.Now()
	checkHistory(ctx, t, t3, tc, "ctenant", "C Tenant")
	checkHistory(ctx, t, t3, tu, "ctenant", "C Tenant 3")

	// try update with wrong kind
	tcr1.Meta.Kind = "WrongKind"
	err = ds.Update(ctx, tcr1)
	assert.Error(t, err)
	assert.EqualError(t, err, "update of type:tenant failed, kind is set to:WrongKind but must be:Tenant")

	// try update with wrong kind
	tcr1.Meta.Kind = "Tenant"
	tcr1.Meta.Apiversion = "v2"
	err = ds.Update(ctx, tcr1)
	assert.Error(t, err)
	assert.EqualError(t, err, "update of type:tenant failed, apiversion must be set to:v1")

	checkHistory(ctx, t, t3, time.Now(), "ctenant", "C Tenant 3")
}

//nolint:unparam
func checkHistoryCreated(ctx context.Context, t *testing.T, id string, name string, desc string) {
	var tgrhc v1.Tenant
	err := ds.GetHistoryCreated(ctx, id, &tgrhc)
	assert.NoError(t, err)
	assert.Equal(t, name, tgrhc.Name)
	assert.Equal(t, desc, tgrhc.GetDescription())
}

func checkHistory(ctx context.Context, t *testing.T, id string, tm time.Time, name string, desc string) {
	var tgrh v1.Tenant
	err := ds.GetHistory(ctx, id, tm, &tgrh)
	assert.NoError(t, err)
	assert.Equal(t, name, tgrh.Name)
	assert.Equal(t, desc, tgrh.GetDescription())
}

func TestGet(t *testing.T) {
	const t4 = "t4"
	assert.NotNil(t, ds, "Datastore must not be nil")
	ctx := context.Background()
	// unregistered type

	ive := &invalidVersionedEntity{}
	err := ds.Get(ctx, "", ive)
	assert.Error(t, err)
	assert.EqualError(t, err, "type:invalid is not registered")

	// unknown id
	var tgr1 v1.Tenant
	err = ds.Get(ctx, "unknown-id", &tgr1)
	assert.Error(t, err)
	assert.EqualError(t, err, "entity of type:tenant with id:unknown-id not found")

	// create a tenant
	tcr1 := &v1.Tenant{
		Meta:        &v1.Meta{Id: t4},
		Name:        "dtenant",
		Description: "D Tenant",
	}
	err = ds.Create(ctx, tcr1)
	assert.NoError(t, err)
	assert.Equal(t, t4, tcr1.GetMeta().GetId())
	assert.Equal(t, "dtenant", tcr1.GetName())
	assert.Equal(t, "D Tenant", tcr1.GetDescription())

	// now get it
	var tgr2 v1.Tenant
	err = ds.Get(ctx, t4, &tgr2)
	assert.NoError(t, err)
	assert.Equal(t, t4, tgr2.GetMeta().GetId())
	assert.Equal(t, "dtenant", tgr2.GetName())
	assert.Equal(t, "D Tenant", tgr2.GetDescription())
}

func TestGetHistory(t *testing.T) {
	const t5 = "t5"
	assert.NotNil(t, ds, "Datastore must not be nil")
	ctx := context.Background()

	tsNow := time.Date(2020, 4, 30, 18, 0, 0, 0, time.UTC)

	ive := &invalidVersionedEntity{}
	err := ds.GetHistory(ctx, "", tsNow, ive)
	assert.Error(t, err)
	assert.EqualError(t, err, "type:invalid is not registered")

	err = ds.GetHistoryCreated(ctx, "", ive)
	assert.Error(t, err)
	assert.EqualError(t, err, "type:invalid is not registered")

	// unknown id
	var tgr1 v1.Tenant
	err = ds.GetHistory(ctx, "unknown-id", tsNow, &tgr1)
	assert.Error(t, err)
	assert.EqualError(t, err, "entity of type:tenant with predicate:[map[id:unknown-id] map[created_at:2020-04-30 18:00:00 +0000 UTC]] not found")

	err = ds.GetHistoryCreated(ctx, "unknown-id", &tgr1)
	assert.Error(t, err)
	assert.EqualError(t, err, "entity of type:tenant with predicate:[map[id:unknown-id] map[op:C]] not found")

	// control time.Now()
	createTS := time.Date(2020, 4, 30, 18, 0, 0, 0, time.UTC)
	setNow(createTS)
	defer resetNow()

	var tgrH v1.Tenant
	err = ds.GetHistory(ctx, t5, createTS, &tgrH)
	assert.Error(t, err)
	assert.EqualError(t, err, "entity of type:tenant with predicate:[map[id:t5] map[created_at:2020-04-30 18:00:00 +0000 UTC]] not found")

	err = ds.GetHistoryCreated(ctx, t5, &tgrH)
	assert.Error(t, err)
	assert.EqualError(t, err, "entity of type:tenant with predicate:[map[id:t5] map[op:C]] not found")

	// create a tenant
	tcr1 := &v1.Tenant{
		Meta:        &v1.Meta{Id: t5},
		Name:        "dtenant",
		Description: "D Tenant",
	}
	err = ds.Create(ctx, tcr1)
	assert.NoError(t, err)

	checkHistoryCreated(ctx, t, t5, "dtenant", "D Tenant")
	checkHistory(ctx, t, t5, createTS, "dtenant", "D Tenant")

	var tcrU v1.Tenant
	err = ds.Get(ctx, t5, &tcrU)
	assert.NoError(t, err)
	assert.Equal(t, createTS, convertToTime(tcrU.Meta.CreatedTime))

	updateTS := time.Date(2020, 4, 30, 20, 0, 0, 0, time.UTC)
	setNow(updateTS)
	tcrU.Name = "dtenant updated"
	err = ds.Update(ctx, &tcrU)
	assert.NoError(t, err)
	assert.Equal(t, updateTS, convertToTime(tcrU.Meta.UpdatedTime))

	checkHistoryCreated(ctx, t, t5, "dtenant", "D Tenant")
	checkHistory(ctx, t, t5, updateTS, "dtenant updated", "D Tenant")

	update2TS := time.Date(2020, 4, 30, 21, 0, 0, 0, time.UTC)
	setNow(update2TS)
	tcrU.Name = "dtenant updated 2"
	err = ds.Update(ctx, &tcrU)
	assert.NoError(t, err)
	assert.Equal(t, update2TS, convertToTime(tcrU.Meta.UpdatedTime))

	checkHistoryCreated(ctx, t, t5, "dtenant", "D Tenant")
	checkHistory(ctx, t, t5, update2TS, "dtenant updated 2", "D Tenant")

	deletedTS := time.Date(2020, 4, 30, 22, 0, 0, 0, time.UTC)
	setNow(deletedTS)
	err = ds.Delete(ctx, tcr1)
	assert.NoError(t, err)

	checkHistoryCreated(ctx, t, t5, "dtenant", "D Tenant")
	checkHistory(ctx, t, t5, deletedTS, "dtenant updated 2", "D Tenant")

	// Check complete history
	// before create
	err = ds.GetHistory(ctx, t5, time.Date(2019, 1, 1, 8, 0, 0, 0, time.UTC), &tgrH)
	assert.Error(t, err)
	assert.EqualError(t, err, "entity of type:tenant with predicate:[map[id:t5] map[created_at:2019-01-01 08:00:00 +0000 UTC]] not found")

	checkHistoryCreated(ctx, t, t5, "dtenant", "D Tenant")
	checkHistory(ctx, t, t5, createTS, "dtenant", "D Tenant")
	checkHistory(ctx, t, t5, updateTS, "dtenant updated", "D Tenant")
	checkHistory(ctx, t, t5, update2TS, "dtenant updated 2", "D Tenant")
	checkHistory(ctx, t, t5, deletedTS, "dtenant updated 2", "D Tenant")
	checkHistory(ctx, t, t5, time.Date(2021, 1, 1, 8, 0, 0, 0, time.UTC), "dtenant updated 2", "D Tenant")
}

func TestFind(t *testing.T) {
	const t6 = "t6"
	assert.NotNil(t, ds, "Datastore must not be nil")
	ctx := context.Background()
	// result not a slice
	var te v1.Tenant
	_, err := ds.Find(ctx, nil, nil, &te)
	assert.Error(t, err)
	assert.EqualError(t, err, "result argument must be a slice address")

	// result is no versionedEntity
	var res []string
	_, err = ds.Find(ctx, nil, nil, &res)
	assert.Error(t, err)
	assert.EqualError(t, err, "result slice element type must implement VersionedJSONEntity-Interface")

	// unregistered type
	var ives = []invalidVersionedEntity{}
	_, err = ds.Find(ctx, nil, nil, &ives)
	assert.Error(t, err)
	assert.EqualError(t, err, "type:invalid is not registered")

	// create a tenant
	tcr1 := &v1.Tenant{
		Meta:        &v1.Meta{Id: t6},
		Name:        "ftenant",
		Description: "F Tenant",
	}
	err = ds.Create(ctx, tcr1)
	assert.NoError(t, err)
	assert.Equal(t, t6, tcr1.GetMeta().GetId())
	assert.Equal(t, "ftenant", tcr1.GetName())
	assert.Equal(t, "F Tenant", tcr1.GetDescription())

	// now search it
	var tfr []v1.Tenant
	filter := make(map[string]interface{})
	filter["id"] = t6
	_, err = ds.Find(ctx, filter, nil, &tfr)
	assert.NoError(t, err)
	assert.NotNil(t, tfr)
	assert.Len(t, tfr, 1)

	// create more
	for i := 0; i < 10; i++ {
		tcr := &v1.Tenant{
			Meta:        &v1.Meta{Id: fmt.Sprintf("ftenant-%d", i)},
			Name:        fmt.Sprintf("tenant-%d", i),
			Description: fmt.Sprintf("Tenant %d", i),
		}
		err = ds.Create(ctx, tcr)
		assert.NoError(t, err)
	}
	// find all
	filter = make(map[string]interface{})
	_, err = ds.Find(ctx, filter, nil, &tfr)
	assert.NoError(t, err)
	assert.NotNil(t, tfr)

	// find one
	var t9 []v1.Tenant
	filter["id"] = "ftenant-9"
	_, err = ds.Find(ctx, filter, nil, &t9)
	assert.NoError(t, err)
	assert.NotNil(t, t9)
	assert.Len(t, t9, 1)

	// find one by name
	var t8 []v1.Tenant
	filter = make(map[string]interface{})
	filter["tenant ->> 'name'"] = "tenant-8"
	_, err = ds.Find(ctx, filter, nil, &t8)
	assert.NoError(t, err)
	assert.NotNil(t, t8)
	assert.Len(t, t8, 1)

	// find one by description
	var t4 []v1.Tenant
	filter = make(map[string]interface{})
	filter["tenant ->> 'description'"] = "Tenant 4"
	_, err = ds.Find(ctx, filter, nil, &t4)
	assert.NoError(t, err)
	assert.NotNil(t, t4)
	assert.Len(t, t4, 1)

}

func TestDelete(t *testing.T) {
	const t9 = "t9"
	assert.NotNil(t, ds, "Datastore must not be nil")
	ctx := context.Background()
	// unregistered type

	ive := &invalidVersionedEntity{}
	err := ds.Delete(ctx, ive)
	assert.Error(t, err)
	assert.EqualError(t, err, "type:invalid is not registered")

	// unknown id
	tdr1 := &v1.Tenant{
		Meta: &v1.Meta{Id: "unknown-id"},
	}
	err = ds.Delete(ctx, tdr1)
	assert.Error(t, err)
	assert.EqualError(t, err, "entity of type:tenant with id:unknown-id not found")

	// create a tenant
	tcr1 := &v1.Tenant{
		Meta:        &v1.Meta{Id: t9},
		Name:        "etenant",
		Description: "E Tenant",
	}
	err = ds.Create(ctx, tcr1)
	assert.NoError(t, err)
	assert.Equal(t, t9, tcr1.GetMeta().GetId())
	assert.Equal(t, "etenant", tcr1.GetName())
	assert.Equal(t, "E Tenant", tcr1.GetDescription())

	// now delete it
	tdr2 := &v1.Tenant{
		Meta: &v1.Meta{Id: t9},
	}
	err = ds.Delete(ctx, tdr2)
	assert.NoError(t, err)

	var tgr v1.Tenant
	err = ds.Get(ctx, t9, &tgr)
	assert.Error(t, err)
	assert.EqualError(t, err, "entity of type:tenant with id:t9 not found")

	var tgrh v1.Tenant
	err = ds.GetHistory(ctx, t9, time.Now(), &tgrh)
	assert.NoError(t, err)
	assert.Equal(t, "etenant", tgrh.Name)
}

func TestAnnotationsAndLabels(t *testing.T) {
	assert.NotNil(t, ds, "Datastore must not be nil")
	ctx := context.Background()
	tcr := &v1.Tenant{
		Meta: &v1.Meta{
			Id: "tenant-3",
			Annotations: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			Labels: []string{
				"color=red",
			},
		},
		Name:        "A Tenant",
		Description: "A very important Tenant",
	}

	err := ds.Create(ctx, tcr)
	assert.NoError(t, err)
	assert.NotNil(t, tcr)
	assert.Equal(t, int64(0), tcr.Meta.Version)
	assert.Equal(t, "A Tenant", tcr.GetName())
	assert.NotNil(t, tcr.GetMeta().GetAnnotations())
	assert.NotNil(t, tcr.GetMeta().GetLabels())
	assert.Equal(t, map[string]string{"key1": "value1", "key2": "value2"}, tcr.GetMeta().GetAnnotations())
	assert.Equal(t, []string{"color=red"}, tcr.GetMeta().GetLabels())

	// get from db
	tget := &v1.Tenant{}
	err = ds.Get(ctx, tcr.Meta.Id, tget)
	require.NoError(t, err)
	require.Equal(t, tcr.Meta.Id, tget.Meta.Id)
	assert.Equal(t, int64(0), tget.Meta.Version)

	// update instance
	tcr.Name = "updated name"
	err = ds.Update(ctx, tcr)
	assert.NoError(t, err)
	assert.NotNil(t, tcr)
	// incremented version after update
	assert.Equal(t, int64(1), tcr.Meta.Version)

	// re-read from db
	var tgr v1.Tenant
	err = ds.Get(ctx, tcr.Meta.GetId(), &tgr)
	assert.NoError(t, err)
	// version is incremented
	assert.Equal(t, int64(1), tgr.Meta.Version)
	// updated data is reflected
	assert.Equal(t, "updated name", tgr.GetName())

	// try to update older version --> optimistic lock error
	tget.Name = "updated older entity"
	err = ds.Update(ctx, tget)
	require.Equal(t, err,
		NewOptimisticLockError(
			fmt.Sprintf("optimistic lock error updating tenant with id %s, existing version 1 mismatches entity version 0", tget.GetMeta().Id),
		),
	)

	// update annotations and labels
	as := tcr.GetMeta().GetAnnotations()
	as["key1"] = "value3"
	ls := []string{"color=red", "size=xlarge"}
	tcr.GetMeta().SetAnnotations(as)
	tcr.GetMeta().SetLabels(ls)
	err = ds.Update(ctx, tcr)
	assert.NoError(t, err)
	assert.NotNil(t, tcr)
	assert.Equal(t, map[string]string{"key1": "value3", "key2": "value2"}, tcr.GetMeta().GetAnnotations())
	assert.Equal(t, []string{"color=red", "size=xlarge"}, tcr.GetMeta().GetLabels())
}

func convertToTime(pbTs *timestamppb.Timestamp) time.Time {
	return pbTs.AsTime()
}

// setNow sets Now
func setNow(t time.Time) {
	Now = func() time.Time {
		return t
	}
}

// resetNow resets the overriden Now to time.Now
func resetNow() {
	Now = time.Now
}
