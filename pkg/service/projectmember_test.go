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
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/runtime/protoimpl"

	"testing"

	"github.com/metal-stack/masterdata-api/pkg/datastore"
	"github.com/metal-stack/masterdata-api/pkg/test/mocks"
)

func TestCreateProjectMember(t *testing.T) {
	storageMock := mocks.NewMockStorage[*v1.ProjectMember](t)
	tenantStorageMock := mocks.NewMockStorage[*v1.Tenant](t)
	projectStorageMock := mocks.NewMockStorage[*v1.Project](t)
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
	resp, err := ts.Create(ctx, pmcr)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.GetProjectMember())
	assert.Equal(t, pmcr.ProjectMember.ProjectId, resp.GetProjectMember().GetProjectId())
}

func TestDeleteProjectMember(t *testing.T) {
	storageMock := mocks.NewMockStorage[*v1.ProjectMember](t)
	tenantStorageMock := mocks.NewMockStorage[*v1.Tenant](t)
	projectStorageMock := mocks.NewMockStorage[*v1.Project](t)
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
	resp, err := ts.Delete(ctx, tdr)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.GetProjectMember())
	assert.Equal(t, tdr.Id, resp.GetProjectMember().GetMeta().GetId())
}

func TestGetProjectMember(t *testing.T) {
	storageMock := mocks.NewMockStorage[*v1.ProjectMember](t)
	tenantStorageMock := mocks.NewMockStorage[*v1.Tenant](t)
	projectStorageMock := mocks.NewMockStorage[*v1.Project](t)
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
	resp, err := ts.Get(ctx, tgr)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.GetProjectMember())
	assert.Equal(t, tgr.Id, resp.GetProjectMember().GetMeta().GetId())
}

func TestFindProjectMember(t *testing.T) {
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
		projectMemberStore = datastore.New(log, db, &v1.ProjectMember{})
		projectStore       = datastore.New(log, db, &v1.Project{})
		tenantStore        = datastore.New(log, db, &v1.Tenant{})

		testTenant1 = &v1.Tenant{
			Meta: &v1.Meta{
				Id:         "tenant-1",
				Kind:       "Tenant",
				Apiversion: "v1",
				Version:    1,
			},
			Name:        "tenant 1",
			Description: "test tenant 1",
		}
		testTenant2 = &v1.Tenant{
			Meta: &v1.Meta{
				Id:         "tenant-2",
				Kind:       "Tenant",
				Apiversion: "v1",
				Version:    1,
			},
			Name:        "tenant 2",
			Description: "test tenant 2",
		}
		testProject1 = &v1.Project{
			Meta: &v1.Meta{
				Id:         "project-1",
				Kind:       "Project",
				Apiversion: "v1",
				Version:    1,
			},
			Name:        "project 1",
			Description: "test project 1",
			TenantId:    "tenant-1",
		}
		testProjectMember1 = &v1.ProjectMember{
			Meta: &v1.Meta{
				Id:         "1",
				Kind:       "ProjectMember",
				Apiversion: "v1",
				Version:    1,
				Annotations: map[string]string{
					"role": "owner",
				},
				Labels: []string{"a", "b"},
			},
			ProjectId: "project-1",
			TenantId:  "tenant-1",
			Namespace: "a",
		}
		testProjectMember2 = &v1.ProjectMember{
			Meta: &v1.Meta{
				Id:         "2",
				Kind:       "ProjectMember",
				Apiversion: "v1",
				Version:    1,
				Annotations: map[string]string{
					"role": "viewer",
				},
				Labels: []string{"c", "d"},
			},
			ProjectId: "project-1",
			TenantId:  "tenant-2",
			Namespace: "a",
		}
		testProjectMember3 = &v1.ProjectMember{
			Meta: &v1.Meta{
				Id:         "3",
				Kind:       "ProjectMember",
				Apiversion: "v1",
				Version:    1,
				Annotations: map[string]string{
					"role": "owner",
				},
				Labels: []string{"e", "f"},
			},
			ProjectId: "project-2",
			TenantId:  "tenant-2",
			Namespace: "a",
		}
		testProjectMember4 = &v1.ProjectMember{
			Meta: &v1.Meta{
				Id:         "4",
				Kind:       "ProjectMember",
				Apiversion: "v1",
				Version:    1,
				Annotations: map[string]string{
					"role": "owner",
				},
			},
			ProjectId: "project-2",
			TenantId:  "tenant-2",
			Namespace: "",
		}

		service = &projectMemberService{
			log:                log,
			projectMemberStore: projectMemberStore,
			tenantStore:        tenantStore,
			projectStore:       projectStore,
		}
	)

	tests := []struct {
		name    string
		prepare func()
		req     *v1.ProjectMemberFindRequest
		want    *v1.ProjectMemberListResponse
		wantErr error
	}{
		{
			name: "find by project",
			req: &v1.ProjectMemberFindRequest{
				ProjectId: pointer.Pointer("project-1"),
				Namespace: "a",
			},
			prepare: func() {
				require.NoError(t, tenantStore.Create(ctx, testTenant1))
				require.NoError(t, tenantStore.Create(ctx, testTenant2))
				require.NoError(t, projectStore.Create(ctx, testProject1))
				require.NoError(t, projectMemberStore.Create(ctx, testProjectMember1))
				require.NoError(t, projectMemberStore.Create(ctx, testProjectMember2))
				require.NoError(t, projectMemberStore.Create(ctx, testProjectMember3))
				require.NoError(t, projectMemberStore.Create(ctx, testProjectMember4))
			},
			want: &v1.ProjectMemberListResponse{
				ProjectMembers: []*v1.ProjectMember{
					testProjectMember1,
					testProjectMember2,
				},
			},
			wantErr: nil,
		},
		{
			name: "find by project id (no results) #1",
			req: &v1.ProjectMemberFindRequest{
				ProjectId: pointer.Pointer("no-result"),
				Namespace: "a",
			},
			prepare: func() {
				require.NoError(t, tenantStore.Create(ctx, testTenant1))
				require.NoError(t, tenantStore.Create(ctx, testTenant2))
				require.NoError(t, projectStore.Create(ctx, testProject1))
				require.NoError(t, projectMemberStore.Create(ctx, testProjectMember1))
				require.NoError(t, projectMemberStore.Create(ctx, testProjectMember2))
				require.NoError(t, projectMemberStore.Create(ctx, testProjectMember3))
				require.NoError(t, projectMemberStore.Create(ctx, testProjectMember4))
			},
			want: &v1.ProjectMemberListResponse{
				ProjectMembers: nil,
			},
			wantErr: nil,
		},
		{
			name: "find by project id (no results) #2",
			req: &v1.ProjectMemberFindRequest{
				ProjectId: pointer.Pointer("project-1"),
				Namespace: "wrong-namespace",
			},
			prepare: func() {
				require.NoError(t, tenantStore.Create(ctx, testTenant1))
				require.NoError(t, tenantStore.Create(ctx, testTenant2))
				require.NoError(t, projectStore.Create(ctx, testProject1))
				require.NoError(t, projectMemberStore.Create(ctx, testProjectMember1))
				require.NoError(t, projectMemberStore.Create(ctx, testProjectMember2))
				require.NoError(t, projectMemberStore.Create(ctx, testProjectMember3))
				require.NoError(t, projectMemberStore.Create(ctx, testProjectMember4))
			},
			want: &v1.ProjectMemberListResponse{
				ProjectMembers: nil,
			},
			wantErr: nil,
		},
		{
			name: "find by tenant",
			req: &v1.ProjectMemberFindRequest{
				TenantId:  pointer.Pointer("tenant-2"),
				Namespace: "a",
			},
			prepare: func() {
				require.NoError(t, tenantStore.Create(ctx, testTenant1))
				require.NoError(t, tenantStore.Create(ctx, testTenant2))
				require.NoError(t, projectStore.Create(ctx, testProject1))
				require.NoError(t, projectMemberStore.Create(ctx, testProjectMember1))
				require.NoError(t, projectMemberStore.Create(ctx, testProjectMember2))
				require.NoError(t, projectMemberStore.Create(ctx, testProjectMember3))
				require.NoError(t, projectMemberStore.Create(ctx, testProjectMember4))
			},
			want: &v1.ProjectMemberListResponse{
				ProjectMembers: []*v1.ProjectMember{
					testProjectMember2,
					testProjectMember3,
				},
			},
			wantErr: nil,
		},
		{
			name: "find by annotation",
			req: &v1.ProjectMemberFindRequest{
				Annotations: map[string]string{"role": "owner"},
				Namespace:   "a",
			},
			prepare: func() {
				require.NoError(t, tenantStore.Create(ctx, testTenant1))
				require.NoError(t, tenantStore.Create(ctx, testTenant2))
				require.NoError(t, projectStore.Create(ctx, testProject1))
				require.NoError(t, projectMemberStore.Create(ctx, testProjectMember1))
				require.NoError(t, projectMemberStore.Create(ctx, testProjectMember2))
				require.NoError(t, projectMemberStore.Create(ctx, testProjectMember3))
				require.NoError(t, projectMemberStore.Create(ctx, testProjectMember4))
			},
			want: &v1.ProjectMemberListResponse{
				ProjectMembers: []*v1.ProjectMember{
					testProjectMember1,
					testProjectMember3,
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

			slices.SortFunc(got.ProjectMembers, func(i, j *v1.ProjectMember) int {
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

func TestUpdateProjectMember(t *testing.T) {
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
		projectMemberStore = datastore.New(log, db, &v1.ProjectMember{})
		projectStore       = datastore.New(log, db, &v1.Project{})
		tenantStore        = datastore.New(log, db, &v1.Tenant{})

		service = &projectMemberService{
			log:                log,
			projectMemberStore: projectMemberStore,
			tenantStore:        tenantStore,
			projectStore:       projectStore,
		}
	)

	tests := []struct {
		name    string
		prepare func()
		req     *v1.ProjectMemberUpdateRequest
		want    *v1.ProjectMemberResponse
		wantErr error
	}{
		{
			name: "update mutable fields",
			req: &v1.ProjectMemberUpdateRequest{
				ProjectMember: &v1.ProjectMember{
					Meta: &v1.Meta{
						Id:      "1",
						Version: 1,
						Annotations: map[string]string{
							"role": "owner",
						},
						Labels: []string{"a", "b"},
					},
					ProjectId: "project-1",
					TenantId:  "tenant-1",
					Namespace: "a",
				},
			},
			prepare: func() {
				require.NoError(t, projectMemberStore.Create(ctx, &v1.ProjectMember{
					Meta: &v1.Meta{
						Id:         "1",
						Kind:       "ProjectMember",
						Apiversion: "v1",
						Version:    1,
					},
					ProjectId: "project-1",
					TenantId:  "tenant-1",
					Namespace: "a",
				}))
			},
			want: &v1.ProjectMemberResponse{
				ProjectMember: &v1.ProjectMember{
					Meta: &v1.Meta{
						Id:         "1",
						Kind:       "ProjectMember",
						Apiversion: "v1",
						Version:    2,
						Annotations: map[string]string{
							"role": "owner",
						},
						Labels: []string{"a", "b"},
					},
					ProjectId: "project-1",
					TenantId:  "tenant-1",
					Namespace: "a",
				},
			},
			wantErr: nil,
		},
		{
			name: "unable to update namespace",
			req: &v1.ProjectMemberUpdateRequest{
				ProjectMember: &v1.ProjectMember{
					Meta: &v1.Meta{
						Id:      "1",
						Version: 1,
					},
					ProjectId: "project-1",
					TenantId:  "tenant-1",
					Namespace: "b",
				},
			},
			prepare: func() {
				require.NoError(t, projectMemberStore.Create(ctx, &v1.ProjectMember{
					Meta: &v1.Meta{
						Id:         "1",
						Kind:       "ProjectMember",
						Apiversion: "v1",
						Version:    1,
					},
					ProjectId: "project-1",
					TenantId:  "tenant-1",
					Namespace: "a",
				}))
			},
			want:    nil,
			wantErr: status.Error(codes.InvalidArgument, "updating the namespace of a project member is not allowed"),
		},
		{
			name: "unable to update project",
			req: &v1.ProjectMemberUpdateRequest{
				ProjectMember: &v1.ProjectMember{
					Meta: &v1.Meta{
						Id:      "1",
						Version: 1,
					},
					ProjectId: "project-2",
					TenantId:  "tenant-1",
					Namespace: "a",
				},
			},
			prepare: func() {
				require.NoError(t, projectMemberStore.Create(ctx, &v1.ProjectMember{
					Meta: &v1.Meta{
						Id:         "1",
						Kind:       "ProjectMember",
						Apiversion: "v1",
						Version:    1,
					},
					ProjectId: "project-1",
					TenantId:  "tenant-1",
					Namespace: "a",
				}))
			},
			want:    nil,
			wantErr: status.Error(codes.InvalidArgument, "updating the project id of a project member is not allowed"),
		},
		{
			name: "unable to update tenant",
			req: &v1.ProjectMemberUpdateRequest{
				ProjectMember: &v1.ProjectMember{
					Meta: &v1.Meta{
						Id:      "1",
						Version: 1,
					},
					ProjectId: "project-1",
					TenantId:  "tenant-2",
					Namespace: "a",
				},
			},
			prepare: func() {
				require.NoError(t, projectMemberStore.Create(ctx, &v1.ProjectMember{
					Meta: &v1.Meta{
						Id:         "1",
						Kind:       "ProjectMember",
						Apiversion: "v1",
						Version:    1,
					},
					ProjectId: "project-1",
					TenantId:  "tenant-1",
					Namespace: "a",
				}))
			},
			want:    nil,
			wantErr: status.Error(codes.InvalidArgument, "updating the tenant id of a project member is not allowed"),
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

			got, err := service.Update(ctx, tt.req)
			if diff := cmp.Diff(err, tt.wantErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
				return
			}

			if err == nil {
				assert.NotNil(t, got.ProjectMember.Meta.UpdatedTime)
			}

			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreTypes(protoimpl.MessageState{}), cmpopts.IgnoreFields(v1.Meta{}, "CreatedTime", "UpdatedTime"), testcommon.IgnoreUnexported()); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
			}
		})
	}
}
