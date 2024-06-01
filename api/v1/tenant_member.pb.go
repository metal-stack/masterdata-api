// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        (unknown)
// source: v1/tenant_member.proto

package v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// TenantMember is the database model
type TenantMember struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Meta *Meta `protobuf:"bytes,1,opt,name=meta,proto3" json:"meta,omitempty"`
	// TenantId is the id of the parent tenant
	TenantId string `protobuf:"bytes,2,opt,name=tenant_id,json=tenantId,proto3" json:"tenant_id,omitempty"`
	// MemberId is the id of the member tenant
	MemberId string `protobuf:"bytes,3,opt,name=member_id,json=memberId,proto3" json:"member_id,omitempty"`
}

func (x *TenantMember) Reset() {
	*x = TenantMember{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_tenant_member_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TenantMember) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TenantMember) ProtoMessage() {}

func (x *TenantMember) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tenant_member_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TenantMember.ProtoReflect.Descriptor instead.
func (*TenantMember) Descriptor() ([]byte, []int) {
	return file_v1_tenant_member_proto_rawDescGZIP(), []int{0}
}

func (x *TenantMember) GetMeta() *Meta {
	if x != nil {
		return x.Meta
	}
	return nil
}

func (x *TenantMember) GetTenantId() string {
	if x != nil {
		return x.TenantId
	}
	return ""
}

func (x *TenantMember) GetMemberId() string {
	if x != nil {
		return x.MemberId
	}
	return ""
}

type TenantMemberCreateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TenantMember *TenantMember `protobuf:"bytes,1,opt,name=tenant_member,json=tenantMember,proto3" json:"tenant_member,omitempty"`
}

func (x *TenantMemberCreateRequest) Reset() {
	*x = TenantMemberCreateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_tenant_member_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TenantMemberCreateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TenantMemberCreateRequest) ProtoMessage() {}

func (x *TenantMemberCreateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tenant_member_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TenantMemberCreateRequest.ProtoReflect.Descriptor instead.
func (*TenantMemberCreateRequest) Descriptor() ([]byte, []int) {
	return file_v1_tenant_member_proto_rawDescGZIP(), []int{1}
}

func (x *TenantMemberCreateRequest) GetTenantMember() *TenantMember {
	if x != nil {
		return x.TenantMember
	}
	return nil
}

type TenantMemberUpdateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TenantMember *TenantMember `protobuf:"bytes,1,opt,name=tenant_member,json=tenantMember,proto3" json:"tenant_member,omitempty"`
}

func (x *TenantMemberUpdateRequest) Reset() {
	*x = TenantMemberUpdateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_tenant_member_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TenantMemberUpdateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TenantMemberUpdateRequest) ProtoMessage() {}

func (x *TenantMemberUpdateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tenant_member_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TenantMemberUpdateRequest.ProtoReflect.Descriptor instead.
func (*TenantMemberUpdateRequest) Descriptor() ([]byte, []int) {
	return file_v1_tenant_member_proto_rawDescGZIP(), []int{2}
}

func (x *TenantMemberUpdateRequest) GetTenantMember() *TenantMember {
	if x != nil {
		return x.TenantMember
	}
	return nil
}

type TenantMemberDeleteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *TenantMemberDeleteRequest) Reset() {
	*x = TenantMemberDeleteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_tenant_member_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TenantMemberDeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TenantMemberDeleteRequest) ProtoMessage() {}

func (x *TenantMemberDeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tenant_member_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TenantMemberDeleteRequest.ProtoReflect.Descriptor instead.
func (*TenantMemberDeleteRequest) Descriptor() ([]byte, []int) {
	return file_v1_tenant_member_proto_rawDescGZIP(), []int{3}
}

func (x *TenantMemberDeleteRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type TenantMemberGetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *TenantMemberGetRequest) Reset() {
	*x = TenantMemberGetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_tenant_member_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TenantMemberGetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TenantMemberGetRequest) ProtoMessage() {}

func (x *TenantMemberGetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tenant_member_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TenantMemberGetRequest.ProtoReflect.Descriptor instead.
func (*TenantMemberGetRequest) Descriptor() ([]byte, []int) {
	return file_v1_tenant_member_proto_rawDescGZIP(), []int{4}
}

func (x *TenantMemberGetRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type TenantMemberFindRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TenantId    *string           `protobuf:"bytes,1,opt,name=tenant_id,json=tenantId,proto3,oneof" json:"tenant_id,omitempty"`
	MemberId    *string           `protobuf:"bytes,2,opt,name=member_id,json=memberId,proto3,oneof" json:"member_id,omitempty"`
	Annotations map[string]string `protobuf:"bytes,6,rep,name=annotations,proto3" json:"annotations,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *TenantMemberFindRequest) Reset() {
	*x = TenantMemberFindRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_tenant_member_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TenantMemberFindRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TenantMemberFindRequest) ProtoMessage() {}

func (x *TenantMemberFindRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tenant_member_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TenantMemberFindRequest.ProtoReflect.Descriptor instead.
func (*TenantMemberFindRequest) Descriptor() ([]byte, []int) {
	return file_v1_tenant_member_proto_rawDescGZIP(), []int{5}
}

func (x *TenantMemberFindRequest) GetTenantId() string {
	if x != nil && x.TenantId != nil {
		return *x.TenantId
	}
	return ""
}

func (x *TenantMemberFindRequest) GetMemberId() string {
	if x != nil && x.MemberId != nil {
		return *x.MemberId
	}
	return ""
}

func (x *TenantMemberFindRequest) GetAnnotations() map[string]string {
	if x != nil {
		return x.Annotations
	}
	return nil
}

type TenantMemberResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TenantMember *TenantMember `protobuf:"bytes,1,opt,name=tenant_member,json=tenantMember,proto3" json:"tenant_member,omitempty"`
}

func (x *TenantMemberResponse) Reset() {
	*x = TenantMemberResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_tenant_member_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TenantMemberResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TenantMemberResponse) ProtoMessage() {}

func (x *TenantMemberResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tenant_member_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TenantMemberResponse.ProtoReflect.Descriptor instead.
func (*TenantMemberResponse) Descriptor() ([]byte, []int) {
	return file_v1_tenant_member_proto_rawDescGZIP(), []int{6}
}

func (x *TenantMemberResponse) GetTenantMember() *TenantMember {
	if x != nil {
		return x.TenantMember
	}
	return nil
}

type TenantMemberListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TenantMembers []*TenantMember `protobuf:"bytes,1,rep,name=tenant_members,json=tenantMembers,proto3" json:"tenant_members,omitempty"`
}

func (x *TenantMemberListResponse) Reset() {
	*x = TenantMemberListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_tenant_member_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TenantMemberListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TenantMemberListResponse) ProtoMessage() {}

func (x *TenantMemberListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tenant_member_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TenantMemberListResponse.ProtoReflect.Descriptor instead.
func (*TenantMemberListResponse) Descriptor() ([]byte, []int) {
	return file_v1_tenant_member_proto_rawDescGZIP(), []int{7}
}

func (x *TenantMemberListResponse) GetTenantMembers() []*TenantMember {
	if x != nil {
		return x.TenantMembers
	}
	return nil
}

var File_v1_tenant_member_proto protoreflect.FileDescriptor

var file_v1_tenant_member_proto_rawDesc = []byte{
	0x0a, 0x16, 0x76, 0x31, 0x2f, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x6d, 0x65, 0x6d, 0x62,
	0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x76, 0x31, 0x1a, 0x0d, 0x76, 0x31,
	0x2f, 0x6d, 0x65, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x66, 0x0a, 0x0c, 0x54,
	0x65, 0x6e, 0x61, 0x6e, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x1c, 0x0a, 0x04, 0x6d,
	0x65, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x76, 0x31, 0x2e, 0x4d,
	0x65, 0x74, 0x61, 0x52, 0x04, 0x6d, 0x65, 0x74, 0x61, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x65, 0x6e,
	0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x65,
	0x6e, 0x61, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x65, 0x6d, 0x62, 0x65,
	0x72, 0x49, 0x64, 0x22, 0x52, 0x0a, 0x19, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x4d, 0x65, 0x6d,
	0x62, 0x65, 0x72, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x35, 0x0a, 0x0d, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x6d, 0x65, 0x6d, 0x62, 0x65,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x65, 0x6e,
	0x61, 0x6e, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x0c, 0x74, 0x65, 0x6e, 0x61, 0x6e,
	0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x22, 0x52, 0x0a, 0x19, 0x54, 0x65, 0x6e, 0x61, 0x6e,
	0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x35, 0x0a, 0x0d, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x6d,
	0x65, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x76, 0x31,
	0x2e, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x0c, 0x74,
	0x65, 0x6e, 0x61, 0x6e, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x22, 0x2b, 0x0a, 0x19, 0x54,
	0x65, 0x6e, 0x61, 0x6e, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x28, 0x0a, 0x16, 0x54, 0x65, 0x6e, 0x61,
	0x6e, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x22, 0x89, 0x02, 0x0a, 0x17, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x4d, 0x65, 0x6d,
	0x62, 0x65, 0x72, 0x46, 0x69, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x20,
	0x0a, 0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x48, 0x00, 0x52, 0x08, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x49, 0x64, 0x88, 0x01, 0x01,
	0x12, 0x20, 0x0a, 0x09, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x08, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x49, 0x64, 0x88,
	0x01, 0x01, 0x12, 0x4e, 0x0a, 0x0b, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x65, 0x6e,
	0x61, 0x6e, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x46, 0x69, 0x6e, 0x64, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x2e, 0x41, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0b, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x1a, 0x3e, 0x0a, 0x10, 0x41, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02,
	0x38, 0x01, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64,
	0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x22, 0x4d,
	0x0a, 0x14, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x35, 0x0a, 0x0d, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74,
	0x5f, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e,
	0x76, 0x31, 0x2e, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x52,
	0x0c, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x22, 0x53, 0x0a,
	0x18, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x4c, 0x69, 0x73,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x37, 0x0a, 0x0e, 0x74, 0x65, 0x6e,
	0x61, 0x6e, 0x74, 0x5f, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x10, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x4d, 0x65, 0x6d,
	0x62, 0x65, 0x72, 0x52, 0x0d, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65,
	0x72, 0x73, 0x32, 0xde, 0x02, 0x0a, 0x13, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x4d, 0x65, 0x6d,
	0x62, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x41, 0x0a, 0x06, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x12, 0x1d, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74,
	0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x4d,
	0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a,
	0x06, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x1d, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x65, 0x6e,
	0x61, 0x6e, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x65, 0x6e, 0x61,
	0x6e, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x41, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x1d, 0x2e, 0x76, 0x31, 0x2e,
	0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x76, 0x31, 0x2e, 0x54,
	0x65, 0x6e, 0x61, 0x6e, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x3b, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x1a, 0x2e, 0x76, 0x31, 0x2e,
	0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x47, 0x65, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x65, 0x6e, 0x61,
	0x6e, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x41, 0x0a, 0x04, 0x46, 0x69, 0x6e, 0x64, 0x12, 0x1b, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x65,
	0x6e, 0x61, 0x6e, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x46, 0x69, 0x6e, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x65, 0x6e, 0x61, 0x6e,
	0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x42, 0x6d, 0x0a, 0x06, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x31, 0x42, 0x11, 0x54,
	0x65, 0x6e, 0x61, 0x6e, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x50, 0x01, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d,
	0x65, 0x74, 0x61, 0x6c, 0x2d, 0x73, 0x74, 0x61, 0x63, 0x6b, 0x2f, 0x6d, 0x61, 0x73, 0x74, 0x65,
	0x72, 0x64, 0x61, 0x74, 0x61, 0x2d, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x56,
	0x58, 0x58, 0xaa, 0x02, 0x02, 0x56, 0x31, 0xca, 0x02, 0x02, 0x56, 0x31, 0xe2, 0x02, 0x0e, 0x56,
	0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x02,
	0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_v1_tenant_member_proto_rawDescOnce sync.Once
	file_v1_tenant_member_proto_rawDescData = file_v1_tenant_member_proto_rawDesc
)

func file_v1_tenant_member_proto_rawDescGZIP() []byte {
	file_v1_tenant_member_proto_rawDescOnce.Do(func() {
		file_v1_tenant_member_proto_rawDescData = protoimpl.X.CompressGZIP(file_v1_tenant_member_proto_rawDescData)
	})
	return file_v1_tenant_member_proto_rawDescData
}

var file_v1_tenant_member_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_v1_tenant_member_proto_goTypes = []interface{}{
	(*TenantMember)(nil),              // 0: v1.TenantMember
	(*TenantMemberCreateRequest)(nil), // 1: v1.TenantMemberCreateRequest
	(*TenantMemberUpdateRequest)(nil), // 2: v1.TenantMemberUpdateRequest
	(*TenantMemberDeleteRequest)(nil), // 3: v1.TenantMemberDeleteRequest
	(*TenantMemberGetRequest)(nil),    // 4: v1.TenantMemberGetRequest
	(*TenantMemberFindRequest)(nil),   // 5: v1.TenantMemberFindRequest
	(*TenantMemberResponse)(nil),      // 6: v1.TenantMemberResponse
	(*TenantMemberListResponse)(nil),  // 7: v1.TenantMemberListResponse
	nil,                               // 8: v1.TenantMemberFindRequest.AnnotationsEntry
	(*Meta)(nil),                      // 9: v1.Meta
}
var file_v1_tenant_member_proto_depIdxs = []int32{
	9,  // 0: v1.TenantMember.meta:type_name -> v1.Meta
	0,  // 1: v1.TenantMemberCreateRequest.tenant_member:type_name -> v1.TenantMember
	0,  // 2: v1.TenantMemberUpdateRequest.tenant_member:type_name -> v1.TenantMember
	8,  // 3: v1.TenantMemberFindRequest.annotations:type_name -> v1.TenantMemberFindRequest.AnnotationsEntry
	0,  // 4: v1.TenantMemberResponse.tenant_member:type_name -> v1.TenantMember
	0,  // 5: v1.TenantMemberListResponse.tenant_members:type_name -> v1.TenantMember
	1,  // 6: v1.TenantMemberService.Create:input_type -> v1.TenantMemberCreateRequest
	2,  // 7: v1.TenantMemberService.Update:input_type -> v1.TenantMemberUpdateRequest
	3,  // 8: v1.TenantMemberService.Delete:input_type -> v1.TenantMemberDeleteRequest
	4,  // 9: v1.TenantMemberService.Get:input_type -> v1.TenantMemberGetRequest
	5,  // 10: v1.TenantMemberService.Find:input_type -> v1.TenantMemberFindRequest
	6,  // 11: v1.TenantMemberService.Create:output_type -> v1.TenantMemberResponse
	6,  // 12: v1.TenantMemberService.Update:output_type -> v1.TenantMemberResponse
	6,  // 13: v1.TenantMemberService.Delete:output_type -> v1.TenantMemberResponse
	6,  // 14: v1.TenantMemberService.Get:output_type -> v1.TenantMemberResponse
	7,  // 15: v1.TenantMemberService.Find:output_type -> v1.TenantMemberListResponse
	11, // [11:16] is the sub-list for method output_type
	6,  // [6:11] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_v1_tenant_member_proto_init() }
func file_v1_tenant_member_proto_init() {
	if File_v1_tenant_member_proto != nil {
		return
	}
	file_v1_meta_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_v1_tenant_member_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TenantMember); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_v1_tenant_member_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TenantMemberCreateRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_v1_tenant_member_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TenantMemberUpdateRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_v1_tenant_member_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TenantMemberDeleteRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_v1_tenant_member_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TenantMemberGetRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_v1_tenant_member_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TenantMemberFindRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_v1_tenant_member_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TenantMemberResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_v1_tenant_member_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TenantMemberListResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_v1_tenant_member_proto_msgTypes[5].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_v1_tenant_member_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_v1_tenant_member_proto_goTypes,
		DependencyIndexes: file_v1_tenant_member_proto_depIdxs,
		MessageInfos:      file_v1_tenant_member_proto_msgTypes,
	}.Build()
	File_v1_tenant_member_proto = out.File
	file_v1_tenant_member_proto_rawDesc = nil
	file_v1_tenant_member_proto_goTypes = nil
	file_v1_tenant_member_proto_depIdxs = nil
}
