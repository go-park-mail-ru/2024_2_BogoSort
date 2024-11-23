// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v5.28.3
// source: survey.proto

package survey

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

type AddAnswerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId     string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	QuestionId string `protobuf:"bytes,2,opt,name=question_id,json=questionId,proto3" json:"question_id,omitempty"`
	Value      int32  `protobuf:"varint,3,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *AddAnswerRequest) Reset() {
	*x = AddAnswerRequest{}
	mi := &file_survey_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AddAnswerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddAnswerRequest) ProtoMessage() {}

func (x *AddAnswerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_survey_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddAnswerRequest.ProtoReflect.Descriptor instead.
func (*AddAnswerRequest) Descriptor() ([]byte, []int) {
	return file_survey_proto_rawDescGZIP(), []int{0}
}

func (x *AddAnswerRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *AddAnswerRequest) GetQuestionId() string {
	if x != nil {
		return x.QuestionId
	}
	return ""
}

func (x *AddAnswerRequest) GetValue() int32 {
	if x != nil {
		return x.Value
	}
	return 0
}

type AddAnswerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *AddAnswerResponse) Reset() {
	*x = AddAnswerResponse{}
	mi := &file_survey_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AddAnswerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddAnswerResponse) ProtoMessage() {}

func (x *AddAnswerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_survey_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddAnswerResponse.ProtoReflect.Descriptor instead.
func (*AddAnswerResponse) Descriptor() ([]byte, []int) {
	return file_survey_proto_rawDescGZIP(), []int{1}
}

func (x *AddAnswerResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type NoContent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *NoContent) Reset() {
	*x = NoContent{}
	mi := &file_survey_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NoContent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NoContent) ProtoMessage() {}

func (x *NoContent) ProtoReflect() protoreflect.Message {
	mi := &file_survey_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NoContent.ProtoReflect.Descriptor instead.
func (*NoContent) Descriptor() ([]byte, []int) {
	return file_survey_proto_rawDescGZIP(), []int{2}
}

var File_survey_proto protoreflect.FileDescriptor

var file_survey_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x73, 0x75, 0x72, 0x76, 0x65, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x73, 0x75, 0x72, 0x76, 0x65, 0x79, 0x22, 0x62, 0x0a, 0x10, 0x41, 0x64, 0x64, 0x41, 0x6e, 0x73,
	0x77, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73,
	0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65,
	0x72, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69,
	0x6f, 0x6e, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x2d, 0x0a, 0x11, 0x41, 0x64,
	0x64, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x0b, 0x0a, 0x09, 0x4e, 0x6f, 0x43,
	0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x32, 0x7f, 0x0a, 0x0d, 0x53, 0x75, 0x72, 0x76, 0x65, 0x79,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x40, 0x0a, 0x09, 0x41, 0x64, 0x64, 0x41, 0x6e,
	0x73, 0x77, 0x65, 0x72, 0x12, 0x18, 0x2e, 0x73, 0x75, 0x72, 0x76, 0x65, 0x79, 0x2e, 0x41, 0x64,
	0x64, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19,
	0x2e, 0x73, 0x75, 0x72, 0x76, 0x65, 0x79, 0x2e, 0x41, 0x64, 0x64, 0x41, 0x6e, 0x73, 0x77, 0x65,
	0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2c, 0x0a, 0x04, 0x50, 0x69, 0x6e,
	0x67, 0x12, 0x11, 0x2e, 0x73, 0x75, 0x72, 0x76, 0x65, 0x79, 0x2e, 0x4e, 0x6f, 0x43, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x1a, 0x11, 0x2e, 0x73, 0x75, 0x72, 0x76, 0x65, 0x79, 0x2e, 0x4e, 0x6f,
	0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x42, 0x0b, 0x5a, 0x09, 0x2e, 0x2f, 0x3b, 0x73, 0x75,
	0x72, 0x76, 0x65, 0x79, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_survey_proto_rawDescOnce sync.Once
	file_survey_proto_rawDescData = file_survey_proto_rawDesc
)

func file_survey_proto_rawDescGZIP() []byte {
	file_survey_proto_rawDescOnce.Do(func() {
		file_survey_proto_rawDescData = protoimpl.X.CompressGZIP(file_survey_proto_rawDescData)
	})
	return file_survey_proto_rawDescData
}

var file_survey_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_survey_proto_goTypes = []any{
	(*AddAnswerRequest)(nil),  // 0: survey.AddAnswerRequest
	(*AddAnswerResponse)(nil), // 1: survey.AddAnswerResponse
	(*NoContent)(nil),         // 2: survey.NoContent
}
var file_survey_proto_depIdxs = []int32{
	0, // 0: survey.SurveyService.AddAnswer:input_type -> survey.AddAnswerRequest
	2, // 1: survey.SurveyService.Ping:input_type -> survey.NoContent
	1, // 2: survey.SurveyService.AddAnswer:output_type -> survey.AddAnswerResponse
	2, // 3: survey.SurveyService.Ping:output_type -> survey.NoContent
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_survey_proto_init() }
func file_survey_proto_init() {
	if File_survey_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_survey_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_survey_proto_goTypes,
		DependencyIndexes: file_survey_proto_depIdxs,
		MessageInfos:      file_survey_proto_msgTypes,
	}.Build()
	File_survey_proto = out.File
	file_survey_proto_rawDesc = nil
	file_survey_proto_goTypes = nil
	file_survey_proto_depIdxs = nil
}
