package service

import (
	"context"
	"log/slog"
	"slices"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/metal-stack/metal-lib/pkg/testcommon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/runtime/protoimpl"

	"testing"

	"github.com/metal-stack/masterdata-api/pkg/datastore"
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
				Max: pointer.Pointer(int32(2)),
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
	storageMock.On("Find", ctx, mock.AnythingOfType("*v1.Paging"), filter).Return(projects, nil, nil)
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

	projectMemberStorageMock.On("Find", ctx, paging, filter).Return([]*v1.ProjectMember{}, nil, nil)
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

func TestFindProject(t *testing.T) {
	ctx := t.Context()
	ves := []datastore.Entity{
		&v1.Project{},
		&v1.ProjectMember{},
		&v1.Tenant{},
		&v1.TenantMember{},
	}

	container, db, err := StartPostgres(ctx, ves...)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, db.Close())
		require.NoError(t, container.Terminate(ctx))
	}()

	var (
		projectStore = datastore.New(log, db, &v1.Project{})
		testProject1 = &v1.Project{
			Meta: &v1.Meta{
				Id:         "1",
				Kind:       "Project",
				Apiversion: "v1",
				Version:    1,
				Annotations: map[string]string{
					"a": "b",
					"c": "d",
				},
				Labels: []string{"e", "f"},
			},
			Name:        "project-1",
			Description: "project 1",
			TenantId:    "tenant-1",
		}
		testProject2 = &v1.Project{
			Meta: &v1.Meta{
				Id:         "2",
				Kind:       "Project",
				Apiversion: "v1",
				Version:    1,
				Annotations: map[string]string{
					"c": "d",
					"e": "f",
				},
				Labels: []string{"f", "g", "h"},
			},
			Name:        "project-2",
			Description: "project 2",
			TenantId:    "tenant-2",
		}

		service = &projectService{
			projectStore: projectStore,
			log:          log,
		}
	)

	tests := []struct {
		name    string
		prepare func()
		req     *v1.ProjectFindRequest
		want    *v1.ProjectListResponse
		wantErr error
	}{
		{
			name: "find by id",
			req: &v1.ProjectFindRequest{
				Id: pointer.Pointer("1"),
			},
			prepare: func() {
				require.NoError(t, projectStore.Create(ctx, testProject1))
				require.NoError(t, projectStore.Create(ctx, testProject2))
			},
			want: &v1.ProjectListResponse{
				Projects: []*v1.Project{
					testProject1,
				},
			},
			wantErr: nil,
		},
		{
			name: "find by id (no results)",
			req: &v1.ProjectFindRequest{
				Id: pointer.Pointer("no-result"),
			},
			prepare: func() {
				require.NoError(t, projectStore.Create(ctx, testProject1))
				require.NoError(t, projectStore.Create(ctx, testProject2))
			},
			want: &v1.ProjectListResponse{
				Projects: nil,
			},
			wantErr: nil,
		},
		{
			name: "find by name",
			req: &v1.ProjectFindRequest{
				Name: pointer.Pointer("project-2"),
			},
			prepare: func() {
				require.NoError(t, projectStore.Create(ctx, testProject1))
				require.NoError(t, projectStore.Create(ctx, testProject2))
			},
			want: &v1.ProjectListResponse{
				Projects: []*v1.Project{
					testProject2,
				},
			},
			wantErr: nil,
		},
		{
			name: "find by tenant",
			req: &v1.ProjectFindRequest{
				TenantId: pointer.Pointer("tenant-2"),
			},
			prepare: func() {
				require.NoError(t, projectStore.Create(ctx, testProject1))
				require.NoError(t, projectStore.Create(ctx, testProject2))
			},
			want: &v1.ProjectListResponse{
				Projects: []*v1.Project{
					testProject2,
				},
			},
			wantErr: nil,
		},
		{
			name: "find by annotation",
			req: &v1.ProjectFindRequest{
				Annotations: map[string]string{
					"a": "b",
				},
			},
			prepare: func() {
				require.NoError(t, projectStore.Create(ctx, testProject1))
				require.NoError(t, projectStore.Create(ctx, testProject2))
			},
			want: &v1.ProjectListResponse{
				Projects: []*v1.Project{
					testProject1,
				},
			},
			wantErr: nil,
		},
		{
			name: "find by annotation #2",
			req: &v1.ProjectFindRequest{
				Annotations: map[string]string{
					"a": "b",
					"c": "d",
				},
			},
			prepare: func() {
				require.NoError(t, projectStore.Create(ctx, testProject1))
				require.NoError(t, projectStore.Create(ctx, testProject2))
			},
			want: &v1.ProjectListResponse{
				Projects: []*v1.Project{
					testProject1,
				},
			},
			wantErr: nil,
		},
		{
			name: "find by annotation #3",
			req: &v1.ProjectFindRequest{
				Annotations: map[string]string{
					"c": "d",
				},
			},
			prepare: func() {
				require.NoError(t, projectStore.Create(ctx, testProject1))
				require.NoError(t, projectStore.Create(ctx, testProject2))
			},
			want: &v1.ProjectListResponse{
				Projects: []*v1.Project{
					testProject1,
					testProject2,
				},
			},
			wantErr: nil,
		},
		{
			name: "find by label",
			req: &v1.ProjectFindRequest{
				Labels: []string{"e"},
			},
			prepare: func() {
				require.NoError(t, projectStore.Create(ctx, testProject1))
				require.NoError(t, projectStore.Create(ctx, testProject2))
			},
			want: &v1.ProjectListResponse{
				Projects: []*v1.Project{
					testProject1,
				},
			},
			wantErr: nil,
		},
		{
			name: "find by label #2",
			req: &v1.ProjectFindRequest{
				Labels: []string{"e", "f"},
			},
			prepare: func() {
				require.NoError(t, projectStore.Create(ctx, testProject1))
				require.NoError(t, projectStore.Create(ctx, testProject2))
			},
			want: &v1.ProjectListResponse{
				Projects: []*v1.Project{
					testProject1,
				},
			},
			wantErr: nil,
		},
		{
			name: "find by label #3",
			req: &v1.ProjectFindRequest{
				Labels: []string{"f"},
			},
			prepare: func() {
				require.NoError(t, projectStore.Create(ctx, testProject1))
				require.NoError(t, projectStore.Create(ctx, testProject2))
			},
			want: &v1.ProjectListResponse{
				Projects: []*v1.Project{
					testProject1,
					testProject2,
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, ve := range ves {
				_, err := db.ExecContext(ctx, "TRUNCATE TABLE "+ve.TableName())
				require.NoError(t, err)
			}

			if tt.prepare != nil {
				tt.prepare()
			}

			got, err := service.Find(ctx, tt.req)
			if diff := cmp.Diff(err, tt.wantErr); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
				return
			}
			slices.SortFunc(got.Projects, func(i, j *v1.Project) int {
				if i.Meta.Id < j.Meta.Id {
					return -1
				} else {
					return 1
				}
			})
			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreTypes(protoimpl.MessageState{}), cmpopts.IgnoreFields(v1.Meta{}, "CreatedTime"), testcommon.IgnoreUnexported()); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
			}
		})
	}
}
