// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v3.12.4
// source: registration/registration.proto

package registration

import (
	platform "github.com/GCFactory/dbo-system/service/registration/proto/platform"
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

// Персональные данные пользователя
type UserData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Inn      string             `protobuf:"bytes,1,opt,name=inn,proto3" json:"inn,omitempty"`           //  ИНН
	Snils    string             `protobuf:"bytes,2,opt,name=snils,proto3" json:"snils,omitempty"`       //  СНИЛС
	Passport *platform.Passport `protobuf:"bytes,3,opt,name=passport,proto3" json:"passport,omitempty"` //  Паспорт
}

func (x *UserData) Reset() {
	*x = UserData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_registration_registration_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserData) ProtoMessage() {}

func (x *UserData) ProtoReflect() protoreflect.Message {
	mi := &file_registration_registration_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserData.ProtoReflect.Descriptor instead.
func (*UserData) Descriptor() ([]byte, []int) {
	return file_registration_registration_proto_rawDescGZIP(), []int{0}
}

func (x *UserData) GetInn() string {
	if x != nil {
		return x.Inn
	}
	return ""
}

func (x *UserData) GetSnils() string {
	if x != nil {
		return x.Snils
	}
	return ""
}

func (x *UserData) GetPassport() *platform.Passport {
	if x != nil {
		return x.Passport
	}
	return nil
}

// Данные для начала регистрации
type RegestrationData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserData     *UserData             `protobuf:"bytes,1,opt,name=user_data,json=userData,proto3" json:"user_data,omitempty"`                                         //  Персональные данные пользователя
	ActivityType platform.ActivityType `protobuf:"varint,2,opt,name=activity_type,json=activityType,proto3,enum=platform.ActivityType" json:"activity_type,omitempty"` //  Тип деятельности
	TaxationType platform.TaxationType `protobuf:"varint,3,opt,name=taxation_type,json=taxationType,proto3,enum=platform.TaxationType" json:"taxation_type,omitempty"` //  Тип налогообложения
}

func (x *RegestrationData) Reset() {
	*x = RegestrationData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_registration_registration_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegestrationData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegestrationData) ProtoMessage() {}

func (x *RegestrationData) ProtoReflect() protoreflect.Message {
	mi := &file_registration_registration_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegestrationData.ProtoReflect.Descriptor instead.
func (*RegestrationData) Descriptor() ([]byte, []int) {
	return file_registration_registration_proto_rawDescGZIP(), []int{1}
}

func (x *RegestrationData) GetUserData() *UserData {
	if x != nil {
		return x.UserData
	}
	return nil
}

func (x *RegestrationData) GetActivityType() platform.ActivityType {
	if x != nil {
		return x.ActivityType
	}
	return platform.ActivityType(0)
}

func (x *RegestrationData) GetTaxationType() platform.TaxationType {
	if x != nil {
		return x.TaxationType
	}
	return platform.TaxationType(0)
}

// Ответ на gateway
type RegistrationResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RequestUuid string `protobuf:"bytes,1,opt,name=request_uuid,json=requestUuid,proto3" json:"request_uuid,omitempty"` //  UUID запроса
	Status      uint32 `protobuf:"varint,2,opt,name=status,proto3" json:"status,omitempty"`                             //  Статус запроса
	Info        string `protobuf:"bytes,3,opt,name=info,proto3" json:"info,omitempty"`                                  //  Дополнительная информация по запросу
}

func (x *RegistrationResult) Reset() {
	*x = RegistrationResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_registration_registration_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegistrationResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegistrationResult) ProtoMessage() {}

func (x *RegistrationResult) ProtoReflect() protoreflect.Message {
	mi := &file_registration_registration_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegistrationResult.ProtoReflect.Descriptor instead.
func (*RegistrationResult) Descriptor() ([]byte, []int) {
	return file_registration_registration_proto_rawDescGZIP(), []int{2}
}

func (x *RegistrationResult) GetRequestUuid() string {
	if x != nil {
		return x.RequestUuid
	}
	return ""
}

func (x *RegistrationResult) GetStatus() uint32 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *RegistrationResult) GetInfo() string {
	if x != nil {
		return x.Info
	}
	return ""
}

// Данные event-а
type EventData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SagaUuid string `protobuf:"bytes,1,opt,name=saga_uuid,json=sagaUuid,proto3" json:"saga_uuid,omitempty"` //  UUID sag-и
	// Types that are assignable to Data:
	//
	//	*EventData_UserData
	Data isEventData_Data `protobuf_oneof:"data"`
}

func (x *EventData) Reset() {
	*x = EventData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_registration_registration_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventData) ProtoMessage() {}

func (x *EventData) ProtoReflect() protoreflect.Message {
	mi := &file_registration_registration_proto_msgTypes[3]
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
	return file_registration_registration_proto_rawDescGZIP(), []int{3}
}

func (x *EventData) GetSagaUuid() string {
	if x != nil {
		return x.SagaUuid
	}
	return ""
}

func (m *EventData) GetData() isEventData_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (x *EventData) GetUserData() *UserData {
	if x, ok := x.GetData().(*EventData_UserData); ok {
		return x.UserData
	}
	return nil
}

type isEventData_Data interface {
	isEventData_Data()
}

type EventData_UserData struct {
	UserData *UserData `protobuf:"bytes,2,opt,name=user_data,json=userData,proto3,oneof"` //  Персональные данные пользователя
}

func (*EventData_UserData) isEventData_Data() {}

// Результат event-а
type EventStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SagaUuid string `protobuf:"bytes,1,opt,name=saga_uuid,json=sagaUuid,proto3" json:"saga_uuid,omitempty"` //  UUID sag-и
	Status   uint32 `protobuf:"varint,2,opt,name=status,proto3" json:"status,omitempty"`                    //  Статус event-а
	Info     string `protobuf:"bytes,3,opt,name=info,proto3" json:"info,omitempty"`                         //  Дополнительная информация по event-у
}

func (x *EventStatus) Reset() {
	*x = EventStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_registration_registration_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventStatus) ProtoMessage() {}

func (x *EventStatus) ProtoReflect() protoreflect.Message {
	mi := &file_registration_registration_proto_msgTypes[4]
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
	return file_registration_registration_proto_rawDescGZIP(), []int{4}
}

func (x *EventStatus) GetSagaUuid() string {
	if x != nil {
		return x.SagaUuid
	}
	return ""
}

func (x *EventStatus) GetStatus() uint32 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *EventStatus) GetInfo() string {
	if x != nil {
		return x.Info
	}
	return ""
}

var File_registration_registration_proto protoreflect.FileDescriptor

var file_registration_registration_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x72,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x0c, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a,
	0x17, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2f, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f,
	0x72, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x62, 0x0a, 0x08, 0x55, 0x73, 0x65, 0x72,
	0x44, 0x61, 0x74, 0x61, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x6e, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x69, 0x6e, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x6e, 0x69, 0x6c, 0x73, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x6e, 0x69, 0x6c, 0x73, 0x12, 0x2e, 0x0a, 0x08,
	0x70, 0x61, 0x73, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12,
	0x2e, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x50, 0x61, 0x73, 0x73, 0x70, 0x6f,
	0x72, 0x74, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x22, 0xc1, 0x01, 0x0a,
	0x10, 0x52, 0x65, 0x67, 0x65, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x61, 0x74,
	0x61, 0x12, 0x33, 0x0a, 0x09, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x52, 0x08, 0x75, 0x73,
	0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x12, 0x3b, 0x0a, 0x0d, 0x61, 0x63, 0x74, 0x69, 0x76, 0x69,
	0x74, 0x79, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x16, 0x2e,
	0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74,
	0x79, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0c, 0x61, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x3b, 0x0a, 0x0d, 0x74, 0x61, 0x78, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x16, 0x2e, 0x70, 0x6c, 0x61,
	0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x54, 0x61, 0x78, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x0c, 0x74, 0x61, 0x78, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65,
	0x22, 0x63, 0x0a, 0x12, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x72, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x55, 0x75, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x12, 0x12, 0x0a, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x69, 0x6e, 0x66, 0x6f, 0x22, 0x67, 0x0a, 0x09, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x44, 0x61,
	0x74, 0x61, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x61, 0x67, 0x61, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x61, 0x67, 0x61, 0x55, 0x75, 0x69, 0x64, 0x12,
	0x35, 0x0a, 0x09, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x16, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x48, 0x00, 0x52, 0x08, 0x75, 0x73,
	0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x42, 0x06, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x56,
	0x0a, 0x0b, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1b, 0x0a,
	0x09, 0x73, 0x61, 0x67, 0x61, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x73, 0x61, 0x67, 0x61, 0x55, 0x75, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x32, 0xf0, 0x01, 0x0a, 0x13, 0x52, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x52,
	0x0a, 0x0e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72,
	0x12, 0x1e, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x52, 0x65, 0x67, 0x65, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x61, 0x74, 0x61,
	0x1a, 0x20, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x12, 0x40, 0x0a, 0x0a, 0x53, 0x74, 0x61, 0x72, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x12, 0x17, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x19, 0x2e, 0x72, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x43, 0x0a, 0x0d, 0x52, 0x6f, 0x6c, 0x6c, 0x42, 0x61, 0x63, 0x6b,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x17, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x19,
	0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x42, 0x10, 0x5a, 0x0e, 0x2e, 0x2f, 0x72,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_registration_registration_proto_rawDescOnce sync.Once
	file_registration_registration_proto_rawDescData = file_registration_registration_proto_rawDesc
)

func file_registration_registration_proto_rawDescGZIP() []byte {
	file_registration_registration_proto_rawDescOnce.Do(func() {
		file_registration_registration_proto_rawDescData = protoimpl.X.CompressGZIP(file_registration_registration_proto_rawDescData)
	})
	return file_registration_registration_proto_rawDescData
}

var file_registration_registration_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_registration_registration_proto_goTypes = []interface{}{
	(*UserData)(nil),           // 0: registration.UserData
	(*RegestrationData)(nil),   // 1: registration.RegestrationData
	(*RegistrationResult)(nil), // 2: registration.RegistrationResult
	(*EventData)(nil),          // 3: registration.EventData
	(*EventStatus)(nil),        // 4: registration.EventStatus
	(*platform.Passport)(nil),  // 5: platform.Passport
	(platform.ActivityType)(0), // 6: platform.ActivityType
	(platform.TaxationType)(0), // 7: platform.TaxationType
}
var file_registration_registration_proto_depIdxs = []int32{
	5, // 0: registration.UserData.passport:type_name -> platform.Passport
	0, // 1: registration.RegestrationData.user_data:type_name -> registration.UserData
	6, // 2: registration.RegestrationData.activity_type:type_name -> platform.ActivityType
	7, // 3: registration.RegestrationData.taxation_type:type_name -> platform.TaxationType
	0, // 4: registration.EventData.user_data:type_name -> registration.UserData
	1, // 5: registration.RegistrationService.RegistrateUser:input_type -> registration.RegestrationData
	3, // 6: registration.RegistrationService.StartEvent:input_type -> registration.EventData
	3, // 7: registration.RegistrationService.RollBackEvent:input_type -> registration.EventData
	2, // 8: registration.RegistrationService.RegistrateUser:output_type -> registration.RegistrationResult
	4, // 9: registration.RegistrationService.StartEvent:output_type -> registration.EventStatus
	4, // 10: registration.RegistrationService.RollBackEvent:output_type -> registration.EventStatus
	8, // [8:11] is the sub-list for method output_type
	5, // [5:8] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_registration_registration_proto_init() }
func file_registration_registration_proto_init() {
	if File_registration_registration_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_registration_registration_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserData); i {
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
		file_registration_registration_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegestrationData); i {
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
		file_registration_registration_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegistrationResult); i {
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
		file_registration_registration_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
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
		file_registration_registration_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
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
	file_registration_registration_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*EventData_UserData)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_registration_registration_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_registration_registration_proto_goTypes,
		DependencyIndexes: file_registration_registration_proto_depIdxs,
		MessageInfos:      file_registration_registration_proto_msgTypes,
	}.Build()
	File_registration_registration_proto = out.File
	file_registration_registration_proto_rawDesc = nil
	file_registration_registration_proto_goTypes = nil
	file_registration_registration_proto_depIdxs = nil
}