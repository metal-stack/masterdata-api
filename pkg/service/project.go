package service

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/masterdata-api/pkg/datastore"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type projectService struct {
	projectStore datastore.Storage[*v1.Project]
	tenantStore  datastore.Storage[*v1.Tenant]
	log          *zap.Logger
}

func NewProjectService(db *sqlx.DB, l *zap.Logger) (*projectService, error) {
	ps, err := datastore.New(l, db, &v1.Project{})
	if err != nil {
		return nil, err
	}
	ts, err := datastore.New(l, db, &v1.Tenant{})
	if err != nil {
		return nil, err
	}
	return &projectService{
		projectStore: NewStorageStatusWrapper(ps),
		tenantStore:  NewStorageStatusWrapper(ts),
		log:          l,
	}, nil
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
	if tenant.GetQuotas() != nil && tenant.GetQuotas().GetProject() != nil && tenant.GetQuotas().GetProject().GetQuota() != nil {
		maxProjects := tenant.GetQuotas().GetProject().GetQuota().GetValue()
		filter := make(map[string]any)
		filter["project ->> 'tenant_id'"] = project.GetTenantId()
		projects, _, err := s.projectStore.Find(ctx, filter, nil)
		if err != nil {
			return nil, err
		}
		if int32(len(projects)) >= maxProjects {
			return nil, status.Error(
				codes.FailedPrecondition,
				fmt.Sprintf("unable to create project, project quota:%d for tenant:%s reached.", maxProjects, project.GetTenantId()))
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
	project := req.Project
	err := s.projectStore.Update(ctx, project)
	return project.NewProjectResponse(), err
}
func (s *projectService) Delete(ctx context.Context, req *v1.ProjectDeleteRequest) (*v1.ProjectResponse, error) {
	project := req.NewProject()
	err := s.projectStore.Delete(ctx, project.Meta.Id)
	return project.NewProjectResponse(), err
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
	filter := make(map[string]any)
	if req.Id != nil {
		filter["id"] = req.Id.GetValue()
	}
	if req.Name != nil {
		filter["project ->> 'name'"] = req.Name.GetValue()
	}
	if req.Description != nil {
		filter["project ->> 'description'"] = req.Description.GetValue()
	}
	if req.TenantId != nil {
		filter["project ->> 'tenant_id'"] = req.TenantId.GetValue()
	}
	for key, value := range req.Annotations {
		// select * from project where project -> 'meta' -> 'annotations' ->>  'metal-stack.io/admitted' = 'true';
		f := fmt.Sprintf("project -> 'meta' -> 'annotations' ->> '%s'", key)
		filter[f] = value
	}
	res, nextPage, err := s.projectStore.Find(ctx, filter, req.Paging)
	if err != nil {
		return nil, err
	}
	resp := new(v1.ProjectListResponse)
	for i := range res {
		p := res[i]
		resp.Projects = append(resp.Projects, p)
	}
	resp.NextPage = nextPage
	return resp, nil
}
