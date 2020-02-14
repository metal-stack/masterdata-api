package service

import (
	"context"
	"os"

	v1 "github.com/metal-stack/masterdata-api/api/v1"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"

	"testing"

	"github.com/metal-stack/masterdata-api/pkg/datastore/mocks"
	"go.uber.org/zap"
)

var log *zap.Logger

func TestMain(m *testing.M) {
	code := 0
	defer func() {
		os.Exit(code)
	}()
	log, _ = zap.NewProduction()
	code = m.Run()
}

func TestCreateTenant(t *testing.T) {
	storageMock := &mocks.Storage{}
	ts := NewTenantService(storageMock, log)
	ctx := context.Background()

	t1 := &v1.Tenant{
		Name:        "First",
		Description: "First Tenant",
	}
	tcr := &v1.TenantCreateRequest{
		Tenant: t1,
	}

	storageMock.On("Create", ctx, t1).Return(nil)
	resp, err := ts.Create(ctx, tcr)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.GetTenant())
	assert.Equal(t, tcr.Tenant.GetName(), resp.GetTenant().GetName())
}

func TestUpdateTenant(t *testing.T) {
	storageMock := &mocks.Storage{}
	ts := NewTenantService(storageMock, log)
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
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.GetTenant())
	assert.Equal(t, tur.Tenant.GetName(), resp.GetTenant().GetName())
}

func TestDeleteTenant(t *testing.T) {
	storageMock := &mocks.Storage{}
	ts := NewTenantService(storageMock, log)
	ctx := context.Background()
	t3 := &v1.Tenant{
		Meta: &v1.Meta{Id: "t3"},
	}
	tdr := &v1.TenantDeleteRequest{
		Id: "t3",
	}

	storageMock.On("Delete", ctx, t3).Return(nil)
	resp, err := ts.Delete(ctx, tdr)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.GetTenant())
	assert.Equal(t, tdr.Id, resp.GetTenant().GetMeta().GetId())
}

// FIXME reenable
// func TestGetTenant(t *testing.T) {
// 	storageMock := &mocks.Storage{}
// 	ts := NewTenantService(storageMock, log)
// 	ctx := context.Background()
// 	t4 := &v1.Tenant{}
// 	tgr := &v1.TenantGetRequest{
// 		Id: "t4",
// 	}

// 	storageMock.On("Get", ctx, "t4", t4).Return(nil)
// 	resp, err := ts.Get(ctx, tgr)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, resp)
// 	assert.NotNil(t, resp.GetTenant())
// 	assert.Equal(t, tgr.Id, resp.GetTenant().GetMeta().GetId())
// }

func TestFindTenantByID(t *testing.T) {
	storageMock := &mocks.Storage{}
	ts := NewTenantService(storageMock, log)
	ctx := context.Background()
	var t5s []v1.Tenant
	// filter by id
	f1 := make(map[string]interface{})
	tfr := &v1.TenantFindRequest{
		Id: &wrappers.StringValue{Value: "t5"},
	}

	f1["id"] = "t5"
	storageMock.On("Find", ctx, f1, &t5s).Return(nil)
	resp, err := ts.Find(ctx, tfr)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestFindTenantByName(t *testing.T) {
	storageMock := &mocks.Storage{}
	ts := NewTenantService(storageMock, log)
	ctx := context.Background()

	// filter by name
	var t6s []v1.Tenant
	tfr := &v1.TenantFindRequest{
		Name: &wrappers.StringValue{Value: "Fifth"},
	}

	f2 := make(map[string]interface{})
	f2["tenant ->> 'name'"] = "Fifth"
	storageMock.On("Find", ctx, f2, &t6s).Return(nil)
	resp, err := ts.Find(ctx, tfr)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
