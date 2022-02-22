// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.4
// source: v1/iam.proto

package v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type IAMConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IssuerConfig *IssuerConfig         `protobuf:"bytes,1,opt,name=issuer_config,json=issuerConfig,proto3" json:"issuer_config,omitempty"`
	IdmConfig    *IDMConfig            `protobuf:"bytes,2,opt,name=idm_config,json=idmConfig,proto3" json:"idm_config,omitempty"`
	GroupConfig  *NamespaceGroupConfig `protobuf:"bytes,3,opt,name=group_config,json=groupConfig,proto3" json:"group_config,omitempty"`
}

func (x *IAMConfig) Reset() {
	*x = IAMConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_iam_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IAMConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IAMConfig) ProtoMessage() {}

func (x *IAMConfig) ProtoReflect() protoreflect.Message {
	mi := &file_v1_iam_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IAMConfig.ProtoReflect.Descriptor instead.
func (*IAMConfig) Descriptor() ([]byte, []int) {
	return file_v1_iam_proto_rawDescGZIP(), []int{0}
}

func (x *IAMConfig) GetIssuerConfig() *IssuerConfig {
	if x != nil {
		return x.IssuerConfig
	}
	return nil
}

func (x *IAMConfig) GetIdmConfig() *IDMConfig {
	if x != nil {
		return x.IdmConfig
	}
	return nil
}

func (x *IAMConfig) GetGroupConfig() *NamespaceGroupConfig {
	if x != nil {
		return x.GroupConfig
	}
	return nil
}

type IssuerConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url      string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	ClientId string `protobuf:"bytes,2,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
}

func (x *IssuerConfig) Reset() {
	*x = IssuerConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_iam_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IssuerConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IssuerConfig) ProtoMessage() {}

func (x *IssuerConfig) ProtoReflect() protoreflect.Message {
	mi := &file_v1_iam_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IssuerConfig.ProtoReflect.Descriptor instead.
func (*IssuerConfig) Descriptor() ([]byte, []int) {
	return file_v1_iam_proto_rawDescGZIP(), []int{1}
}

func (x *IssuerConfig) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *IssuerConfig) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

// mandatory config
type IDMConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IdmType string `protobuf:"bytes,1,opt,name=idm_type,json=idmType,proto3" json:"idm_type,omitempty"`
	// optional
	ConnectorConfig *ConnectorConfig `protobuf:"bytes,2,opt,name=connector_config,json=connectorConfig,proto3" json:"connector_config,omitempty"`
}

func (x *IDMConfig) Reset() {
	*x = IDMConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_iam_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IDMConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IDMConfig) ProtoMessage() {}

func (x *IDMConfig) ProtoReflect() protoreflect.Message {
	mi := &file_v1_iam_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IDMConfig.ProtoReflect.Descriptor instead.
func (*IDMConfig) Descriptor() ([]byte, []int) {
	return file_v1_iam_proto_rawDescGZIP(), []int{2}
}

func (x *IDMConfig) GetIdmType() string {
	if x != nil {
		return x.IdmType
	}
	return ""
}

func (x *IDMConfig) GetConnectorConfig() *ConnectorConfig {
	if x != nil {
		return x.ConnectorConfig
	}
	return nil
}

// Config for group-rolebinding-controller
type NamespaceGroupConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// no action is taken or any namespace in this list
	ExcludedNamespaces string `protobuf:"bytes,1,opt,name=excluded_namespaces,json=excludedNamespaces,proto3" json:"excluded_namespaces,omitempty"`
	// for each element a RoleBinding is created in any Namespace - ClusterRoles are bound with this name
	// admin,edit,view
	ExpectedGroupsList string `protobuf:"bytes,2,opt,name=expected_groups_list,json=expectedGroupsList,proto3" json:"expected_groups_list,omitempty"`
	// Maximum length of namespace-part in clusterGroupname and therefore in the corresponding groupname in the directory.
	// 20 chars für AD, given the naming-conventions
	NamespaceMaxLength int32 `protobuf:"varint,3,opt,name=namespace_max_length,json=namespaceMaxLength,proto3" json:"namespace_max_length,omitempty"`
	// The created RoleBindings will reference this group (from token).
	// oidc:{{ .Namespace }}-{{ .Group }}
	ClusterGroupnameTemplate string `protobuf:"bytes,4,opt,name=cluster_groupname_template,json=clusterGroupnameTemplate,proto3" json:"cluster_groupname_template,omitempty"`
	// The RoleBindings will created with this name.
	// oidc-{{ .Namespace }}-{{ .Group }}
	RolebindingNameTemplate string `protobuf:"bytes,5,opt,name=rolebinding_name_template,json=rolebindingNameTemplate,proto3" json:"rolebinding_name_template,omitempty"`
}

func (x *NamespaceGroupConfig) Reset() {
	*x = NamespaceGroupConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_iam_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NamespaceGroupConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NamespaceGroupConfig) ProtoMessage() {}

func (x *NamespaceGroupConfig) ProtoReflect() protoreflect.Message {
	mi := &file_v1_iam_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NamespaceGroupConfig.ProtoReflect.Descriptor instead.
func (*NamespaceGroupConfig) Descriptor() ([]byte, []int) {
	return file_v1_iam_proto_rawDescGZIP(), []int{3}
}

func (x *NamespaceGroupConfig) GetExcludedNamespaces() string {
	if x != nil {
		return x.ExcludedNamespaces
	}
	return ""
}

func (x *NamespaceGroupConfig) GetExpectedGroupsList() string {
	if x != nil {
		return x.ExpectedGroupsList
	}
	return ""
}

func (x *NamespaceGroupConfig) GetNamespaceMaxLength() int32 {
	if x != nil {
		return x.NamespaceMaxLength
	}
	return 0
}

func (x *NamespaceGroupConfig) GetClusterGroupnameTemplate() string {
	if x != nil {
		return x.ClusterGroupnameTemplate
	}
	return ""
}

func (x *NamespaceGroupConfig) GetRolebindingNameTemplate() string {
	if x != nil {
		return x.RolebindingNameTemplate
	}
	return ""
}

// optional config if idm webhook is used to automatically create/delete groups/roles in the tenant idm
type ConnectorConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// the following are all mandatory
	IdmApiUrl      string `protobuf:"bytes,1,opt,name=idm_api_url,json=idmApiUrl,proto3" json:"idm_api_url,omitempty"`
	IdmApiUser     string `protobuf:"bytes,2,opt,name=idm_api_user,json=idmApiUser,proto3" json:"idm_api_user,omitempty"`
	IdmApiPassword string `protobuf:"bytes,3,opt,name=idm_api_password,json=idmApiPassword,proto3" json:"idm_api_password,omitempty"`
	IdmSystemId    string `protobuf:"bytes,4,opt,name=idm_system_id,json=idmSystemId,proto3" json:"idm_system_id,omitempty"`
	IdmAccessCode  string `protobuf:"bytes,5,opt,name=idm_access_code,json=idmAccessCode,proto3" json:"idm_access_code,omitempty"`
	IdmCustomerId  string `protobuf:"bytes,6,opt,name=idm_customer_id,json=idmCustomerId,proto3" json:"idm_customer_id,omitempty"`
	IdmGroupOu     string `protobuf:"bytes,7,opt,name=idm_group_ou,json=idmGroupOu,proto3" json:"idm_group_ou,omitempty"`
	// optional
	IdmGroupnameTemplate *wrapperspb.StringValue `protobuf:"bytes,8,opt,name=idm_groupname_template,json=idmGroupnameTemplate,proto3" json:"idm_groupname_template,omitempty"`
	IdmDomainName        string                  `protobuf:"bytes,9,opt,name=idm_domain_name,json=idmDomainName,proto3" json:"idm_domain_name,omitempty"`
	IdmTenantPrefix      string                  `protobuf:"bytes,10,opt,name=idm_tenant_prefix,json=idmTenantPrefix,proto3" json:"idm_tenant_prefix,omitempty"`
	IdmSubmitter         string                  `protobuf:"bytes,11,opt,name=idm_submitter,json=idmSubmitter,proto3" json:"idm_submitter,omitempty"`
	IdmJobInfo           string                  `protobuf:"bytes,12,opt,name=idm_job_info,json=idmJobInfo,proto3" json:"idm_job_info,omitempty"`
	IdmReqSystem         string                  `protobuf:"bytes,13,opt,name=idm_req_system,json=idmReqSystem,proto3" json:"idm_req_system,omitempty"`
	IdmReqUser           string                  `protobuf:"bytes,14,opt,name=idm_req_user,json=idmReqUser,proto3" json:"idm_req_user,omitempty"`
	IdmReqEmail          string                  `protobuf:"bytes,15,opt,name=idm_req_email,json=idmReqEmail,proto3" json:"idm_req_email,omitempty"`
}

func (x *ConnectorConfig) Reset() {
	*x = ConnectorConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_iam_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConnectorConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnectorConfig) ProtoMessage() {}

func (x *ConnectorConfig) ProtoReflect() protoreflect.Message {
	mi := &file_v1_iam_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnectorConfig.ProtoReflect.Descriptor instead.
func (*ConnectorConfig) Descriptor() ([]byte, []int) {
	return file_v1_iam_proto_rawDescGZIP(), []int{4}
}

func (x *ConnectorConfig) GetIdmApiUrl() string {
	if x != nil {
		return x.IdmApiUrl
	}
	return ""
}

func (x *ConnectorConfig) GetIdmApiUser() string {
	if x != nil {
		return x.IdmApiUser
	}
	return ""
}

func (x *ConnectorConfig) GetIdmApiPassword() string {
	if x != nil {
		return x.IdmApiPassword
	}
	return ""
}

func (x *ConnectorConfig) GetIdmSystemId() string {
	if x != nil {
		return x.IdmSystemId
	}
	return ""
}

func (x *ConnectorConfig) GetIdmAccessCode() string {
	if x != nil {
		return x.IdmAccessCode
	}
	return ""
}

func (x *ConnectorConfig) GetIdmCustomerId() string {
	if x != nil {
		return x.IdmCustomerId
	}
	return ""
}

func (x *ConnectorConfig) GetIdmGroupOu() string {
	if x != nil {
		return x.IdmGroupOu
	}
	return ""
}

func (x *ConnectorConfig) GetIdmGroupnameTemplate() *wrapperspb.StringValue {
	if x != nil {
		return x.IdmGroupnameTemplate
	}
	return nil
}

func (x *ConnectorConfig) GetIdmDomainName() string {
	if x != nil {
		return x.IdmDomainName
	}
	return ""
}

func (x *ConnectorConfig) GetIdmTenantPrefix() string {
	if x != nil {
		return x.IdmTenantPrefix
	}
	return ""
}

func (x *ConnectorConfig) GetIdmSubmitter() string {
	if x != nil {
		return x.IdmSubmitter
	}
	return ""
}

func (x *ConnectorConfig) GetIdmJobInfo() string {
	if x != nil {
		return x.IdmJobInfo
	}
	return ""
}

func (x *ConnectorConfig) GetIdmReqSystem() string {
	if x != nil {
		return x.IdmReqSystem
	}
	return ""
}

func (x *ConnectorConfig) GetIdmReqUser() string {
	if x != nil {
		return x.IdmReqUser
	}
	return ""
}

func (x *ConnectorConfig) GetIdmReqEmail() string {
	if x != nil {
		return x.IdmReqEmail
	}
	return ""
}

var File_v1_iam_proto protoreflect.FileDescriptor

var file_v1_iam_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x76, 0x31, 0x2f, 0x69, 0x61, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02,
	0x76, 0x31, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xad, 0x01, 0x0a, 0x09, 0x49, 0x41, 0x4d, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x12, 0x35, 0x0a, 0x0d, 0x69, 0x73, 0x73, 0x75, 0x65, 0x72, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x73, 0x73,
	0x75, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x0c, 0x69, 0x73, 0x73, 0x75, 0x65,
	0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x2c, 0x0a, 0x0a, 0x69, 0x64, 0x6d, 0x5f, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x76, 0x31,
	0x2e, 0x49, 0x44, 0x4d, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x09, 0x69, 0x64, 0x6d, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x3b, 0x0a, 0x0c, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x76, 0x31,
	0x2e, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x0b, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x22, 0x3d, 0x0a, 0x0c, 0x49, 0x73, 0x73, 0x75, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x75, 0x72, 0x6c, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49,
	0x64, 0x22, 0x66, 0x0a, 0x09, 0x49, 0x44, 0x4d, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x19,
	0x0a, 0x08, 0x69, 0x64, 0x6d, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x69, 0x64, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x12, 0x3e, 0x0a, 0x10, 0x63, 0x6f, 0x6e,
	0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x0f, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63,
	0x74, 0x6f, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22, 0xa5, 0x02, 0x0a, 0x14, 0x4e, 0x61,
	0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x12, 0x2f, 0x0a, 0x13, 0x65, 0x78, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x64, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x12, 0x65, 0x78, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61,
	0x63, 0x65, 0x73, 0x12, 0x30, 0x0a, 0x14, 0x65, 0x78, 0x70, 0x65, 0x63, 0x74, 0x65, 0x64, 0x5f,
	0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x12, 0x65, 0x78, 0x70, 0x65, 0x63, 0x74, 0x65, 0x64, 0x47, 0x72, 0x6f, 0x75, 0x70,
	0x73, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x30, 0x0a, 0x14, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61,
	0x63, 0x65, 0x5f, 0x6d, 0x61, 0x78, 0x5f, 0x6c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x12, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x4d, 0x61,
	0x78, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x12, 0x3c, 0x0a, 0x1a, 0x63, 0x6c, 0x75, 0x73, 0x74,
	0x65, 0x72, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x6e, 0x61, 0x6d, 0x65, 0x5f, 0x74, 0x65, 0x6d,
	0x70, 0x6c, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x18, 0x63, 0x6c, 0x75,
	0x73, 0x74, 0x65, 0x72, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x6e, 0x61, 0x6d, 0x65, 0x54, 0x65, 0x6d,
	0x70, 0x6c, 0x61, 0x74, 0x65, 0x12, 0x3a, 0x0a, 0x19, 0x72, 0x6f, 0x6c, 0x65, 0x62, 0x69, 0x6e,
	0x64, 0x69, 0x6e, 0x67, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x5f, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61,
	0x74, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x17, 0x72, 0x6f, 0x6c, 0x65, 0x62, 0x69,
	0x6e, 0x64, 0x69, 0x6e, 0x67, 0x4e, 0x61, 0x6d, 0x65, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74,
	0x65, 0x22, 0xee, 0x04, 0x0a, 0x0f, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x1e, 0x0a, 0x0b, 0x69, 0x64, 0x6d, 0x5f, 0x61, 0x70, 0x69,
	0x5f, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x69, 0x64, 0x6d, 0x41,
	0x70, 0x69, 0x55, 0x72, 0x6c, 0x12, 0x20, 0x0a, 0x0c, 0x69, 0x64, 0x6d, 0x5f, 0x61, 0x70, 0x69,
	0x5f, 0x75, 0x73, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x69, 0x64, 0x6d,
	0x41, 0x70, 0x69, 0x55, 0x73, 0x65, 0x72, 0x12, 0x28, 0x0a, 0x10, 0x69, 0x64, 0x6d, 0x5f, 0x61,
	0x70, 0x69, 0x5f, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0e, 0x69, 0x64, 0x6d, 0x41, 0x70, 0x69, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72,
	0x64, 0x12, 0x22, 0x0a, 0x0d, 0x69, 0x64, 0x6d, 0x5f, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x5f,
	0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x69, 0x64, 0x6d, 0x53, 0x79, 0x73,
	0x74, 0x65, 0x6d, 0x49, 0x64, 0x12, 0x26, 0x0a, 0x0f, 0x69, 0x64, 0x6d, 0x5f, 0x61, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d,
	0x69, 0x64, 0x6d, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x26, 0x0a,
	0x0f, 0x69, 0x64, 0x6d, 0x5f, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x69, 0x64, 0x6d, 0x43, 0x75, 0x73, 0x74, 0x6f,
	0x6d, 0x65, 0x72, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0c, 0x69, 0x64, 0x6d, 0x5f, 0x67, 0x72, 0x6f,
	0x75, 0x70, 0x5f, 0x6f, 0x75, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x69, 0x64, 0x6d,
	0x47, 0x72, 0x6f, 0x75, 0x70, 0x4f, 0x75, 0x12, 0x52, 0x0a, 0x16, 0x69, 0x64, 0x6d, 0x5f, 0x67,
	0x72, 0x6f, 0x75, 0x70, 0x6e, 0x61, 0x6d, 0x65, 0x5f, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74,
	0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x14, 0x69, 0x64, 0x6d, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x6e,
	0x61, 0x6d, 0x65, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x12, 0x26, 0x0a, 0x0f, 0x69,
	0x64, 0x6d, 0x5f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x09,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x69, 0x64, 0x6d, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x2a, 0x0a, 0x11, 0x69, 0x64, 0x6d, 0x5f, 0x74, 0x65, 0x6e, 0x61, 0x6e,
	0x74, 0x5f, 0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f,
	0x69, 0x64, 0x6d, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x50, 0x72, 0x65, 0x66, 0x69, 0x78, 0x12,
	0x23, 0x0a, 0x0d, 0x69, 0x64, 0x6d, 0x5f, 0x73, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x74, 0x65, 0x72,
	0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x69, 0x64, 0x6d, 0x53, 0x75, 0x62, 0x6d, 0x69,
	0x74, 0x74, 0x65, 0x72, 0x12, 0x20, 0x0a, 0x0c, 0x69, 0x64, 0x6d, 0x5f, 0x6a, 0x6f, 0x62, 0x5f,
	0x69, 0x6e, 0x66, 0x6f, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x69, 0x64, 0x6d, 0x4a,
	0x6f, 0x62, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x24, 0x0a, 0x0e, 0x69, 0x64, 0x6d, 0x5f, 0x72, 0x65,
	0x71, 0x5f, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c,
	0x69, 0x64, 0x6d, 0x52, 0x65, 0x71, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x12, 0x20, 0x0a, 0x0c,
	0x69, 0x64, 0x6d, 0x5f, 0x72, 0x65, 0x71, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x18, 0x0e, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x69, 0x64, 0x6d, 0x52, 0x65, 0x71, 0x55, 0x73, 0x65, 0x72, 0x12, 0x22,
	0x0a, 0x0d, 0x69, 0x64, 0x6d, 0x5f, 0x72, 0x65, 0x71, 0x5f, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18,
	0x0f, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x69, 0x64, 0x6d, 0x52, 0x65, 0x71, 0x45, 0x6d, 0x61,
	0x69, 0x6c, 0x42, 0x06, 0x5a, 0x04, 0x2e, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_v1_iam_proto_rawDescOnce sync.Once
	file_v1_iam_proto_rawDescData = file_v1_iam_proto_rawDesc
)

func file_v1_iam_proto_rawDescGZIP() []byte {
	file_v1_iam_proto_rawDescOnce.Do(func() {
		file_v1_iam_proto_rawDescData = protoimpl.X.CompressGZIP(file_v1_iam_proto_rawDescData)
	})
	return file_v1_iam_proto_rawDescData
}

var file_v1_iam_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_v1_iam_proto_goTypes = []interface{}{
	(*IAMConfig)(nil),              // 0: v1.IAMConfig
	(*IssuerConfig)(nil),           // 1: v1.IssuerConfig
	(*IDMConfig)(nil),              // 2: v1.IDMConfig
	(*NamespaceGroupConfig)(nil),   // 3: v1.NamespaceGroupConfig
	(*ConnectorConfig)(nil),        // 4: v1.ConnectorConfig
	(*wrapperspb.StringValue)(nil), // 5: google.protobuf.StringValue
}
var file_v1_iam_proto_depIdxs = []int32{
	1, // 0: v1.IAMConfig.issuer_config:type_name -> v1.IssuerConfig
	2, // 1: v1.IAMConfig.idm_config:type_name -> v1.IDMConfig
	3, // 2: v1.IAMConfig.group_config:type_name -> v1.NamespaceGroupConfig
	4, // 3: v1.IDMConfig.connector_config:type_name -> v1.ConnectorConfig
	5, // 4: v1.ConnectorConfig.idm_groupname_template:type_name -> google.protobuf.StringValue
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_v1_iam_proto_init() }
func file_v1_iam_proto_init() {
	if File_v1_iam_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_v1_iam_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IAMConfig); i {
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
		file_v1_iam_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IssuerConfig); i {
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
		file_v1_iam_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IDMConfig); i {
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
		file_v1_iam_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NamespaceGroupConfig); i {
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
		file_v1_iam_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConnectorConfig); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_v1_iam_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_v1_iam_proto_goTypes,
		DependencyIndexes: file_v1_iam_proto_depIdxs,
		MessageInfos:      file_v1_iam_proto_msgTypes,
	}.Build()
	File_v1_iam_proto = out.File
	file_v1_iam_proto_rawDesc = nil
	file_v1_iam_proto_goTypes = nil
	file_v1_iam_proto_depIdxs = nil
}
