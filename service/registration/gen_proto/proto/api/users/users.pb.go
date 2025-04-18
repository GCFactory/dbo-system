// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v3.21.12
// source: users/users.proto

package users

import (
	platform "github.com/GCFactory/dbo-system/service/registration/gen_proto/proto/platform"
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

// Дополнительные сведения для операции
type OperationDetails struct {
	state    protoimpl.MessageState `protogen:"open.v1"`
	UserUuid string                 `protobuf:"bytes,1,opt,name=user_uuid,json=userUuid,proto3" json:"user_uuid,omitempty"` //  UUID счёта
	// Types that are valid to be assigned to AdditionalData:
	//
	//	*OperationDetails_Passport
	//	*OperationDetails_SomeData
	AdditionalData isOperationDetails_AdditionalData `protobuf_oneof:"AdditionalData"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *OperationDetails) Reset() {
	*x = OperationDetails{}
	mi := &file_users_users_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *OperationDetails) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OperationDetails) ProtoMessage() {}

func (x *OperationDetails) ProtoReflect() protoreflect.Message {
	mi := &file_users_users_proto_msgTypes[0]
	if x != nil {
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
	return file_users_users_proto_rawDescGZIP(), []int{0}
}

func (x *OperationDetails) GetUserUuid() string {
	if x != nil {
		return x.UserUuid
	}
	return ""
}

func (x *OperationDetails) GetAdditionalData() isOperationDetails_AdditionalData {
	if x != nil {
		return x.AdditionalData
	}
	return nil
}

func (x *OperationDetails) GetPassport() *platform.Passport {
	if x != nil {
		if x, ok := x.AdditionalData.(*OperationDetails_Passport); ok {
			return x.Passport
		}
	}
	return nil
}

func (x *OperationDetails) GetSomeData() string {
	if x != nil {
		if x, ok := x.AdditionalData.(*OperationDetails_SomeData); ok {
			return x.SomeData
		}
	}
	return ""
}

type isOperationDetails_AdditionalData interface {
	isOperationDetails_AdditionalData()
}

type OperationDetails_Passport struct {
	Passport *platform.Passport `protobuf:"bytes,2,opt,name=passport,proto3,oneof"`
}

type OperationDetails_SomeData struct {
	SomeData string `protobuf:"bytes,3,opt,name=some_data,json=someData,proto3,oneof"`
}

func (*OperationDetails_Passport) isOperationDetails_AdditionalData() {}

func (*OperationDetails_SomeData) isOperationDetails_AdditionalData() {}

type UserInfo struct {
	state         protoimpl.MessageState      `protogen:"open.v1"`
	Passport      *platform.Passport          `protobuf:"bytes,1,opt,name=passport,proto3" json:"passport,omitempty"` //  Passport
	UserInn       string                      `protobuf:"bytes,2,opt,name=user_inn,json=userInn,proto3" json:"user_inn,omitempty"`
	UserData      *platform.UserLoginPassword `protobuf:"bytes,3,opt,name=user_data,json=userData,proto3" json:"user_data,omitempty"` // Логин и пароль пользователя
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UserInfo) Reset() {
	*x = UserInfo{}
	mi := &file_users_users_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserInfo) ProtoMessage() {}

func (x *UserInfo) ProtoReflect() protoreflect.Message {
	mi := &file_users_users_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserInfo.ProtoReflect.Descriptor instead.
func (*UserInfo) Descriptor() ([]byte, []int) {
	return file_users_users_proto_rawDescGZIP(), []int{1}
}

func (x *UserInfo) GetPassport() *platform.Passport {
	if x != nil {
		return x.Passport
	}
	return nil
}

func (x *UserInfo) GetUserInn() string {
	if x != nil {
		return x.UserInn
	}
	return ""
}

func (x *UserInfo) GetUserData() *platform.UserLoginPassword {
	if x != nil {
		return x.UserData
	}
	return nil
}

// Данные event-а
type EventData struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	SagaUuid      string                 `protobuf:"bytes,1,opt,name=saga_uuid,json=sagaUuid,proto3" json:"saga_uuid,omitempty"`                //  UUID sag-и
	EventUuid     string                 `protobuf:"bytes,2,opt,name=event_uuid,json=eventUuid,proto3" json:"event_uuid,omitempty"`             //  UUID evnet-a
	OperationName string                 `protobuf:"bytes,3,opt,name=operation_name,json=operationName,proto3" json:"operation_name,omitempty"` // Тип операции
	// Types that are valid to be assigned to Data:
	//
	//	*EventData_UserInfo
	//	*EventData_AdditionalInfo
	Data          isEventData_Data `protobuf_oneof:"data"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EventData) Reset() {
	*x = EventData{}
	mi := &file_users_users_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EventData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventData) ProtoMessage() {}

func (x *EventData) ProtoReflect() protoreflect.Message {
	mi := &file_users_users_proto_msgTypes[2]
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
	return file_users_users_proto_rawDescGZIP(), []int{2}
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

func (x *EventData) GetData() isEventData_Data {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *EventData) GetUserInfo() *UserInfo {
	if x != nil {
		if x, ok := x.Data.(*EventData_UserInfo); ok {
			return x.UserInfo
		}
	}
	return nil
}

func (x *EventData) GetAdditionalInfo() *OperationDetails {
	if x != nil {
		if x, ok := x.Data.(*EventData_AdditionalInfo); ok {
			return x.AdditionalInfo
		}
	}
	return nil
}

type isEventData_Data interface {
	isEventData_Data()
}

type EventData_UserInfo struct {
	UserInfo *UserInfo `protobuf:"bytes,4,opt,name=user_info,json=userInfo,proto3,oneof"`
}

type EventData_AdditionalInfo struct {
	AdditionalInfo *OperationDetails `protobuf:"bytes,5,opt,name=additional_info,json=additionalInfo,proto3,oneof"` //  Дополнительная информация по операции
}

func (*EventData_UserInfo) isEventData_Data() {}

func (*EventData_AdditionalInfo) isEventData_Data() {}

type ListOfAccounts struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Accounts      []string               `protobuf:"bytes,1,rep,name=accounts,proto3" json:"accounts,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListOfAccounts) Reset() {
	*x = ListOfAccounts{}
	mi := &file_users_users_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListOfAccounts) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListOfAccounts) ProtoMessage() {}

func (x *ListOfAccounts) ProtoReflect() protoreflect.Message {
	mi := &file_users_users_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListOfAccounts.ProtoReflect.Descriptor instead.
func (*ListOfAccounts) Descriptor() ([]byte, []int) {
	return file_users_users_proto_rawDescGZIP(), []int{3}
}

func (x *ListOfAccounts) GetAccounts() []string {
	if x != nil {
		return x.Accounts
	}
	return nil
}

type FullData struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Passport      *platform.Passport     `protobuf:"bytes,1,opt,name=passport,proto3" json:"passport,omitempty"`
	UserId        string                 `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	UserInn       string                 `protobuf:"bytes,3,opt,name=user_inn,json=userInn,proto3" json:"user_inn,omitempty"`
	UserLogin     string                 `protobuf:"bytes,4,opt,name=user_login,json=userLogin,proto3" json:"user_login,omitempty"`
	Accounts      *ListOfAccounts        `protobuf:"bytes,5,opt,name=accounts,proto3" json:"accounts,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FullData) Reset() {
	*x = FullData{}
	mi := &file_users_users_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FullData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FullData) ProtoMessage() {}

func (x *FullData) ProtoReflect() protoreflect.Message {
	mi := &file_users_users_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FullData.ProtoReflect.Descriptor instead.
func (*FullData) Descriptor() ([]byte, []int) {
	return file_users_users_proto_rawDescGZIP(), []int{4}
}

func (x *FullData) GetPassport() *platform.Passport {
	if x != nil {
		return x.Passport
	}
	return nil
}

func (x *FullData) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *FullData) GetUserInn() string {
	if x != nil {
		return x.UserInn
	}
	return ""
}

func (x *FullData) GetUserLogin() string {
	if x != nil {
		return x.UserLogin
	}
	return ""
}

func (x *FullData) GetAccounts() *ListOfAccounts {
	if x != nil {
		return x.Accounts
	}
	return nil
}

// Результат event-а
type EventSuccess struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	SagaUuid      string                 `protobuf:"bytes,1,opt,name=saga_uuid,json=sagaUuid,proto3" json:"saga_uuid,omitempty"`                //  UUID sag-и
	EventUuid     string                 `protobuf:"bytes,2,opt,name=event_uuid,json=eventUuid,proto3" json:"event_uuid,omitempty"`             //  UUID evnet-a
	OperationName string                 `protobuf:"bytes,3,opt,name=operation_name,json=operationName,proto3" json:"operation_name,omitempty"` //  Тип выполняемой операции
	// Types that are valid to be assigned to Result:
	//
	//	*EventSuccess_Info
	//	*EventSuccess_FullData
	//	*EventSuccess_Accounts
	Result        isEventSuccess_Result `protobuf_oneof:"result"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EventSuccess) Reset() {
	*x = EventSuccess{}
	mi := &file_users_users_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EventSuccess) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventSuccess) ProtoMessage() {}

func (x *EventSuccess) ProtoReflect() protoreflect.Message {
	mi := &file_users_users_proto_msgTypes[5]
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
	return file_users_users_proto_rawDescGZIP(), []int{5}
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

func (x *EventSuccess) GetResult() isEventSuccess_Result {
	if x != nil {
		return x.Result
	}
	return nil
}

func (x *EventSuccess) GetInfo() string {
	if x != nil {
		if x, ok := x.Result.(*EventSuccess_Info); ok {
			return x.Info
		}
	}
	return ""
}

func (x *EventSuccess) GetFullData() *FullData {
	if x != nil {
		if x, ok := x.Result.(*EventSuccess_FullData); ok {
			return x.FullData
		}
	}
	return nil
}

func (x *EventSuccess) GetAccounts() *ListOfAccounts {
	if x != nil {
		if x, ok := x.Result.(*EventSuccess_Accounts); ok {
			return x.Accounts
		}
	}
	return nil
}

type isEventSuccess_Result interface {
	isEventSuccess_Result()
}

type EventSuccess_Info struct {
	Info string `protobuf:"bytes,4,opt,name=info,proto3,oneof"` //  Дополнительная информация по event-у
}

type EventSuccess_FullData struct {
	FullData *FullData `protobuf:"bytes,5,opt,name=full_data,json=fullData,proto3,oneof"` //  Full user data
}

type EventSuccess_Accounts struct {
	Accounts *ListOfAccounts `protobuf:"bytes,6,opt,name=accounts,proto3,oneof"`
}

func (*EventSuccess_Info) isEventSuccess_Result() {}

func (*EventSuccess_FullData) isEventSuccess_Result() {}

func (*EventSuccess_Accounts) isEventSuccess_Result() {}

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
	mi := &file_users_users_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EventError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventError) ProtoMessage() {}

func (x *EventError) ProtoReflect() protoreflect.Message {
	mi := &file_users_users_proto_msgTypes[6]
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
	return file_users_users_proto_rawDescGZIP(), []int{6}
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

var File_users_users_proto protoreflect.FileDescriptor

var file_users_users_proto_rawDesc = string([]byte{
	0x0a, 0x11, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x1a, 0x17, 0x70, 0x6c, 0x61, 0x74,
	0x66, 0x6f, 0x72, 0x6d, 0x2f, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x92, 0x01, 0x0a, 0x10, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x75, 0x73, 0x65, 0x72,
	0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65,
	0x72, 0x55, 0x75, 0x69, 0x64, 0x12, 0x30, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x70, 0x6f, 0x72,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f,
	0x72, 0x6d, 0x2e, 0x50, 0x61, 0x73, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x48, 0x00, 0x52, 0x08, 0x70,
	0x61, 0x73, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x1d, 0x0a, 0x09, 0x73, 0x6f, 0x6d, 0x65, 0x5f,
	0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x08, 0x73, 0x6f,
	0x6d, 0x65, 0x44, 0x61, 0x74, 0x61, 0x42, 0x10, 0x0a, 0x0e, 0x41, 0x64, 0x64, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x61, 0x6c, 0x44, 0x61, 0x74, 0x61, 0x22, 0x8f, 0x01, 0x0a, 0x08, 0x55, 0x73, 0x65,
	0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x2e, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x70, 0x6f, 0x72,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f,
	0x72, 0x6d, 0x2e, 0x50, 0x61, 0x73, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x08, 0x70, 0x61, 0x73,
	0x73, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x6e,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x75, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x6e,
	0x12, 0x38, 0x0a, 0x09, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x55,
	0x73, 0x65, 0x72, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x22, 0xea, 0x01, 0x0a, 0x09, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x61, 0x67, 0x61,
	0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x61, 0x67,
	0x61, 0x55, 0x75, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x75,
	0x75, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x55, 0x75, 0x69, 0x64, 0x12, 0x25, 0x0a, 0x0e, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6f, 0x70,
	0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x2e, 0x0a, 0x09, 0x75,
	0x73, 0x65, 0x72, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f,
	0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x48,
	0x00, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x42, 0x0a, 0x0f, 0x61,
	0x64, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x4f, 0x70, 0x65,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x48, 0x00, 0x52,
	0x0e, 0x61, 0x64, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x49, 0x6e, 0x66, 0x6f, 0x42,
	0x06, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x2c, 0x0a, 0x0e, 0x4c, 0x69, 0x73, 0x74, 0x4f,
	0x66, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x61, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x61, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x73, 0x22, 0xc0, 0x01, 0x0a, 0x08, 0x46, 0x75, 0x6c, 0x6c, 0x44, 0x61,
	0x74, 0x61, 0x12, 0x2e, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e,
	0x50, 0x61, 0x73, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x70, 0x6f,
	0x72, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x75,
	0x73, 0x65, 0x72, 0x5f, 0x69, 0x6e, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x75,
	0x73, 0x65, 0x72, 0x49, 0x6e, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x6c,
	0x6f, 0x67, 0x69, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x75, 0x73, 0x65, 0x72,
	0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x31, 0x0a, 0x08, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e,
	0x4c, 0x69, 0x73, 0x74, 0x4f, 0x66, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x52, 0x08,
	0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x22, 0xf6, 0x01, 0x0a, 0x0c, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x61, 0x67,
	0x61, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x61,
	0x67, 0x61, 0x55, 0x75, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f,
	0x75, 0x75, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x55, 0x75, 0x69, 0x64, 0x12, 0x25, 0x0a, 0x0e, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6f,
	0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x04,
	0x69, 0x6e, 0x66, 0x6f, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x04, 0x69, 0x6e,
	0x66, 0x6f, 0x12, 0x2e, 0x0a, 0x09, 0x66, 0x75, 0x6c, 0x6c, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x46, 0x75,
	0x6c, 0x6c, 0x44, 0x61, 0x74, 0x61, 0x48, 0x00, 0x52, 0x08, 0x66, 0x75, 0x6c, 0x6c, 0x44, 0x61,
	0x74, 0x61, 0x12, 0x33, 0x0a, 0x08, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x4c, 0x69, 0x73,
	0x74, 0x4f, 0x66, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x48, 0x00, 0x52, 0x08, 0x61,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x42, 0x08, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x22, 0x9b, 0x01, 0x0a, 0x0a, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x12, 0x1b, 0x0a, 0x09, 0x73, 0x61, 0x67, 0x61, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x61, 0x67, 0x61, 0x55, 0x75, 0x69, 0x64, 0x12, 0x1d, 0x0a,
	0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x55, 0x75, 0x69, 0x64, 0x12, 0x25, 0x0a, 0x0e,
	0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x69,
	0x6e, 0x66, 0x6f, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x42,
	0x13, 0x5a, 0x11, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x75,
	0x73, 0x65, 0x72, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_users_users_proto_rawDescOnce sync.Once
	file_users_users_proto_rawDescData []byte
)

func file_users_users_proto_rawDescGZIP() []byte {
	file_users_users_proto_rawDescOnce.Do(func() {
		file_users_users_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_users_users_proto_rawDesc), len(file_users_users_proto_rawDesc)))
	})
	return file_users_users_proto_rawDescData
}

var file_users_users_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_users_users_proto_goTypes = []any{
	(*OperationDetails)(nil),           // 0: users.OperationDetails
	(*UserInfo)(nil),                   // 1: users.UserInfo
	(*EventData)(nil),                  // 2: users.EventData
	(*ListOfAccounts)(nil),             // 3: users.ListOfAccounts
	(*FullData)(nil),                   // 4: users.FullData
	(*EventSuccess)(nil),               // 5: users.EventSuccess
	(*EventError)(nil),                 // 6: users.EventError
	(*platform.Passport)(nil),          // 7: platform.Passport
	(*platform.UserLoginPassword)(nil), // 8: platform.UserLoginPassword
}
var file_users_users_proto_depIdxs = []int32{
	7, // 0: users.OperationDetails.passport:type_name -> platform.Passport
	7, // 1: users.UserInfo.passport:type_name -> platform.Passport
	8, // 2: users.UserInfo.user_data:type_name -> platform.UserLoginPassword
	1, // 3: users.EventData.user_info:type_name -> users.UserInfo
	0, // 4: users.EventData.additional_info:type_name -> users.OperationDetails
	7, // 5: users.FullData.passport:type_name -> platform.Passport
	3, // 6: users.FullData.accounts:type_name -> users.ListOfAccounts
	4, // 7: users.EventSuccess.full_data:type_name -> users.FullData
	3, // 8: users.EventSuccess.accounts:type_name -> users.ListOfAccounts
	9, // [9:9] is the sub-list for method output_type
	9, // [9:9] is the sub-list for method input_type
	9, // [9:9] is the sub-list for extension type_name
	9, // [9:9] is the sub-list for extension extendee
	0, // [0:9] is the sub-list for field type_name
}

func init() { file_users_users_proto_init() }
func file_users_users_proto_init() {
	if File_users_users_proto != nil {
		return
	}
	file_users_users_proto_msgTypes[0].OneofWrappers = []any{
		(*OperationDetails_Passport)(nil),
		(*OperationDetails_SomeData)(nil),
	}
	file_users_users_proto_msgTypes[2].OneofWrappers = []any{
		(*EventData_UserInfo)(nil),
		(*EventData_AdditionalInfo)(nil),
	}
	file_users_users_proto_msgTypes[5].OneofWrappers = []any{
		(*EventSuccess_Info)(nil),
		(*EventSuccess_FullData)(nil),
		(*EventSuccess_Accounts)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_users_users_proto_rawDesc), len(file_users_users_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_users_users_proto_goTypes,
		DependencyIndexes: file_users_users_proto_depIdxs,
		MessageInfos:      file_users_users_proto_msgTypes,
	}.Build()
	File_users_users_proto = out.File
	file_users_users_proto_goTypes = nil
	file_users_users_proto_depIdxs = nil
}
