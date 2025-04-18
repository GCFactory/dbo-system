// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v3.21.12
// source: notification/notification.proto

package notification_api

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type AdditionalInfo struct {
	state             protoimpl.MessageState `protogen:"open.v1"`
	UserId            string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Email             *string                `protobuf:"bytes,2,opt,name=email,proto3,oneof" json:"email,omitempty"`
	EmailNotification *bool                  `protobuf:"varint,3,opt,name=email_notification,json=emailNotification,proto3,oneof" json:"email_notification,omitempty"`
	unknownFields     protoimpl.UnknownFields
	sizeCache         protoimpl.SizeCache
}

func (x *AdditionalInfo) Reset() {
	*x = AdditionalInfo{}
	mi := &file_notification_notification_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AdditionalInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AdditionalInfo) ProtoMessage() {}

func (x *AdditionalInfo) ProtoReflect() protoreflect.Message {
	mi := &file_notification_notification_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AdditionalInfo.ProtoReflect.Descriptor instead.
func (*AdditionalInfo) Descriptor() ([]byte, []int) {
	return file_notification_notification_proto_rawDescGZIP(), []int{0}
}

func (x *AdditionalInfo) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *AdditionalInfo) GetEmail() string {
	if x != nil && x.Email != nil {
		return *x.Email
	}
	return ""
}

func (x *AdditionalInfo) GetEmailNotification() bool {
	if x != nil && x.EmailNotification != nil {
		return *x.EmailNotification
	}
	return false
}

// Данные event-а
type EventData struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	SagaUuid       string                 `protobuf:"bytes,1,opt,name=saga_uuid,json=sagaUuid,proto3" json:"saga_uuid,omitempty"`                //  UUID sag-и
	EventUuid      string                 `protobuf:"bytes,2,opt,name=event_uuid,json=eventUuid,proto3" json:"event_uuid,omitempty"`             //  UUID evnet-a
	OperationName  string                 `protobuf:"bytes,3,opt,name=operation_name,json=operationName,proto3" json:"operation_name,omitempty"` // Тип операции
	AdditionalInfo *AdditionalInfo        `protobuf:"bytes,4,opt,name=additional_info,json=additionalInfo,proto3" json:"additional_info,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *EventData) Reset() {
	*x = EventData{}
	mi := &file_notification_notification_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EventData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventData) ProtoMessage() {}

func (x *EventData) ProtoReflect() protoreflect.Message {
	mi := &file_notification_notification_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventData.ProtoReflect.Descriptor instead.
func (*EventData) Descriptor() ([]byte, []int) {
	return file_notification_notification_proto_rawDescGZIP(), []int{1}
}

func (x *EventData) GetSagaUuid() string {
	if x != nil {
		return x.SagaUuid
	}
	return ""
}

func (x *EventData) GetEventUuid() string {
	if x != nil {
		return x.EventUuid
	}
	return ""
}

func (x *EventData) GetOperationName() string {
	if x != nil {
		return x.OperationName
	}
	return ""
}

func (x *EventData) GetAdditionalInfo() *AdditionalInfo {
	if x != nil {
		return x.AdditionalInfo
	}
	return nil
}

// Результат event-а
type EventSuccess struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	SagaUuid      string                 `protobuf:"bytes,1,opt,name=saga_uuid,json=sagaUuid,proto3" json:"saga_uuid,omitempty"`                //  UUID sag-и
	EventUuid     string                 `protobuf:"bytes,2,opt,name=event_uuid,json=eventUuid,proto3" json:"event_uuid,omitempty"`             //  UUID evnet-a
	OperationName string                 `protobuf:"bytes,3,opt,name=operation_name,json=operationName,proto3" json:"operation_name,omitempty"` //  Тип выполняемой операции
	Info          string                 `protobuf:"bytes,5,opt,name=info,proto3" json:"info,omitempty"`                                        //  Дополнительная информация по event-у
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EventSuccess) Reset() {
	*x = EventSuccess{}
	mi := &file_notification_notification_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EventSuccess) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventSuccess) ProtoMessage() {}

func (x *EventSuccess) ProtoReflect() protoreflect.Message {
	mi := &file_notification_notification_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventSuccess.ProtoReflect.Descriptor instead.
func (*EventSuccess) Descriptor() ([]byte, []int) {
	return file_notification_notification_proto_rawDescGZIP(), []int{2}
}

func (x *EventSuccess) GetSagaUuid() string {
	if x != nil {
		return x.SagaUuid
	}
	return ""
}

func (x *EventSuccess) GetEventUuid() string {
	if x != nil {
		return x.EventUuid
	}
	return ""
}

func (x *EventSuccess) GetOperationName() string {
	if x != nil {
		return x.OperationName
	}
	return ""
}

func (x *EventSuccess) GetInfo() string {
	if x != nil {
		return x.Info
	}
	return ""
}

type EventError struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	SagaUuid      string                 `protobuf:"bytes,1,opt,name=saga_uuid,json=sagaUuid,proto3" json:"saga_uuid,omitempty"`                //  UUID sag-и
	EventUuid     string                 `protobuf:"bytes,2,opt,name=event_uuid,json=eventUuid,proto3" json:"event_uuid,omitempty"`             //  UUID evnet-a
	OperationName string                 `protobuf:"bytes,3,opt,name=operation_name,json=operationName,proto3" json:"operation_name,omitempty"` //  Тип выполняемой операции
	Status        uint32                 `protobuf:"varint,4,opt,name=status,proto3" json:"status,omitempty"`                                   //  Event error status
	Info          string                 `protobuf:"bytes,5,opt,name=info,proto3" json:"info,omitempty"`                                        //  Дополнительная информация по event-у
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EventError) Reset() {
	*x = EventError{}
	mi := &file_notification_notification_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EventError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventError) ProtoMessage() {}

func (x *EventError) ProtoReflect() protoreflect.Message {
	mi := &file_notification_notification_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventError.ProtoReflect.Descriptor instead.
func (*EventError) Descriptor() ([]byte, []int) {
	return file_notification_notification_proto_rawDescGZIP(), []int{3}
}

func (x *EventError) GetSagaUuid() string {
	if x != nil {
		return x.SagaUuid
	}
	return ""
}

func (x *EventError) GetEventUuid() string {
	if x != nil {
		return x.EventUuid
	}
	return ""
}

func (x *EventError) GetOperationName() string {
	if x != nil {
		return x.OperationName
	}
	return ""
}

func (x *EventError) GetStatus() uint32 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *EventError) GetInfo() string {
	if x != nil {
		return x.Info
	}
	return ""
}

var File_notification_notification_proto protoreflect.FileDescriptor

var file_notification_notification_proto_rawDesc = string([]byte{
	0x0a, 0x1f, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x6e,
	0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x0c, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22,
	0x99, 0x01, 0x0a, 0x0e, 0x41, 0x64, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x49, 0x6e,
	0x66, 0x6f, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x05, 0x65,
	0x6d, 0x61, 0x69, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x05, 0x65, 0x6d,
	0x61, 0x69, 0x6c, 0x88, 0x01, 0x01, 0x12, 0x32, 0x0a, 0x12, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x5f,
	0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x08, 0x48, 0x01, 0x52, 0x11, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x4e, 0x6f, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x65,
	0x6d, 0x61, 0x69, 0x6c, 0x42, 0x15, 0x0a, 0x13, 0x5f, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x5f, 0x6e,
	0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0xb5, 0x01, 0x0a, 0x09,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x61, 0x67,
	0x61, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x61,
	0x67, 0x61, 0x55, 0x75, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f,
	0x75, 0x75, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x55, 0x75, 0x69, 0x64, 0x12, 0x25, 0x0a, 0x0e, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6f,
	0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x45, 0x0a, 0x0f,
	0x61, 0x64, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x41, 0x64, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x49,
	0x6e, 0x66, 0x6f, 0x52, 0x0e, 0x61, 0x64, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x49,
	0x6e, 0x66, 0x6f, 0x22, 0x85, 0x01, 0x0a, 0x0c, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x53, 0x75, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x61, 0x67, 0x61, 0x5f, 0x75, 0x75, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x61, 0x67, 0x61, 0x55, 0x75, 0x69,
	0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x55, 0x75, 0x69, 0x64,
	0x12, 0x25, 0x0a, 0x0e, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x22, 0x9b, 0x01, 0x0a, 0x0a,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x61,
	0x67, 0x61, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73,
	0x61, 0x67, 0x61, 0x55, 0x75, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x55, 0x75, 0x69, 0x64, 0x12, 0x25, 0x0a, 0x0e, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d,
	0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x42, 0x1e, 0x5a, 0x1c, 0x2e, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
})

var (
	file_notification_notification_proto_rawDescOnce sync.Once
	file_notification_notification_proto_rawDescData []byte
)

func file_notification_notification_proto_rawDescGZIP() []byte {
	file_notification_notification_proto_rawDescOnce.Do(func() {
		file_notification_notification_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_notification_notification_proto_rawDesc), len(file_notification_notification_proto_rawDesc)))
	})
	return file_notification_notification_proto_rawDescData
}

var file_notification_notification_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_notification_notification_proto_goTypes = []any{
	(*AdditionalInfo)(nil), // 0: notification.AdditionalInfo
	(*EventData)(nil),      // 1: notification.EventData
	(*EventSuccess)(nil),   // 2: notification.EventSuccess
	(*EventError)(nil),     // 3: notification.EventError
}
var file_notification_notification_proto_depIdxs = []int32{
	0, // 0: notification.EventData.additional_info:type_name -> notification.AdditionalInfo
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_notification_notification_proto_init() }
func file_notification_notification_proto_init() {
	if File_notification_notification_proto != nil {
		return
	}
	file_notification_notification_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_notification_notification_proto_rawDesc), len(file_notification_notification_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_notification_notification_proto_goTypes,
		DependencyIndexes: file_notification_notification_proto_depIdxs,
		MessageInfos:      file_notification_notification_proto_msgTypes,
	}.Build()
	File_notification_notification_proto = out.File
	file_notification_notification_proto_goTypes = nil
	file_notification_notification_proto_depIdxs = nil
}
