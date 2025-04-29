package service

import (
	"context"
	"fmt"
	"log/slog"

	"connectrpc.com/connect"
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

func NewProjectMemberService(l *slog.Logger, pds ProjectDataStore, pmds ProjectMemberDataStore, tds TenantDataStore) *projectMemberService {
	return &projectMemberService{
		projectMemberStore: NewStorageStatusWrapper(pmds),
		tenantStore:        NewStorageStatusWrapper(tds),
		projectStore:       NewStorageStatusWrapper(pds),
		log:                l,
	}
}

func (s *projectMemberService) Create(ctx context.Context, rq *connect.Request[v1.ProjectMemberCreateRequest]) (*connect.Response[v1.ProjectMemberResponse], error) {
	req := rq.Msg
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
	return connect.NewResponse(projectMember.NewProjectMemberResponse()), err
}
func (s *projectMemberService) Update(ctx context.Context, rq *connect.Request[v1.ProjectMemberUpdateRequest]) (*connect.Response[v1.ProjectMemberResponse], error) {
	req := rq.Msg
	projectMember := req.ProjectMember
	err := s.projectMemberStore.Update(ctx, projectMember)
	return connect.NewResponse(projectMember.NewProjectMemberResponse()), err
}
func (s *projectMemberService) Delete(ctx context.Context, rq *connect.Request[v1.ProjectMemberDeleteRequest]) (*connect.Response[v1.ProjectMemberResponse], error) {
	req := rq.Msg
	projectMember := req.NewProjectMember()
	err := s.projectMemberStore.Delete(ctx, projectMember.Meta.Id)
	return connect.NewResponse(projectMember.NewProjectMemberResponse()), err
}
func (s *projectMemberService) Get(ctx context.Context, rq *connect.Request[v1.ProjectMemberGetRequest]) (*connect.Response[v1.ProjectMemberResponse], error) {
	req := rq.Msg
	projectMember, err := s.projectMemberStore.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(projectMember.NewProjectMemberResponse()), nil
}
func (s *projectMemberService) Find(ctx context.Context, rq *connect.Request[v1.ProjectMemberFindRequest]) (*connect.Response[v1.ProjectMemberListResponse], error) {
	req := rq.Msg
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
	return connect.NewResponse(resp), nil
}
