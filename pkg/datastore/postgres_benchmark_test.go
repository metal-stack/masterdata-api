package datastore

import (
	"context"
	"testing"

	"github.com/google/uuid"
	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var (
	ds Storage[*v1.Tenant]
)

func init() {
	db, _ = createPostgresConnection()
	ds, _ = NewPostgresStorage(zap.NewNop(), db, &v1.Tenant{})
}

func BenchmarkGetTenant(b *testing.B) {

	t1 := &v1.Tenant{
		Meta: &v1.Meta{
			Id: "t1",
		},
	}
	err := ds.Create(context.Background(), t1)
	assert.NoError(b, err)
	defer func() {
		ds.Delete(context.Background(), "t1")
	}()

	for n := 0; n < b.N; n++ {
		t, err := ds.Get(context.Background(), "t1")
		assert.NoError(b, err)
		assert.NotNil(b, t)
	}
}

func BenchmarkCreateTenant(b *testing.B) {
	for n := 0; n < b.N; n++ {
		err := ds.Create(context.Background(), &v1.Tenant{
			Meta: &v1.Meta{
				Id: uuid.NewString(),
			},
		})
		assert.NoError(b, err)
	}
}

func BenchmarkFindTenant(b *testing.B) {
	err := ds.Create(context.Background(), &v1.Tenant{
		Meta: &v1.Meta{
			Id: "t1",
		},
		Name: "tenant-1",
	})
	assert.NoError(b, err)
	defer func() {
		ds.Delete(context.Background(), "t1")
	}()

	for n := 0; n < b.N; n++ {
		f := make(map[string]any)
		f["tenant ->> 'name'"] = "tenant-1"

		t, _, err := ds.Find(context.Background(), f, nil)
		assert.NoError(b, err)
		assert.NotNil(b, t)
		assert.Len(b, t, 1)
	}
}
