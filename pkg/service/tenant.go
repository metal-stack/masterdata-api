package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/davecgh/go-spew/spew"
	"github.com/jmoiron/sqlx"
	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/masterdata-api/pkg/datastore"
)

type tenantService struct {
	db                *sqlx.DB
	tenantStore       datastore.Storage[*v1.Tenant]
	tenantMemberStore datastore.Storage[*v1.TenantMember]
	log               *slog.Logger
}

func NewTenantService(db *sqlx.DB, l *slog.Logger) (*tenantService, error) {
	ts, err := datastore.New(l, db, &v1.Tenant{})
	if err != nil {
		return nil, err
	}
	tms, err := datastore.New(l, db, &v1.TenantMember{})
	if err != nil {
		return nil, err
	}
	return &tenantService{
		db:                db,
		tenantStore:       NewStorageStatusWrapper(ts),
		tenantMemberStore: NewStorageStatusWrapper(tms),
		log:               l,
	}, nil
}

func (s *tenantService) Create(ctx context.Context, req *v1.TenantCreateRequest) (*v1.TenantResponse, error) {
	tenant := req.Tenant
	// allow create without sending Meta
	if tenant.Meta == nil {
		tenant.Meta = &v1.Meta{}
	}
	err := s.tenantStore.Create(ctx, tenant)
	return tenant.NewTenantResponse(), err
}
func (s *tenantService) Update(ctx context.Context, req *v1.TenantUpdateRequest) (*v1.TenantResponse, error) {
	tenant := req.Tenant
	err := s.tenantStore.Update(ctx, tenant)
	return tenant.NewTenantResponse(), err
}

func (s *tenantService) Delete(ctx context.Context, req *v1.TenantDeleteRequest) (*v1.TenantResponse, error) {
	tenant := req.NewTenant()
	tenantFilter := map[string]any{
		"tenantmember ->> 'tenant_id'": tenant.Meta.Id,
	}
	memberFilter := map[string]any{
		"tenantmember ->> 'member_id'": tenant.Meta.Id,
	}
	tenantMemberships, _, err := s.tenantMemberStore.Find(ctx, tenantFilter, nil)
	if err != nil {
		return nil, err
	}
	memberMemberships, _, err := s.tenantMemberStore.Find(ctx, memberFilter, nil)
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, m := range tenantMemberships {
		ids = append(ids, m.Meta.Id)
	}
	for _, m := range memberMemberships {
		ids = append(ids, m.Meta.Id)
	}

	if len(ids) > 0 {
		err = s.tenantMemberStore.DeleteAll(ctx, ids...)
		if err != nil {
			return nil, err
		}
	}
	err = s.tenantStore.Delete(ctx, tenant.Meta.Id)
	if err != nil {
		return nil, err
	}
	return tenant.NewTenantResponse(), nil
}

func (s *tenantService) Get(ctx context.Context, req *v1.TenantGetRequest) (*v1.TenantResponse, error) {
	tenant, err := s.tenantStore.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// response with entity, no error
	return tenant.NewTenantResponse(), nil
}

func (s *tenantService) GetHistory(ctx context.Context, req *v1.TenantGetHistoryRequest) (*v1.TenantResponse, error) {
	tenant := &v1.Tenant{}
	at := req.At.AsTime()
	s.log.Info("getHistory", "id", req.Id, "at", at)
	err := s.tenantStore.GetHistory(ctx, req.Id, at, tenant)
	if err != nil {
		return nil, err
	}

	// response with entity, no error
	return tenant.NewTenantResponse(), nil
}

func (s *tenantService) Find(ctx context.Context, req *v1.TenantFindRequest) (*v1.TenantListResponse, error) {
	filter := make(map[string]any)
	if req.Id != nil {
		filter["id"] = req.GetId().GetValue()
	}
	if req.Name != nil {
		filter["tenant ->> 'name'"] = req.GetName().GetValue()
	}
	for key, value := range req.Annotations {
		// select * from tenants where tenant -> 'meta' -> 'annotations' ->>  'metal-stack.io/admitted' = 'true';
		f := fmt.Sprintf("tenant -> 'meta' -> 'annotations' ->> '%s'", key)
		filter[f] = value
	}
	res, nextPage, err := s.tenantStore.Find(ctx, filter, req.Paging)
	if err != nil {
		return nil, err
	}
	resp := new(v1.TenantListResponse)
	resp.Tenants = append(resp.Tenants, res...)
	resp.NextPage = nextPage
	return resp, nil
}

func (s *tenantService) ProjectsFromMemberships(ctx context.Context, req *v1.ProjectsFromMembershipsRequest) (*v1.ProjectsFromMembershipsResponse, error) {
	pm := datastore.Entity(&v1.ProjectMember{})
	p := datastore.Entity(&v1.Project{})

	dm := sq.
		Select(
			p.JSONField(),
			pm.JSONField()+"->'meta'->>'annotations' AS annotations",
		).
		From(pm.TableName()).
		Join(p.TableName() + " ON " + p.TableName() + ".id = " + pm.JSONField() + "->>'project_id'").
		Where(pm.JSONField() + "->>'tenant_id' = :tenantId")

	query, _, err := dm.ToSql()
	if err != nil {
		return nil, err
	}

	spew.Dump(query)

	rows, err := s.db.NamedQueryContext(ctx, query, map[string]any{"tenantId": req.TenantId})
	if err != nil {
		return nil, err
	}

	var projects []*v1.ProjectMembershipWithAnnotations
	for rows.Next() {
		type annotations map[string]string
		type res struct {
			Project     *v1.Project
			Annotations string
		}
		var r res
		err = rows.Scan(&r.Project, &r.Annotations)
		if err != nil {
			return nil, err
		}

		var unmarshalled map[string]string
		err = json.Unmarshal([]byte(r.Annotations), &unmarshalled)
		if err != nil {
			return nil, err
		}

		projects = append(projects, &v1.ProjectMembershipWithAnnotations{
			Project:            r.Project,
			ProjectAnnotations: unmarshalled,
			TenantAnnotations:  make(map[string]string),
		})
	}

	spew.Dump(projects)

	return &v1.ProjectsFromMembershipsResponse{}, nil
}

func (s *tenantService) TenantsFromMemberships(context.Context, *v1.TenantsFromMembershipsRequest) (*v1.TenantsFromMembershipsResponse, error) {
	panic("unimplemented")
}
