package service

import (
	"context"
	"log/slog"
	"os"

	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"testing"

	"github.com/metal-stack/masterdata-api/pkg/datastore/mocks"
)

var log *slog.Logger

func TestMain(m *testing.M) {
	code := 0
	defer func() {
		os.Exit(code)
	}()
	log = slog.Default()
	code = m.Run()
}

func TestCreateTenant(t *testing.T) {
	storageMock := &mocks.Storage[*v1.Tenant]{}
	ts := &tenantService{
		tenantStore: storageMock,
		log:         slog.Default(),
	}
	ctx := context.Background()

	t1 := &v1.Tenant{
		Name:        "First",
		Description: "First Tenant",
		Meta: &v1.Meta{
			Annotations: map[string]string{
				"metal-stack.io/contract": "2345",
			},
			Labels: []string{
				"color=blue",
			},
		},
	}
	tcr := &v1.TenantCreateRequest{
		Tenant: t1,
	}

	storageMock.On("Create", ctx, t1).Return(nil)
	resp, err := ts.Create(ctx, tcr)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.GetTenant())
	assert.Equal(t, tcr.Tenant.GetName(), resp.GetTenant().GetName())
}

func TestUpdateTenant(t *testing.T) {
	storageMock := &mocks.Storage[*v1.Tenant]{}
	ts := &tenantService{
		tenantStore: storageMock,
		log:         slog.Default(),
	}
	ctx := context.Background()

	t1 := &v1.Tenant{
		Name:        "Second",
		Description: "Second Tenant",
	}
	tur := &v1.TenantUpdateRequest{
		Tenant: t1,
	}

	storageMock.On("Update", ctx, t1).Return(nil)
	resp, err := ts.Update(ctx, tur)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.GetTenant())
	assert.Equal(t, tur.Tenant.GetName(), resp.GetTenant().GetName())
}

func TestDeleteTenant(t *testing.T) {
	storageMock := &mocks.Storage[*v1.Tenant]{}
	ts := &tenantService{
		tenantStore: storageMock,
		log:         slog.Default(),
	}
	ctx := context.Background()
	t3 := &v1.Tenant{
		Meta: &v1.Meta{Id: "t3"},
	}
	tdr := &v1.TenantDeleteRequest{
		Id: "t3",
	}

	storageMock.On("Delete", ctx, t3.Meta.Id).Return(nil)
	resp, err := ts.Delete(ctx, tdr)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.GetTenant())
	assert.Equal(t, tdr.Id, resp.GetTenant().GetMeta().GetId())
}

func TestGetTenant(t *testing.T) {
	storageMock := &mocks.Storage[*v1.Tenant]{}
	ts := &tenantService{
		tenantStore: storageMock,
		log:         slog.Default(),
	}
	ctx := context.Background()
	t4 := &v1.Tenant{
		Meta: &v1.Meta{Id: "t4"},
	}
	tgr := &v1.TenantGetRequest{
		Id: "t4",
	}

	storageMock.On("Get", ctx, "t4").Return(t4, nil)
	resp, err := ts.Get(ctx, tgr)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.GetTenant())
	assert.Equal(t, tgr.Id, resp.GetTenant().GetMeta().GetId())
}

func TestFindTenantByID(t *testing.T) {
	storageMock := &mocks.Storage[*v1.Tenant]{}
	ts := &tenantService{
		tenantStore: storageMock,
		log:         slog.Default(),
	}
	ctx := context.Background()
	var t5s []*v1.Tenant
	// filter by id
	f1 := make(map[string]any)
	tfr := &v1.TenantFindRequest{
		Id: &wrapperspb.StringValue{Value: "t5"},
	}

	f1["id"] = "t5"
	storageMock.On("Find", ctx, f1, mock.AnythingOfType("*v1.Paging")).Return(t5s, nil, nil)
	resp, err := ts.Find(ctx, tfr)
	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestFindTenantByName(t *testing.T) {
	storageMock := &mocks.Storage[*v1.Tenant]{}
	ts := &tenantService{
		tenantStore: storageMock,
		log:         slog.Default(),
	}
	ctx := context.Background()

	// filter by name
	var t6s []*v1.Tenant
	tfr := &v1.TenantFindRequest{
		Name: &wrapperspb.StringValue{Value: "Fifth"},
	}

	f2 := make(map[string]any)
	f2["tenant ->> 'name'"] = "Fifth"
	storageMock.On("Find", ctx, f2, mock.AnythingOfType("*v1.Paging")).Return(t6s, nil, nil)
	resp, err := ts.Find(ctx, tfr)
	require.NoError(t, err)
	assert.NotNil(t, resp)
}
