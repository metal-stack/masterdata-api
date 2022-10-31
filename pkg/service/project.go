package service

import (
	"context"
	"fmt"

	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/masterdata-api/pkg/datastore"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProjectService struct {
	Storage datastore.Storage
	log     *zap.Logger
}

func NewProjectService(s datastore.Storage, l *zap.Logger) *ProjectService {
	return &ProjectService{
		Storage: NewStorageStatusWrapper(s),
		log:     l,
	}
}

func (s *ProjectService) Create(ctx context.Context, req *v1.ProjectCreateRequest) (*v1.ProjectResponse, error) {
	project := req.Project

	tenant := &v1.Tenant{}
	err := s.Storage.Get(ctx, project.GetTenantId(), tenant)
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
		filter := make(map[string]interface{})
		filter["project ->> 'tenant_id'"] = project.GetTenantId()
		var projects []v1.Project
		_, err = s.Storage.Find(ctx, filter, nil, &projects)
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
	err = s.Storage.Create(ctx, project)
	return project.NewProjectResponse(), err
}
func (s *ProjectService) Update(ctx context.Context, req *v1.ProjectUpdateRequest) (*v1.ProjectResponse, error) {
	project := req.Project
	err := s.Storage.Update(ctx, project)
	return project.NewProjectResponse(), err
}
func (s *ProjectService) Delete(ctx context.Context, req *v1.ProjectDeleteRequest) (*v1.ProjectResponse, error) {
	project := req.NewProject()
	err := s.Storage.Delete(ctx, project)
	return project.NewProjectResponse(), err
}
func (s *ProjectService) Get(ctx context.Context, req *v1.ProjectGetRequest) (*v1.ProjectResponse, error) {
	project := &v1.Project{}
	err := s.Storage.Get(ctx, req.Id, project)
	if err != nil {
		return nil, err
	}
	return project.NewProjectResponse(), nil
}
func (s *ProjectService) GetHistory(ctx context.Context, req *v1.ProjectGetHistoryRequest) (*v1.ProjectResponse, error) {
	project := &v1.Project{}
	at := req.At.AsTime()
	err := s.Storage.GetHistory(ctx, req.Id, at, project)
	if err != nil {
		return nil, err
	}
	return project.NewProjectResponse(), nil
}
func (s *ProjectService) Find(ctx context.Context, req *v1.ProjectFindRequest) (*v1.ProjectListResponse, error) {
	var res []v1.Project
	filter := make(map[string]interface{})
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
	nextPage, err := s.Storage.Find(ctx, filter, req.Paging, &res)
	if err != nil {
		return nil, err
	}
	resp := new(v1.ProjectListResponse)
	for i := range res {
		p := &res[i]
		resp.Projects = append(resp.Projects, p)
	}
	resp.NextPage = nextPage
	return resp, nil
}
