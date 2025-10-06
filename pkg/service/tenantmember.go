package service

import (
	"context"
	"fmt"
	"log/slog"

	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/masterdata-api/pkg/datastore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type tenantMemberService struct {
	tenantMemberStore datastore.Storage[*v1.TenantMember]
	tenantStore       datastore.Storage[*v1.Tenant]
	log               *slog.Logger
}

func NewTenantMemberService(l *slog.Logger, tds TenantDataStore, tmds TenantMemberDataStore) *tenantMemberService {
	return &tenantMemberService{
		tenantMemberStore: NewStorageStatusWrapper(tmds),
		tenantStore:       NewStorageStatusWrapper(tds),
		log:               l,
	}
}

func (s *tenantMemberService) Create(ctx context.Context, rq *v1.TenantMemberCreateRequest) (*v1.TenantMemberResponse, error) {
	tenantMember := rq.TenantMember

	_, err := s.tenantStore.Get(ctx, tenantMember.GetTenantId())
	if err != nil && v1.IsNotFound(err) {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("unable to find tenant:%s for tenantMember", tenantMember.GetTenantId()))
	}
	if err != nil {
		return nil, err
	}

	// allow create without sending Meta
	if tenantMember.Meta == nil {
		tenantMember.Meta = &v1.Meta{}
	}

	err = s.tenantMemberStore.Create(ctx, tenantMember)

	return tenantMember.NewTenantMemberResponse(), err
}

func (s *tenantMemberService) Update(ctx context.Context, rq *v1.TenantMemberUpdateRequest) (*v1.TenantMemberResponse, error) {
	tenantMember := rq.TenantMember

	old, err := s.tenantMemberStore.Get(ctx, tenantMember.Meta.Id)
	if err != nil {
		return nil, err
	}

	if old.TenantId != tenantMember.TenantId {
		return nil, status.Error(codes.InvalidArgument, "updating the tenant id of a tenant member is not allowed")
	}
	if old.MemberId != tenantMember.MemberId {
		return nil, status.Error(codes.InvalidArgument, "updating the member id of a tenant member is not allowed")
	}
	if old.Namespace != tenantMember.Namespace {
		return nil, status.Error(codes.InvalidArgument, "updating the namespace of a tenant member is not allowed")
	}

	err = s.tenantMemberStore.Update(ctx, tenantMember)

	return tenantMember.NewTenantMemberResponse(), err
}

func (s *tenantMemberService) Delete(ctx context.Context, rq *v1.TenantMemberDeleteRequest) (*v1.TenantMemberResponse, error) {
	tenantMember := rq.NewTenantMember()

	err := s.tenantMemberStore.Delete(ctx, tenantMember.Meta.Id)

	return tenantMember.NewTenantMemberResponse(), err
}

func (s *tenantMemberService) Get(ctx context.Context, rq *v1.TenantMemberGetRequest) (*v1.TenantMemberResponse, error) {
	tenantMember, err := s.tenantMemberStore.Get(ctx, rq.Id)
	if err != nil {
		return nil, err
	}

	return tenantMember.NewTenantMemberResponse(), nil
}

func (s *tenantMemberService) Find(ctx context.Context, rq *v1.TenantMemberFindRequest) (*v1.TenantMemberListResponse, error) {
	filter := map[string]any{
		"COALESCE(tenantmember ->> 'namespace', '')": rq.Namespace,
	}

	if rq.TenantId != nil {
		filter["tenantmember ->> 'tenant_id'"] = rq.TenantId
	}
	if rq.MemberId != nil {
		filter["tenantmember ->> 'member_id'"] = rq.MemberId
	}
	for key, value := range rq.Annotations {
		// select * from tenantMember where tenantMember -> 'meta' -> 'annotations' ->>  'metal-stack.io/role' = 'owner';
		f := fmt.Sprintf("tenantmember -> 'meta' -> 'annotations' ->> '%s'", key)
		filter[f] = value
	}

	res, _, err := s.tenantMemberStore.Find(ctx, nil, filter)
	if err != nil {
		return nil, err
	}

	resp := new(v1.TenantMemberListResponse)
	resp.TenantMembers = append(resp.TenantMembers, res...)

	return resp, nil
}
