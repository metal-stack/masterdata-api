// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: v1/project_member.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	ProjectMemberService_Create_FullMethodName = "/v1.ProjectMemberService/Create"
	ProjectMemberService_Update_FullMethodName = "/v1.ProjectMemberService/Update"
	ProjectMemberService_Delete_FullMethodName = "/v1.ProjectMemberService/Delete"
	ProjectMemberService_Get_FullMethodName    = "/v1.ProjectMemberService/Get"
	ProjectMemberService_Find_FullMethodName   = "/v1.ProjectMemberService/Find"
)

// ProjectMemberServiceClient is the client API for ProjectMemberService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProjectMemberServiceClient interface {
	Create(ctx context.Context, in *ProjectMemberCreateRequest, opts ...grpc.CallOption) (*ProjectMemberResponse, error)
	Update(ctx context.Context, in *ProjectMemberUpdateRequest, opts ...grpc.CallOption) (*ProjectMemberResponse, error)
	Delete(ctx context.Context, in *ProjectMemberDeleteRequest, opts ...grpc.CallOption) (*ProjectMemberResponse, error)
	Get(ctx context.Context, in *ProjectMemberGetRequest, opts ...grpc.CallOption) (*ProjectMemberResponse, error)
	Find(ctx context.Context, in *ProjectMemberFindRequest, opts ...grpc.CallOption) (*ProjectMemberListResponse, error)
}

type projectMemberServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewProjectMemberServiceClient(cc grpc.ClientConnInterface) ProjectMemberServiceClient {
	return &projectMemberServiceClient{cc}
}

func (c *projectMemberServiceClient) Create(ctx context.Context, in *ProjectMemberCreateRequest, opts ...grpc.CallOption) (*ProjectMemberResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ProjectMemberResponse)
	err := c.cc.Invoke(ctx, ProjectMemberService_Create_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectMemberServiceClient) Update(ctx context.Context, in *ProjectMemberUpdateRequest, opts ...grpc.CallOption) (*ProjectMemberResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ProjectMemberResponse)
	err := c.cc.Invoke(ctx, ProjectMemberService_Update_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectMemberServiceClient) Delete(ctx context.Context, in *ProjectMemberDeleteRequest, opts ...grpc.CallOption) (*ProjectMemberResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ProjectMemberResponse)
	err := c.cc.Invoke(ctx, ProjectMemberService_Delete_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectMemberServiceClient) Get(ctx context.Context, in *ProjectMemberGetRequest, opts ...grpc.CallOption) (*ProjectMemberResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ProjectMemberResponse)
	err := c.cc.Invoke(ctx, ProjectMemberService_Get_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectMemberServiceClient) Find(ctx context.Context, in *ProjectMemberFindRequest, opts ...grpc.CallOption) (*ProjectMemberListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ProjectMemberListResponse)
	err := c.cc.Invoke(ctx, ProjectMemberService_Find_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProjectMemberServiceServer is the server API for ProjectMemberService service.
// All implementations should embed UnimplementedProjectMemberServiceServer
// for forward compatibility.
type ProjectMemberServiceServer interface {
	Create(context.Context, *ProjectMemberCreateRequest) (*ProjectMemberResponse, error)
	Update(context.Context, *ProjectMemberUpdateRequest) (*ProjectMemberResponse, error)
	Delete(context.Context, *ProjectMemberDeleteRequest) (*ProjectMemberResponse, error)
	Get(context.Context, *ProjectMemberGetRequest) (*ProjectMemberResponse, error)
	Find(context.Context, *ProjectMemberFindRequest) (*ProjectMemberListResponse, error)
}

// UnimplementedProjectMemberServiceServer should be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedProjectMemberServiceServer struct{}

func (UnimplementedProjectMemberServiceServer) Create(context.Context, *ProjectMemberCreateRequest) (*ProjectMemberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedProjectMemberServiceServer) Update(context.Context, *ProjectMemberUpdateRequest) (*ProjectMemberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedProjectMemberServiceServer) Delete(context.Context, *ProjectMemberDeleteRequest) (*ProjectMemberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedProjectMemberServiceServer) Get(context.Context, *ProjectMemberGetRequest) (*ProjectMemberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedProjectMemberServiceServer) Find(context.Context, *ProjectMemberFindRequest) (*ProjectMemberListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Find not implemented")
}
func (UnimplementedProjectMemberServiceServer) testEmbeddedByValue() {}

// UnsafeProjectMemberServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProjectMemberServiceServer will
// result in compilation errors.
type UnsafeProjectMemberServiceServer interface {
	mustEmbedUnimplementedProjectMemberServiceServer()
}

func RegisterProjectMemberServiceServer(s grpc.ServiceRegistrar, srv ProjectMemberServiceServer) {
	// If the following call pancis, it indicates UnimplementedProjectMemberServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ProjectMemberService_ServiceDesc, srv)
}

func _ProjectMemberService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProjectMemberCreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectMemberServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectMemberService_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectMemberServiceServer).Create(ctx, req.(*ProjectMemberCreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectMemberService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProjectMemberUpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectMemberServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectMemberService_Update_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectMemberServiceServer).Update(ctx, req.(*ProjectMemberUpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectMemberService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProjectMemberDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectMemberServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectMemberService_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectMemberServiceServer).Delete(ctx, req.(*ProjectMemberDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectMemberService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProjectMemberGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectMemberServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectMemberService_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectMemberServiceServer).Get(ctx, req.(*ProjectMemberGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectMemberService_Find_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProjectMemberFindRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectMemberServiceServer).Find(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectMemberService_Find_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectMemberServiceServer).Find(ctx, req.(*ProjectMemberFindRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ProjectMemberService_ServiceDesc is the grpc.ServiceDesc for ProjectMemberService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProjectMemberService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.ProjectMemberService",
	HandlerType: (*ProjectMemberServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _ProjectMemberService_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _ProjectMemberService_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _ProjectMemberService_Delete_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _ProjectMemberService_Get_Handler,
		},
		{
			MethodName: "Find",
			Handler:    _ProjectMemberService_Find_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/project_member.proto",
}
