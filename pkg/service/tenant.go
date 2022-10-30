package service

import (
	"context"

	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/masterdata-api/pkg/datastore"
	"go.uber.org/zap"
)

type TenantService struct {
	Storage datastore.Storage
	log     *zap.Logger
}

func NewTenantService(s datastore.Storage, l *zap.Logger) *TenantService {
	return &TenantService{
		Storage: NewStorageStatusWrapper(s),
		log:     l,
	}
}

func (s *TenantService) Create(ctx context.Context, req *v1.TenantCreateRequest) (*v1.TenantResponse, error) {
	tenant := req.Tenant
	// allow create without sending Meta
	if tenant.Meta == nil {
		tenant.Meta = &v1.Meta{}
	}
	err := s.Storage.Create(ctx, tenant)
	return tenant.NewTenantResponse(), err
}
func (s *TenantService) Update(ctx context.Context, req *v1.TenantUpdateRequest) (*v1.TenantResponse, error) {
	tenant := req.Tenant
	err := s.Storage.Update(ctx, tenant)
	return tenant.NewTenantResponse(), err
}

func (s *TenantService) Delete(ctx context.Context, req *v1.TenantDeleteRequest) (*v1.TenantResponse, error) {
	tenant := req.NewTenant()
	err := s.Storage.Delete(ctx, tenant)
	return tenant.NewTenantResponse(), err
}

func (s *TenantService) Get(ctx context.Context, req *v1.TenantGetRequest) (*v1.TenantResponse, error) {
	tenant := &v1.Tenant{}
	err := s.Storage.Get(ctx, req.Id, tenant)
	if err != nil {
		return nil, err
	}

	// response with entity, no error
	return tenant.NewTenantResponse(), nil
}

func (s *TenantService) GetHistory(ctx context.Context, req *v1.TenantGetHistoryRequest) (*v1.TenantResponse, error) {
	tenant := &v1.Tenant{}
	at := req.At.AsTime()
	s.log.Info("getHistory", zap.String("id", req.Id), zap.Time("at", at))
	err := s.Storage.GetHistory(ctx, req.Id, at, tenant)
	if err != nil {
		return nil, err
	}

	// response with entity, no error
	return tenant.NewTenantResponse(), nil
}

func (s *TenantService) Find(ctx context.Context, req *v1.TenantFindRequest) (*v1.TenantListResponse, error) {
	var res []v1.Tenant
	filter := make(map[string]interface{})
	if req.Id != nil {
		filter["id"] = req.GetId().GetValue()
	}
	if req.Name != nil {
		filter["tenant ->> 'name'"] = req.GetName().GetValue()
	}
	nextPage, err := s.Storage.Find(ctx, filter, req.Paging, &res)
	if err != nil {
		return nil, err
	}
	resp := new(v1.TenantListResponse)
	for i := range res {
		t := &res[i]
		resp.Tenants = append(resp.Tenants, t)
	}
	resp.NextPage = nextPage
	return resp, nil
}
