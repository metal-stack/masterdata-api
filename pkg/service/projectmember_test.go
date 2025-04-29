package service

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"
	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"testing"

	"github.com/metal-stack/masterdata-api/pkg/datastore/mocks"
)

func TestCreateProjectMember(t *testing.T) {
	storageMock := &mocks.Storage[*v1.ProjectMember]{}
	tenantStorageMock := &mocks.Storage[*v1.Tenant]{}
	projectStorageMock := &mocks.Storage[*v1.Project]{}
	ts := &projectMemberService{
		projectMemberStore: storageMock,
		tenantStore:        tenantStorageMock,
		projectStore:       projectStorageMock,
		log:                slog.Default(),
	}
	ctx := context.Background()

	t1 := &v1.Tenant{}
	p1 := &v1.Project{}
	pm1 := &v1.ProjectMember{
		ProjectId: "p1",
		TenantId:  "t1",
	}
	pmcr := &v1.ProjectMemberCreateRequest{
		ProjectMember: pm1,
	}
	tenantStorageMock.On("Get", ctx, pm1.GetTenantId()).Return(t1, nil)
	projectStorageMock.On("Get", ctx, pm1.GetProjectId()).Return(p1, nil)
	storageMock.On("Create", ctx, pm1).Return(nil)
	resp, err := ts.Create(ctx, connect.NewRequest(pmcr))
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Msg.ProjectMember)
	assert.Equal(t, pmcr.ProjectMember.ProjectId, resp.Msg.ProjectMember.GetProjectId())
}

func TestUpdateProjectMember(t *testing.T) {
	storageMock := &mocks.Storage[*v1.ProjectMember]{}
	tenantStorageMock := &mocks.Storage[*v1.Tenant]{}
	projectStorageMock := &mocks.Storage[*v1.Project]{}
	ts := &projectMemberService{
		projectMemberStore: storageMock,
		tenantStore:        tenantStorageMock,
		projectStore:       projectStorageMock,
		log:                slog.Default(),
	}
	ctx := context.Background()

	meta := &v1.Meta{Id: "p2", Annotations: map[string]string{"key": "value"}}
	pm1 := &v1.ProjectMember{
		Meta:      meta,
		ProjectId: "p1",
		TenantId:  "t1",
	}
	meta.Annotations = map[string]string{"key": "value2"}
	pmur := &v1.ProjectMemberUpdateRequest{
		ProjectMember: &v1.ProjectMember{
			Meta:      meta,
			ProjectId: "p1",
			TenantId:  "t1",
		},
	}

	storageMock.On("Update", ctx, pm1).Return(nil)
	resp, err := ts.Update(ctx, connect.NewRequest(pmur))
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Msg.ProjectMember)
	assert.Equal(t, pmur.ProjectMember.Meta.Annotations, resp.Msg.ProjectMember.Meta.Annotations)
}

func TestDeleteProjectMember(t *testing.T) {
	storageMock := &mocks.Storage[*v1.ProjectMember]{}
	tenantStorageMock := &mocks.Storage[*v1.Tenant]{}
	projectStorageMock := &mocks.Storage[*v1.Project]{}
	ts := &projectMemberService{
		projectMemberStore: storageMock,
		tenantStore:        tenantStorageMock,
		projectStore:       projectStorageMock,
		log:                slog.Default(),
	}
	ctx := context.Background()
	t3 := &v1.ProjectMember{
		Meta: &v1.Meta{Id: "p3"},
	}
	tdr := &v1.ProjectMemberDeleteRequest{
		Id: "p3",
	}

	storageMock.On("Delete", ctx, t3.Meta.Id).Return(nil)
	resp, err := ts.Delete(ctx, connect.NewRequest(tdr))
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Msg.ProjectMember)
	assert.Equal(t, tdr.Id, resp.Msg.ProjectMember.GetMeta().GetId())
}

func TestGetProjectMember(t *testing.T) {
	storageMock := &mocks.Storage[*v1.ProjectMember]{}
	tenantStorageMock := &mocks.Storage[*v1.Tenant]{}
	projectStorageMock := &mocks.Storage[*v1.Project]{}
	ts := &projectMemberService{
		projectMemberStore: storageMock,
		tenantStore:        tenantStorageMock,
		projectStore:       projectStorageMock,
		log:                slog.Default(),
	}
	ctx := context.Background()
	t4 := &v1.ProjectMember{
		Meta: &v1.Meta{Id: "p4"},
	}
	tgr := &v1.ProjectMemberGetRequest{
		Id: "p4",
	}

	storageMock.On("Get", ctx, "p4").Return(t4, nil)
	resp, err := ts.Get(ctx, connect.NewRequest(tgr))
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Msg.ProjectMember)
	assert.Equal(t, tgr.Id, resp.Msg.ProjectMember.GetMeta().GetId())
}

func TestFindProjectMemberByProject(t *testing.T) {
	storageMock := &mocks.Storage[*v1.ProjectMember]{}
	tenantStorageMock := &mocks.Storage[*v1.Tenant]{}
	projectStorageMock := &mocks.Storage[*v1.Project]{}
	ts := &projectMemberService{
		projectMemberStore: storageMock,
		tenantStore:        tenantStorageMock,
		projectStore:       projectStorageMock,
		log:                slog.Default(),
	}
	ctx := context.Background()

	// filter by name
	var t6s []*v1.ProjectMember
	tfr := &v1.ProjectMemberFindRequest{
		ProjectId: pointer.Pointer("p1"),
	}

	f2 := make(map[string]any)
	f2["projectmember ->> 'project_id'"] = pointer.Pointer("p1")
	storageMock.On("Find", ctx, f2, mock.AnythingOfType("*v1.Paging")).Return(t6s, nil, nil)
	resp, err := ts.Find(ctx, connect.NewRequest(tfr))
	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestFindProjectMemberByTenant(t *testing.T) {
	storageMock := &mocks.Storage[*v1.ProjectMember]{}
	tenantStorageMock := &mocks.Storage[*v1.Tenant]{}
	projectStorageMock := &mocks.Storage[*v1.Project]{}
	ts := &projectMemberService{
		projectMemberStore: storageMock,
		tenantStore:        tenantStorageMock,
		projectStore:       projectStorageMock,
		log:                slog.Default(),
	}
	ctx := context.Background()

	// filter by name
	var t6s []*v1.ProjectMember
	tfr := &v1.ProjectMemberFindRequest{
		TenantId: pointer.Pointer("t1"),
	}

	f2 := make(map[string]any)
	f2["projectmember ->> 'tenant_id'"] = pointer.Pointer("t1")
	storageMock.On("Find", ctx, f2, mock.AnythingOfType("*v1.Paging")).Return(t6s, nil, nil)
	resp, err := ts.Find(ctx, connect.NewRequest(tfr))
	require.NoError(t, err)
	assert.NotNil(t, resp)
}
