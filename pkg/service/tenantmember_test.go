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

func TestFindTenantMember(t *testing.T) {
	ctx := t.Context()
	ves := []datastore.Entity{
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
		tenantMemberStore = datastore.New(log, db, &v1.TenantMember{})
		tenantStore       = datastore.New(log, db, &v1.Tenant{})

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
		testTenantMember1 = &v1.TenantMember{
			Meta: &v1.Meta{
				Id:         "1",
				Kind:       "TenantMember",
				Apiversion: "v1",
				Version:    1,
				Annotations: map[string]string{
					"role": "owner",
				},
				Labels: []string{"a", "b"},
			},
			TenantId:  "tenant-1",
			MemberId:  "tenant-1",
			Namespace: "a",
		}
		testTenantMember2 = &v1.TenantMember{
			Meta: &v1.Meta{
				Id:         "2",
				Kind:       "TenantMember",
				Apiversion: "v1",
				Version:    1,
				Annotations: map[string]string{
					"role": "owner",
				},
				Labels: []string{"c", "d"},
			},
			TenantId:  "tenant-2",
			MemberId:  "tenant-2",
			Namespace: "a",
		}
		testTenantMember3 = &v1.TenantMember{
			Meta: &v1.Meta{
				Id:         "3",
				Kind:       "TenantMember",
				Apiversion: "v1",
				Version:    1,
				Annotations: map[string]string{
					"role": "editor",
				},
				Labels: []string{"e", "f"},
			},
			TenantId:  "tenant-1",
			MemberId:  "tenant-2",
			Namespace: "a",
		}
		testTenantMember4 = &v1.TenantMember{
			Meta: &v1.Meta{
				Id:         "4",
				Kind:       "TenantMember",
				Apiversion: "v1",
				Version:    1,
				Annotations: map[string]string{
					"role": "editor",
				},
				Labels: []string{"e", "f"},
			},
			TenantId:  "tenant-1",
			MemberId:  "tenant-2",
			Namespace: "",
		}

		service = &tenantMemberService{
			log:               log,
			tenantMemberStore: tenantMemberStore,
			tenantStore:       tenantStore,
		}
	)

	tests := []struct {
		name    string
		prepare func()
		req     *v1.TenantMemberFindRequest
		want    *v1.TenantMemberListResponse
		wantErr error
	}{
		{
			name: "find by tenant",
			req: &v1.TenantMemberFindRequest{
				TenantId:  pointer.Pointer("tenant-1"),
				Namespace: "a",
			},
			prepare: func() {
				require.NoError(t, tenantStore.Create(ctx, testTenant1))
				require.NoError(t, tenantStore.Create(ctx, testTenant2))
				require.NoError(t, tenantMemberStore.Create(ctx, testTenantMember1))
				require.NoError(t, tenantMemberStore.Create(ctx, testTenantMember2))
				require.NoError(t, tenantMemberStore.Create(ctx, testTenantMember3))
				require.NoError(t, tenantMemberStore.Create(ctx, testTenantMember4))
			},
			want: &v1.TenantMemberListResponse{
				TenantMembers: []*v1.TenantMember{
					testTenantMember1,
					testTenantMember3,
				},
			},
			wantErr: nil,
		},
		{
			name: "find by tenant id (no results) #1",
			req: &v1.TenantMemberFindRequest{
				TenantId:  pointer.Pointer("no-result"),
				Namespace: "a",
			},
			prepare: func() {
				require.NoError(t, tenantStore.Create(ctx, testTenant1))
				require.NoError(t, tenantStore.Create(ctx, testTenant2))
				require.NoError(t, tenantMemberStore.Create(ctx, testTenantMember1))
				require.NoError(t, tenantMemberStore.Create(ctx, testTenantMember2))
				require.NoError(t, tenantMemberStore.Create(ctx, testTenantMember3))
				require.NoError(t, tenantMemberStore.Create(ctx, testTenantMember4))
			},
			want: &v1.TenantMemberListResponse{
				TenantMembers: nil,
			},
			wantErr: nil,
		},
		{
			name: "find by tenant id (no results) #2",
			req: &v1.TenantMemberFindRequest{
				TenantId:  pointer.Pointer("tenant-1"),
				Namespace: "wrong-namespace",
			},
			prepare: func() {
				require.NoError(t, tenantStore.Create(ctx, testTenant1))
				require.NoError(t, tenantStore.Create(ctx, testTenant2))
				require.NoError(t, tenantMemberStore.Create(ctx, testTenantMember1))
				require.NoError(t, tenantMemberStore.Create(ctx, testTenantMember2))
				require.NoError(t, tenantMemberStore.Create(ctx, testTenantMember3))
				require.NoError(t, tenantMemberStore.Create(ctx, testTenantMember4))
			},
			want: &v1.TenantMemberListResponse{
				TenantMembers: nil,
			},
			wantErr: nil,
		},
		{
			name: "find by tenant",
			req: &v1.TenantMemberFindRequest{
				TenantId:  pointer.Pointer("tenant-2"),
				Namespace: "a",
			},
			prepare: func() {
				require.NoError(t, tenantStore.Create(ctx, testTenant1))
				require.NoError(t, tenantStore.Create(ctx, testTenant2))
				require.NoError(t, tenantMemberStore.Create(ctx, testTenantMember1))
				require.NoError(t, tenantMemberStore.Create(ctx, testTenantMember2))
				require.NoError(t, tenantMemberStore.Create(ctx, testTenantMember3))
				require.NoError(t, tenantMemberStore.Create(ctx, testTenantMember4))
			},
			want: &v1.TenantMemberListResponse{
				TenantMembers: []*v1.TenantMember{
					testTenantMember2,
				},
			},
			wantErr: nil,
		},
		{
			name: "find by annotation",
			req: &v1.TenantMemberFindRequest{
				Annotations: map[string]string{"role": "owner"},
				Namespace:   "a",
			},
			prepare: func() {
				require.NoError(t, tenantStore.Create(ctx, testTenant1))
				require.NoError(t, tenantStore.Create(ctx, testTenant2))
				require.NoError(t, tenantMemberStore.Create(ctx, testTenantMember1))
				require.NoError(t, tenantMemberStore.Create(ctx, testTenantMember2))
				require.NoError(t, tenantMemberStore.Create(ctx, testTenantMember3))
				require.NoError(t, tenantMemberStore.Create(ctx, testTenantMember4))
			},
			want: &v1.TenantMemberListResponse{
				TenantMembers: []*v1.TenantMember{
					testTenantMember1,
					testTenantMember2,
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

			slices.SortFunc(got.TenantMembers, func(i, j *v1.TenantMember) int {
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

func TestUpdateTenantMember(t *testing.T) {
	ctx := t.Context()
	ves := []datastore.Entity{
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
		tenantMemberStore = datastore.New(log, db, &v1.TenantMember{})
		tenantStore       = datastore.New(log, db, &v1.Tenant{})

		service = &tenantMemberService{
			log:               log,
			tenantMemberStore: tenantMemberStore,
			tenantStore:       tenantStore,
		}
	)

	tests := []struct {
		name    string
		prepare func()
		req     *v1.TenantMemberUpdateRequest
		want    *v1.TenantMemberResponse
		wantErr error
	}{
		{
			name: "update mutable fields",
			req: &v1.TenantMemberUpdateRequest{
				TenantMember: &v1.TenantMember{
					Meta: &v1.Meta{
						Id:      "1",
						Version: 1,
						Annotations: map[string]string{
							"role": "owner",
						},
						Labels: []string{"a", "b"},
					},
					TenantId:  "tenant-1",
					MemberId:  "tenant-1",
					Namespace: "a",
				},
			},
			prepare: func() {
				require.NoError(t, tenantMemberStore.Create(ctx, &v1.TenantMember{
					Meta: &v1.Meta{
						Id:         "1",
						Kind:       "TenantMember",
						Apiversion: "v1",
						Version:    1,
					},
					TenantId:  "tenant-1",
					MemberId:  "tenant-1",
					Namespace: "a",
				}))
			},
			want: &v1.TenantMemberResponse{
				TenantMember: &v1.TenantMember{
					Meta: &v1.Meta{
						Id:         "1",
						Kind:       "TenantMember",
						Apiversion: "v1",
						Version:    2,
						Annotations: map[string]string{
							"role": "owner",
						},
						Labels: []string{"a", "b"},
					},
					TenantId:  "tenant-1",
					MemberId:  "tenant-1",
					Namespace: "a",
				},
			},
			wantErr: nil,
		},
		{
			name: "unable to update namespace",
			req: &v1.TenantMemberUpdateRequest{
				TenantMember: &v1.TenantMember{
					Meta: &v1.Meta{
						Id:      "1",
						Version: 1,
					},
					TenantId:  "tenant-1",
					MemberId:  "tenant-1",
					Namespace: "b",
				},
			},
			prepare: func() {
				require.NoError(t, tenantMemberStore.Create(ctx, &v1.TenantMember{
					Meta: &v1.Meta{
						Id:         "1",
						Kind:       "TenantMember",
						Apiversion: "v1",
						Version:    1,
					},
					TenantId:  "tenant-1",
					MemberId:  "tenant-1",
					Namespace: "a",
				}))
			},
			want:    nil,
			wantErr: status.Error(codes.InvalidArgument, "updating the namespace of a tenant member is not allowed"),
		},
		{
			name: "unable to update tenant id",
			req: &v1.TenantMemberUpdateRequest{
				TenantMember: &v1.TenantMember{
					Meta: &v1.Meta{
						Id:      "1",
						Version: 1,
					},
					TenantId:  "tenant-2",
					MemberId:  "tenant-1",
					Namespace: "a",
				},
			},
			prepare: func() {
				require.NoError(t, tenantMemberStore.Create(ctx, &v1.TenantMember{
					Meta: &v1.Meta{
						Id:         "1",
						Kind:       "TenantMember",
						Apiversion: "v1",
						Version:    1,
					},
					TenantId:  "tenant-1",
					MemberId:  "tenant-1",
					Namespace: "a",
				}))
			},
			want:    nil,
			wantErr: status.Error(codes.InvalidArgument, "updating the tenant id of a tenant member is not allowed"),
		},
		{
			name: "unable to update member id",
			req: &v1.TenantMemberUpdateRequest{
				TenantMember: &v1.TenantMember{
					Meta: &v1.Meta{
						Id:      "1",
						Version: 1,
					},
					TenantId:  "tenant-1",
					MemberId:  "tenant-2",
					Namespace: "a",
				},
			},
			prepare: func() {
				require.NoError(t, tenantMemberStore.Create(ctx, &v1.TenantMember{
					Meta: &v1.Meta{
						Id:         "1",
						Kind:       "TenantMember",
						Apiversion: "v1",
						Version:    1,
					},
					TenantId:  "tenant-1",
					MemberId:  "tenant-1",
					Namespace: "a",
				}))
			},
			want:    nil,
			wantErr: status.Error(codes.InvalidArgument, "updating the member id of a tenant member is not allowed"),
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
				assert.NotNil(t, got.TenantMember.Meta.UpdatedTime)
			}

			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreTypes(protoimpl.MessageState{}), cmpopts.IgnoreFields(v1.Meta{}, "CreatedTime", "UpdatedTime"), testcommon.IgnoreUnexported()); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
			}
		})
	}
}
