// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        (unknown)
// source: v1/quota.proto

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

// QuotaSet defines the types of possible Quotas
// might be specified by project or tenant
// whatever quota is reached first counts
// it always defines the max amount of this type
type QuotaSet struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// cluster the amount of clusters
	Cluster *Quota `protobuf:"bytes,1,opt,name=cluster,proto3" json:"cluster,omitempty"`
	// machine the amount of machines
	Machine *Quota `protobuf:"bytes,2,opt,name=machine,proto3" json:"machine,omitempty"`
	// ip the amount of aquired ip´s
	Ip *Quota `protobuf:"bytes,3,opt,name=ip,proto3" json:"ip,omitempty"`
	// project the amount of projects of a tenant
	Project *Quota `protobuf:"bytes,4,opt,name=project,proto3" json:"project,omitempty"`
}

func (x *QuotaSet) Reset() {
	*x = QuotaSet{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_quota_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QuotaSet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QuotaSet) ProtoMessage() {}

func (x *QuotaSet) ProtoReflect() protoreflect.Message {
	mi := &file_v1_quota_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QuotaSet.ProtoReflect.Descriptor instead.
func (*QuotaSet) Descriptor() ([]byte, []int) {
	return file_v1_quota_proto_rawDescGZIP(), []int{0}
}

func (x *QuotaSet) GetCluster() *Quota {
	if x != nil {
		return x.Cluster
	}
	return nil
}

func (x *QuotaSet) GetMachine() *Quota {
	if x != nil {
		return x.Machine
	}
	return nil
}

func (x *QuotaSet) GetIp() *Quota {
	if x != nil {
		return x.Ip
	}
	return nil
}

func (x *QuotaSet) GetProject() *Quota {
	if x != nil {
		return x.Project
	}
	return nil
}

// Quota is the actual maximum amount
type Quota struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// quota is the maximum amount for the current QuotaSet, can be nil
	Quota *wrapperspb.Int32Value `protobuf:"bytes,1,opt,name=quota,proto3" json:"quota,omitempty"`
}

func (x *Quota) Reset() {
	*x = Quota{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_quota_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Quota) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Quota) ProtoMessage() {}

func (x *Quota) ProtoReflect() protoreflect.Message {
	mi := &file_v1_quota_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Quota.ProtoReflect.Descriptor instead.
func (*Quota) Descriptor() ([]byte, []int) {
	return file_v1_quota_proto_rawDescGZIP(), []int{1}
}

func (x *Quota) GetQuota() *wrapperspb.Int32Value {
	if x != nil {
		return x.Quota
	}
	return nil
}

var File_v1_quota_proto protoreflect.FileDescriptor

var file_v1_quota_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x76, 0x31, 0x2f, 0x71, 0x75, 0x6f, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x02, 0x76, 0x31, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x94, 0x01, 0x0a, 0x08, 0x51, 0x75, 0x6f, 0x74, 0x61, 0x53, 0x65,
	0x74, 0x12, 0x23, 0x0a, 0x07, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x09, 0x2e, 0x76, 0x31, 0x2e, 0x51, 0x75, 0x6f, 0x74, 0x61, 0x52, 0x07, 0x63,
	0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x12, 0x23, 0x0a, 0x07, 0x6d, 0x61, 0x63, 0x68, 0x69, 0x6e,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x76, 0x31, 0x2e, 0x51, 0x75, 0x6f,
	0x74, 0x61, 0x52, 0x07, 0x6d, 0x61, 0x63, 0x68, 0x69, 0x6e, 0x65, 0x12, 0x19, 0x0a, 0x02, 0x69,
	0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x76, 0x31, 0x2e, 0x51, 0x75, 0x6f,
	0x74, 0x61, 0x52, 0x02, 0x69, 0x70, 0x12, 0x23, 0x0a, 0x07, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x76, 0x31, 0x2e, 0x51, 0x75, 0x6f,
	0x74, 0x61, 0x52, 0x07, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x22, 0x3a, 0x0a, 0x05, 0x51,
	0x75, 0x6f, 0x74, 0x61, 0x12, 0x31, 0x0a, 0x05, 0x71, 0x75, 0x6f, 0x74, 0x61, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x49, 0x6e, 0x74, 0x33, 0x32, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x52, 0x05, 0x71, 0x75, 0x6f, 0x74, 0x61, 0x42, 0x66, 0x0a, 0x06, 0x63, 0x6f, 0x6d, 0x2e, 0x76,
	0x31, 0x42, 0x0a, 0x51, 0x75, 0x6f, 0x74, 0x61, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a,
	0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x65, 0x74, 0x61,
	0x6c, 0x2d, 0x73, 0x74, 0x61, 0x63, 0x6b, 0x2f, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x64, 0x61,
	0x74, 0x61, 0x2d, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x56, 0x58, 0x58, 0xaa,
	0x02, 0x02, 0x56, 0x31, 0xca, 0x02, 0x02, 0x56, 0x31, 0xe2, 0x02, 0x0e, 0x56, 0x31, 0x5c, 0x47,
	0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x02, 0x56, 0x31, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_v1_quota_proto_rawDescOnce sync.Once
	file_v1_quota_proto_rawDescData = file_v1_quota_proto_rawDesc
)

func file_v1_quota_proto_rawDescGZIP() []byte {
	file_v1_quota_proto_rawDescOnce.Do(func() {
		file_v1_quota_proto_rawDescData = protoimpl.X.CompressGZIP(file_v1_quota_proto_rawDescData)
	})
	return file_v1_quota_proto_rawDescData
}

var file_v1_quota_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_v1_quota_proto_goTypes = []interface{}{
	(*QuotaSet)(nil),              // 0: v1.QuotaSet
	(*Quota)(nil),                 // 1: v1.Quota
	(*wrapperspb.Int32Value)(nil), // 2: google.protobuf.Int32Value
}
var file_v1_quota_proto_depIdxs = []int32{
	1, // 0: v1.QuotaSet.cluster:type_name -> v1.Quota
	1, // 1: v1.QuotaSet.machine:type_name -> v1.Quota
	1, // 2: v1.QuotaSet.ip:type_name -> v1.Quota
	1, // 3: v1.QuotaSet.project:type_name -> v1.Quota
	2, // 4: v1.Quota.quota:type_name -> google.protobuf.Int32Value
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_v1_quota_proto_init() }
func file_v1_quota_proto_init() {
	if File_v1_quota_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_v1_quota_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QuotaSet); i {
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
		file_v1_quota_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Quota); i {
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
			RawDescriptor: file_v1_quota_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_v1_quota_proto_goTypes,
		DependencyIndexes: file_v1_quota_proto_depIdxs,
		MessageInfos:      file_v1_quota_proto_msgTypes,
	}.Build()
	File_v1_quota_proto = out.File
	file_v1_quota_proto_rawDesc = nil
	file_v1_quota_proto_goTypes = nil
	file_v1_quota_proto_depIdxs = nil
}
