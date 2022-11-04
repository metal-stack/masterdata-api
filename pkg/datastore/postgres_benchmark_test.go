package datastore

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func BenchMain(m *testing.B) {
	fmt.Println("benchmain")
}

func BenchmarkGetTenant(b *testing.B) {
	db, err := createPostgresConnection()
	assert.NoError(b, err)
	ds, err := NewPostgresStorage(zap.NewNop(), db, &v1.Tenant{})
	require.NoError(b, err)

	err = ds.Create(context.Background(), &v1.Tenant{
		Meta: &v1.Meta{
			Id: "t1",
		},
	})
	assert.NoError(b, err)

	for n := 0; n < b.N; n++ {
		t, err := ds.Get(context.Background(), "t1")
		assert.NoError(b, err)
		assert.NotNil(b, t)
	}
}

func BenchmarkCreateTenant(b *testing.B) {
	db, err := createPostgresConnection()
	assert.NoError(b, err)
	ds, err := NewPostgresStorage(zap.NewNop(), db, &v1.Tenant{})
	require.NoError(b, err)

	for n := 0; n < b.N; n++ {
		err = ds.Create(context.Background(), &v1.Tenant{
			Meta: &v1.Meta{
				Id: uuid.NewString(),
			},
		})
		assert.NoError(b, err)
	}
}

func BenchmarkFindTenant(b *testing.B) {
	db, err := createPostgresConnection()
	assert.NoError(b, err)
	ds, err := NewPostgresStorage(zap.NewNop(), db, &v1.Tenant{})
	require.NoError(b, err)

	err = ds.Create(context.Background(), &v1.Tenant{
		Meta: &v1.Meta{
			Id: "t1",
		},
		Name: "tenant-1",
	})
	assert.NoError(b, err)

	for n := 0; n < b.N; n++ {
		f := make(map[string]any)
		f["tenant ->> 'name'"] = "tenant-1"

		t, _, err := ds.Find(context.Background(), f, nil)
		assert.NoError(b, err)
		assert.NotNil(b, t)
		assert.Len(b, t, 1)
	}
}
