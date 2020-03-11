// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.20.0
// 	protoc        v3.11.4
// source: v1/meta.proto

package v1

import (
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Meta struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string               `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Kind        string               `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"`
	Apiversion  string               `protobuf:"bytes,3,opt,name=apiversion,proto3" json:"apiversion,omitempty"`
	Version     int64                `protobuf:"varint,4,opt,name=version,proto3" json:"version,omitempty"`
	CreatedTime *timestamp.Timestamp `protobuf:"bytes,5,opt,name=created_time,json=createdTime,proto3" json:"created_time,omitempty"`
	UpdatedTime *timestamp.Timestamp `protobuf:"bytes,6,opt,name=updated_time,json=updatedTime,proto3" json:"updated_time,omitempty"`
}

func (x *Meta) Reset() {
	*x = Meta{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1_meta_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Meta) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Meta) ProtoMessage() {}

func (x *Meta) ProtoReflect() protoreflect.Message {
	mi := &file_v1_meta_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Meta.ProtoReflect.Descriptor instead.
func (*Meta) Descriptor() ([]byte, []int) {
	return file_v1_meta_proto_rawDescGZIP(), []int{0}
}

func (x *Meta) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Meta) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *Meta) GetApiversion() string {
	if x != nil {
		return x.Apiversion
	}
	return ""
}

func (x *Meta) GetVersion() int64 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *Meta) GetCreatedTime() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedTime
	}
	return nil
}

func (x *Meta) GetUpdatedTime() *timestamp.Timestamp {
	if x != nil {
		return x.UpdatedTime
	}
	return nil
}

var File_v1_meta_proto protoreflect.FileDescriptor

var file_v1_meta_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x76, 0x31, 0x2f, 0x6d, 0x65, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x02, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe2, 0x01, 0x0a, 0x04, 0x4d, 0x65, 0x74, 0x61, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x6b, 0x69, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x69, 0x6e,
	0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x61, 0x70, 0x69, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x61, 0x70, 0x69, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x3d, 0x0a, 0x0c, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0b, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x3d, 0x0a, 0x0c, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0b, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x42, 0x2e, 0x5a, 0x2c, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x65, 0x74, 0x61, 0x6c, 0x2d, 0x73, 0x74,
	0x61, 0x63, 0x6b, 0x2f, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x64, 0x61, 0x74, 0x61, 0x2d, 0x61,
	0x70, 0x69, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_v1_meta_proto_rawDescOnce sync.Once
	file_v1_meta_proto_rawDescData = file_v1_meta_proto_rawDesc
)

func file_v1_meta_proto_rawDescGZIP() []byte {
	file_v1_meta_proto_rawDescOnce.Do(func() {
		file_v1_meta_proto_rawDescData = protoimpl.X.CompressGZIP(file_v1_meta_proto_rawDescData)
	})
	return file_v1_meta_proto_rawDescData
}

var file_v1_meta_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_v1_meta_proto_goTypes = []interface{}{
	(*Meta)(nil),                // 0: v1.Meta
	(*timestamp.Timestamp)(nil), // 1: google.protobuf.Timestamp
}
var file_v1_meta_proto_depIdxs = []int32{
	1, // 0: v1.Meta.created_time:type_name -> google.protobuf.Timestamp
	1, // 1: v1.Meta.updated_time:type_name -> google.protobuf.Timestamp
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_v1_meta_proto_init() }
func file_v1_meta_proto_init() {
	if File_v1_meta_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_v1_meta_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Meta); i {
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
			RawDescriptor: file_v1_meta_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_v1_meta_proto_goTypes,
		DependencyIndexes: file_v1_meta_proto_depIdxs,
		MessageInfos:      file_v1_meta_proto_msgTypes,
	}.Build()
	File_v1_meta_proto = out.File
	file_v1_meta_proto_rawDesc = nil
	file_v1_meta_proto_goTypes = nil
	file_v1_meta_proto_depIdxs = nil
}
