// Code generated by protoc-gen-go. DO NOT EDIT.
// source: v1/tenant.proto

package v1

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Tenant struct {
	Meta                 *Meta      `protobuf:"bytes,1,opt,name=meta,proto3" json:"meta,omitempty"`
	Name                 string     `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description          string     `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	DefaultQuotas        *QuotaSet  `protobuf:"bytes,4,opt,name=default_quotas,json=defaultQuotas,proto3" json:"default_quotas,omitempty"`
	Quotas               *QuotaSet  `protobuf:"bytes,5,opt,name=quotas,proto3" json:"quotas,omitempty"`
	IamConfig            *IAMConfig `protobuf:"bytes,6,opt,name=iam_config,json=iamConfig,proto3" json:"iam_config,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Tenant) Reset()         { *m = Tenant{} }
func (m *Tenant) String() string { return proto.CompactTextString(m) }
func (*Tenant) ProtoMessage()    {}
func (*Tenant) Descriptor() ([]byte, []int) {
	return fileDescriptor_941e7e5149062005, []int{0}
}

func (m *Tenant) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Tenant.Unmarshal(m, b)
}
func (m *Tenant) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Tenant.Marshal(b, m, deterministic)
}
func (m *Tenant) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Tenant.Merge(m, src)
}
func (m *Tenant) XXX_Size() int {
	return xxx_messageInfo_Tenant.Size(m)
}
func (m *Tenant) XXX_DiscardUnknown() {
	xxx_messageInfo_Tenant.DiscardUnknown(m)
}

var xxx_messageInfo_Tenant proto.InternalMessageInfo

func (m *Tenant) GetMeta() *Meta {
	if m != nil {
		return m.Meta
	}
	return nil
}

func (m *Tenant) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Tenant) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Tenant) GetDefaultQuotas() *QuotaSet {
	if m != nil {
		return m.DefaultQuotas
	}
	return nil
}

func (m *Tenant) GetQuotas() *QuotaSet {
	if m != nil {
		return m.Quotas
	}
	return nil
}

func (m *Tenant) GetIamConfig() *IAMConfig {
	if m != nil {
		return m.IamConfig
	}
	return nil
}

type TenantCreateRequest struct {
	Tenant               *Tenant  `protobuf:"bytes,1,opt,name=tenant,proto3" json:"tenant,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TenantCreateRequest) Reset()         { *m = TenantCreateRequest{} }
func (m *TenantCreateRequest) String() string { return proto.CompactTextString(m) }
func (*TenantCreateRequest) ProtoMessage()    {}
func (*TenantCreateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_941e7e5149062005, []int{1}
}

func (m *TenantCreateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TenantCreateRequest.Unmarshal(m, b)
}
func (m *TenantCreateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TenantCreateRequest.Marshal(b, m, deterministic)
}
func (m *TenantCreateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TenantCreateRequest.Merge(m, src)
}
func (m *TenantCreateRequest) XXX_Size() int {
	return xxx_messageInfo_TenantCreateRequest.Size(m)
}
func (m *TenantCreateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_TenantCreateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TenantCreateRequest proto.InternalMessageInfo

func (m *TenantCreateRequest) GetTenant() *Tenant {
	if m != nil {
		return m.Tenant
	}
	return nil
}

type TenantUpdateRequest struct {
	Tenant               *Tenant  `protobuf:"bytes,1,opt,name=tenant,proto3" json:"tenant,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TenantUpdateRequest) Reset()         { *m = TenantUpdateRequest{} }
func (m *TenantUpdateRequest) String() string { return proto.CompactTextString(m) }
func (*TenantUpdateRequest) ProtoMessage()    {}
func (*TenantUpdateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_941e7e5149062005, []int{2}
}

func (m *TenantUpdateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TenantUpdateRequest.Unmarshal(m, b)
}
func (m *TenantUpdateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TenantUpdateRequest.Marshal(b, m, deterministic)
}
func (m *TenantUpdateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TenantUpdateRequest.Merge(m, src)
}
func (m *TenantUpdateRequest) XXX_Size() int {
	return xxx_messageInfo_TenantUpdateRequest.Size(m)
}
func (m *TenantUpdateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_TenantUpdateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TenantUpdateRequest proto.InternalMessageInfo

func (m *TenantUpdateRequest) GetTenant() *Tenant {
	if m != nil {
		return m.Tenant
	}
	return nil
}

type TenantDeleteRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TenantDeleteRequest) Reset()         { *m = TenantDeleteRequest{} }
func (m *TenantDeleteRequest) String() string { return proto.CompactTextString(m) }
func (*TenantDeleteRequest) ProtoMessage()    {}
func (*TenantDeleteRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_941e7e5149062005, []int{3}
}

func (m *TenantDeleteRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TenantDeleteRequest.Unmarshal(m, b)
}
func (m *TenantDeleteRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TenantDeleteRequest.Marshal(b, m, deterministic)
}
func (m *TenantDeleteRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TenantDeleteRequest.Merge(m, src)
}
func (m *TenantDeleteRequest) XXX_Size() int {
	return xxx_messageInfo_TenantDeleteRequest.Size(m)
}
func (m *TenantDeleteRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_TenantDeleteRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TenantDeleteRequest proto.InternalMessageInfo

func (m *TenantDeleteRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type TenantGetRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TenantGetRequest) Reset()         { *m = TenantGetRequest{} }
func (m *TenantGetRequest) String() string { return proto.CompactTextString(m) }
func (*TenantGetRequest) ProtoMessage()    {}
func (*TenantGetRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_941e7e5149062005, []int{4}
}

func (m *TenantGetRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TenantGetRequest.Unmarshal(m, b)
}
func (m *TenantGetRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TenantGetRequest.Marshal(b, m, deterministic)
}
func (m *TenantGetRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TenantGetRequest.Merge(m, src)
}
func (m *TenantGetRequest) XXX_Size() int {
	return xxx_messageInfo_TenantGetRequest.Size(m)
}
func (m *TenantGetRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_TenantGetRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TenantGetRequest proto.InternalMessageInfo

func (m *TenantGetRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type TenantFindRequest struct {
	Id                   *wrappers.StringValue `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                 *wrappers.StringValue `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *TenantFindRequest) Reset()         { *m = TenantFindRequest{} }
func (m *TenantFindRequest) String() string { return proto.CompactTextString(m) }
func (*TenantFindRequest) ProtoMessage()    {}
func (*TenantFindRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_941e7e5149062005, []int{5}
}

func (m *TenantFindRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TenantFindRequest.Unmarshal(m, b)
}
func (m *TenantFindRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TenantFindRequest.Marshal(b, m, deterministic)
}
func (m *TenantFindRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TenantFindRequest.Merge(m, src)
}
func (m *TenantFindRequest) XXX_Size() int {
	return xxx_messageInfo_TenantFindRequest.Size(m)
}
func (m *TenantFindRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_TenantFindRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TenantFindRequest proto.InternalMessageInfo

func (m *TenantFindRequest) GetId() *wrappers.StringValue {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *TenantFindRequest) GetName() *wrappers.StringValue {
	if m != nil {
		return m.Name
	}
	return nil
}

type TenantResponse struct {
	Tenant               *Tenant  `protobuf:"bytes,1,opt,name=tenant,proto3" json:"tenant,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TenantResponse) Reset()         { *m = TenantResponse{} }
func (m *TenantResponse) String() string { return proto.CompactTextString(m) }
func (*TenantResponse) ProtoMessage()    {}
func (*TenantResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_941e7e5149062005, []int{6}
}

func (m *TenantResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TenantResponse.Unmarshal(m, b)
}
func (m *TenantResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TenantResponse.Marshal(b, m, deterministic)
}
func (m *TenantResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TenantResponse.Merge(m, src)
}
func (m *TenantResponse) XXX_Size() int {
	return xxx_messageInfo_TenantResponse.Size(m)
}
func (m *TenantResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_TenantResponse.DiscardUnknown(m)
}

var xxx_messageInfo_TenantResponse proto.InternalMessageInfo

func (m *TenantResponse) GetTenant() *Tenant {
	if m != nil {
		return m.Tenant
	}
	return nil
}

type TenantListResponse struct {
	Tenants              []*Tenant `protobuf:"bytes,1,rep,name=tenants,proto3" json:"tenants,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *TenantListResponse) Reset()         { *m = TenantListResponse{} }
func (m *TenantListResponse) String() string { return proto.CompactTextString(m) }
func (*TenantListResponse) ProtoMessage()    {}
func (*TenantListResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_941e7e5149062005, []int{7}
}

func (m *TenantListResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TenantListResponse.Unmarshal(m, b)
}
func (m *TenantListResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TenantListResponse.Marshal(b, m, deterministic)
}
func (m *TenantListResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TenantListResponse.Merge(m, src)
}
func (m *TenantListResponse) XXX_Size() int {
	return xxx_messageInfo_TenantListResponse.Size(m)
}
func (m *TenantListResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_TenantListResponse.DiscardUnknown(m)
}

var xxx_messageInfo_TenantListResponse proto.InternalMessageInfo

func (m *TenantListResponse) GetTenants() []*Tenant {
	if m != nil {
		return m.Tenants
	}
	return nil
}

func init() {
	proto.RegisterType((*Tenant)(nil), "v1.Tenant")
	proto.RegisterType((*TenantCreateRequest)(nil), "v1.TenantCreateRequest")
	proto.RegisterType((*TenantUpdateRequest)(nil), "v1.TenantUpdateRequest")
	proto.RegisterType((*TenantDeleteRequest)(nil), "v1.TenantDeleteRequest")
	proto.RegisterType((*TenantGetRequest)(nil), "v1.TenantGetRequest")
	proto.RegisterType((*TenantFindRequest)(nil), "v1.TenantFindRequest")
	proto.RegisterType((*TenantResponse)(nil), "v1.TenantResponse")
	proto.RegisterType((*TenantListResponse)(nil), "v1.TenantListResponse")
}

func init() {
	proto.RegisterFile("v1/tenant.proto", fileDescriptor_941e7e5149062005)
}

var fileDescriptor_941e7e5149062005 = []byte{
	// 462 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x91, 0xdd, 0x6e, 0xd3, 0x30,
	0x14, 0xc7, 0x95, 0xb4, 0x04, 0x7a, 0xba, 0x16, 0x38, 0x7c, 0x45, 0xd5, 0x84, 0xaa, 0x68, 0x48,
	0xbb, 0x98, 0x12, 0xd2, 0xb1, 0x0b, 0xb8, 0x43, 0x43, 0x4c, 0x48, 0xec, 0x82, 0x14, 0xb8, 0x9d,
	0xbc, 0xe6, 0xb4, 0xb2, 0xd4, 0xc4, 0x59, 0xec, 0x84, 0x97, 0xe1, 0xf5, 0x78, 0x0f, 0x14, 0x3b,
	0x99, 0x53, 0xb4, 0x0a, 0xb8, 0x8b, 0xff, 0x1f, 0xc7, 0x39, 0xfe, 0xc1, 0xc3, 0x3a, 0x8e, 0x14,
	0xe5, 0x2c, 0x57, 0x61, 0x51, 0x0a, 0x25, 0xd0, 0xad, 0xe3, 0xd9, 0xa4, 0x8e, 0xa3, 0x8c, 0x14,
	0x33, 0xd2, 0x6c, 0x5a, 0xc7, 0xd1, 0x4d, 0x25, 0x6e, 0xcf, 0x07, 0x75, 0x1c, 0x71, 0x96, 0xb5,
	0xa7, 0x97, 0x1b, 0x21, 0x36, 0x5b, 0x8a, 0xf4, 0xe9, 0xba, 0x5a, 0x47, 0x3f, 0x4a, 0x56, 0x14,
	0x54, 0x4a, 0xe3, 0x07, 0xbf, 0x1c, 0xf0, 0xbe, 0xea, 0x1b, 0xf0, 0x10, 0x86, 0xcd, 0x58, 0xdf,
	0x99, 0x3b, 0xc7, 0xe3, 0xc5, 0x83, 0xb0, 0x8e, 0xc3, 0x4b, 0x52, 0x2c, 0xd1, 0x2a, 0x22, 0x0c,
	0x73, 0x96, 0x91, 0xef, 0xce, 0x9d, 0xe3, 0x51, 0xa2, 0xbf, 0x71, 0x0e, 0xe3, 0x94, 0xe4, 0xaa,
	0xe4, 0x85, 0xe2, 0x22, 0xf7, 0x07, 0xda, 0xea, 0x4b, 0x78, 0x0a, 0xd3, 0x94, 0xd6, 0xac, 0xda,
	0xaa, 0x2b, 0xfd, 0x8f, 0xd2, 0x1f, 0xea, 0xe9, 0x07, 0xcd, 0xf4, 0x2f, 0x8d, 0xb2, 0x24, 0x95,
	0x4c, 0xda, 0x8c, 0x16, 0x24, 0x1e, 0x81, 0xd7, 0x86, 0xef, 0xdd, 0x11, 0x6e, 0x3d, 0x3c, 0x01,
	0xe0, 0x2c, 0xbb, 0x5a, 0x89, 0x7c, 0xcd, 0x37, 0xbe, 0xa7, 0x93, 0x93, 0x26, 0xf9, 0xe9, 0xfd,
	0xe5, 0xb9, 0x16, 0x93, 0x11, 0x67, 0x99, 0xf9, 0x0c, 0xde, 0xc2, 0x13, 0xb3, 0xe6, 0x79, 0x49,
	0x4c, 0x51, 0x42, 0x37, 0x15, 0x49, 0x85, 0x01, 0x78, 0xe6, 0x7d, 0xdb, 0xad, 0xa1, 0x19, 0x60,
	0x82, 0x49, 0xeb, 0xd8, 0xea, 0xb7, 0x22, 0xfd, 0xcf, 0xea, 0xab, 0xae, 0xfa, 0x81, 0xb6, 0x64,
	0xab, 0x53, 0x70, 0x79, 0xaa, 0x6b, 0xa3, 0xc4, 0xe5, 0x69, 0x10, 0xc0, 0x23, 0x13, 0xbb, 0x20,
	0xb5, 0x2f, 0x23, 0xe1, 0xb1, 0xc9, 0x7c, 0xe4, 0x79, 0xda, 0x85, 0x4e, 0x6e, 0x43, 0xe3, 0xc5,
	0x61, 0x68, 0x50, 0x87, 0x1d, 0xea, 0x70, 0xa9, 0x4a, 0x9e, 0x6f, 0xbe, 0xb3, 0x6d, 0x45, 0xcd,
	0x08, 0x7c, 0xdd, 0x43, 0xf8, 0xb7, 0xbc, 0x4e, 0x06, 0x6f, 0x60, 0xda, 0x6e, 0x44, 0xb2, 0x10,
	0xb9, 0xa4, 0x7f, 0xda, 0xfa, 0x1d, 0xa0, 0x51, 0x3e, 0x73, 0x69, 0x9b, 0x47, 0x70, 0xdf, 0xf8,
	0xd2, 0x77, 0xe6, 0x83, 0x3f, 0xaa, 0x9d, 0xb5, 0xf8, 0xe9, 0xc2, 0xc4, 0x68, 0x4b, 0x2a, 0x6b,
	0xbe, 0x22, 0x3c, 0x03, 0xcf, 0x30, 0xc3, 0x17, 0xb6, 0xb0, 0x43, 0x71, 0x86, 0xbd, 0x49, 0xdd,
	0x75, 0x67, 0xe0, 0x19, 0x5e, 0xfd, 0xda, 0x0e, 0xc1, 0x7d, 0x35, 0xc3, 0xaa, 0x5f, 0xdb, 0xa1,
	0x77, 0x67, 0x2d, 0x82, 0xc1, 0x05, 0x29, 0x7c, 0x6a, 0x2d, 0x8b, 0x72, 0xcf, 0x3d, 0xc3, 0x06,
	0x24, 0x3e, 0xb3, 0x5e, 0x0f, 0xec, 0xec, 0xb9, 0x95, 0xfb, 0x8f, 0x78, 0xed, 0x69, 0x58, 0xa7,
	0xbf, 0x03, 0x00, 0x00, 0xff, 0xff, 0x0c, 0x90, 0xd8, 0xb4, 0x19, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// TenantServiceClient is the client API for TenantService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TenantServiceClient interface {
	Create(ctx context.Context, in *TenantCreateRequest, opts ...grpc.CallOption) (*TenantResponse, error)
	Update(ctx context.Context, in *TenantUpdateRequest, opts ...grpc.CallOption) (*TenantResponse, error)
	Delete(ctx context.Context, in *TenantDeleteRequest, opts ...grpc.CallOption) (*TenantResponse, error)
	Get(ctx context.Context, in *TenantGetRequest, opts ...grpc.CallOption) (*TenantResponse, error)
	Find(ctx context.Context, in *TenantFindRequest, opts ...grpc.CallOption) (*TenantListResponse, error)
}

type tenantServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTenantServiceClient(cc grpc.ClientConnInterface) TenantServiceClient {
	return &tenantServiceClient{cc}
}

func (c *tenantServiceClient) Create(ctx context.Context, in *TenantCreateRequest, opts ...grpc.CallOption) (*TenantResponse, error) {
	out := new(TenantResponse)
	err := c.cc.Invoke(ctx, "/v1.TenantService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tenantServiceClient) Update(ctx context.Context, in *TenantUpdateRequest, opts ...grpc.CallOption) (*TenantResponse, error) {
	out := new(TenantResponse)
	err := c.cc.Invoke(ctx, "/v1.TenantService/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tenantServiceClient) Delete(ctx context.Context, in *TenantDeleteRequest, opts ...grpc.CallOption) (*TenantResponse, error) {
	out := new(TenantResponse)
	err := c.cc.Invoke(ctx, "/v1.TenantService/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tenantServiceClient) Get(ctx context.Context, in *TenantGetRequest, opts ...grpc.CallOption) (*TenantResponse, error) {
	out := new(TenantResponse)
	err := c.cc.Invoke(ctx, "/v1.TenantService/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tenantServiceClient) Find(ctx context.Context, in *TenantFindRequest, opts ...grpc.CallOption) (*TenantListResponse, error) {
	out := new(TenantListResponse)
	err := c.cc.Invoke(ctx, "/v1.TenantService/Find", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TenantServiceServer is the server API for TenantService service.
type TenantServiceServer interface {
	Create(context.Context, *TenantCreateRequest) (*TenantResponse, error)
	Update(context.Context, *TenantUpdateRequest) (*TenantResponse, error)
	Delete(context.Context, *TenantDeleteRequest) (*TenantResponse, error)
	Get(context.Context, *TenantGetRequest) (*TenantResponse, error)
	Find(context.Context, *TenantFindRequest) (*TenantListResponse, error)
}

// UnimplementedTenantServiceServer can be embedded to have forward compatible implementations.
type UnimplementedTenantServiceServer struct {
}

func (*UnimplementedTenantServiceServer) Create(ctx context.Context, req *TenantCreateRequest) (*TenantResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (*UnimplementedTenantServiceServer) Update(ctx context.Context, req *TenantUpdateRequest) (*TenantResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (*UnimplementedTenantServiceServer) Delete(ctx context.Context, req *TenantDeleteRequest) (*TenantResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (*UnimplementedTenantServiceServer) Get(ctx context.Context, req *TenantGetRequest) (*TenantResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (*UnimplementedTenantServiceServer) Find(ctx context.Context, req *TenantFindRequest) (*TenantListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Find not implemented")
}

func RegisterTenantServiceServer(s *grpc.Server, srv TenantServiceServer) {
	s.RegisterService(&_TenantService_serviceDesc, srv)
}

func _TenantService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TenantCreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TenantServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.TenantService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TenantServiceServer).Create(ctx, req.(*TenantCreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TenantService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TenantUpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TenantServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.TenantService/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TenantServiceServer).Update(ctx, req.(*TenantUpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TenantService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TenantDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TenantServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.TenantService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TenantServiceServer).Delete(ctx, req.(*TenantDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TenantService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TenantGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TenantServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.TenantService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TenantServiceServer).Get(ctx, req.(*TenantGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TenantService_Find_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TenantFindRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TenantServiceServer).Find(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.TenantService/Find",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TenantServiceServer).Find(ctx, req.(*TenantFindRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _TenantService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "v1.TenantService",
	HandlerType: (*TenantServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _TenantService_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _TenantService_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _TenantService_Delete_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _TenantService_Get_Handler,
		},
		{
			MethodName: "Find",
			Handler:    _TenantService_Find_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/tenant.proto",
}
