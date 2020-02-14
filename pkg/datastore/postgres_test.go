package datastore

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"runtime/debug"
	"testing"
	"time"

	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"go.uber.org/zap"
)

var (
	ds  *Datastore
	log *zap.Logger
)

// to test unregistered type checks
type invalidVersionedEntity struct{}

func (v *invalidVersionedEntity) JSONField() string { return "invalid" }
func (v *invalidVersionedEntity) TableName() string { return "" }
func (v *invalidVersionedEntity) Schema() string    { return "" }
func (v *invalidVersionedEntity) GetMeta() *v1.Meta { return nil }

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
		ds, err = NewPostgresStorage(log, ip, port.Port(), "postgres", "password", "postgres", "disable", &v1.Project{}, &v1.Tenant{})
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			// log.Sugar().Errorw("cannot connect to postgres", "err", err)
			continue
		}
		err = ds.db.Ping()
		if err != nil {
			log.Sugar().Errorw("Could not connect to postgres server", "err", err)
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

	// get existing
	var tgr v1.Tenant
	err = ds.Get(ctx, tcr.Meta.GetId(), &tgr)
	assert.NoError(t, err)
	assert.NotNil(t, tgr)
	assert.Equal(t, "tenant-1", tgr.Meta.Id)
	assert.Equal(t, "A Tenant", tgr.GetName())
	assert.Equal(t, "A very important Tenant", tgr.GetDescription())

	// get unknown
	var tgr2 v1.Tenant
	err = ds.Get(ctx, "unknown-id", &tgr2)
	assert.Error(t, err)
	assert.EqualError(t, err, "entity of type:tenant with id:unknown-id not found")
	assert.NotNil(t, tgr2)

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
	err = ds.Find(ctx, filter, &tenants)
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
	assert.EqualError(t, err, "not found: delete of tenant with id tenant-1 affected 0 rows")

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
		Meta:        &v1.Meta{Id: "t1"},
		Name:        "atenant",
		Description: "A Tenant",
	}
	err = ds.Create(ctx, tcr1)
	assert.NoError(t, err)
	// specified id is persisted
	assert.Equal(t, "t1", tcr1.Meta.Id)
	// initial version is set
	assert.Equal(t, int64(0), tcr1.Meta.Version)
	assert.Equal(t, "atenant", tcr1.GetName())
	assert.Equal(t, "A Tenant", tcr1.GetDescription())

	// create with same id
	tcr2 := &v1.Tenant{
		Meta:        &v1.Meta{Id: "t1"},
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

}

func TestUpdate(t *testing.T) {
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
	assert.EqualError(t, err, "entity of type:tenant has no id, cannot update: meta:<> name:\"ctenant\" description:\"C Tenant\" ")

	// tenant with id is not found
	tcr1 = &v1.Tenant{
		Meta:        &v1.Meta{Id: "t3"},
		Name:        "ctenant",
		Description: "C Tenant",
	}
	err = ds.Update(ctx, tcr1)
	assert.Error(t, err)
	assert.EqualError(t, err, "update - no entity of type:tenant with id:t3 found")
	// create tenant
	err = ds.Create(ctx, tcr1)
	assert.NoError(t, err)
	assert.Equal(t, "t3", tcr1.GetMeta().GetId())
	assert.Equal(t, "ctenant", tcr1.GetName())
	assert.Equal(t, "C Tenant", tcr1.GetDescription())
	// now update existing
	tcr1.Description = "C Tenant 3"
	err = ds.Update(ctx, tcr1)
	assert.NoError(t, err)
	assert.Equal(t, "t3", tcr1.GetMeta().GetId())
	assert.Equal(t, "ctenant", tcr1.GetName())
	assert.Equal(t, "C Tenant 3", tcr1.GetDescription())

}

func TestGet(t *testing.T) {
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
		Meta:        &v1.Meta{Id: "t4"},
		Name:        "dtenant",
		Description: "D Tenant",
	}
	err = ds.Create(ctx, tcr1)
	assert.NoError(t, err)
	assert.Equal(t, "t4", tcr1.GetMeta().GetId())
	assert.Equal(t, "dtenant", tcr1.GetName())
	assert.Equal(t, "D Tenant", tcr1.GetDescription())

	// now get it
	var tgr2 v1.Tenant
	err = ds.Get(ctx, "t4", &tgr2)
	assert.NoError(t, err)
	assert.Equal(t, "t4", tgr2.GetMeta().GetId())
	assert.Equal(t, "dtenant", tgr2.GetName())
	assert.Equal(t, "D Tenant", tgr2.GetDescription())

}

func TestFind(t *testing.T) {
	assert.NotNil(t, ds, "Datastore must not be nil")
	ctx := context.Background()
	// result not a slice
	var te v1.Tenant
	err := ds.Find(ctx, nil, te)
	assert.Error(t, err)
	assert.EqualError(t, err, "result argument must be a slice address")

	// result is no versionedEntity
	var res []string
	err = ds.Find(ctx, nil, &res)
	assert.Error(t, err)
	assert.EqualError(t, err, "result slice element type must implement VersionedJSONEntity-Interface")

	// unregistered type
	var ives = []invalidVersionedEntity{}
	err = ds.Find(ctx, nil, &ives)
	assert.Error(t, err)
	assert.EqualError(t, err, "type:invalid is not registered")

	// create a tenant
	tcr1 := &v1.Tenant{
		Meta:        &v1.Meta{Id: "t6"},
		Name:        "ftenant",
		Description: "F Tenant",
	}
	err = ds.Create(ctx, tcr1)
	assert.NoError(t, err)
	assert.Equal(t, "t6", tcr1.GetMeta().GetId())
	assert.Equal(t, "ftenant", tcr1.GetName())
	assert.Equal(t, "F Tenant", tcr1.GetDescription())

	// now search it
	var tfr []v1.Tenant
	filter := make(map[string]interface{})
	filter["id"] = "t6"
	err = ds.Find(ctx, filter, &tfr)
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
	err = ds.Find(ctx, filter, &tfr)
	assert.NoError(t, err)
	assert.NotNil(t, tfr)

	// find one
	var t9 []v1.Tenant
	filter["id"] = "ftenant-9"
	err = ds.Find(ctx, filter, &t9)
	assert.NoError(t, err)
	assert.NotNil(t, t9)
	assert.Len(t, t9, 1)

	// find one by name
	var t8 []v1.Tenant
	filter = make(map[string]interface{})
	filter["tenant ->> 'name'"] = "tenant-8"
	err = ds.Find(ctx, filter, &t8)
	assert.NoError(t, err)
	assert.NotNil(t, t8)
	assert.Len(t, t8, 1)

	// find one by description
	var t4 []v1.Tenant
	filter = make(map[string]interface{})
	filter["tenant ->> 'description'"] = "Tenant 4"
	err = ds.Find(ctx, filter, &t4)
	assert.NoError(t, err)
	assert.NotNil(t, t4)
	assert.Len(t, t4, 1)

}

func TestDelete(t *testing.T) {
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
	assert.EqualError(t, err, "not found: delete of tenant with id unknown-id affected 0 rows")

	// create a tenant
	tcr1 := &v1.Tenant{
		Meta:        &v1.Meta{Id: "t5"},
		Name:        "etenant",
		Description: "E Tenant",
	}
	err = ds.Create(ctx, tcr1)
	assert.NoError(t, err)
	assert.Equal(t, "t5", tcr1.GetMeta().GetId())
	assert.Equal(t, "etenant", tcr1.GetName())
	assert.Equal(t, "E Tenant", tcr1.GetDescription())

	// now delete it
	tdr2 := &v1.Tenant{
		Meta: &v1.Meta{Id: "t5"},
	}
	err = ds.Delete(ctx, tdr2)
	assert.NoError(t, err)

	var tgr v1.Tenant
	err = ds.Get(ctx, "t5", &tgr)
	assert.Error(t, err)
	assert.EqualError(t, err, "entity of type:tenant with id:t5 not found")

}
