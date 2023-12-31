// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.4
// source: file-api.proto

package v1

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type IsAliveRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *IsAliveRequest) Reset() {
	*x = IsAliveRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_file_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IsAliveRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IsAliveRequest) ProtoMessage() {}

func (x *IsAliveRequest) ProtoReflect() protoreflect.Message {
	mi := &file_file_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IsAliveRequest.ProtoReflect.Descriptor instead.
func (*IsAliveRequest) Descriptor() ([]byte, []int) {
	return file_file_api_proto_rawDescGZIP(), []int{0}
}

type IsAliveResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *IsAliveResponse) Reset() {
	*x = IsAliveResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_file_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IsAliveResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IsAliveResponse) ProtoMessage() {}

func (x *IsAliveResponse) ProtoReflect() protoreflect.Message {
	mi := &file_file_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IsAliveResponse.ProtoReflect.Descriptor instead.
func (*IsAliveResponse) Descriptor() ([]byte, []int) {
	return file_file_api_proto_rawDescGZIP(), []int{1}
}

type DeliveryGetFile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UUID string `protobuf:"bytes,101,opt,name=UUID,proto3" json:"UUID,omitempty"`
}

func (x *DeliveryGetFile) Reset() {
	*x = DeliveryGetFile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_file_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeliveryGetFile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeliveryGetFile) ProtoMessage() {}

func (x *DeliveryGetFile) ProtoReflect() protoreflect.Message {
	mi := &file_file_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeliveryGetFile.ProtoReflect.Descriptor instead.
func (*DeliveryGetFile) Descriptor() ([]byte, []int) {
	return file_file_api_proto_rawDescGZIP(), []int{2}
}

func (x *DeliveryGetFile) GetUUID() string {
	if x != nil {
		return x.UUID
	}
	return ""
}

type DeliveryPutFile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	File *File `protobuf:"bytes,21,opt,name=File,proto3" json:"File,omitempty"`
}

func (x *DeliveryPutFile) Reset() {
	*x = DeliveryPutFile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_file_api_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeliveryPutFile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeliveryPutFile) ProtoMessage() {}

func (x *DeliveryPutFile) ProtoReflect() protoreflect.Message {
	mi := &file_file_api_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeliveryPutFile.ProtoReflect.Descriptor instead.
func (*DeliveryPutFile) Descriptor() ([]byte, []int) {
	return file_file_api_proto_rawDescGZIP(), []int{3}
}

func (x *DeliveryPutFile) GetFile() *File {
	if x != nil {
		return x.File
	}
	return nil
}

type File struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UUID         string `protobuf:"bytes,21,opt,name=UUID,proto3" json:"UUID,omitempty"`
	BusinessType string `protobuf:"bytes,22,opt,name=BusinessType,proto3" json:"BusinessType,omitempty"`
	FileName     string `protobuf:"bytes,30,opt,name=FileName,proto3" json:"FileName,omitempty"`
	Hash         string `protobuf:"bytes,40,opt,name=Hash,proto3" json:"Hash,omitempty"`
	Size         uint64 `protobuf:"varint,41,opt,name=Size,proto3" json:"Size,omitempty"`
	Content      []byte `protobuf:"bytes,42,opt,name=Content,proto3,oneof" json:"Content,omitempty"`
}

func (x *File) Reset() {
	*x = File{}
	if protoimpl.UnsafeEnabled {
		mi := &file_file_api_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *File) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*File) ProtoMessage() {}

func (x *File) ProtoReflect() protoreflect.Message {
	mi := &file_file_api_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use File.ProtoReflect.Descriptor instead.
func (*File) Descriptor() ([]byte, []int) {
	return file_file_api_proto_rawDescGZIP(), []int{4}
}

func (x *File) GetUUID() string {
	if x != nil {
		return x.UUID
	}
	return ""
}

func (x *File) GetBusinessType() string {
	if x != nil {
		return x.BusinessType
	}
	return ""
}

func (x *File) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *File) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *File) GetSize() uint64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *File) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

var File_file_api_proto protoreflect.FileDescriptor

var file_file_api_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x66, 0x69, 0x6c, 0x65, 0x2d, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e,
	0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x10,
	0x0a, 0x0e, 0x49, 0x73, 0x41, 0x6c, 0x69, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x22, 0x11, 0x0a, 0x0f, 0x49, 0x73, 0x41, 0x6c, 0x69, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x2f, 0x0a, 0x0f, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x79, 0x47,
	0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x55, 0x55, 0x49, 0x44, 0x18, 0x65,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x55, 0x55, 0x49, 0x44, 0x4a, 0x04, 0x08, 0x01, 0x10, 0x65,
	0x52, 0x02, 0x49, 0x64, 0x22, 0x32, 0x0a, 0x0f, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x79,
	0x50, 0x75, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x19, 0x0a, 0x04, 0x46, 0x69, 0x6c, 0x65, 0x18,
	0x15, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x04, 0x46, 0x69,
	0x6c, 0x65, 0x4a, 0x04, 0x08, 0x01, 0x10, 0x15, 0x22, 0xb9, 0x01, 0x0a, 0x04, 0x46, 0x69, 0x6c,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x55, 0x55, 0x49, 0x44, 0x18, 0x15, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x55, 0x55, 0x49, 0x44, 0x12, 0x22, 0x0a, 0x0c, 0x42, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73,
	0x73, 0x54, 0x79, 0x70, 0x65, 0x18, 0x16, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x42, 0x75, 0x73,
	0x69, 0x6e, 0x65, 0x73, 0x73, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x46, 0x69, 0x6c,
	0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x1e, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x46, 0x69, 0x6c,
	0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x48, 0x61, 0x73, 0x68, 0x18, 0x28, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x48, 0x61, 0x73, 0x68, 0x12, 0x12, 0x0a, 0x04, 0x53, 0x69, 0x7a,
	0x65, 0x18, 0x29, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x1d, 0x0a,
	0x07, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x2a, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x00,
	0x52, 0x07, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x88, 0x01, 0x01, 0x42, 0x0a, 0x0a, 0x08,
	0x5f, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x4a, 0x04, 0x08, 0x01, 0x10, 0x15, 0x4a, 0x04,
	0x08, 0x1f, 0x10, 0x28, 0x32, 0x90, 0x01, 0x0a, 0x0e, 0x46, 0x69, 0x6c, 0x65, 0x41, 0x70, 0x69,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3b, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x46, 0x69,
	0x6c, 0x65, 0x12, 0x10, 0x2e, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x79, 0x47, 0x65, 0x74,
	0x46, 0x69, 0x6c, 0x65, 0x1a, 0x05, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x22, 0x17, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x11, 0x12, 0x0f, 0x2f, 0x67, 0x65, 0x74, 0x66, 0x69, 0x6c, 0x65, 0x2f, 0x7b, 0x55,
	0x55, 0x49, 0x44, 0x7d, 0x12, 0x41, 0x0a, 0x07, 0x49, 0x73, 0x41, 0x6c, 0x69, 0x76, 0x65, 0x12,
	0x0f, 0x2e, 0x49, 0x73, 0x41, 0x6c, 0x69, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x10, 0x2e, 0x49, 0x73, 0x41, 0x6c, 0x69, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x13, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0d, 0x12, 0x0b, 0x2f, 0x72, 0x65, 0x61,
	0x64, 0x79, 0x2f, 0x6c, 0x69, 0x76, 0x65, 0x42, 0x08, 0x5a, 0x06, 0x61, 0x70, 0x69, 0x2f, 0x76,
	0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_file_api_proto_rawDescOnce sync.Once
	file_file_api_proto_rawDescData = file_file_api_proto_rawDesc
)

func file_file_api_proto_rawDescGZIP() []byte {
	file_file_api_proto_rawDescOnce.Do(func() {
		file_file_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_file_api_proto_rawDescData)
	})
	return file_file_api_proto_rawDescData
}

var file_file_api_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_file_api_proto_goTypes = []interface{}{
	(*IsAliveRequest)(nil),  // 0: IsAliveRequest
	(*IsAliveResponse)(nil), // 1: IsAliveResponse
	(*DeliveryGetFile)(nil), // 2: DeliveryGetFile
	(*DeliveryPutFile)(nil), // 3: DeliveryPutFile
	(*File)(nil),            // 4: File
}
var file_file_api_proto_depIdxs = []int32{
	4, // 0: DeliveryPutFile.File:type_name -> File
	2, // 1: FileApiService.GetFile:input_type -> DeliveryGetFile
	0, // 2: FileApiService.IsAlive:input_type -> IsAliveRequest
	4, // 3: FileApiService.GetFile:output_type -> File
	1, // 4: FileApiService.IsAlive:output_type -> IsAliveResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_file_api_proto_init() }
func file_file_api_proto_init() {
	if File_file_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_file_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IsAliveRequest); i {
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
		file_file_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IsAliveResponse); i {
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
		file_file_api_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeliveryGetFile); i {
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
		file_file_api_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeliveryPutFile); i {
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
		file_file_api_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*File); i {
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
	file_file_api_proto_msgTypes[4].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_file_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_file_api_proto_goTypes,
		DependencyIndexes: file_file_api_proto_depIdxs,
		MessageInfos:      file_file_api_proto_msgTypes,
	}.Build()
	File_file_api_proto = out.File
	file_file_api_proto_rawDesc = nil
	file_file_api_proto_goTypes = nil
	file_file_api_proto_depIdxs = nil
}
