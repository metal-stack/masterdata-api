package service

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/masterdata-api/pkg/datastore"
	"github.com/metal-stack/masterdata-api/pkg/datastore/mocks"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/metal-stack/metal-lib/pkg/testcommon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/runtime/protoimpl"
	"google.golang.org/protobuf/types/known/wrapperspb"
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
	memberStorageMock.On("Find", ctx, tfilter, paging).Return([]*v1.TenantMember{}, nil, nil)
	memberStorageMock.On("Find", ctx, mfilter, paging).Return([]*v1.TenantMember{}, nil, nil)
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

func Test_tenantService_ProjectsFromMemberships(t *testing.T) {
	ctx := context.Background()
	ves := []datastore.Entity{
		&v1.Project{},
		&v1.ProjectMember{},
		&v1.Tenant{},
		&v1.TenantMember{},
	}

	container, db, err := StartPostgres(ves...)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, container.Stop(ctx, pointer.Pointer(3*time.Second)))
	}()
	defer func() {
		require.NoError(t, db.Close())
	}()

	s := &tenantService{
		db:  db,
		log: slog.Default(),
	}

	var (
		projectStore, _       = datastore.New(log, db, &v1.Project{})
		tenantMemberStore, _  = datastore.New(log, db, &v1.TenantMember{})
		projectMemberStore, _ = datastore.New(log, db, &v1.ProjectMember{})
	)

	tests := []struct {
		name    string
		prepare func()
		req     *v1.ProjectsFromMembershipsRequest
		want    *v1.ProjectsFromMembershipsResponse
		wantErr error
	}{
		{
			name: "direct membership",
			req: &v1.ProjectsFromMembershipsRequest{
				TenantId: "a",
			},
			prepare: func() {
				err := projectStore.Create(ctx, &v1.Project{Meta: &v1.Meta{Id: "1"}})
				require.NoError(t, err)
				err = projectMemberStore.Create(ctx, &v1.ProjectMember{Meta: &v1.Meta{Annotations: map[string]string{"role": "admin"}}, ProjectId: "1", TenantId: "a"})
				require.NoError(t, err)
			},
			want: &v1.ProjectsFromMembershipsResponse{
				Projects: []*v1.ProjectMembershipWithAnnotations{{
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
			req: &v1.ProjectsFromMembershipsRequest{
				TenantId: "a",
			},
			prepare: func() {
				err := projectStore.Create(ctx, &v1.Project{Meta: &v1.Meta{Id: "1"}, TenantId: "b"})
				require.NoError(t, err)
				err = tenantMemberStore.Create(ctx, &v1.TenantMember{Meta: &v1.Meta{Annotations: map[string]string{"role": "viewer"}}, TenantId: "b", MemberId: "a"})
				require.NoError(t, err)
			},
			want: &v1.ProjectsFromMembershipsResponse{
				Projects: []*v1.ProjectMembershipWithAnnotations{{
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

			got, err := s.ProjectsFromMemberships(ctx, tt.req)
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

func Test_tenantService_TenantsFromMemberships(t *testing.T) {
	ctx := context.Background()
	ves := []datastore.Entity{
		&v1.Project{},
		&v1.ProjectMember{},
		&v1.Tenant{},
		&v1.TenantMember{},
	}

	container, db, err := StartPostgres(ves...)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, container.Stop(ctx, pointer.Pointer(3*time.Second)))
	}()
	defer func() {
		require.NoError(t, db.Close())
	}()

	s := &tenantService{
		db:  db,
		log: slog.Default(),
	}

	var (
		tenantStore, _       = datastore.New(log, db, &v1.Tenant{})
		tenantMemberStore, _ = datastore.New(log, db, &v1.TenantMember{})
	)

	tests := []struct {
		name    string
		prepare func()
		req     *v1.TenantsFromMembershipsRequest
		want    *v1.TenantsFromMembershipsResponse
		wantErr error
	}{
		{
			name: "direct membership",
			req: &v1.TenantsFromMembershipsRequest{
				TenantId: "a",
			},
			prepare: func() {
				err := tenantStore.Create(ctx, &v1.Tenant{Meta: &v1.Meta{Id: "b"}})
				require.NoError(t, err)
				err = tenantMemberStore.Create(ctx, &v1.TenantMember{Meta: &v1.Meta{Annotations: map[string]string{"role": "admin"}}, MemberId: "a", TenantId: "b"})
				require.NoError(t, err)
			},
			want: &v1.TenantsFromMembershipsResponse{
				Tenants: []*v1.Tenant{
					{
						Meta: &v1.Meta{
							Kind:       "Tenant",
							Apiversion: "v1",
							Id:         "b",
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

			got, err := s.TenantsFromMemberships(ctx, tt.req)
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
