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

func (s *tenantMemberService) Create(ctx context.Context, req *v1.TenantMemberCreateRequest) (*v1.TenantMemberResponse, error) {
	tenantMember := req.TenantMember

	_, err := s.tenantStore.Get(ctx, tenantMember.GetTenantId())
	if err != nil && v1.IsNotFound(err) {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("unable to find tenant:%s for tenantMember", tenantMember.GetTenantId()))
	}
	if err != nil {
		return nil, err
	}

	_, err = s.tenantStore.Get(ctx, tenantMember.GetMemberId())
	if err != nil && v1.IsNotFound(err) {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("unable to find member:%s for tenantMember", tenantMember.GetMemberId()))
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
func (s *tenantMemberService) Update(ctx context.Context, req *v1.TenantMemberUpdateRequest) (*v1.TenantMemberResponse, error) {
	tenantMember := req.TenantMember
	err := s.tenantMemberStore.Update(ctx, tenantMember)
	return tenantMember.NewTenantMemberResponse(), err
}
func (s *tenantMemberService) Delete(ctx context.Context, req *v1.TenantMemberDeleteRequest) (*v1.TenantMemberResponse, error) {
	tenantMember := req.NewTenantMember()
	err := s.tenantMemberStore.Delete(ctx, tenantMember.Meta.Id)
	return tenantMember.NewTenantMemberResponse(), err
}
func (s *tenantMemberService) Get(ctx context.Context, req *v1.TenantMemberGetRequest) (*v1.TenantMemberResponse, error) {
	tenantMember, err := s.tenantMemberStore.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return tenantMember.NewTenantMemberResponse(), nil
}
func (s *tenantMemberService) Find(ctx context.Context, req *v1.TenantMemberFindRequest) (*v1.TenantMemberListResponse, error) {
	filter := make(map[string]any)
	if req.TenantId != nil {
		filter["tenantmember ->> 'tenant_id'"] = req.TenantId
	}
	if req.MemberId != nil {
		filter["tenantmember ->> 'member_id'"] = req.MemberId
	}
	for key, value := range req.Annotations {
		// select * from tenantMember where tenantMember -> 'meta' -> 'annotations' ->>  'metal-stack.io/role' = 'owner';
		f := fmt.Sprintf("tenantmember -> 'meta' -> 'annotations' ->> '%s'", key)
		filter[f] = value
	}
	res, _, err := s.tenantMemberStore.Find(ctx, filter, nil)
	if err != nil {
		return nil, err
	}
	resp := new(v1.TenantMemberListResponse)
	resp.TenantMembers = append(resp.TenantMembers, res...)
	return resp, nil
}
