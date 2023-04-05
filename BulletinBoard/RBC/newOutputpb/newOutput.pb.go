// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.6.1
// source: newOutput.proto

package newOutputpb

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

// The request message containing the user's name.
type NewOutput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type   int64  `protobuf:"varint,1,opt,name=type,proto3" json:"type,omitempty"`
	Round  string `protobuf:"bytes,2,opt,name=round,proto3" json:"round,omitempty"`
	View   string `protobuf:"bytes,3,opt,name=view,proto3" json:"view,omitempty"`
	Sender string `protobuf:"bytes,4,opt,name=sender,proto3" json:"sender,omitempty"`
	Output string `protobuf:"bytes,5,opt,name=output,proto3" json:"output,omitempty"`
	Sig    string `protobuf:"bytes,6,opt,name=sig,proto3" json:"sig,omitempty"`
}

func (x *NewOutput) Reset() {
	*x = NewOutput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_newOutput_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NewOutput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NewOutput) ProtoMessage() {}

func (x *NewOutput) ProtoReflect() protoreflect.Message {
	mi := &file_newOutput_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NewOutput.ProtoReflect.Descriptor instead.
func (*NewOutput) Descriptor() ([]byte, []int) {
	return file_newOutput_proto_rawDescGZIP(), []int{0}
}

func (x *NewOutput) GetType() int64 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *NewOutput) GetRound() string {
	if x != nil {
		return x.Round
	}
	return ""
}

func (x *NewOutput) GetView() string {
	if x != nil {
		return x.View
	}
	return ""
}

func (x *NewOutput) GetSender() string {
	if x != nil {
		return x.Sender
	}
	return ""
}

func (x *NewOutput) GetOutput() string {
	if x != nil {
		return x.Output
	}
	return ""
}

func (x *NewOutput) GetSig() string {
	if x != nil {
		return x.Sig
	}
	return ""
}

type NewOutputResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *NewOutputResponse) Reset() {
	*x = NewOutputResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_newOutput_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NewOutputResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NewOutputResponse) ProtoMessage() {}

func (x *NewOutputResponse) ProtoReflect() protoreflect.Message {
	mi := &file_newOutput_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NewOutputResponse.ProtoReflect.Descriptor instead.
func (*NewOutputResponse) Descriptor() ([]byte, []int) {
	return file_newOutput_proto_rawDescGZIP(), []int{1}
}

var File_newOutput_proto protoreflect.FileDescriptor

var file_newOutput_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x6e, 0x65, 0x77, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x0b, 0x6e, 0x65, 0x77, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x70, 0x62, 0x22, 0x8b,
	0x01, 0x0a, 0x09, 0x4e, 0x65, 0x77, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x12, 0x12, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x76, 0x69, 0x65, 0x77, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x76, 0x69, 0x65, 0x77, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65,
	0x6e, 0x64, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x65, 0x6e, 0x64,
	0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x69,
	0x67, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x73, 0x69, 0x67, 0x22, 0x13, 0x0a, 0x11,
	0x4e, 0x65, 0x77, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x32, 0x5f, 0x0a, 0x0f, 0x4e, 0x65, 0x77, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x48, 0x61,
	0x6e, 0x64, 0x6c, 0x65, 0x12, 0x4c, 0x0a, 0x10, 0x4e, 0x65, 0x77, 0x4f, 0x75, 0x74, 0x70, 0x75,
	0x74, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x12, 0x16, 0x2e, 0x6e, 0x65, 0x77, 0x4f, 0x75,
	0x74, 0x70, 0x75, 0x74, 0x70, 0x62, 0x2e, 0x4e, 0x65, 0x77, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74,
	0x1a, 0x1e, 0x2e, 0x6e, 0x65, 0x77, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x70, 0x62, 0x2e, 0x4e,
	0x65, 0x77, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x42, 0x10, 0x5a, 0x0e, 0x2e, 0x2e, 0x2f, 0x6e, 0x65, 0x77, 0x4f, 0x75, 0x74, 0x70,
	0x75, 0x74, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_newOutput_proto_rawDescOnce sync.Once
	file_newOutput_proto_rawDescData = file_newOutput_proto_rawDesc
)

func file_newOutput_proto_rawDescGZIP() []byte {
	file_newOutput_proto_rawDescOnce.Do(func() {
		file_newOutput_proto_rawDescData = protoimpl.X.CompressGZIP(file_newOutput_proto_rawDescData)
	})
	return file_newOutput_proto_rawDescData
}

var file_newOutput_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_newOutput_proto_goTypes = []interface{}{
	(*NewOutput)(nil),         // 0: newOutputpb.NewOutput
	(*NewOutputResponse)(nil), // 1: newOutputpb.NewOutputResponse
}
var file_newOutput_proto_depIdxs = []int32{
	0, // 0: newOutputpb.NewOutputHandle.NewOutputReceive:input_type -> newOutputpb.NewOutput
	1, // 1: newOutputpb.NewOutputHandle.NewOutputReceive:output_type -> newOutputpb.NewOutputResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_newOutput_proto_init() }
func file_newOutput_proto_init() {
	if File_newOutput_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_newOutput_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NewOutput); i {
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
		file_newOutput_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NewOutputResponse); i {
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
			RawDescriptor: file_newOutput_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_newOutput_proto_goTypes,
		DependencyIndexes: file_newOutput_proto_depIdxs,
		MessageInfos:      file_newOutput_proto_msgTypes,
	}.Build()
	File_newOutput_proto = out.File
	file_newOutput_proto_rawDesc = nil
	file_newOutput_proto_goTypes = nil
	file_newOutput_proto_depIdxs = nil
}