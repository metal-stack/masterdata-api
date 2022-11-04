package datastore

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/stretchr/testify/require"

	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

var db *sqlx.DB

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

	var err error
	db, err = createPostgresConnection()
	if err != nil {
		panic(err)
	}
	code = m.Run()
}

func TestCRUD(t *testing.T) {
	tenantDS, err := NewPostgresStorage(zaptest.NewLogger(t), db, &v1.Tenant{})
	require.NoError(t, err)
	assert.NotNil(t, tenantDS, "Datastore must not be nil")
	ctx := context.Background()
	tcr := &v1.Tenant{
		Meta:        &v1.Meta{Id: "tenant-1"},
		Name:        "A Tenant",
		Description: "A very important Tenant",
	}

	err = tenantDS.Create(ctx, tcr)
	assert.NoError(t, err)
	assert.NotNil(t, tcr)
	// specified id is persisted
	assert.Equal(t, "tenant-1", tcr.Meta.Id)
	// initial version is set
	assert.Equal(t, int64(0), tcr.Meta.Version)
	assert.Equal(t, "A Tenant", tcr.GetName())
	assert.Equal(t, "A very important Tenant", tcr.GetDescription())

	err = tenantDS.Create(ctx, tcr)
	assert.EqualError(t, err, "an entity of type:tenant with the id:tenant-1 already exists")

	// get existing
	tgr, err := tenantDS.Get(ctx, tcr.Meta.GetId())
	assert.NoError(t, err)
	assert.NotNil(t, &tgr)
	assert.Equal(t, "tenant-1", tgr.Meta.Id)
	assert.Equal(t, "A Tenant", tgr.GetName())
	assert.Equal(t, "A very important Tenant", tgr.GetDescription())

	// get unknown
	tgr2, err := tenantDS.Get(ctx, "unknown-id")
	assert.Error(t, err)
	assert.EqualError(t, err, "entity of type:tenant with id:unknown-id not found")
	assert.NotNil(t, &tgr2)

	// update without meta and id
	err = tenantDS.Update(ctx, tgr2)
	assert.Error(t, err)
	assert.EqualError(t, err, "update of type:tenant failed, meta is nil")

	// update with unknown id
	tcr2 := &v1.Tenant{
		Meta:        &v1.Meta{Id: "tenant-2"},
		Name:        "A second Tenant",
		Description: "A not so important Tenant",
	}
	err = tenantDS.Update(ctx, tcr2)
	assert.Error(t, err)
	assert.EqualError(t, err, "update - no entity of type:tenant with id:tenant-2 found")

	// update name
	tcr.Name = "Important Tenant"
	err = tenantDS.Update(ctx, tcr)
	assert.NoError(t, err)
	assert.Equal(t, "Important Tenant", tcr.GetName())

	// find existing
	filter := make(map[string]any)
	// filter["tenant->>name"] = "Important Tenant"
	filter["id"] = "tenant-1"
	tenants, _, err := tenantDS.Find(ctx, filter, nil)
	assert.NoError(t, err)
	assert.NotNil(t, tenants)
	assert.Len(t, tenants, 1)
	assert.Equal(t, "Important Tenant", tenants[0].Name)

	// delete existing
	err = tenantDS.Delete(ctx, tcr)
	assert.NoError(t, err)
	assert.NotNil(t, tcr)

	// delete not existing
	err = tenantDS.Delete(ctx, tcr)
	assert.Error(t, err)
	assert.EqualError(t, err, "entity of type:tenant with id:tenant-1 not found")

}

func TestUpdateOptimisticLock(t *testing.T) {
	tenantDS, err := NewPostgresStorage(zaptest.NewLogger(t), db, &v1.Tenant{})
	require.NoError(t, err)
	assert.NotNil(t, tenantDS, "Datastore must not be nil")
	ctx := context.Background()
	tcr := &v1.Tenant{
		Meta:        &v1.Meta{Id: "tenant-2"},
		Name:        "A Tenant",
		Description: "A very important Tenant",
	}

	err = tenantDS.Create(ctx, tcr)
	assert.NoError(t, err)
	assert.NotNil(t, tcr)
	assert.Equal(t, int64(0), tcr.Meta.Version)
	assert.Equal(t, "A Tenant", tcr.GetName())

	// get from db
	tget, err := tenantDS.Get(ctx, tcr.Meta.Id)
	require.NoError(t, err)
	require.Equal(t, tcr.Meta.Id, tget.Meta.Id)
	assert.Equal(t, int64(0), tget.Meta.Version)

	// update instance
	tcr.Name = "updated name"
	err = tenantDS.Update(ctx, tcr)
	assert.NoError(t, err)
	assert.NotNil(t, tcr)
	// incremented version after update
	assert.Equal(t, int64(1), tcr.Meta.Version)

	// re-read from db
	tgr, err := tenantDS.Get(ctx, tcr.Meta.GetId())
	assert.NoError(t, err)
	// version is incremented
	assert.Equal(t, int64(1), tgr.Meta.Version)
	// updated data is reflected
	assert.Equal(t, "updated name", tgr.GetName())

	// try to update older version --> optimistic lock error
	tget.Name = "updated older entity"
	err = tenantDS.Update(ctx, tget)
	require.Equal(t, err, NewOptimisticLockError(fmt.Sprintf("optimistic lock error updating tenant with id %s, existing version 1 mismatches entity version 0", tget.GetMeta().Id)))
}

func TestCreate(t *testing.T) {
	const t1 = "t1"
	tenantDS, err := NewPostgresStorage(zaptest.NewLogger(t), db, &v1.Tenant{})
	require.NoError(t, err)
	assert.NotNil(t, tenantDS, "Datastore must not be nil")
	ctx := context.Background()

	tcr1 := &v1.Tenant{
		Name:        "atenant",
		Description: "A Tenant",
	}

	// meta is nil
	err = tenantDS.Create(ctx, tcr1)
	assert.Error(t, err)
	assert.EqualError(t, err, "create of type:tenant failed, meta is nil")

	// valid entity
	tcr1 = &v1.Tenant{
		Meta:        &v1.Meta{Id: t1},
		Name:        "atenant",
		Description: "A Tenant",
	}
	err = tenantDS.Create(ctx, tcr1)
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
	err = tenantDS.Create(ctx, tcr2)
	assert.Error(t, err)
	assert.EqualError(t, err, "an entity of type:tenant with the id:t1 already exists")

	// create with empty id
	tcr3 := &v1.Tenant{
		Meta:        &v1.Meta{},
		Name:        "ctenant",
		Description: "C Tenant",
	}
	err = tenantDS.Create(ctx, tcr3)
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
	err = tenantDS.Create(ctx, tcr4)
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
	err = tenantDS.Create(ctx, tcr5)
	assert.Error(t, err)
	assert.EqualError(t, err, "create of type:tenant failed, kind is set to:Project but must be:Tenant")

	// create with wrong apiversion
	tcr6 := &v1.Tenant{
		Meta:        &v1.Meta{Apiversion: "v2"},
		Name:        "ftenant",
		Description: "F Tenant",
	}
	err = tenantDS.Create(ctx, tcr6)
	assert.Error(t, err)
	assert.EqualError(t, err, "create of type:tenant failed, apiversion must be set to:v1")
}

func TestUpdate(t *testing.T) {
	const t3 = "t3"
	tenantDS, err := NewPostgresStorage(zaptest.NewLogger(t), db, &v1.Tenant{})
	require.NoError(t, err)
	assert.NotNil(t, tenantDS, "Datastore must not be nil")
	ctx := context.Background()

	// meta is nil
	tcr1 := &v1.Tenant{
		Name:        "ctenant",
		Description: "C Tenant",
	}
	err = tenantDS.Update(ctx, tcr1)
	assert.Error(t, err)
	assert.EqualError(t, err, "update of type:tenant failed, meta is nil")

	// id is empty
	tcr1 = &v1.Tenant{
		Meta:        &v1.Meta{Id: ""},
		Name:        "ctenant",
		Description: "C Tenant",
	}
	err = tenantDS.Update(ctx, tcr1)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "entity of type:tenant has no id, cannot update: meta:{}")

	// tenant with id is not found
	tcr1 = &v1.Tenant{
		Meta:        &v1.Meta{Id: t3},
		Name:        "ctenant",
		Description: "C Tenant",
	}
	err = tenantDS.Update(ctx, tcr1)
	assert.Error(t, err)
	assert.EqualError(t, err, "update - no entity of type:tenant with id:t3 found")
	// create tenant
	err = tenantDS.Create(ctx, tcr1)
	assert.NoError(t, err)
	assert.Equal(t, t3, tcr1.GetMeta().GetId())
	assert.Equal(t, "ctenant", tcr1.GetName())
	assert.Equal(t, "C Tenant", tcr1.GetDescription())

	tc := time.Now()
	checkHistory(ctx, t, t3, tc, "ctenant", "C Tenant")

	// now update existing
	tcr1.Description = "C Tenant 3"
	err = tenantDS.Update(ctx, tcr1)
	assert.NoError(t, err)
	assert.Equal(t, t3, tcr1.GetMeta().GetId())
	assert.Equal(t, "ctenant", tcr1.GetName())
	assert.Equal(t, "C Tenant 3", tcr1.GetDescription())

	tu := time.Now()
	checkHistory(ctx, t, t3, tc, "ctenant", "C Tenant")
	checkHistory(ctx, t, t3, tu, "ctenant", "C Tenant 3")

	// try update with wrong kind
	tcr1.Meta.Kind = "WrongKind"
	err = tenantDS.Update(ctx, tcr1)
	assert.Error(t, err)
	assert.EqualError(t, err, "update of type:tenant failed, kind is set to:WrongKind but must be:Tenant")

	// try update with wrong kind
	tcr1.Meta.Kind = "Tenant"
	tcr1.Meta.Apiversion = "v2"
	err = tenantDS.Update(ctx, tcr1)
	assert.Error(t, err)
	assert.EqualError(t, err, "update of type:tenant failed, apiversion must be set to:v1")

	checkHistory(ctx, t, t3, time.Now(), "ctenant", "C Tenant 3")
}

//nolint:unparam
func checkHistoryCreated(ctx context.Context, t *testing.T, id string, name string, desc string) {
	tenantDS, err := NewPostgresStorage(zaptest.NewLogger(t), db, &v1.Tenant{})
	require.NoError(t, err)
	var tgrhc v1.Tenant
	err = tenantDS.GetHistoryCreated(ctx, id, &tgrhc)
	assert.NoError(t, err)
	assert.Equal(t, name, tgrhc.Name)
	assert.Equal(t, desc, tgrhc.GetDescription())
}

func checkHistory(ctx context.Context, t *testing.T, id string, tm time.Time, name string, desc string) {
	tenantDS, err := NewPostgresStorage(zaptest.NewLogger(t), db, &v1.Tenant{})
	require.NoError(t, err)
	var tgrh v1.Tenant
	err = tenantDS.GetHistory(ctx, id, tm, &tgrh)
	assert.NoError(t, err)
	assert.Equal(t, name, tgrh.Name)
	assert.Equal(t, desc, tgrh.GetDescription())
}

func TestGet(t *testing.T) {
	const t4 = "t4"
	tenantDS, err := NewPostgresStorage(zaptest.NewLogger(t), db, &v1.Tenant{})
	require.NoError(t, err)
	assert.NotNil(t, tenantDS, "Datastore must not be nil")
	ctx := context.Background()
	// unknown id
	_, err = tenantDS.Get(ctx, "unknown-id")
	assert.Error(t, err)
	assert.EqualError(t, err, "entity of type:tenant with id:unknown-id not found")

	// create a tenant
	tcr1 := &v1.Tenant{
		Meta:        &v1.Meta{Id: t4},
		Name:        "dtenant",
		Description: "D Tenant",
	}
	err = tenantDS.Create(ctx, tcr1)
	assert.NoError(t, err)
	assert.Equal(t, t4, tcr1.GetMeta().GetId())
	assert.Equal(t, "dtenant", tcr1.GetName())
	assert.Equal(t, "D Tenant", tcr1.GetDescription())

	// now get it
	tgr2, err := tenantDS.Get(ctx, t4)
	assert.NoError(t, err)
	assert.Equal(t, t4, tgr2.GetMeta().GetId())
	assert.Equal(t, "dtenant", tgr2.GetName())
	assert.Equal(t, "D Tenant", tgr2.GetDescription())
}

func TestGetHistory(t *testing.T) {
	const t5 = "t5"
	tenantDS, err := NewPostgresStorage(zaptest.NewLogger(t), db, &v1.Tenant{})
	require.NoError(t, err)
	assert.NotNil(t, tenantDS, "Datastore must not be nil")
	ctx := context.Background()

	tsNow := time.Date(2020, 4, 30, 18, 0, 0, 0, time.UTC)

	// unknown id
	var tgr1 v1.Tenant
	err = tenantDS.GetHistory(ctx, "unknown-id", tsNow, &tgr1)
	assert.Error(t, err)
	assert.EqualError(t, err, "entity of type:tenant with predicate:[map[id:unknown-id] map[created_at:2020-04-30 18:00:00 +0000 UTC]] not found")

	err = tenantDS.GetHistoryCreated(ctx, "unknown-id", &tgr1)
	assert.Error(t, err)
	assert.EqualError(t, err, "entity of type:tenant with predicate:[map[id:unknown-id] map[op:C]] not found")

	// control time.Now()
	createTS := time.Date(2020, 4, 30, 18, 0, 0, 0, time.UTC)
	setNow(createTS)
	defer resetNow()

	var tgrH v1.Tenant
	err = tenantDS.GetHistory(ctx, t5, createTS, &tgrH)
	assert.Error(t, err)
	assert.EqualError(t, err, "entity of type:tenant with predicate:[map[id:t5] map[created_at:2020-04-30 18:00:00 +0000 UTC]] not found")

	err = tenantDS.GetHistoryCreated(ctx, t5, &tgrH)
	assert.Error(t, err)
	assert.EqualError(t, err, "entity of type:tenant with predicate:[map[id:t5] map[op:C]] not found")

	// create a tenant
	tcr1 := &v1.Tenant{
		Meta:        &v1.Meta{Id: t5},
		Name:        "dtenant",
		Description: "D Tenant",
	}
	err = tenantDS.Create(ctx, tcr1)
	assert.NoError(t, err)

	checkHistoryCreated(ctx, t, t5, "dtenant", "D Tenant")
	checkHistory(ctx, t, t5, createTS, "dtenant", "D Tenant")

	tcrU, err := tenantDS.Get(ctx, t5)
	assert.NoError(t, err)
	assert.Equal(t, createTS, convertToTime(tcrU.Meta.CreatedTime))

	updateTS := time.Date(2020, 4, 30, 20, 0, 0, 0, time.UTC)
	setNow(updateTS)
	tcrU.Name = "dtenant updated"
	err = tenantDS.Update(ctx, tcrU)
	assert.NoError(t, err)
	assert.Equal(t, updateTS, convertToTime(tcrU.Meta.UpdatedTime))

	checkHistoryCreated(ctx, t, t5, "dtenant", "D Tenant")
	checkHistory(ctx, t, t5, updateTS, "dtenant updated", "D Tenant")

	update2TS := time.Date(2020, 4, 30, 21, 0, 0, 0, time.UTC)
	setNow(update2TS)
	tcrU.Name = "dtenant updated 2"
	err = tenantDS.Update(ctx, tcrU)
	assert.NoError(t, err)
	assert.Equal(t, update2TS, convertToTime(tcrU.Meta.UpdatedTime))

	checkHistoryCreated(ctx, t, t5, "dtenant", "D Tenant")
	checkHistory(ctx, t, t5, update2TS, "dtenant updated 2", "D Tenant")

	deletedTS := time.Date(2020, 4, 30, 22, 0, 0, 0, time.UTC)
	setNow(deletedTS)
	err = tenantDS.Delete(ctx, tcr1)
	assert.NoError(t, err)

	checkHistoryCreated(ctx, t, t5, "dtenant", "D Tenant")
	checkHistory(ctx, t, t5, deletedTS, "dtenant updated 2", "D Tenant")

	// Check complete history
	// before create
	err = tenantDS.GetHistory(ctx, t5, time.Date(2019, 1, 1, 8, 0, 0, 0, time.UTC), &tgrH)
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
	tenantDS, err := NewPostgresStorage(zaptest.NewLogger(t), db, &v1.Tenant{})
	require.NoError(t, err)
	assert.NotNil(t, tenantDS, "Datastore must not be nil")
	ctx := context.Background()

	// create a tenant
	tcr1 := &v1.Tenant{
		Meta:        &v1.Meta{Id: t6},
		Name:        "ftenant",
		Description: "F Tenant",
	}
	err = tenantDS.Create(ctx, tcr1)
	assert.NoError(t, err)
	assert.Equal(t, t6, tcr1.GetMeta().GetId())
	assert.Equal(t, "ftenant", tcr1.GetName())
	assert.Equal(t, "F Tenant", tcr1.GetDescription())

	// now search it
	filter := make(map[string]any)
	filter["id"] = t6
	tfr, _, err := tenantDS.Find(ctx, filter, nil)
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
		err = tenantDS.Create(ctx, tcr)
		assert.NoError(t, err)
	}
	// find all
	filter = make(map[string]any)
	tfr, _, err = tenantDS.Find(ctx, filter, nil)
	assert.NoError(t, err)
	assert.NotNil(t, tfr)

	// find one
	filter["id"] = "ftenant-9"
	t9, _, err := tenantDS.Find(ctx, filter, nil)
	assert.NoError(t, err)
	assert.NotNil(t, t9)
	assert.Len(t, t9, 1)

	// find one by name
	filter = make(map[string]any)
	filter["tenant ->> 'name'"] = "tenant-8"
	t8, _, err := tenantDS.Find(ctx, filter, nil)
	assert.NoError(t, err)
	assert.NotNil(t, t8)
	assert.Len(t, t8, 1)

	// find one by description
	filter = make(map[string]any)
	filter["tenant ->> 'description'"] = "Tenant 4"
	t4, _, err := tenantDS.Find(ctx, filter, nil)
	assert.NoError(t, err)
	assert.NotNil(t, t4)
	assert.Len(t, t4, 1)

}

func TestDelete(t *testing.T) {
	const t9 = "t9"
	tenantDS, err := NewPostgresStorage(zaptest.NewLogger(t), db, &v1.Tenant{})
	require.NoError(t, err)
	assert.NotNil(t, tenantDS, "Datastore must not be nil")
	ctx := context.Background()

	// unknown id
	tdr1 := &v1.Tenant{
		Meta: &v1.Meta{Id: "unknown-id"},
	}
	err = tenantDS.Delete(ctx, tdr1)
	assert.Error(t, err)
	assert.EqualError(t, err, "entity of type:tenant with id:unknown-id not found")

	// create a tenant
	tcr1 := &v1.Tenant{
		Meta:        &v1.Meta{Id: t9},
		Name:        "etenant",
		Description: "E Tenant",
	}
	err = tenantDS.Create(ctx, tcr1)
	assert.NoError(t, err)
	assert.Equal(t, t9, tcr1.GetMeta().GetId())
	assert.Equal(t, "etenant", tcr1.GetName())
	assert.Equal(t, "E Tenant", tcr1.GetDescription())

	// now delete it
	tdr2 := &v1.Tenant{
		Meta: &v1.Meta{Id: t9},
	}
	err = tenantDS.Delete(ctx, tdr2)
	assert.NoError(t, err)

	_, err = tenantDS.Get(ctx, t9)
	assert.Error(t, err)
	assert.EqualError(t, err, "entity of type:tenant with id:t9 not found")

	var tgrh v1.Tenant
	err = tenantDS.GetHistory(ctx, t9, time.Now(), &tgrh)
	assert.NoError(t, err)
	assert.Equal(t, "etenant", tgrh.Name)
}

func TestAnnotationsAndLabels(t *testing.T) {
	tenantDS, err := NewPostgresStorage(zaptest.NewLogger(t), db, &v1.Tenant{})
	require.NoError(t, err)
	assert.NotNil(t, tenantDS, "Datastore must not be nil")
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

	err = tenantDS.Create(ctx, tcr)
	assert.NoError(t, err)
	assert.NotNil(t, tcr)
	assert.Equal(t, int64(0), tcr.Meta.Version)
	assert.Equal(t, "A Tenant", tcr.GetName())
	assert.NotNil(t, tcr.GetMeta().GetAnnotations())
	assert.NotNil(t, tcr.GetMeta().GetLabels())
	assert.Equal(t, map[string]string{"key1": "value1", "key2": "value2"}, tcr.GetMeta().GetAnnotations())
	assert.Equal(t, []string{"color=red"}, tcr.GetMeta().GetLabels())

	// get from db
	tget, err := tenantDS.Get(ctx, tcr.Meta.Id)
	require.NoError(t, err)
	require.Equal(t, tcr.Meta.Id, tget.Meta.Id)
	assert.Equal(t, int64(0), tget.Meta.Version)

	// update instance
	tcr.Name = "updated name"
	err = tenantDS.Update(ctx, tcr)
	assert.NoError(t, err)
	assert.NotNil(t, tcr)
	// incremented version after update
	assert.Equal(t, int64(1), tcr.Meta.Version)

	// re-read from db
	tgr, err := tenantDS.Get(ctx, tcr.Meta.GetId())
	assert.NoError(t, err)
	// version is incremented
	assert.Equal(t, int64(1), tgr.Meta.Version)
	// updated data is reflected
	assert.Equal(t, "updated name", tgr.GetName())

	// try to update older version --> optimistic lock error
	tget.Name = "updated older entity"
	err = tenantDS.Update(ctx, tget)
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
	err = tenantDS.Update(ctx, tcr)
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

func createPostgresConnection() (*sqlx.DB, error) {

	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
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
		return nil, err
	}

	ip, err := postgres.Host(ctx)
	if err != nil {
		return nil, err
	}
	port, err := postgres.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return nil, err
	}

	log := zap.NewNop()
	var db *sqlx.DB
	for {
		var err error
		ves := []Entity{
			&v1.Project{},
			&v1.Tenant{},
		}
		db, err = NewPostgresDB(log, ip, port.Port(), "postgres", "password", "postgres", "disable", ves...)
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		err = db.Ping()
		if err == nil {
			break
		}
	}
	return db, nil
}
