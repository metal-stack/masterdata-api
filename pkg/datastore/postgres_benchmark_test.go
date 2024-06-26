package datastore

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	ds Storage[*v1.Tenant]
)

func init() {
	_, db, _ = createPostgresConnection()
	ds = New(slog.Default(), db, &v1.Tenant{})
}

func BenchmarkGetTenant(b *testing.B) {

	t1 := &v1.Tenant{
		Meta: &v1.Meta{
			Id: "t1",
		},
	}
	err := ds.Create(context.Background(), t1)
	require.NoError(b, err)
	defer func() {
		_ = ds.Delete(context.Background(), "t1")
	}()

	for n := 0; n < b.N; n++ {
		t, err := ds.Get(context.Background(), "t1")
		require.NoError(b, err)
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
		require.NoError(b, err)
	}
}

func BenchmarkUpdateTenant(b *testing.B) {
	t1 := &v1.Tenant{
		Meta: &v1.Meta{
			Id: "t1-update",
		},
	}
	err := ds.Create(context.Background(), t1)
	require.NoError(b, err)
	defer func() {
		_ = ds.Delete(context.Background(), "t1-update")
	}()

	for n := 0; n < b.N; n++ {
		t1, err := ds.Get(context.Background(), t1.Meta.Id)
		require.NoError(b, err)
		t1.Name = fmt.Sprintf("t1-create-%d", n)
		t1.Meta.Version = int64(t1.Meta.Version)
		err = ds.Update(context.Background(), t1)
		require.NoError(b, err)
	}
}

func BenchmarkFindTenant(b *testing.B) {
	err := ds.Create(context.Background(), &v1.Tenant{
		Meta: &v1.Meta{
			Id: "t1",
		},
		Name: "tenant-1",
	})
	require.NoError(b, err)
	defer func() {
		_ = ds.Delete(context.Background(), "t1")
	}()

	for n := 0; n < b.N; n++ {
		f := make(map[string]any)
		f["tenant ->> 'name'"] = "tenant-1"

		t, _, err := ds.Find(context.Background(), f, nil)
		require.NoError(b, err)
		assert.NotNil(b, t)
		assert.Len(b, t, 1)
	}
}
