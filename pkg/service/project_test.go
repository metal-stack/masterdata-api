package service

import (
	"context"
	"log/slog"

	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"testing"

	"github.com/metal-stack/masterdata-api/pkg/datastore/mocks"
)

func TestCreateProject(t *testing.T) {
	storageMock := &mocks.Storage[*v1.Project]{}
	tenantStorageMock := &mocks.Storage[*v1.Tenant]{}
	ts := &projectService{
		projectStore: storageMock,
		tenantStore:  tenantStorageMock,
		log:          slog.Default(),
	}
	ctx := context.Background()

	t1 := &v1.Tenant{}
	p1 := &v1.Project{
		Name:        "FirstP",
		Description: "First Project",
		TenantId:    "t1",
		Meta: &v1.Meta{
			Annotations: map[string]string{
				"metal-stack.io/contract": "1234",
			},
			Labels: []string{
				"color=green",
			},
		},
	}
	tcr := &v1.ProjectCreateRequest{
		Project: p1,
	}
	tenantStorageMock.On("Get", ctx, p1.GetTenantId()).Return(t1, nil)
	storageMock.On("Create", ctx, p1).Return(nil)
	resp, err := ts.Create(ctx, tcr)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.GetProject())
	assert.Equal(t, tcr.Project.GetName(), resp.GetProject().GetName())
}

func TestCreateProjectWithQuotaCheck(t *testing.T) {
	storageMock := &mocks.Storage[*v1.Project]{}
	tenantStorageMock := &mocks.Storage[*v1.Tenant]{}
	ts := &projectService{
		projectStore: storageMock,
		tenantStore:  tenantStorageMock,
		log:          slog.Default(),
	}
	ctx := context.Background()

	t1 := &v1.Tenant{
		Quotas: &v1.QuotaSet{
			Project: &v1.Quota{
				Quota: &wrapperspb.Int32Value{Value: 2},
			},
		},
	}
	p1 := &v1.Project{
		Name:        "FirstP",
		Description: "First Project",
		TenantId:    "t1",
	}
	tcr := &v1.ProjectCreateRequest{
		Project: p1,
	}
	filter := make(map[string]any)
	filter["project ->> 'tenant_id'"] = p1.TenantId
	var projects []*v1.Project
	// see: https://github.com/stretchr/testify/blob/master/mock/mock.go#L149-L162
	tenantStorageMock.On("Get", ctx, p1.GetTenantId()).Return(t1, nil)
	storageMock.On("Find", ctx, filter, mock.AnythingOfType("*v1.Paging")).Return(projects, nil, nil)
	storageMock.On("Create", ctx, p1).Return(nil)
	resp, err := ts.Create(ctx, tcr)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.GetProject())
	assert.Equal(t, tcr.Project.GetName(), resp.GetProject().GetName())
}

func TestUpdateProject(t *testing.T) {
	storageMock := &mocks.Storage[*v1.Project]{}
	tenantStorageMock := &mocks.Storage[*v1.Tenant]{}
	ts := &projectService{
		projectStore: storageMock,
		tenantStore:  tenantStorageMock,
		log:          slog.Default(),
	}
	ctx := context.Background()

	t1 := &v1.Project{
		Meta:        &v1.Meta{Id: "p2"},
		Name:        "SecondP",
		Description: "Second Project",
	}
	tur := &v1.ProjectUpdateRequest{
		Project: &v1.Project{
			Meta:        &v1.Meta{Id: "p2"},
			Name:        "SecondP",
			Description: "Second Project",
		},
	}

	storageMock.On("Get", ctx, t1.Meta.Id).Return(t1, nil)

	storageMock.On("Update", ctx, t1).Return(nil)
	resp, err := ts.Update(ctx, tur)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.GetProject())
	assert.Equal(t, tur.GetProject().GetName(), resp.GetProject().GetName())
}

func TestDeleteProject(t *testing.T) {
	storageMock := &mocks.Storage[*v1.Project]{}
	projectMemberStorageMock := &mocks.Storage[*v1.ProjectMember]{}
	tenantStorageMock := &mocks.Storage[*v1.Tenant]{}
	ps := &projectService{
		projectStore:       storageMock,
		projectMemberStore: projectMemberStorageMock,
		tenantStore:        tenantStorageMock,
		log:                slog.Default(),
	}
	ctx := context.Background()
	p3 := &v1.Project{
		Meta: &v1.Meta{Id: "p3"},
	}
	pdr := &v1.ProjectDeleteRequest{
		Id: "p3",
	}
	filter := map[string]any{
		"projectmember ->> 'project_id'": p3.Meta.Id,
	}
	var paging *v1.Paging

	projectMemberStorageMock.On("Find", ctx, filter, paging).Return([]*v1.ProjectMember{}, nil, nil)
	storageMock.On("DeleteAll", ctx, p3.Meta.Id).Return(nil)
	storageMock.On("Delete", ctx, p3.Meta.Id).Return(nil)
	resp, err := ps.Delete(ctx, pdr)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.GetProject())
	assert.Equal(t, pdr.Id, resp.GetProject().GetMeta().GetId())
}

func TestGetProject(t *testing.T) {
	storageMock := &mocks.Storage[*v1.Project]{}
	tenantStorageMock := &mocks.Storage[*v1.Tenant]{}
	ts := &projectService{
		projectStore: storageMock,
		tenantStore:  tenantStorageMock,
		log:          slog.Default(),
	}
	ctx := context.Background()
	t4 := &v1.Project{
		Meta: &v1.Meta{Id: "p4"},
	}
	tgr := &v1.ProjectGetRequest{
		Id: "p4",
	}

	storageMock.On("Get", ctx, "p4").Return(t4, nil)
	resp, err := ts.Get(ctx, tgr)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.GetProject())
	assert.Equal(t, tgr.Id, resp.GetProject().GetMeta().GetId())
}

func TestFindProjectByID(t *testing.T) {
	storageMock := &mocks.Storage[*v1.Project]{}
	tenantStorageMock := &mocks.Storage[*v1.Tenant]{}
	ts := &projectService{
		projectStore: storageMock,
		tenantStore:  tenantStorageMock,
		log:          slog.Default(),
	}
	ctx := context.Background()
	var t5s []*v1.Project
	// filter by id
	f1 := make(map[string]any)
	tfr := &v1.ProjectFindRequest{
		Id: &wrapperspb.StringValue{Value: "p5"},
	}

	f1["id"] = "p5"
	storageMock.On("Find", ctx, f1, mock.AnythingOfType("*v1.Paging")).Return(t5s, nil, nil)
	resp, err := ts.Find(ctx, tfr)
	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestFindProjectByName(t *testing.T) {
	storageMock := &mocks.Storage[*v1.Project]{}
	tenantStorageMock := &mocks.Storage[*v1.Tenant]{}
	ts := &projectService{
		projectStore: storageMock,
		tenantStore:  tenantStorageMock,
		log:          slog.Default(),
	}
	ctx := context.Background()

	// filter by name
	var t6s []*v1.Project
	tfr := &v1.ProjectFindRequest{
		Name: &wrapperspb.StringValue{Value: "Sixth"},
	}

	f2 := make(map[string]any)
	f2["project ->> 'name'"] = "Sixth"
	storageMock.On("Find", ctx, f2, mock.AnythingOfType("*v1.Paging")).Return(t6s, nil, nil)
	resp, err := ts.Find(ctx, tfr)
	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestFindProjectByTenant(t *testing.T) {
	storageMock := &mocks.Storage[*v1.Project]{}
	tenantStorageMock := &mocks.Storage[*v1.Tenant]{}
	ts := &projectService{
		projectStore: storageMock,
		tenantStore:  tenantStorageMock,
		log:          slog.Default(),
	}
	ctx := context.Background()

	// filter by name
	var t6s []*v1.Project
	tfr := &v1.ProjectFindRequest{
		TenantId: &wrapperspb.StringValue{Value: "p1"},
	}

	f2 := make(map[string]any)
	f2["project ->> 'tenant_id'"] = "p1"
	storageMock.On("Find", ctx, f2, mock.AnythingOfType("*v1.Paging")).Return(t6s, nil, nil)
	resp, err := ts.Find(ctx, tfr)
	require.NoError(t, err)
	assert.NotNil(t, resp)
}
