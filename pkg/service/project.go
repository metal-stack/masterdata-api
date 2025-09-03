package service

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/masterdata-api/pkg/datastore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type projectService struct {
	projectStore       datastore.Storage[*v1.Project]
	projectMemberStore datastore.Storage[*v1.ProjectMember]
	tenantStore        datastore.Storage[*v1.Tenant]
	log                *slog.Logger
}

func NewProjectService(l *slog.Logger, pds ProjectDataStore, pmds ProjectMemberDataStore, tds TenantDataStore) *projectService {
	return &projectService{
		projectStore:       NewStorageStatusWrapper(pds),
		projectMemberStore: NewStorageStatusWrapper(pmds),
		tenantStore:        NewStorageStatusWrapper(tds),
		log:                l,
	}
}

func (s *projectService) Create(ctx context.Context, req *v1.ProjectCreateRequest) (*v1.ProjectResponse, error) {
	project := req.Project

	tenant, err := s.tenantStore.Get(ctx, project.GetTenantId())
	if err != nil && v1.IsNotFound(err) {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("unable to find tenant:%s for project", project.GetTenantId()))
	}
	if err != nil {
		return nil, err
	}
	if tenant.GetDefaultQuotas() != nil && project.GetQuotas() == nil {
		project.Quotas = tenant.GetDefaultQuotas()
	}

	// Check if tenant defines project quotas
	if tenant.GetQuotas() != nil && tenant.GetQuotas().GetProject() != nil && tenant.GetQuotas().GetProject().Max != nil {
		maxProjects := tenant.GetQuotas().GetProject().Max
		filter := make(map[string]any)
		filter["project ->> 'tenant_id'"] = project.GetTenantId()
		projects, _, err := s.projectStore.Find(ctx, nil, filter)
		if err != nil {
			return nil, err
		}

		if maxProjects != nil {
			if len(projects) >= int(*maxProjects) {
				return nil, status.Error(
					codes.FailedPrecondition,
					fmt.Sprintf("unable to create project, project quota:%d for tenant:%s reached.", maxProjects, project.GetTenantId()))
			}
		}
	}

	// allow create without sending Meta
	if project.Meta == nil {
		project.Meta = &v1.Meta{}
	}
	err = s.projectStore.Create(ctx, project)
	return project.NewProjectResponse(), err
}
func (s *projectService) Update(ctx context.Context, req *v1.ProjectUpdateRequest) (*v1.ProjectResponse, error) {
	old, err := s.projectStore.Get(ctx, req.Project.Meta.Id)
	if err != nil {
		return nil, err
	}
	project := req.Project
	if old.TenantId != project.TenantId {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("update tenant of project:%s is not allowed", project.Meta.Id))
	}
	err = s.projectStore.Update(ctx, project)
	return project.NewProjectResponse(), err
}
func (s *projectService) Delete(ctx context.Context, req *v1.ProjectDeleteRequest) (*v1.ProjectResponse, error) {
	project := req.NewProject()
	filter := map[string]any{
		"projectmember ->> 'project_id'": project.Meta.Id,
	}
	memberships, _, err := s.projectMemberStore.Find(ctx, nil, filter)
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, m := range memberships {
		ids = append(ids, m.Meta.Id)
	}

	if len(ids) > 0 {
		err = s.projectMemberStore.DeleteAll(ctx, ids...)
		if err != nil {
			return nil, err
		}
	}
	err = s.projectStore.Delete(ctx, project.Meta.Id)
	if err != nil {
		return nil, err
	}
	return project.NewProjectResponse(), nil
}
func (s *projectService) Get(ctx context.Context, req *v1.ProjectGetRequest) (*v1.ProjectResponse, error) {
	project, err := s.projectStore.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return project.NewProjectResponse(), nil
}
func (s *projectService) GetHistory(ctx context.Context, req *v1.ProjectGetHistoryRequest) (*v1.ProjectResponse, error) {
	project := &v1.Project{}
	at := req.At.AsTime()
	err := s.projectStore.GetHistory(ctx, req.Id, at, project)
	if err != nil {
		return nil, err
	}
	return project.NewProjectResponse(), nil
}
func (s *projectService) Find(ctx context.Context, req *v1.ProjectFindRequest) (*v1.ProjectListResponse, error) {

	var filters []any

	mapFilter := make(map[string]any)
	if req.Id != nil {
		mapFilter["id"] = req.Id
	}
	if req.Name != nil {
		mapFilter["project ->> 'name'"] = req.Name
	}
	if req.Description != nil {
		mapFilter["project ->> 'description'"] = req.Description
	}
	if req.TenantId != nil {
		mapFilter["project ->> 'tenant_id'"] = req.TenantId
	}
	for key, value := range req.Annotations {
		// select * from project where project -> 'meta' -> 'annotations' ->>  'metal-stack.io/admitted' = 'true';
		f := fmt.Sprintf("project -> 'meta' -> 'annotations' ->> '%s'", key)
		mapFilter[f] = value
	}

	if len(mapFilter) > 0 {
		filters = append(filters, mapFilter)
	}

	if len(req.Labels) > 0 {
		var contains []string

		for _, label := range req.Labels {
			contains = append(contains, strconv.Quote(label))
		}

		// select * from projects where project -> 'meta' -> 'labels' @> '["a=b","c=d"]';
		labelFilter := fmt.Sprintf(`project -> 'meta' -> 'labels' @> '[%s]'`, strings.Join(contains, ","))

		filters = append(filters, labelFilter)
	}

	res, nextPage, err := s.projectStore.Find(ctx, req.Paging, filters...)
	if err != nil {
		return nil, err
	}
	resp := new(v1.ProjectListResponse)
	resp.Projects = append(resp.Projects, res...)
	resp.NextPage = nextPage
	return resp, nil
}
