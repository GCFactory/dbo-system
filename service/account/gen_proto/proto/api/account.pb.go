// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.12.4
// source: account/account.proto

package api

import (
	platform "github.com/GCFactory/dbo-system/service/account/gen_proto/proto/platform"
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

// Дополнительные сведения для операции
type OperationDetails struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccUuid        string  `protobuf:"bytes,1,opt,name=acc_uuid,json=accUuid,proto3" json:"acc_uuid,omitempty"`                        //  UUID счёта
	AdditionalData float32 `protobuf:"fixed32,2,opt,name=additional_data,json=additionalData,proto3" json:"additional_data,omitempty"` //  Дополнительные данные
}

func (x *OperationDetails) Reset() {
	*x = OperationDetails{}
	if protoimpl.UnsafeEnabled {
		mi := &file_account_account_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OperationDetails) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OperationDetails) ProtoMessage() {}

func (x *OperationDetails) ProtoReflect() protoreflect.Message {
	mi := &file_account_account_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OperationDetails.ProtoReflect.Descriptor instead.
func (*OperationDetails) Descriptor() ([]byte, []int) {
	return file_account_account_proto_rawDescGZIP(), []int{0}
}

func (x *OperationDetails) GetAccUuid() string {
	if x != nil {
		return x.AccUuid
	}
	return ""
}

func (x *OperationDetails) GetAdditionalData() float32 {
	if x != nil {
		return x.AdditionalData
	}
	return 0
}

// Данные event-а
type EventData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SagaUuid      string `protobuf:"bytes,1,opt,name=saga_uuid,json=sagaUuid,proto3" json:"saga_uuid,omitempty"`                //  UUID sag-и
	OperationName string `protobuf:"bytes,2,opt,name=operation_name,json=operationName,proto3" json:"operation_name,omitempty"` // Тип операции
	// Types that are assignable to Data:
	//
	//	*EventData_AccountData
	//	*EventData_AdditionalInfo
	Data isEventData_Data `protobuf_oneof:"data"`
}

func (x *EventData) Reset() {
	*x = EventData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_account_account_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventData) ProtoMessage() {}

func (x *EventData) ProtoReflect() protoreflect.Message {
	mi := &file_account_account_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
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
	return file_account_account_proto_rawDescGZIP(), []int{1}
}

func (x *EventData) GetSagaUuid() string {
	if x != nil {
		return x.SagaUuid
	}
	return ""
}

func (x *EventData) GetOperationName() string {
	if x != nil {
		return x.OperationName
	}
	return ""
}

func (m *EventData) GetData() isEventData_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (x *EventData) GetAccountData() *platform.AccountDetails {
	if x, ok := x.GetData().(*EventData_AccountData); ok {
		return x.AccountData
	}
	return nil
}

func (x *EventData) GetAdditionalInfo() *OperationDetails {
	if x, ok := x.GetData().(*EventData_AdditionalInfo); ok {
		return x.AdditionalInfo
	}
	return nil
}

type isEventData_Data interface {
	isEventData_Data()
}

type EventData_AccountData struct {
	AccountData *platform.AccountDetails `protobuf:"bytes,3,opt,name=account_data,json=accountData,proto3,oneof"` //  Реквизиты счёта
}

type EventData_AdditionalInfo struct {
	AdditionalInfo *OperationDetails `protobuf:"bytes,4,opt,name=additional_info,json=additionalInfo,proto3,oneof"` //  Дополнительная информация по операции
}

func (*EventData_AccountData) isEventData_Data() {}

func (*EventData_AdditionalInfo) isEventData_Data() {}

// Результат event-а
type EventStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SagaUuid      string `protobuf:"bytes,1,opt,name=saga_uuid,json=sagaUuid,proto3" json:"saga_uuid,omitempty"`                //  UUID sag-и
	OperationName string `protobuf:"bytes,2,opt,name=operation_name,json=operationName,proto3" json:"operation_name,omitempty"` //  Тип выполняемой операции
	Status        uint32 `protobuf:"varint,3,opt,name=status,proto3" json:"status,omitempty"`                                   //  Статус event-а
	// Types that are assignable to Result:
	//
	//	*EventStatus_Info
	//	*EventStatus_AccData
	Result isEventStatus_Result `protobuf_oneof:"result"`
}

func (x *EventStatus) Reset() {
	*x = EventStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_account_account_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventStatus) ProtoMessage() {}

func (x *EventStatus) ProtoReflect() protoreflect.Message {
	mi := &file_account_account_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventStatus.ProtoReflect.Descriptor instead.
func (*EventStatus) Descriptor() ([]byte, []int) {
	return file_account_account_proto_rawDescGZIP(), []int{2}
}

func (x *EventStatus) GetSagaUuid() string {
	if x != nil {
		return x.SagaUuid
	}
	return ""
}

func (x *EventStatus) GetOperationName() string {
	if x != nil {
		return x.OperationName
	}
	return ""
}

func (x *EventStatus) GetStatus() uint32 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (m *EventStatus) GetResult() isEventStatus_Result {
	if m != nil {
		return m.Result
	}
	return nil
}

func (x *EventStatus) GetInfo() string {
	if x, ok := x.GetResult().(*EventStatus_Info); ok {
		return x.Info
	}
	return ""
}

func (x *EventStatus) GetAccData() *platform.FullAccountData {
	if x, ok := x.GetResult().(*EventStatus_AccData); ok {
		return x.AccData
	}
	return nil
}

type isEventStatus_Result interface {
	isEventStatus_Result()
}

type EventStatus_Info struct {
	Info string `protobuf:"bytes,4,opt,name=info,proto3,oneof"` //  Дополнительная информация по event-у
}

type EventStatus_AccData struct {
	AccData *platform.FullAccountData `protobuf:"bytes,5,opt,name=acc_data,json=accData,proto3,oneof"` //  Данные счета
}

func (*EventStatus_Info) isEventStatus_Result() {}

func (*EventStatus_AccData) isEventStatus_Result() {}

var File_account_account_proto protoreflect.FileDescriptor

var file_account_account_proto_rawDesc = []byte{
	0x0a, 0x15, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x1a, 0x17, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2f, 0x70, 0x6c, 0x61, 0x74, 0x66,
	0x6f, 0x72, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x56, 0x0a, 0x10, 0x4f, 0x70, 0x65,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x12, 0x19, 0x0a,
	0x08, 0x61, 0x63, 0x63, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x61, 0x63, 0x63, 0x55, 0x75, 0x69, 0x64, 0x12, 0x27, 0x0a, 0x0f, 0x61, 0x64, 0x64, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x02, 0x52, 0x0e, 0x61, 0x64, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x44, 0x61, 0x74,
	0x61, 0x22, 0xdc, 0x01, 0x0a, 0x09, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x44, 0x61, 0x74, 0x61, 0x12,
	0x1b, 0x0a, 0x09, 0x73, 0x61, 0x67, 0x61, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x73, 0x61, 0x67, 0x61, 0x55, 0x75, 0x69, 0x64, 0x12, 0x25, 0x0a, 0x0e,
	0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x3d, 0x0a, 0x0c, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x70, 0x6c, 0x61, 0x74,
	0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x44, 0x65, 0x74, 0x61,
	0x69, 0x6c, 0x73, 0x48, 0x00, 0x52, 0x0b, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x44, 0x61,
	0x74, 0x61, 0x12, 0x44, 0x0a, 0x0f, 0x61, 0x64, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c,
	0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x61, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x44,
	0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x48, 0x00, 0x52, 0x0e, 0x61, 0x64, 0x64, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x61, 0x6c, 0x49, 0x6e, 0x66, 0x6f, 0x42, 0x06, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x22, 0xc1, 0x01, 0x0a, 0x0b, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x1b, 0x0a, 0x09, 0x73, 0x61, 0x67, 0x61, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x61, 0x67, 0x61, 0x55, 0x75, 0x69, 0x64, 0x12, 0x25, 0x0a,
	0x0e, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x14, 0x0a, 0x04,
	0x69, 0x6e, 0x66, 0x6f, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x04, 0x69, 0x6e,
	0x66, 0x6f, 0x12, 0x36, 0x0a, 0x08, 0x61, 0x63, 0x63, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e,
	0x46, 0x75, 0x6c, 0x6c, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x44, 0x61, 0x74, 0x61, 0x48,
	0x00, 0x52, 0x07, 0x61, 0x63, 0x63, 0x44, 0x61, 0x74, 0x61, 0x42, 0x08, 0x0a, 0x06, 0x72, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x32, 0x4c, 0x0a, 0x0e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3a, 0x0a, 0x0e, 0x52, 0x65, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x12, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x14, 0x2e, 0x61,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x42, 0x0d, 0x5a, 0x0b, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x70,
	0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_account_account_proto_rawDescOnce sync.Once
	file_account_account_proto_rawDescData = file_account_account_proto_rawDesc
)

func file_account_account_proto_rawDescGZIP() []byte {
	file_account_account_proto_rawDescOnce.Do(func() {
		file_account_account_proto_rawDescData = protoimpl.X.CompressGZIP(file_account_account_proto_rawDescData)
	})
	return file_account_account_proto_rawDescData
}

var file_account_account_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_account_account_proto_goTypes = []any{
	(*OperationDetails)(nil),         // 0: account.OperationDetails
	(*EventData)(nil),                // 1: account.EventData
	(*EventStatus)(nil),              // 2: account.EventStatus
	(*platform.AccountDetails)(nil),  // 3: platform.AccountDetails
	(*platform.FullAccountData)(nil), // 4: platform.FullAccountData
}
var file_account_account_proto_depIdxs = []int32{
	3, // 0: account.EventData.account_data:type_name -> platform.AccountDetails
	0, // 1: account.EventData.additional_info:type_name -> account.OperationDetails
	4, // 2: account.EventStatus.acc_data:type_name -> platform.FullAccountData
	1, // 3: account.AccountService.ReserveAccount:input_type -> account.EventData
	2, // 4: account.AccountService.ReserveAccount:output_type -> account.EventStatus
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_account_account_proto_init() }
func file_account_account_proto_init() {
	if File_account_account_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_account_account_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*OperationDetails); i {
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
		file_account_account_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*EventData); i {
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
		file_account_account_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*EventStatus); i {
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
	file_account_account_proto_msgTypes[1].OneofWrappers = []any{
		(*EventData_AccountData)(nil),
		(*EventData_AdditionalInfo)(nil),
	}
	file_account_account_proto_msgTypes[2].OneofWrappers = []any{
		(*EventStatus_Info)(nil),
		(*EventStatus_AccData)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_account_account_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_account_account_proto_goTypes,
		DependencyIndexes: file_account_account_proto_depIdxs,
		MessageInfos:      file_account_account_proto_msgTypes,
	}.Build()
	File_account_account_proto = out.File
	file_account_account_proto_rawDesc = nil
	file_account_account_proto_goTypes = nil
	file_account_account_proto_depIdxs = nil
}
