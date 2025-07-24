package service

import (
	"context"
	"log/slog"
	"os"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/metal-stack/metal-lib/pkg/testcommon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/runtime/protoimpl"

	"github.com/metal-stack/masterdata-api/pkg/datastore"
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
	memberStorageMock := &mocks.Storage[*v1.TenantMember]{}
	ts := &tenantService{
		tenantStore:       storageMock,
		tenantMemberStore: memberStorageMock,
		log:               slog.Default(),
	}
	ctx := context.Background()
	t3 := &v1.Tenant{
		Meta: &v1.Meta{Id: "t3"},
	}
	tdr := &v1.TenantDeleteRequest{
		Id: "t3",
	}
	tfilter := map[string]any{
		"tenantmember ->> 'tenant_id'": t3.Meta.Id,
	}
	mfilter := map[string]any{
		"tenantmember ->> 'member_id'": t3.Meta.Id,
	}
	var paging *v1.Paging

	storageMock.On("Delete", ctx, t3.Meta.Id).Return(nil)
	memberStorageMock.On("Find", ctx, tfilter, paging).Return([]*v1.TenantMember{
		{
			Meta: &v1.Meta{
				Id: "t3",
			},
			TenantId: t3.Meta.Id,
			MemberId: t3.Meta.Id,
		},
	}, nil, nil)
	memberStorageMock.On("Find", ctx, mfilter, paging).Return([]*v1.TenantMember{
		{
			Meta: &v1.Meta{
				Id: "t3",
			},
			TenantId: t3.Meta.Id,
			MemberId: t3.Meta.Id,
		},
	}, nil, nil)
	memberStorageMock.On("DeleteAll", ctx, "t3").Return(nil)
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
		Id: pointer.Pointer("t5"),
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
		Name: pointer.Pointer("Fifth"),
	}

	f2 := make(map[string]any)
	f2["tenant ->> 'name'"] = "Fifth"
	storageMock.On("Find", ctx, f2, mock.AnythingOfType("*v1.Paging")).Return(t6s, nil, nil)
	resp, err := ts.Find(ctx, tfr)
	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func Test_tenantService_FindParticipatingProjects(t *testing.T) {
	ctx := context.Background()
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

	s := &tenantService{
		db:  db,
		log: slog.Default(),
	}

	var (
		projectStore       = datastore.New(log, db, &v1.Project{})
		tenantMemberStore  = datastore.New(log, db, &v1.TenantMember{})
		projectMemberStore = datastore.New(log, db, &v1.ProjectMember{})
	)

	tests := []struct {
		name    string
		prepare func()
		req     *v1.FindParticipatingProjectsRequest
		want    *v1.FindParticipatingProjectsResponse
		wantErr error
	}{
		{
			name: "no memberships",
			req: &v1.FindParticipatingProjectsRequest{
				TenantId:         "a",
				IncludeInherited: pointer.Pointer(true),
			},
			prepare: func() {
			},
			want:    &v1.FindParticipatingProjectsResponse{},
			wantErr: nil,
		},
		{
			name: "ignores foreign memberships",
			req: &v1.FindParticipatingProjectsRequest{
				TenantId:         "a",
				IncludeInherited: pointer.Pointer(true),
			},
			prepare: func() {
				err := projectStore.Create(ctx, &v1.Project{Meta: &v1.Meta{Id: "1"}})
				require.NoError(t, err)
				err = projectMemberStore.Create(ctx, &v1.ProjectMember{Meta: &v1.Meta{Annotations: map[string]string{"role": "admin"}}, ProjectId: "1", TenantId: "someone else"})
				require.NoError(t, err)
			},
			want:    &v1.FindParticipatingProjectsResponse{},
			wantErr: nil,
		},
		{
			name: "direct membership including 0 inherited",
			req: &v1.FindParticipatingProjectsRequest{
				TenantId:         "a",
				IncludeInherited: pointer.Pointer(true),
			},
			prepare: func() {
				err := projectStore.Create(ctx, &v1.Project{Meta: &v1.Meta{Id: "1"}})
				require.NoError(t, err)
				err = projectMemberStore.Create(ctx, &v1.ProjectMember{Meta: &v1.Meta{Annotations: map[string]string{"role": "admin"}}, ProjectId: "1", TenantId: "a"})
				require.NoError(t, err)
			},
			want: &v1.FindParticipatingProjectsResponse{
				Projects: []*v1.ProjectWithMembershipAnnotations{{
					Project: &v1.Project{
						Meta: &v1.Meta{
							Kind:       "Project",
							Apiversion: "v1",
							Id:         "1",
						},
					},
					ProjectAnnotations: map[string]string{"role": "admin"},
					TenantAnnotations:  nil,
				}},
			},
			wantErr: nil,
		},
		{
			name: "direct membership excluding inherited",
			req: &v1.FindParticipatingProjectsRequest{
				TenantId:         "a",
				IncludeInherited: pointer.Pointer(false),
			},
			prepare: func() {
				err := projectStore.Create(ctx, &v1.Project{Meta: &v1.Meta{Id: "1"}})
				require.NoError(t, err)
				err = projectStore.Create(ctx, &v1.Project{Meta: &v1.Meta{Id: "2"}, TenantId: "b"})
				require.NoError(t, err)
				err = projectMemberStore.Create(ctx, &v1.ProjectMember{Meta: &v1.Meta{Annotations: map[string]string{"role": "admin"}}, ProjectId: "1", TenantId: "a"})
				require.NoError(t, err)
				err = projectMemberStore.Create(ctx, &v1.ProjectMember{Meta: &v1.Meta{Annotations: map[string]string{"role": "admin"}}, ProjectId: "2", TenantId: "b"})
				require.NoError(t, err)
				err = tenantMemberStore.Create(ctx, &v1.TenantMember{Meta: &v1.Meta{Annotations: map[string]string{"role": "editor"}}, MemberId: "a", TenantId: "b"})
				require.NoError(t, err)
			},
			want: &v1.FindParticipatingProjectsResponse{
				Projects: []*v1.ProjectWithMembershipAnnotations{{
					Project: &v1.Project{
						Meta: &v1.Meta{
							Kind:       "Project",
							Apiversion: "v1",
							Id:         "1",
						},
					},
					ProjectAnnotations: map[string]string{"role": "admin"},
					TenantAnnotations:  nil,
				}},
			},
			wantErr: nil,
		},
		{
			name: "inherited membership",
			req: &v1.FindParticipatingProjectsRequest{
				TenantId:         "a",
				IncludeInherited: pointer.Pointer(true),
			},
			prepare: func() {
				err := projectStore.Create(ctx, &v1.Project{Meta: &v1.Meta{Id: "1"}, TenantId: "b"})
				require.NoError(t, err)
				err = tenantMemberStore.Create(ctx, &v1.TenantMember{Meta: &v1.Meta{Annotations: map[string]string{"role": "viewer"}}, TenantId: "b", MemberId: "a"})
				require.NoError(t, err)
			},
			want: &v1.FindParticipatingProjectsResponse{
				Projects: []*v1.ProjectWithMembershipAnnotations{{
					Project: &v1.Project{
						Meta: &v1.Meta{
							Kind:       "Project",
							Apiversion: "v1",
							Id:         "1",
						},
						TenantId: "b",
					},
					ProjectAnnotations: nil,
					TenantAnnotations:  map[string]string{"role": "viewer"},
				}},
			},
			wantErr: nil,
		},
		{
			name: "direct and indirect memberships including inherited",
			req: &v1.FindParticipatingProjectsRequest{
				TenantId:         "req-tenant",
				IncludeInherited: pointer.Pointer(true),
			},
			prepare: func() {
				err := projectStore.Create(ctx, &v1.Project{
					Meta:     &v1.Meta{Id: "direct-1"},
					TenantId: "req-tenant",
				})
				require.NoError(t, err)
				err = projectMemberStore.Create(ctx, &v1.ProjectMember{
					Meta:      &v1.Meta{Annotations: map[string]string{"role": "owner"}},
					ProjectId: "direct-1",
					TenantId:  "req-tenant",
				})
				require.NoError(t, err)
				err = tenantMemberStore.Create(ctx, &v1.TenantMember{
					Meta:     &v1.Meta{Annotations: map[string]string{"role": "editor"}},
					MemberId: "req-tenant",
					TenantId: "parent",
				})
				require.NoError(t, err)
				err = projectStore.Create(ctx, &v1.Project{
					Meta:     &v1.Meta{Id: "indirect-2"},
					TenantId: "parent",
				})
				require.NoError(t, err)
				err = projectMemberStore.Create(ctx, &v1.ProjectMember{
					Meta:      &v1.Meta{Annotations: map[string]string{"role": "admin"}},
					ProjectId: "indirect-2",
					TenantId:  "parent",
				})
				require.NoError(t, err)
			},
			want: &v1.FindParticipatingProjectsResponse{
				Projects: []*v1.ProjectWithMembershipAnnotations{
					{
						Project: &v1.Project{
							Meta: &v1.Meta{
								Kind:       "Project",
								Apiversion: "v1",
								Id:         "direct-1",
							},
							TenantId: "req-tenant",
						},
						ProjectAnnotations: map[string]string{"role": "owner"},
						TenantAnnotations:  nil,
					},
					{
						Project: &v1.Project{
							Meta: &v1.Meta{
								Kind:       "Project",
								Apiversion: "v1",
								Id:         "indirect-2",
							},
							TenantId: "parent",
						},
						ProjectAnnotations: nil,
						TenantAnnotations:  map[string]string{"role": "editor"},
					},
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

			got, err := s.FindParticipatingProjects(ctx, tt.req)
			if diff := cmp.Diff(err, tt.wantErr); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
				return
			}
			slices.SortFunc(got.Projects, func(i, j *v1.ProjectWithMembershipAnnotations) int {
				if i.Project.Meta.Id < j.Project.Meta.Id {
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

func Test_tenantService_FindParticipatingTenants(t *testing.T) {
	ctx := context.Background()
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

	s := &tenantService{
		db:  db,
		log: slog.Default(),
	}

	var (
		projectStore       = datastore.New(log, db, &v1.Project{})
		tenantMemberStore  = datastore.New(log, db, &v1.TenantMember{})
		projectMemberStore = datastore.New(log, db, &v1.ProjectMember{})
		tenantStore        = datastore.New(log, db, &v1.Tenant{})
	)

	tests := []struct {
		name    string
		req     *v1.FindParticipatingTenantsRequest
		prepare func()
		want    *v1.FindParticipatingTenantsResponse
		wantErr error
	}{
		{
			name: "no memberships",
			req: &v1.FindParticipatingTenantsRequest{
				TenantId:         "a",
				IncludeInherited: pointer.Pointer(true),
			},
			prepare: func() {},
			want:    &v1.FindParticipatingTenantsResponse{},
			wantErr: nil,
		},
		{
			name: "ignore foreign memberships",
			req: &v1.FindParticipatingTenantsRequest{
				TenantId:         "a",
				IncludeInherited: pointer.Pointer(true),
			},
			prepare: func() {
				err := tenantStore.Create(ctx, &v1.Tenant{Meta: &v1.Meta{Id: "a"}})
				require.NoError(t, err)
				err = tenantStore.Create(ctx, &v1.Tenant{Meta: &v1.Meta{Id: "b"}})
				require.NoError(t, err)
				err = tenantStore.Create(ctx, &v1.Tenant{Meta: &v1.Meta{Id: "c"}})
				require.NoError(t, err)
				err = tenantMemberStore.Create(ctx, &v1.TenantMember{Meta: &v1.Meta{Annotations: map[string]string{"role": "admin"}}, MemberId: "c", TenantId: "b"})
				require.NoError(t, err)
			},
			want:    &v1.FindParticipatingTenantsResponse{},
			wantErr: err,
		},
		{
			name: "direct membership",
			req: &v1.FindParticipatingTenantsRequest{
				TenantId:         "a",
				IncludeInherited: pointer.Pointer(true),
			},
			prepare: func() {
				err := tenantStore.Create(ctx, &v1.Tenant{Meta: &v1.Meta{Id: "b"}})
				require.NoError(t, err)
				err = tenantMemberStore.Create(ctx, &v1.TenantMember{Meta: &v1.Meta{Annotations: map[string]string{"role": "admin"}}, MemberId: "a", TenantId: "b"})
				require.NoError(t, err)
			},
			want: &v1.FindParticipatingTenantsResponse{
				Tenants: []*v1.TenantWithMembershipAnnotations{
					{
						Tenant: &v1.Tenant{
							Meta: &v1.Meta{
								Kind:       "Tenant",
								Apiversion: "v1",
								Id:         "b",
							},
						},
						TenantAnnotations: map[string]string{"role": "admin"},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "indirect membership",
			req: &v1.FindParticipatingTenantsRequest{
				TenantId:         "a",
				IncludeInherited: pointer.Pointer(true),
			},
			prepare: func() {
				err := projectStore.Create(ctx, &v1.Project{Meta: &v1.Meta{Id: "1"}, TenantId: "b"})
				require.NoError(t, err)
				err = tenantStore.Create(ctx, &v1.Tenant{Meta: &v1.Meta{Id: "b"}})
				require.NoError(t, err)
				err = projectMemberStore.Create(ctx, &v1.ProjectMember{Meta: &v1.Meta{Annotations: map[string]string{"role": "admin"}}, ProjectId: "1", TenantId: "a"})
				require.NoError(t, err)
			},
			want: &v1.FindParticipatingTenantsResponse{
				Tenants: []*v1.TenantWithMembershipAnnotations{
					{
						Tenant: &v1.Tenant{
							Meta: &v1.Meta{
								Kind:       "Tenant",
								Apiversion: "v1",
								Id:         "b",
							},
						},
						ProjectAnnotations: map[string]string{"role": "admin"},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "exclude inherited",
			req: &v1.FindParticipatingTenantsRequest{
				TenantId:         "a",
				IncludeInherited: pointer.Pointer(false),
			},
			prepare: func() {
				err := projectStore.Create(ctx, &v1.Project{Meta: &v1.Meta{Id: "1"}, TenantId: "b"})
				require.NoError(t, err)
				err = tenantStore.Create(ctx, &v1.Tenant{Meta: &v1.Meta{Id: "b"}})
				require.NoError(t, err)
				err = projectMemberStore.Create(ctx, &v1.ProjectMember{Meta: &v1.Meta{Annotations: map[string]string{"role": "admin"}}, ProjectId: "1", TenantId: "a"})
				require.NoError(t, err)
			},
			want:    &v1.FindParticipatingTenantsResponse{},
			wantErr: nil,
		},
		{
			name: "direct and indirect memberships",
			req: &v1.FindParticipatingTenantsRequest{
				TenantId:         "req-tnt",
				IncludeInherited: pointer.Pointer(true),
			},
			prepare: func() {
				err = tenantStore.Create(ctx, &v1.Tenant{Meta: &v1.Meta{Id: "indirect-tnt"}})
				require.NoError(t, err)
				err := projectStore.Create(ctx, &v1.Project{
					Meta:     &v1.Meta{Id: "indirect"},
					TenantId: "indirect-tnt",
				})
				require.NoError(t, err)
				err = projectMemberStore.Create(ctx, &v1.ProjectMember{
					Meta:      &v1.Meta{Annotations: map[string]string{"role": "admin"}},
					ProjectId: "indirect",
					TenantId:  "req-tnt",
				})
				require.NoError(t, err)

				err = tenantStore.Create(ctx, &v1.Tenant{Meta: &v1.Meta{Id: "direct-tnt"}})
				require.NoError(t, err)
				err = tenantMemberStore.Create(ctx, &v1.TenantMember{
					Meta:     &v1.Meta{Annotations: map[string]string{"relation": "direct"}},
					TenantId: "direct-tnt",
					MemberId: "req-tnt",
				})
				require.NoError(t, err)
			},
			want: &v1.FindParticipatingTenantsResponse{
				Tenants: []*v1.TenantWithMembershipAnnotations{
					{
						Tenant: &v1.Tenant{
							Meta: &v1.Meta{
								Kind:       "Tenant",
								Apiversion: "v1",
								Id:         "direct-tnt",
							},
						},
						TenantAnnotations: map[string]string{"relation": "direct"},
					},
					{
						Tenant: &v1.Tenant{
							Meta: &v1.Meta{
								Kind:       "Tenant",
								Apiversion: "v1",
								Id:         "indirect-tnt",
							},
						},
						ProjectAnnotations: map[string]string{"role": "admin"},
					},
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

			got, err := s.FindParticipatingTenants(ctx, tt.req)
			if diff := cmp.Diff(err, tt.wantErr); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
				return
			}

			slices.SortFunc(got.Tenants, func(i, j *v1.TenantWithMembershipAnnotations) int {
				if i.Tenant.Meta.Id < j.Tenant.Meta.Id {
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

func Test_tenantService_ListTenantMembers(t *testing.T) {
	ctx := context.Background()
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

	s := &tenantService{
		db:  db,
		log: slog.Default(),
	}

	var (
		projectStore       = datastore.New(log, db, &v1.Project{})
		tenantMemberStore  = datastore.New(log, db, &v1.TenantMember{})
		projectMemberStore = datastore.New(log, db, &v1.ProjectMember{})
		tenantStore        = datastore.New(log, db, &v1.Tenant{})
	)

	tests := []struct {
		name    string
		req     *v1.ListTenantMembersRequest
		prepare func()
		want    *v1.ListTenantMembersResponse
		wantErr error
	}{
		{
			name: "no members",
			req: &v1.ListTenantMembersRequest{
				TenantId:         "acme",
				IncludeInherited: pointer.Pointer(true),
			},
			prepare: func() {
			},
			want:    &v1.ListTenantMembersResponse{},
			wantErr: err,
		},
		{
			name: "ignore foreign members",
			req: &v1.ListTenantMembersRequest{
				TenantId:         "acme",
				IncludeInherited: pointer.Pointer(true),
			},
			prepare: func() {
				err := tenantStore.Create(ctx, &v1.Tenant{Meta: &v1.Meta{Id: "acme"}})
				require.NoError(t, err)
				err = tenantStore.Create(ctx, &v1.Tenant{Meta: &v1.Meta{Id: "azure"}})
				require.NoError(t, err)
				err = tenantStore.Create(ctx, &v1.Tenant{Meta: &v1.Meta{Id: "google"}})
				require.NoError(t, err)
				err = tenantMemberStore.Create(ctx, &v1.TenantMember{Meta: &v1.Meta{Annotations: map[string]string{"role": "admin"}}, MemberId: "azure", TenantId: "google"})
				require.NoError(t, err)
			},
			want:    &v1.ListTenantMembersResponse{},
			wantErr: err,
		},
		{
			name: "direct membership",
			req: &v1.ListTenantMembersRequest{
				TenantId:         "acme",
				IncludeInherited: pointer.Pointer(true),
			},
			prepare: func() {
				err := tenantStore.Create(ctx, &v1.Tenant{Meta: &v1.Meta{Id: "azure"}})
				require.NoError(t, err)
				err = tenantMemberStore.Create(ctx, &v1.TenantMember{Meta: &v1.Meta{Annotations: map[string]string{"role": "admin"}}, MemberId: "azure", TenantId: "acme"})
				require.NoError(t, err)
			},
			want: &v1.ListTenantMembersResponse{
				Tenants: []*v1.TenantWithMembershipAnnotations{
					{
						Tenant: &v1.Tenant{
							Meta: &v1.Meta{
								Kind:       "Tenant",
								Apiversion: "v1",
								Id:         "azure",
							},
						},
						TenantAnnotations: map[string]string{"role": "admin"},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "indirect membership",
			req: &v1.ListTenantMembersRequest{
				TenantId:         "acme",
				IncludeInherited: pointer.Pointer(true),
			},
			prepare: func() {
				err := projectStore.Create(ctx, &v1.Project{Meta: &v1.Meta{Id: "1"}, TenantId: "acme"})
				require.NoError(t, err)
				err = tenantStore.Create(ctx, &v1.Tenant{Meta: &v1.Meta{Id: "google"}})
				require.NoError(t, err)
				err = projectMemberStore.Create(ctx, &v1.ProjectMember{Meta: &v1.Meta{Annotations: map[string]string{"role": "editor"}}, ProjectId: "1", TenantId: "google"})
				require.NoError(t, err)
			},
			want: &v1.ListTenantMembersResponse{
				Tenants: []*v1.TenantWithMembershipAnnotations{
					{
						Tenant: &v1.Tenant{
							Meta: &v1.Meta{
								Kind:       "Tenant",
								Apiversion: "v1",
								Id:         "google",
							},
						},
						ProjectIds: []string{
							"1",
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "exclude inherited",
			req: &v1.ListTenantMembersRequest{
				TenantId:         "acme",
				IncludeInherited: pointer.Pointer(false),
			},
			prepare: func() {
				err := projectStore.Create(ctx, &v1.Project{Meta: &v1.Meta{Id: "1"}, TenantId: "acme"})
				require.NoError(t, err)
				err = tenantStore.Create(ctx, &v1.Tenant{Meta: &v1.Meta{Id: "google"}})
				require.NoError(t, err)
				err = projectMemberStore.Create(ctx, &v1.ProjectMember{Meta: &v1.Meta{Annotations: map[string]string{"role": "editor"}}, ProjectId: "1", TenantId: "google"})
				require.NoError(t, err)
			},
			want:    &v1.ListTenantMembersResponse{},
			wantErr: nil,
		},
		{
			name: "indirect membership in multiple projects",
			req: &v1.ListTenantMembersRequest{
				TenantId:         "github",
				IncludeInherited: pointer.Pointer(true),
			},
			prepare: func() {
				err := tenantStore.Create(ctx, &v1.Tenant{Meta: &v1.Meta{Id: "github"}})
				require.NoError(t, err)
				err = tenantStore.Create(ctx, &v1.Tenant{Meta: &v1.Meta{Id: "azure"}})
				require.NoError(t, err)
				err = projectStore.Create(ctx, &v1.Project{Meta: &v1.Meta{Id: "1"}, TenantId: "github"})
				require.NoError(t, err)
				err = projectStore.Create(ctx, &v1.Project{Meta: &v1.Meta{Id: "2"}, TenantId: "github"})
				require.NoError(t, err)
				err = projectMemberStore.Create(ctx, &v1.ProjectMember{Meta: &v1.Meta{Annotations: map[string]string{"project-role": "owner"}}, ProjectId: "1", TenantId: "github"})
				require.NoError(t, err)
				err = projectMemberStore.Create(ctx, &v1.ProjectMember{Meta: &v1.Meta{Annotations: map[string]string{"project-role": "owner"}}, ProjectId: "2", TenantId: "github"})
				require.NoError(t, err)
				err = projectMemberStore.Create(ctx, &v1.ProjectMember{Meta: &v1.Meta{Annotations: map[string]string{"project-role": "viewer"}}, ProjectId: "2", TenantId: "azure"})
				require.NoError(t, err)
				err = tenantMemberStore.Create(ctx, &v1.TenantMember{Meta: &v1.Meta{Annotations: map[string]string{"tenant-role": "owner"}}, MemberId: "github", TenantId: "github"})
				require.NoError(t, err)
			},
			want: &v1.ListTenantMembersResponse{
				Tenants: []*v1.TenantWithMembershipAnnotations{
					{
						Tenant: &v1.Tenant{
							Meta: &v1.Meta{
								Kind:       "Tenant",
								Apiversion: "v1",
								Id:         "github",
							},
						},
						TenantAnnotations: map[string]string{"tenant-role": "owner"},
						ProjectIds: []string{
							"1",
							"2",
						},
					},
					{
						Tenant: &v1.Tenant{
							Meta: &v1.Meta{
								Kind:       "Tenant",
								Apiversion: "v1",
								Id:         "azure",
							},
						},
						ProjectIds: []string{
							"2",
						},
					},
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

			got, err := s.ListTenantMembers(ctx, tt.req)
			if diff := cmp.Diff(err, tt.wantErr); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
				return
			}
			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreTypes(protoimpl.MessageState{}), cmpopts.IgnoreFields(v1.Meta{}, "CreatedTime"), testcommon.IgnoreUnexported()); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
			}
		})
	}
}
