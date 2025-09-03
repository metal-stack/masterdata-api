package service

import (
	"context"
	"log/slog"

	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"testing"

	"github.com/metal-stack/masterdata-api/pkg/test/mocks"
)

func TestCreateTenantMember(t *testing.T) {
	storageMock := mocks.NewMockStorage[*v1.TenantMember](t)
	tenantStorageMock := mocks.NewMockStorage[*v1.Tenant](t)
	ts := &tenantMemberService{
		tenantMemberStore: storageMock,
		tenantStore:       tenantStorageMock,
		log:               slog.Default(),
	}
	ctx := context.Background()

	t1 := &v1.Tenant{}
	m1 := &v1.Tenant{}
	pm1 := &v1.TenantMember{
		TenantId: "t1",
		MemberId: "m1",
	}
	pmcr := &v1.TenantMemberCreateRequest{
		TenantMember: pm1,
	}
	tenantStorageMock.On("Get", ctx, pm1.GetTenantId()).Return(t1, nil)
	tenantStorageMock.On("Get", ctx, pm1.GetMemberId()).Return(m1, nil)
	storageMock.On("Create", ctx, pm1).Return(nil)
	resp, err := ts.Create(ctx, pmcr)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.GetTenantMember())
	assert.Equal(t, pmcr.TenantMember.TenantId, resp.GetTenantMember().GetTenantId())
}

func TestUpdateTenantMember(t *testing.T) {
	storageMock := mocks.NewMockStorage[*v1.TenantMember](t)
	tenantStorageMock := mocks.NewMockStorage[*v1.Tenant](t)
	ts := &tenantMemberService{
		tenantMemberStore: storageMock,
		tenantStore:       tenantStorageMock,
		log:               slog.Default(),
	}
	ctx := context.Background()

	meta := &v1.Meta{Id: "p2", Annotations: map[string]string{"key": "value"}}
	pm1 := &v1.TenantMember{
		Meta:     meta,
		TenantId: "p1",
		MemberId: "t1",
	}
	meta.Annotations = map[string]string{"key": "value2"}
	pmur := &v1.TenantMemberUpdateRequest{
		TenantMember: &v1.TenantMember{
			Meta:     meta,
			TenantId: "p1",
			MemberId: "t1",
		},
	}

	storageMock.On("Get", ctx, pm1.Meta.Id).Return(pm1, nil)
	storageMock.On("Update", ctx, pm1).Return(nil)
	resp, err := ts.Update(ctx, pmur)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.GetTenantMember())
	assert.Equal(t, pmur.GetTenantMember().Meta.Annotations, resp.GetTenantMember().Meta.Annotations)
}

func TestDeleteTenantMember(t *testing.T) {
	storageMock := mocks.NewMockStorage[*v1.TenantMember](t)
	tenantStorageMock := mocks.NewMockStorage[*v1.Tenant](t)
	ts := &tenantMemberService{
		tenantMemberStore: storageMock,
		tenantStore:       tenantStorageMock,
		log:               slog.Default(),
	}
	ctx := context.Background()
	t3 := &v1.TenantMember{
		Meta: &v1.Meta{Id: "p3"},
	}
	tdr := &v1.TenantMemberDeleteRequest{
		Id: "p3",
	}

	storageMock.On("Delete", ctx, t3.Meta.Id).Return(nil)
	resp, err := ts.Delete(ctx, tdr)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.GetTenantMember())
	assert.Equal(t, tdr.Id, resp.GetTenantMember().GetMeta().GetId())
}

func TestGetTenantMember(t *testing.T) {
	storageMock := mocks.NewMockStorage[*v1.TenantMember](t)
	tenantStorageMock := mocks.NewMockStorage[*v1.Tenant](t)
	ts := &tenantMemberService{
		tenantMemberStore: storageMock,
		tenantStore:       tenantStorageMock,
		log:               slog.Default(),
	}
	ctx := context.Background()
	t4 := &v1.TenantMember{
		Meta: &v1.Meta{Id: "p4"},
	}
	tgr := &v1.TenantMemberGetRequest{
		Id: "p4",
	}

	storageMock.On("Get", ctx, "p4").Return(t4, nil)
	resp, err := ts.Get(ctx, tgr)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.GetTenantMember())
	assert.Equal(t, tgr.Id, resp.GetTenantMember().GetMeta().GetId())
}

func TestFindTenantMemberByTenant(t *testing.T) {
	storageMock := mocks.NewMockStorage[*v1.TenantMember](t)
	tenantStorageMock := mocks.NewMockStorage[*v1.Tenant](t)
	ts := &tenantMemberService{
		tenantMemberStore: storageMock,
		tenantStore:       tenantStorageMock,
		log:               slog.Default(),
	}
	ctx := context.Background()

	// filter by name
	var t6s []*v1.TenantMember
	tfr := &v1.TenantMemberFindRequest{
		TenantId:  pointer.Pointer("p1"),
		Namespace: "a",
	}

	f2 := make(map[string]any)
	f2["tenantmember ->> 'tenant_id'"] = pointer.Pointer("p1")
	f2["tenantmember ->> 'namespace'"] = "a"
	storageMock.On("Find", ctx, mock.AnythingOfType("*v1.Paging"), []any{f2}).Return(t6s, nil, nil)
	resp, err := ts.Find(ctx, tfr)
	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestFindTenantMemberByMember(t *testing.T) {
	storageMock := mocks.NewMockStorage[*v1.TenantMember](t)
	tenantStorageMock := mocks.NewMockStorage[*v1.Tenant](t)
	ts := &tenantMemberService{
		tenantMemberStore: storageMock,
		tenantStore:       tenantStorageMock,
		log:               slog.Default(),
	}
	ctx := context.Background()

	// filter by name
	var t6s []*v1.TenantMember
	tfr := &v1.TenantMemberFindRequest{
		MemberId:  pointer.Pointer("t1"),
		Namespace: "a",
	}

	f2 := make(map[string]any)
	f2["tenantmember ->> 'member_id'"] = pointer.Pointer("t1")
	f2["tenantmember ->> 'namespace'"] = "a"
	storageMock.On("Find", ctx, mock.AnythingOfType("*v1.Paging"), []any{f2}).Return(t6s, nil, nil)
	resp, err := ts.Find(ctx, tfr)
	require.NoError(t, err)
	assert.NotNil(t, resp)
}
