// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: v1/tenant_member.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	TenantMemberService_Create_FullMethodName = "/v1.TenantMemberService/Create"
	TenantMemberService_Update_FullMethodName = "/v1.TenantMemberService/Update"
	TenantMemberService_Delete_FullMethodName = "/v1.TenantMemberService/Delete"
	TenantMemberService_Get_FullMethodName    = "/v1.TenantMemberService/Get"
	TenantMemberService_Find_FullMethodName   = "/v1.TenantMemberService/Find"
)

// TenantMemberServiceClient is the client API for TenantMemberService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TenantMemberServiceClient interface {
	Create(ctx context.Context, in *TenantMemberCreateRequest, opts ...grpc.CallOption) (*TenantMemberResponse, error)
	Update(ctx context.Context, in *TenantMemberUpdateRequest, opts ...grpc.CallOption) (*TenantMemberResponse, error)
	Delete(ctx context.Context, in *TenantMemberDeleteRequest, opts ...grpc.CallOption) (*TenantMemberResponse, error)
	Get(ctx context.Context, in *TenantMemberGetRequest, opts ...grpc.CallOption) (*TenantMemberResponse, error)
	Find(ctx context.Context, in *TenantMemberFindRequest, opts ...grpc.CallOption) (*TenantMemberListResponse, error)
}

type tenantMemberServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTenantMemberServiceClient(cc grpc.ClientConnInterface) TenantMemberServiceClient {
	return &tenantMemberServiceClient{cc}
}

func (c *tenantMemberServiceClient) Create(ctx context.Context, in *TenantMemberCreateRequest, opts ...grpc.CallOption) (*TenantMemberResponse, error) {
	out := new(TenantMemberResponse)
	err := c.cc.Invoke(ctx, TenantMemberService_Create_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tenantMemberServiceClient) Update(ctx context.Context, in *TenantMemberUpdateRequest, opts ...grpc.CallOption) (*TenantMemberResponse, error) {
	out := new(TenantMemberResponse)
	err := c.cc.Invoke(ctx, TenantMemberService_Update_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tenantMemberServiceClient) Delete(ctx context.Context, in *TenantMemberDeleteRequest, opts ...grpc.CallOption) (*TenantMemberResponse, error) {
	out := new(TenantMemberResponse)
	err := c.cc.Invoke(ctx, TenantMemberService_Delete_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tenantMemberServiceClient) Get(ctx context.Context, in *TenantMemberGetRequest, opts ...grpc.CallOption) (*TenantMemberResponse, error) {
	out := new(TenantMemberResponse)
	err := c.cc.Invoke(ctx, TenantMemberService_Get_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tenantMemberServiceClient) Find(ctx context.Context, in *TenantMemberFindRequest, opts ...grpc.CallOption) (*TenantMemberListResponse, error) {
	out := new(TenantMemberListResponse)
	err := c.cc.Invoke(ctx, TenantMemberService_Find_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TenantMemberServiceServer is the server API for TenantMemberService service.
// All implementations should embed UnimplementedTenantMemberServiceServer
// for forward compatibility
type TenantMemberServiceServer interface {
	Create(context.Context, *TenantMemberCreateRequest) (*TenantMemberResponse, error)
	Update(context.Context, *TenantMemberUpdateRequest) (*TenantMemberResponse, error)
	Delete(context.Context, *TenantMemberDeleteRequest) (*TenantMemberResponse, error)
	Get(context.Context, *TenantMemberGetRequest) (*TenantMemberResponse, error)
	Find(context.Context, *TenantMemberFindRequest) (*TenantMemberListResponse, error)
}

// UnimplementedTenantMemberServiceServer should be embedded to have forward compatible implementations.
type UnimplementedTenantMemberServiceServer struct {
}

func (UnimplementedTenantMemberServiceServer) Create(context.Context, *TenantMemberCreateRequest) (*TenantMemberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedTenantMemberServiceServer) Update(context.Context, *TenantMemberUpdateRequest) (*TenantMemberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedTenantMemberServiceServer) Delete(context.Context, *TenantMemberDeleteRequest) (*TenantMemberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedTenantMemberServiceServer) Get(context.Context, *TenantMemberGetRequest) (*TenantMemberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedTenantMemberServiceServer) Find(context.Context, *TenantMemberFindRequest) (*TenantMemberListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Find not implemented")
}

// UnsafeTenantMemberServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TenantMemberServiceServer will
// result in compilation errors.
type UnsafeTenantMemberServiceServer interface {
	mustEmbedUnimplementedTenantMemberServiceServer()
}

func RegisterTenantMemberServiceServer(s grpc.ServiceRegistrar, srv TenantMemberServiceServer) {
	s.RegisterService(&TenantMemberService_ServiceDesc, srv)
}

func _TenantMemberService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TenantMemberCreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TenantMemberServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TenantMemberService_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TenantMemberServiceServer).Create(ctx, req.(*TenantMemberCreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TenantMemberService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TenantMemberUpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TenantMemberServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TenantMemberService_Update_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TenantMemberServiceServer).Update(ctx, req.(*TenantMemberUpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TenantMemberService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TenantMemberDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TenantMemberServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TenantMemberService_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TenantMemberServiceServer).Delete(ctx, req.(*TenantMemberDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TenantMemberService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TenantMemberGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TenantMemberServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TenantMemberService_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TenantMemberServiceServer).Get(ctx, req.(*TenantMemberGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TenantMemberService_Find_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TenantMemberFindRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TenantMemberServiceServer).Find(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TenantMemberService_Find_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TenantMemberServiceServer).Find(ctx, req.(*TenantMemberFindRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TenantMemberService_ServiceDesc is the grpc.ServiceDesc for TenantMemberService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TenantMemberService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.TenantMemberService",
	HandlerType: (*TenantMemberServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _TenantMemberService_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _TenantMemberService_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _TenantMemberService_Delete_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _TenantMemberService_Get_Handler,
		},
		{
			MethodName: "Find",
			Handler:    _TenantMemberService_Find_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/tenant_member.proto",
}
