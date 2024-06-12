// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: v1/project_member.proto

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

// ProjectMember is the database model
type ProjectMember struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Meta      *Meta  `protobuf:"bytes,1,opt,name=meta,proto3" json:"meta,omitempty"`
	ProjectId string `protobuf:"bytes,2,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	TenantId  string `protobuf:"bytes,4,opt,name=tenant_id,json=tenantId,proto3" json:"tenant_id,omitempty"`
}

func (x *ProjectMember) Reset() {
	*x = ProjectMember{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_project_member_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProjectMember) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProjectMember) ProtoMessage() {}

func (x *ProjectMember) ProtoReflect() protoreflect.Message {
	mi := &file_v1_project_member_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProjectMember.ProtoReflect.Descriptor instead.
func (*ProjectMember) Descriptor() ([]byte, []int) {
	return file_v1_project_member_proto_rawDescGZIP(), []int{0}
}

func (x *ProjectMember) GetMeta() *Meta {
	if x != nil {
		return x.Meta
	}
	return nil
}

func (x *ProjectMember) GetProjectId() string {
	if x != nil {
		return x.ProjectId
	}
	return ""
}

func (x *ProjectMember) GetTenantId() string {
	if x != nil {
		return x.TenantId
	}
	return ""
}

type ProjectMemberCreateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectMember *ProjectMember `protobuf:"bytes,1,opt,name=project_member,json=projectMember,proto3" json:"project_member,omitempty"`
}

func (x *ProjectMemberCreateRequest) Reset() {
	*x = ProjectMemberCreateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_project_member_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProjectMemberCreateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProjectMemberCreateRequest) ProtoMessage() {}

func (x *ProjectMemberCreateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_project_member_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProjectMemberCreateRequest.ProtoReflect.Descriptor instead.
func (*ProjectMemberCreateRequest) Descriptor() ([]byte, []int) {
	return file_v1_project_member_proto_rawDescGZIP(), []int{1}
}

func (x *ProjectMemberCreateRequest) GetProjectMember() *ProjectMember {
	if x != nil {
		return x.ProjectMember
	}
	return nil
}

type ProjectMemberUpdateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectMember *ProjectMember `protobuf:"bytes,1,opt,name=project_member,json=projectMember,proto3" json:"project_member,omitempty"`
}

func (x *ProjectMemberUpdateRequest) Reset() {
	*x = ProjectMemberUpdateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_project_member_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProjectMemberUpdateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProjectMemberUpdateRequest) ProtoMessage() {}

func (x *ProjectMemberUpdateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_project_member_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProjectMemberUpdateRequest.ProtoReflect.Descriptor instead.
func (*ProjectMemberUpdateRequest) Descriptor() ([]byte, []int) {
	return file_v1_project_member_proto_rawDescGZIP(), []int{2}
}

func (x *ProjectMemberUpdateRequest) GetProjectMember() *ProjectMember {
	if x != nil {
		return x.ProjectMember
	}
	return nil
}

type ProjectMemberDeleteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *ProjectMemberDeleteRequest) Reset() {
	*x = ProjectMemberDeleteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_project_member_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProjectMemberDeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProjectMemberDeleteRequest) ProtoMessage() {}

func (x *ProjectMemberDeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_project_member_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProjectMemberDeleteRequest.ProtoReflect.Descriptor instead.
func (*ProjectMemberDeleteRequest) Descriptor() ([]byte, []int) {
	return file_v1_project_member_proto_rawDescGZIP(), []int{3}
}

func (x *ProjectMemberDeleteRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type ProjectMemberGetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *ProjectMemberGetRequest) Reset() {
	*x = ProjectMemberGetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_project_member_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProjectMemberGetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProjectMemberGetRequest) ProtoMessage() {}

func (x *ProjectMemberGetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_project_member_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProjectMemberGetRequest.ProtoReflect.Descriptor instead.
func (*ProjectMemberGetRequest) Descriptor() ([]byte, []int) {
	return file_v1_project_member_proto_rawDescGZIP(), []int{4}
}

func (x *ProjectMemberGetRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type ProjectMemberFindRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectId   *string           `protobuf:"bytes,1,opt,name=project_id,json=projectId,proto3,oneof" json:"project_id,omitempty"`
	TenantId    *string           `protobuf:"bytes,2,opt,name=tenant_id,json=tenantId,proto3,oneof" json:"tenant_id,omitempty"`
	Annotations map[string]string `protobuf:"bytes,6,rep,name=annotations,proto3" json:"annotations,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *ProjectMemberFindRequest) Reset() {
	*x = ProjectMemberFindRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_project_member_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProjectMemberFindRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProjectMemberFindRequest) ProtoMessage() {}

func (x *ProjectMemberFindRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_project_member_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProjectMemberFindRequest.ProtoReflect.Descriptor instead.
func (*ProjectMemberFindRequest) Descriptor() ([]byte, []int) {
	return file_v1_project_member_proto_rawDescGZIP(), []int{5}
}

func (x *ProjectMemberFindRequest) GetProjectId() string {
	if x != nil && x.ProjectId != nil {
		return *x.ProjectId
	}
	return ""
}

func (x *ProjectMemberFindRequest) GetTenantId() string {
	if x != nil && x.TenantId != nil {
		return *x.TenantId
	}
	return ""
}

func (x *ProjectMemberFindRequest) GetAnnotations() map[string]string {
	if x != nil {
		return x.Annotations
	}
	return nil
}

type ProjectMemberResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectMember *ProjectMember `protobuf:"bytes,1,opt,name=project_member,json=projectMember,proto3" json:"project_member,omitempty"`
}

func (x *ProjectMemberResponse) Reset() {
	*x = ProjectMemberResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_project_member_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProjectMemberResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProjectMemberResponse) ProtoMessage() {}

func (x *ProjectMemberResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v1_project_member_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProjectMemberResponse.ProtoReflect.Descriptor instead.
func (*ProjectMemberResponse) Descriptor() ([]byte, []int) {
	return file_v1_project_member_proto_rawDescGZIP(), []int{6}
}

func (x *ProjectMemberResponse) GetProjectMember() *ProjectMember {
	if x != nil {
		return x.ProjectMember
	}
	return nil
}

type ProjectMemberListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectMembers []*ProjectMember `protobuf:"bytes,1,rep,name=project_members,json=projectMembers,proto3" json:"project_members,omitempty"`
}

func (x *ProjectMemberListResponse) Reset() {
	*x = ProjectMemberListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_project_member_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProjectMemberListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProjectMemberListResponse) ProtoMessage() {}

func (x *ProjectMemberListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v1_project_member_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProjectMemberListResponse.ProtoReflect.Descriptor instead.
func (*ProjectMemberListResponse) Descriptor() ([]byte, []int) {
	return file_v1_project_member_proto_rawDescGZIP(), []int{7}
}

func (x *ProjectMemberListResponse) GetProjectMembers() []*ProjectMember {
	if x != nil {
		return x.ProjectMembers
	}
	return nil
}

var File_v1_project_member_proto protoreflect.FileDescriptor

var file_v1_project_member_proto_rawDesc = []byte{
	0x0a, 0x17, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x6d, 0x65, 0x6d,
	0x62, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x76, 0x31, 0x1a, 0x0d, 0x76,
	0x31, 0x2f, 0x6d, 0x65, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x69, 0x0a, 0x0d,
	0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x1c, 0x0a,
	0x04, 0x6d, 0x65, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x76, 0x31,
	0x2e, 0x4d, 0x65, 0x74, 0x61, 0x52, 0x04, 0x6d, 0x65, 0x74, 0x61, 0x12, 0x1d, 0x0a, 0x0a, 0x70,
	0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x65,
	0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74,
	0x65, 0x6e, 0x61, 0x6e, 0x74, 0x49, 0x64, 0x22, 0x56, 0x0a, 0x1a, 0x50, 0x72, 0x6f, 0x6a, 0x65,
	0x63, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x38, 0x0a, 0x0e, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x5f, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e,
	0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72,
	0x52, 0x0d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x22,
	0x56, 0x0a, 0x1a, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x38, 0x0a,
	0x0e, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x6a, 0x65,
	0x63, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x0d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x22, 0x2c, 0x0a, 0x1a, 0x50, 0x72, 0x6f, 0x6a, 0x65,
	0x63, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x29, 0x0a, 0x17, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x22, 0x8e, 0x02, 0x0a, 0x18, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x4d, 0x65, 0x6d, 0x62,
	0x65, 0x72, 0x46, 0x69, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x22, 0x0a,
	0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x48, 0x00, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x88, 0x01,
	0x01, 0x12, 0x20, 0x0a, 0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x08, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x49, 0x64,
	0x88, 0x01, 0x01, 0x12, 0x4f, 0x0a, 0x0b, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72,
	0x6f, 0x6a, 0x65, 0x63, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x46, 0x69, 0x6e, 0x64, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x41, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0b, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x1a, 0x3e, 0x0a, 0x10, 0x41, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x3a, 0x02, 0x38, 0x01, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x5f, 0x69, 0x64, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69,
	0x64, 0x22, 0x51, 0x0a, 0x15, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x4d, 0x65, 0x6d, 0x62,
	0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x38, 0x0a, 0x0e, 0x70, 0x72,
	0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x11, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x4d,
	0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x0d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x4d, 0x65,
	0x6d, 0x62, 0x65, 0x72, 0x22, 0x57, 0x0a, 0x19, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x4d,
	0x65, 0x6d, 0x62, 0x65, 0x72, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x3a, 0x0a, 0x0f, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x6d, 0x65, 0x6d,
	0x62, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x76, 0x31, 0x2e,
	0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x0e, 0x70,
	0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x32, 0xe9, 0x02,
	0x0a, 0x14, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x43, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x12, 0x1e, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x4d, 0x65, 0x6d,
	0x62, 0x65, 0x72, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x19, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x4d, 0x65, 0x6d,
	0x62, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x43, 0x0a, 0x06, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x1e, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x6a, 0x65,
	0x63, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x6a, 0x65,
	0x63, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x43, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x1e, 0x2e, 0x76, 0x31, 0x2e,
	0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x76, 0x31, 0x2e,
	0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3d, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x1b, 0x2e, 0x76,
	0x31, 0x2e, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x47,
	0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x76, 0x31, 0x2e, 0x50,
	0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x43, 0x0a, 0x04, 0x46, 0x69, 0x6e, 0x64, 0x12, 0x1c, 0x2e, 0x76,
	0x31, 0x2e, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x46,
	0x69, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x76, 0x31, 0x2e,
	0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x4c, 0x69, 0x73,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x6e, 0x0a, 0x06, 0x63, 0x6f, 0x6d,
	0x2e, 0x76, 0x31, 0x42, 0x12, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x4d, 0x65, 0x6d, 0x62,
	0x65, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x65, 0x74, 0x61, 0x6c, 0x2d, 0x73, 0x74, 0x61, 0x63,
	0x6b, 0x2f, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x64, 0x61, 0x74, 0x61, 0x2d, 0x61, 0x70, 0x69,
	0x2f, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x56, 0x58, 0x58, 0xaa, 0x02, 0x02, 0x56, 0x31, 0xca, 0x02,
	0x02, 0x56, 0x31, 0xe2, 0x02, 0x0e, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x02, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_v1_project_member_proto_rawDescOnce sync.Once
	file_v1_project_member_proto_rawDescData = file_v1_project_member_proto_rawDesc
)

func file_v1_project_member_proto_rawDescGZIP() []byte {
	file_v1_project_member_proto_rawDescOnce.Do(func() {
		file_v1_project_member_proto_rawDescData = protoimpl.X.CompressGZIP(file_v1_project_member_proto_rawDescData)
	})
	return file_v1_project_member_proto_rawDescData
}

var file_v1_project_member_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_v1_project_member_proto_goTypes = []any{
	(*ProjectMember)(nil),              // 0: v1.ProjectMember
	(*ProjectMemberCreateRequest)(nil), // 1: v1.ProjectMemberCreateRequest
	(*ProjectMemberUpdateRequest)(nil), // 2: v1.ProjectMemberUpdateRequest
	(*ProjectMemberDeleteRequest)(nil), // 3: v1.ProjectMemberDeleteRequest
	(*ProjectMemberGetRequest)(nil),    // 4: v1.ProjectMemberGetRequest
	(*ProjectMemberFindRequest)(nil),   // 5: v1.ProjectMemberFindRequest
	(*ProjectMemberResponse)(nil),      // 6: v1.ProjectMemberResponse
	(*ProjectMemberListResponse)(nil),  // 7: v1.ProjectMemberListResponse
	nil,                                // 8: v1.ProjectMemberFindRequest.AnnotationsEntry
	(*Meta)(nil),                       // 9: v1.Meta
}
var file_v1_project_member_proto_depIdxs = []int32{
	9,  // 0: v1.ProjectMember.meta:type_name -> v1.Meta
	0,  // 1: v1.ProjectMemberCreateRequest.project_member:type_name -> v1.ProjectMember
	0,  // 2: v1.ProjectMemberUpdateRequest.project_member:type_name -> v1.ProjectMember
	8,  // 3: v1.ProjectMemberFindRequest.annotations:type_name -> v1.ProjectMemberFindRequest.AnnotationsEntry
	0,  // 4: v1.ProjectMemberResponse.project_member:type_name -> v1.ProjectMember
	0,  // 5: v1.ProjectMemberListResponse.project_members:type_name -> v1.ProjectMember
	1,  // 6: v1.ProjectMemberService.Create:input_type -> v1.ProjectMemberCreateRequest
	2,  // 7: v1.ProjectMemberService.Update:input_type -> v1.ProjectMemberUpdateRequest
	3,  // 8: v1.ProjectMemberService.Delete:input_type -> v1.ProjectMemberDeleteRequest
	4,  // 9: v1.ProjectMemberService.Get:input_type -> v1.ProjectMemberGetRequest
	5,  // 10: v1.ProjectMemberService.Find:input_type -> v1.ProjectMemberFindRequest
	6,  // 11: v1.ProjectMemberService.Create:output_type -> v1.ProjectMemberResponse
	6,  // 12: v1.ProjectMemberService.Update:output_type -> v1.ProjectMemberResponse
	6,  // 13: v1.ProjectMemberService.Delete:output_type -> v1.ProjectMemberResponse
	6,  // 14: v1.ProjectMemberService.Get:output_type -> v1.ProjectMemberResponse
	7,  // 15: v1.ProjectMemberService.Find:output_type -> v1.ProjectMemberListResponse
	11, // [11:16] is the sub-list for method output_type
	6,  // [6:11] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_v1_project_member_proto_init() }
func file_v1_project_member_proto_init() {
	if File_v1_project_member_proto != nil {
		return
	}
	file_v1_meta_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_v1_project_member_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*ProjectMember); i {
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
		file_v1_project_member_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*ProjectMemberCreateRequest); i {
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
		file_v1_project_member_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*ProjectMemberUpdateRequest); i {
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
		file_v1_project_member_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*ProjectMemberDeleteRequest); i {
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
		file_v1_project_member_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*ProjectMemberGetRequest); i {
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
		file_v1_project_member_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*ProjectMemberFindRequest); i {
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
		file_v1_project_member_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*ProjectMemberResponse); i {
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
		file_v1_project_member_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*ProjectMemberListResponse); i {
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
	file_v1_project_member_proto_msgTypes[5].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_v1_project_member_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_v1_project_member_proto_goTypes,
		DependencyIndexes: file_v1_project_member_proto_depIdxs,
		MessageInfos:      file_v1_project_member_proto_msgTypes,
	}.Build()
	File_v1_project_member_proto = out.File
	file_v1_project_member_proto_rawDesc = nil
	file_v1_project_member_proto_goTypes = nil
	file_v1_project_member_proto_depIdxs = nil
}
