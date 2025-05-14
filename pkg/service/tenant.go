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
	log               *slog.Logger
	tenantStore       datastore.Storage[*v1.Tenant]
	tenantMemberStore datastore.Storage[*v1.TenantMember]
	db                *sqlx.DB
}

var (
	projectMembers = datastore.Entity(&v1.ProjectMember{})
	tenantMembers  = datastore.Entity(&v1.TenantMember{})
	projects       = datastore.Entity(&v1.Project{})
	tenants        = datastore.Entity(&v1.Tenant{})
)

func NewTenantService(db *sqlx.DB, l *slog.Logger, tds TenantDataStore, tmds TenantMemberDataStore) *tenantService {
	return &tenantService{
		db:                db,
		tenantStore:       NewStorageStatusWrapper(tds),
		tenantMemberStore: NewStorageStatusWrapper(tmds),
		log:               l,
	}
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
	hostFilter := map[string]any{
		"tenantmember ->> 'tenant_id'": tenant.Meta.Id,
	}
	memberFilter := map[string]any{
		"tenantmember ->> 'member_id'": tenant.Meta.Id,
	}
	tenantIsHostMemberships, _, err := s.tenantMemberStore.Find(ctx, hostFilter, nil)
	if err != nil {
		return nil, err
	}
	tenantIsMemberMemberships, _, err := s.tenantMemberStore.Find(ctx, memberFilter, nil)
	if err != nil {
		return nil, err
	}

	unionMap := make(map[string]bool)
	for _, m := range tenantIsHostMemberships {
		unionMap[m.Meta.Id] = true
	}
	for _, m := range tenantIsMemberMemberships {
		unionMap[m.Meta.Id] = true
	}

	var ids []string
	for k := range unionMap {
		ids = append(ids, k)
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

var (
	queryDirectProjectParticipations = sq.Select(
		projects.JSONField(),
		projectMembers.JSONField()+"->'meta'->>'annotations' AS project_membership_annotations",
	).
		From(projectMembers.TableName()).
		Join(projects.TableName() + " ON " + projects.TableName() + ".id = " + projectMembers.JSONField() + "->>'project_id'").
		Where(projectMembers.JSONField() + "->>'tenant_id' = :tenantId")

	queryInheritedProjectParticipations = sq.Select(
		projects.JSONField(),
		tenantMembers.JSONField()+"->'meta'->>'annotations' AS tenant_membership_annotations",
	).
		From(tenantMembers.TableName()).
		Join(projects.TableName() + " ON " + projects.JSONField() + "->>'tenant_id' = " + tenantMembers.JSONField() + "->>'tenant_id'").
		Where(tenantMembers.JSONField() + "->>'member_id' = :tenantId")
)

// FindParticipatingProjects returns all projects in which a member participates.
// This includes projects in which the member is explicitly participating through a project membership but may also
// include memberships, which are inherited by the tenant membership.
func (s *tenantService) FindParticipatingProjects(ctx context.Context, req *v1.FindParticipatingProjectsRequest) (*v1.FindParticipatingProjectsResponse, error) {
	type result struct {
		Project                      *v1.Project
		TenantMembershipAnnotations  []byte `db:"tenant_membership_annotations"`
		ProjectMembershipAnnotations []byte `db:"project_membership_annotations"`
	}

	var (
		res       []*v1.ProjectWithMembershipAnnotations
		resultMap = map[string]*v1.ProjectWithMembershipAnnotations{}

		input = map[string]any{"tenantId": req.TenantId}

		resultFn = func(e result) error {
			p, ok := resultMap[e.Project.Meta.Id]
			if !ok {
				p = &v1.ProjectWithMembershipAnnotations{
					Project: e.Project,
				}
			}

			if e.TenantMembershipAnnotations != nil {
				err := json.Unmarshal(e.TenantMembershipAnnotations, &p.TenantAnnotations)
				if err != nil {
					return err
				}
			}

			if e.ProjectMembershipAnnotations != nil {
				err := json.Unmarshal(e.ProjectMembershipAnnotations, &p.ProjectAnnotations)
				if err != nil {
					return err
				}
			}

			resultMap[e.Project.Meta.Id] = p

			return nil
		}
	)

	err := datastore.RunQuery(ctx, s.log, s.db, queryDirectProjectParticipations, input, resultFn)
	if err != nil {
		return nil, err
	}

	includeInherited := true
	if req.IncludeInherited != nil {
		includeInherited = *req.IncludeInherited
	}

	if includeInherited {
		err := datastore.RunQuery(ctx, s.log, s.db, queryInheritedProjectParticipations, input, resultFn)
		if err != nil {
			return nil, err
		}
	}

	for _, p := range resultMap {
		res = append(res, p)
	}

	return &v1.FindParticipatingProjectsResponse{Projects: res}, nil
}

var (
	queryDirectTenantParticipations = sq.Select(
		tenants.JSONField(),
		tenantMembers.JSONField()+"->'meta'->>'annotations' AS tenant_membership_annotations",
	).
		From(tenantMembers.TableName()).
		Join(tenants.TableName() + " ON " + tenants.TableName() + ".id = " + tenantMembers.JSONField() + "->>'tenant_id'").
		Where(tenantMembers.JSONField() + "->>'member_id' = :tenantId")

	queryInheritedTenantParticipations = sq.Select(
		tenants.JSONField(),
		projectMembers.JSONField()+"->'meta'->>'annotations' AS project_membership_annotations",
	).
		From(projectMembers.TableName()).
		Join(projects.TableName() + " ON " + projects.TableName() + ".id = " + projectMembers.JSONField() + "->>'project_id'").
		Join(tenants.TableName() + " ON " + tenants.TableName() + ".id = " + projects.JSONField() + "->>'tenant_id'").
		Where(projectMembers.JSONField() + "->>'tenant_id' = :tenantId")
)

// FindParticipatingTenants returns all tenants in which a member participates.
// This includes tenants in which the member is explicitly participating through a tenant membership but may also
// include memberships, which are inherited by the project memberships (e.g. through project invites).
func (s *tenantService) FindParticipatingTenants(ctx context.Context, req *v1.FindParticipatingTenantsRequest) (*v1.FindParticipatingTenantsResponse, error) {
	type result struct {
		Tenant                       *v1.Tenant
		TenantMembershipAnnotations  []byte `db:"tenant_membership_annotations"`
		ProjectMembershipAnnotations []byte `db:"project_membership_annotations"`
	}

	var (
		input = map[string]any{"tenantId": req.TenantId}

		res       []*v1.TenantWithMembershipAnnotations
		resultMap = map[string]*v1.TenantWithMembershipAnnotations{}

		resultFn = func(e result) error {
			t, ok := resultMap[e.Tenant.Meta.Id]
			if !ok {
				t = &v1.TenantWithMembershipAnnotations{
					Tenant: e.Tenant,
				}
			}

			if e.TenantMembershipAnnotations != nil {
				err := json.Unmarshal(e.TenantMembershipAnnotations, &t.TenantAnnotations)
				if err != nil {
					return err
				}
			}

			if e.ProjectMembershipAnnotations != nil {
				err := json.Unmarshal(e.ProjectMembershipAnnotations, &t.ProjectAnnotations)
				if err != nil {
					return err
				}
			}

			resultMap[e.Tenant.Meta.Id] = t

			return nil
		}
	)

	err := datastore.RunQuery(ctx, s.log, s.db, queryDirectTenantParticipations, input, resultFn)
	if err != nil {
		return nil, err
	}

	includeInherited := true
	if req.IncludeInherited != nil {
		includeInherited = *req.IncludeInherited
	}

	if includeInherited {
		err = datastore.RunQuery(ctx, s.log, s.db, queryInheritedTenantParticipations, input, resultFn)
		if err != nil {
			return nil, err
		}
	}

	for _, t := range resultMap {
		res = append(res, t)
	}

	return &v1.FindParticipatingTenantsResponse{Tenants: res}, nil
}

var (
	queryDirectTenantMembers = sq.Select(
		tenants.JSONField(),
		tenantMembers.JSONField()+"->'meta'->>'annotations' AS tenant_membership_annotations",
	).
		From(tenantMembers.TableName()).
		Join(tenants.TableName() + " ON " + tenants.TableName() + ".id = " + tenantMembers.JSONField() + "->>'member_id'").
		Where(tenantMembers.JSONField() + "->>'tenant_id' = :tenantId")

	queryInheritedTenantMembers = sq.Select(
		tenants.JSONField(),
		projects.JSONField(),
	).
		From(projectMembers.TableName()).
		Join(projects.TableName() + " ON " + projects.TableName() + ".id = " + projectMembers.JSONField() + "->>'project_id'").
		Join(tenants.TableName() + " ON " + tenants.TableName() + ".id = " + projectMembers.JSONField() + "->>'tenant_id'").
		Where(projects.JSONField() + "->>'tenant_id' = :tenantId")
)

// ListTenantMembers returns all members of a tenant.
// This includes members which are explicitly participating through a tenant membership but may also
// include memberships, which are inherited by the project memberships (e.g. through project invites).
func (s *tenantService) ListTenantMembers(ctx context.Context, req *v1.ListTenantMembersRequest) (*v1.ListTenantMembersResponse, error) {
	type result struct {
		Tenant                      *v1.Tenant
		TenantMembershipAnnotations []byte `db:"tenant_membership_annotations"`
		Project                     *v1.Project
	}

	var (
		res       []*v1.TenantWithMembershipAnnotations
		resultMap = map[string]*v1.TenantWithMembershipAnnotations{}

		input = map[string]any{"tenantId": req.TenantId}

		resultFn = func(e result) error {
			t, ok := resultMap[e.Tenant.Meta.Id]
			if !ok {
				t = &v1.TenantWithMembershipAnnotations{
					Tenant: e.Tenant,
				}
			}

			if e.TenantMembershipAnnotations != nil {
				err := json.Unmarshal(e.TenantMembershipAnnotations, &t.TenantAnnotations)
				if err != nil {
					return err
				}
			}

			if e.Project != nil {
				t.ProjectIds = append(t.ProjectIds, e.Project.Meta.Id)
			}

			resultMap[e.Tenant.Meta.Id] = t

			return nil
		}
	)

	err := datastore.RunQuery(ctx, s.log, s.db, queryDirectTenantMembers, input, resultFn)
	if err != nil {
		return nil, err
	}

	includeInherited := true
	if req.IncludeInherited != nil {
		includeInherited = *req.IncludeInherited
	}

	if includeInherited {
		err = datastore.RunQuery(ctx, s.log, s.db, queryInheritedTenantMembers, input, resultFn)
		if err != nil {
			return nil, err
		}
	}

	for _, t := range resultMap {
		res = append(res, t)
	}

	return &v1.ListTenantMembersResponse{Tenants: res}, nil
}
