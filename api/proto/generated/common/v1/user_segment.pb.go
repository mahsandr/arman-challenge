// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        (unknown)
// source: common/v1/user_segment.proto

package commonv1

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

type UserSegment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId    string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Segment   string `protobuf:"bytes,2,opt,name=segment,proto3" json:"segment,omitempty"`
	Timestamp int64  `protobuf:"varint,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *UserSegment) Reset() {
	*x = UserSegment{}
	mi := &file_common_v1_user_segment_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserSegment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserSegment) ProtoMessage() {}

func (x *UserSegment) ProtoReflect() protoreflect.Message {
	mi := &file_common_v1_user_segment_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserSegment.ProtoReflect.Descriptor instead.
func (*UserSegment) Descriptor() ([]byte, []int) {
	return file_common_v1_user_segment_proto_rawDescGZIP(), []int{0}
}

func (x *UserSegment) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *UserSegment) GetSegment() string {
	if x != nil {
		return x.Segment
	}
	return ""
}

func (x *UserSegment) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

var File_common_v1_user_segment_proto protoreflect.FileDescriptor

var file_common_v1_user_segment_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x75, 0x73, 0x65, 0x72,
	0x5f, 0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x22, 0x5e, 0x0a, 0x0b, 0x55, 0x73, 0x65,
	0x72, 0x53, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49,
	0x64, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0xaf, 0x01, 0x0a, 0x0d, 0x63, 0x6f,
	0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x42, 0x10, 0x55, 0x73, 0x65,
	0x72, 0x53, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a,
	0x47, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x61, 0x68, 0x73,
	0x61, 0x6e, 0x64, 0x72, 0x2f, 0x61, 0x72, 0x6d, 0x61, 0x6e, 0x2d, 0x63, 0x68, 0x61, 0x6c, 0x6c,
	0x65, 0x6e, 0x67, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x67, 0x65, 0x6e, 0x65,
	0x72, 0x61, 0x74, 0x65, 0x64, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x3b,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x43, 0x58, 0x58, 0xaa, 0x02,
	0x09, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x09, 0x43, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x15, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5c,
	0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02,
	0x0a, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_common_v1_user_segment_proto_rawDescOnce sync.Once
	file_common_v1_user_segment_proto_rawDescData = file_common_v1_user_segment_proto_rawDesc
)

func file_common_v1_user_segment_proto_rawDescGZIP() []byte {
	file_common_v1_user_segment_proto_rawDescOnce.Do(func() {
		file_common_v1_user_segment_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_v1_user_segment_proto_rawDescData)
	})
	return file_common_v1_user_segment_proto_rawDescData
}

var file_common_v1_user_segment_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_common_v1_user_segment_proto_goTypes = []any{
	(*UserSegment)(nil), // 0: common.v1.UserSegment
}
var file_common_v1_user_segment_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_common_v1_user_segment_proto_init() }
func file_common_v1_user_segment_proto_init() {
	if File_common_v1_user_segment_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_common_v1_user_segment_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_common_v1_user_segment_proto_goTypes,
		DependencyIndexes: file_common_v1_user_segment_proto_depIdxs,
		MessageInfos:      file_common_v1_user_segment_proto_msgTypes,
	}.Build()
	File_common_v1_user_segment_proto = out.File
	file_common_v1_user_segment_proto_rawDesc = nil
	file_common_v1_user_segment_proto_goTypes = nil
	file_common_v1_user_segment_proto_depIdxs = nil
}
