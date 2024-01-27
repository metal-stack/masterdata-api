package service

import (
	"context"

	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"testing"

	"github.com/metal-stack/masterdata-api/pkg/datastore/mocks"
)

func TestCreateProjectMember(t *testing.T) {
	storageMock := &mocks.Storage[*v1.ProjectMember]{}
	tenantStorageMock := &mocks.Storage[*v1.Tenant]{}
	ts := &projectMemberService{
		projectMemberStore: storageMock,
		tenantStore:        tenantStorageMock,
		log:                zaptest.NewLogger(t),
	}
	ctx := context.Background()

	t1 := &v1.Tenant{}
	p1 := &v1.ProjectMember{
		ProjectId: "p1",
		TenantId:  "t1",
	}
	tcr := &v1.ProjectMemberCreateRequest{
		ProjectMember: p1,
	}
	tenantStorageMock.On("Get", ctx, p1.GetTenantId()).Return(t1, nil)
	storageMock.On("Create", ctx, p1).Return(nil)
	resp, err := ts.Create(ctx, tcr)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.GetProjectMember())
	assert.Equal(t, tcr.ProjectMember.ProjectId, resp.GetProjectMember().GetProjectId())
}

func TestUpdateProjectMember(t *testing.T) {
	storageMock := &mocks.Storage[*v1.ProjectMember]{}
	tenantStorageMock := &mocks.Storage[*v1.Tenant]{}
	ts := &projectMemberService{
		projectMemberStore: storageMock,
		tenantStore:        tenantStorageMock,
		log:                zaptest.NewLogger(t),
	}
	ctx := context.Background()

	t1 := &v1.ProjectMember{
		Meta:      &v1.Meta{Id: "p2", Annotations: map[string]string{"key": "value"}},
		ProjectId: "p1",
		TenantId:  "t1",
	}
	tur := &v1.ProjectMemberUpdateRequest{
		ProjectMember: &v1.ProjectMember{
			Meta: &v1.Meta{Id: "p2", Annotations: map[string]string{"key": "value2"}},
		},
	}

	storageMock.On("Update", ctx, t1).Return(nil)
	resp, err := ts.Update(ctx, tur)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.GetProjectMember())
	assert.Equal(t, tur.GetProjectMember().Meta.Annotations, resp.GetProjectMember().Meta.Annotations)
}

func TestDeleteProjectMember(t *testing.T) {
	storageMock := &mocks.Storage[*v1.ProjectMember]{}
	tenantStorageMock := &mocks.Storage[*v1.Tenant]{}
	ts := &projectMemberService{
		projectMemberStore: storageMock,
		tenantStore:        tenantStorageMock,
		log:                zaptest.NewLogger(t),
	}
	ctx := context.Background()
	t3 := &v1.ProjectMember{
		Meta: &v1.Meta{Id: "p3"},
	}
	tdr := &v1.ProjectMemberDeleteRequest{
		Id: "p3",
	}

	storageMock.On("Delete", ctx, t3.Meta.Id).Return(nil)
	resp, err := ts.Delete(ctx, tdr)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.GetProjectMember())
	assert.Equal(t, tdr.Id, resp.GetProjectMember().GetMeta().GetId())
}

func TestGetProjectMember(t *testing.T) {
	storageMock := &mocks.Storage[*v1.ProjectMember]{}
	tenantStorageMock := &mocks.Storage[*v1.Tenant]{}
	ts := &projectMemberService{
		projectMemberStore: storageMock,
		tenantStore:        tenantStorageMock,
		log:                zaptest.NewLogger(t),
	}
	ctx := context.Background()
	t4 := &v1.ProjectMember{
		Meta: &v1.Meta{Id: "p4"},
	}
	tgr := &v1.ProjectMemberGetRequest{
		Id: "p4",
	}

	storageMock.On("Get", ctx, "p4").Return(t4, nil)
	resp, err := ts.Get(ctx, tgr)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.GetProjectMember())
	assert.Equal(t, tgr.Id, resp.GetProjectMember().GetMeta().GetId())
}

func TestFindProjectMemberByProject(t *testing.T) {
	storageMock := &mocks.Storage[*v1.ProjectMember]{}
	tenantStorageMock := &mocks.Storage[*v1.Tenant]{}
	ts := &projectMemberService{
		projectMemberStore: storageMock,
		tenantStore:        tenantStorageMock,
		log:                zaptest.NewLogger(t),
	}
	ctx := context.Background()

	// filter by name
	var t6s []*v1.ProjectMember
	tfr := &v1.ProjectMemberFindRequest{
		ProjectId: pointer.Pointer("p1"),
	}

	f2 := make(map[string]any)
	f2["projectmember ->> 'project_id'"] = pointer.Pointer("p1")
	storageMock.On("Find", ctx, f2).Return(t6s, nil, nil)
	resp, err := ts.Find(ctx, tfr)
	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestFindProjectMemberByTenant(t *testing.T) {
	storageMock := &mocks.Storage[*v1.ProjectMember]{}
	tenantStorageMock := &mocks.Storage[*v1.Tenant]{}
	ts := &projectMemberService{
		projectMemberStore: storageMock,
		tenantStore:        tenantStorageMock,
		log:                zaptest.NewLogger(t),
	}
	ctx := context.Background()

	// filter by name
	var t6s []*v1.ProjectMember
	tfr := &v1.ProjectMemberFindRequest{
		TenantId: pointer.Pointer("t1"),
	}

	f2 := make(map[string]any)
	f2["projectmember ->> 'tenant_id'"] = pointer.Pointer("t1")
	storageMock.On("Find", ctx, f2).Return(t6s, nil, nil)
	resp, err := ts.Find(ctx, tfr)
	require.NoError(t, err)
	assert.NotNil(t, resp)
}
