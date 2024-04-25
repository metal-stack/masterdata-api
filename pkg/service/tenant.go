package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
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
	var (
		pm = datastore.Entity(&v1.ProjectMember{})
		tm = datastore.Entity(&v1.TenantMember{})
		p  = datastore.Entity(&v1.Project{})

		res       []*v1.ProjectMembershipWithAnnotations
		resultMap = map[string]*v1.ProjectMembershipWithAnnotations{}

		// all projects with direct memberships
		directProjects = sq.Select(
			p.JSONField(),
			pm.JSONField()+"->'meta'->>'annotations' AS annotations",
		).
			From(pm.TableName()).
			Join(p.TableName() + " ON " + p.TableName() + ".id = " + pm.JSONField() + "->>'project_id'").
			Where(pm.JSONField() + "->>'tenant_id' = :tenantId")

		// all projects with no direct membership because inherited from tenant membership
		inheritedProjects = sq.Select(
			p.JSONField(),
			tm.JSONField()+"->'meta'->>'annotations' AS annotations",
		).
			From(tm.TableName()).
			Join(p.TableName() + " ON " + p.JSONField() + "->>'tenant_id' = " + tm.JSONField() + "->>'tenant_id'").
			Where(tm.JSONField() + "->>'member_id' = :tenantId")

		runQuery = func(builder sq.SelectBuilder, callback func(*v1.ProjectMembershipWithAnnotations, map[string]string)) error {
			query, vals, err := builder.ToSql()
			if err != nil {
				return err
			}

			if s.log.Enabled(ctx, slog.LevelDebug) {
				s.log.Debug("query", "sql", query, "values", vals)
			}

			rows, err := s.db.NamedQueryContext(ctx, query, map[string]any{"tenantId": req.TenantId})
			if err != nil {
				return err
			}
			defer func() {
				err = rows.Close()
				if err != nil {
					s.log.Error("error closing result rows", "error", err)
				}
			}()

			for rows.Next() {
				var (
					project     *v1.Project
					raw         []byte
					annotations map[string]string
				)

				err = rows.Scan(&project, &raw)
				if err != nil {
					return err
				}

				err = json.Unmarshal(raw, &annotations)
				if err != nil {
					return err
				}

				p, ok := resultMap[project.Meta.Id]
				if !ok {
					p = &v1.ProjectMembershipWithAnnotations{
						Project: project,
					}
				}

				callback(p, annotations)

				resultMap[project.Meta.Id] = p
			}

			return nil
		}
	)

	err := runQuery(directProjects, func(pmwa *v1.ProjectMembershipWithAnnotations, annotations map[string]string) {
		pmwa.ProjectAnnotations = annotations
	})
	if err != nil {
		return nil, err
	}

	includeInherited := true
	if req.IncludeInherited != nil {
		includeInherited = *req.IncludeInherited
	}

	if includeInherited {
		err = runQuery(inheritedProjects, func(pmwa *v1.ProjectMembershipWithAnnotations, annotations map[string]string) {
			pmwa.TenantAnnotations = annotations
		})
		if err != nil {
			return nil, err
		}
	}

	for _, p := range resultMap {
		res = append(res, p)
	}

	return &v1.ProjectsFromMembershipsResponse{Projects: res}, nil
}

func (s *tenantService) TenantsFromMemberships(ctx context.Context, req *v1.TenantsFromMembershipsRequest) (*v1.TenantsFromMembershipsResponse, error) {
	var (
		pm = datastore.Entity(&v1.ProjectMember{})
		tm = datastore.Entity(&v1.TenantMember{})
		p  = datastore.Entity(&v1.Project{})
		t  = datastore.Entity(&v1.Tenant{})

		res       []*v1.TenantMembershipWithAnnotations
		resultMap = map[string]*v1.TenantMembershipWithAnnotations{}

		directTenants = sq.Select(
			t.JSONField(),
			tm.JSONField()+"->'meta'->>'annotations' AS annotations",
		).
			From(tm.TableName()).
			Join(t.TableName() + " ON " + t.TableName() + ".id = " + tm.JSONField() + "->>'tenant_id'").
			Where(tm.JSONField() + "->>'member_id' = :tenantId")

		inheritedTenants = sq.Select(
			t.JSONField(),
			pm.JSONField()+"->'meta'->>'annotations' AS annotations",
		).
			From(pm.TableName()).
			Join(p.TableName() + " ON " + p.TableName() + ".id = " + pm.JSONField() + "->>'project_id'").
			Join(t.TableName() + " ON " + t.TableName() + ".id = " + p.JSONField() + "->>'tenant_id'").
			Where(pm.JSONField() + "->>'tenant_id' = :tenantId")

		runQuery = func(builder sq.SelectBuilder, callback func(*v1.TenantMembershipWithAnnotations, map[string]string)) error {
			query, vals, err := builder.ToSql()
			if err != nil {
				return err
			}

			if s.log.Enabled(ctx, slog.LevelDebug) {
				s.log.Debug("query", "sql", query, "values", vals)
			}

			rows, err := s.db.NamedQueryContext(ctx, query, map[string]any{"tenantId": req.TenantId})
			if err != nil {
				return err
			}
			defer func() {
				err = rows.Close()
				if err != nil {
					s.log.Error("error closing result rows", "error", err)
				}
			}()

			for rows.Next() {
				var (
					tenant      *v1.Tenant
					raw         []byte
					annotations map[string]string
				)

				err = rows.Scan(&tenant, &raw)
				if err != nil {
					return err
				}

				err = json.Unmarshal(raw, &annotations)
				if err != nil {
					return err
				}

				t, ok := resultMap[tenant.Meta.Id]
				if !ok {
					t = &v1.TenantMembershipWithAnnotations{
						Tenant: tenant,
					}
				}

				callback(t, annotations)

				resultMap[tenant.Meta.Id] = t
			}

			return nil
		}
	)

	err := runQuery(directTenants, func(tmwa *v1.TenantMembershipWithAnnotations, annotations map[string]string) {
		tmwa.TenantAnnotations = annotations
	})
	if err != nil {
		return nil, err
	}

	includeInherited := true
	if req.IncludeInherited != nil {
		includeInherited = *req.IncludeInherited
	}

	if includeInherited {
		err = runQuery(inheritedTenants, func(tmwa *v1.TenantMembershipWithAnnotations, annotations map[string]string) {
			tmwa.ProjectAnnotations = annotations
		})
		if err != nil {
			return nil, err
		}
	}

	for _, p := range resultMap {
		res = append(res, p)
	}

	return &v1.TenantsFromMembershipsResponse{Tenants: res}, nil
}

func (s *tenantService) GetAllTenants(ctx context.Context, req *v1.GetAllTenantsRequest) (*v1.GetAllTenantsResponse, error) {
	var (
		pm = datastore.Entity(&v1.ProjectMember{})
		tm = datastore.Entity(&v1.TenantMember{})
		p  = datastore.Entity(&v1.Project{})
		t  = datastore.Entity(&v1.Tenant{})

		res       []*v1.TenantMembershipWithAnnotations
		resultMap = map[string]*v1.TenantMembershipWithAnnotations{}

		directTenants = sq.Select(
			t.JSONField(),
			tm.JSONField()+"->'meta'->>'annotations' AS annotations",
		).
			From(tm.TableName()).
			Join(t.TableName() + " ON " + t.TableName() + ".id = " + tm.JSONField() + "->>'member_id'").
			Where(tm.JSONField() + "->>'tenant_id' = :tenantId")

		inheritedTenants = sq.Select(
			t.JSONField(),
			pm.JSONField()+"->'meta'->>'annotations' AS annotations",
		).
			From(pm.TableName()).
			Join(p.TableName() + " ON " + p.TableName() + ".id = " + pm.JSONField() + "->>'project_id'").
			Join(t.TableName() + " ON " + t.TableName() + ".id = " + pm.JSONField() + "->>'tenant_id'").
			Where(p.JSONField() + "->>'tenant_id' = :tenantId")

		runQuery = func(builder sq.SelectBuilder, callback func(*v1.TenantMembershipWithAnnotations, map[string]string)) error {
			query, vals, err := builder.ToSql()
			if err != nil {
				return err
			}

			if s.log.Enabled(ctx, slog.LevelDebug) {
				s.log.Debug("query", "sql", query, "values", vals)
			}

			rows, err := s.db.NamedQueryContext(ctx, query, map[string]any{"tenantId": req.TenantId})
			if err != nil {
				return err
			}
			defer func() {
				err = rows.Close()
				if err != nil {
					s.log.Error("error closing result rows", "error", err)
				}
			}()

			for rows.Next() {
				var (
					tenant      *v1.Tenant
					raw         []byte
					annotations map[string]string
				)

				err = rows.Scan(&tenant, &raw)
				if err != nil {
					return err
				}

				err = json.Unmarshal(raw, &annotations)
				if err != nil {
					return err
				}

				t, ok := resultMap[tenant.Meta.Id]
				if !ok {
					t = &v1.TenantMembershipWithAnnotations{
						Tenant: tenant,
					}
				}

				callback(t, annotations)

				resultMap[tenant.Meta.Id] = t
			}

			return nil
		}
	)

	err := runQuery(directTenants, func(tmwa *v1.TenantMembershipWithAnnotations, annotations map[string]string) {
		tmwa.TenantAnnotations = annotations
	})
	if err != nil {
		return nil, err
	}

	includeInherited := true
	if req.IncludeInherited != nil {
		includeInherited = *req.IncludeInherited
	}

	if includeInherited {
		err = runQuery(inheritedTenants, func(tmwa *v1.TenantMembershipWithAnnotations, annotations map[string]string) {
			tmwa.ProjectAnnotations = annotations
		})
		if err != nil {
			return nil, err
		}
	}

	for _, p := range resultMap {
		res = append(res, p)
	}

	return &v1.GetAllTenantsResponse{Tenants: res}, nil
}
