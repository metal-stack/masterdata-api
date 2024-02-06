package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	v1 "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/masterdata-api/pkg/datastore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type projectMemberService struct {
	projectMemberStore datastore.Storage[*v1.ProjectMember]
	tenantStore        datastore.Storage[*v1.Tenant]
	projectStore       datastore.Storage[*v1.Project]
	log                *slog.Logger
}

func NewProjectMemberService(db *sqlx.DB, l *slog.Logger) (*projectMemberService, error) {
	pms, err := datastore.New(l, db, &v1.ProjectMember{})
	if err != nil {
		return nil, err
	}
	ts, err := datastore.New(l, db, &v1.Tenant{})
	if err != nil {
		return nil, err
	}
	ps, err := datastore.New(l, db, &v1.Project{})
	if err != nil {
		return nil, err
	}
	return &projectMemberService{
		projectMemberStore: NewStorageStatusWrapper(pms),
		tenantStore:        NewStorageStatusWrapper(ts),
		projectStore:       NewStorageStatusWrapper(ps),
		log:                l,
	}, nil
}

func (s *projectMemberService) Create(ctx context.Context, req *v1.ProjectMemberCreateRequest) (*v1.ProjectMemberResponse, error) {
	projectMember := req.ProjectMember

	_, err := s.tenantStore.Get(ctx, projectMember.GetTenantId())
	if err != nil && v1.IsNotFound(err) {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("unable to find tenant:%s for projectMember", projectMember.GetTenantId()))
	}
	if err != nil {
		return nil, err
	}

	_, err = s.projectStore.Get(ctx, projectMember.GetProjectId())
	if err != nil && v1.IsNotFound(err) {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("unable to find project:%s for projectMember", projectMember.GetProjectId()))
	}
	if err != nil {
		return nil, err
	}

	// allow create without sending Meta
	if projectMember.Meta == nil {
		projectMember.Meta = &v1.Meta{}
	}
	err = s.projectMemberStore.Create(ctx, projectMember)
	return projectMember.NewProjectMemberResponse(), err
}
func (s *projectMemberService) Update(ctx context.Context, req *v1.ProjectMemberUpdateRequest) (*v1.ProjectMemberResponse, error) {
	projectMember := req.ProjectMember
	err := s.projectMemberStore.Update(ctx, projectMember)
	return projectMember.NewProjectMemberResponse(), err
}
func (s *projectMemberService) Delete(ctx context.Context, req *v1.ProjectMemberDeleteRequest) (*v1.ProjectMemberResponse, error) {
	projectMember := req.NewProjectMember()
	err := s.projectMemberStore.Delete(ctx, projectMember.Meta.Id)
	return projectMember.NewProjectMemberResponse(), err
}
func (s *projectMemberService) Get(ctx context.Context, req *v1.ProjectMemberGetRequest) (*v1.ProjectMemberResponse, error) {
	projectMember, err := s.projectMemberStore.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return projectMember.NewProjectMemberResponse(), nil
}
func (s *projectMemberService) Find(ctx context.Context, req *v1.ProjectMemberFindRequest) (*v1.ProjectMemberListResponse, error) {
	filter := make(map[string]any)
	if req.ProjectId != nil {
		filter["projectmember ->> 'project_id'"] = req.ProjectId
	}
	if req.TenantId != nil {
		filter["projectmember ->> 'tenant_id'"] = req.TenantId
	}
	for key, value := range req.Annotations {
		// select * from projectMember where projectMember -> 'meta' -> 'annotations' ->>  'metal-stack.io/role' = 'owner';
		f := fmt.Sprintf("projectmember -> 'meta' -> 'annotations' ->> '%s'", key)
		filter[f] = value
	}
	res, _, err := s.projectMemberStore.Find(ctx, filter, nil)
	if err != nil {
		return nil, err
	}
	resp := new(v1.ProjectMemberListResponse)
	resp.ProjectMembers = append(resp.ProjectMembers, res...)
	return resp, nil
}
