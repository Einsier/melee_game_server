// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1-devel
// 	protoc        v3.6.1
// source: response.proto

package proto

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

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResponseCode            ResponseCode             `protobuf:"varint,1,opt,name=responseCode,proto3,enum=proto.ResponseCode" json:"responseCode,omitempty"`
	PlayerEnterGameResponse *PlayerEnterGameResponse `protobuf:"bytes,2,opt,name=playerEnterGameResponse,proto3" json:"playerEnterGameResponse,omitempty"`
	PlayerQuitGameResponse  *PlayerQuitGameResponse  `protobuf:"bytes,3,opt,name=playerQuitGameResponse,proto3" json:"playerQuitGameResponse,omitempty"`
	HeroGetPropResponse     *HeroGetPropResponse     `protobuf:"bytes,4,opt,name=heroGetPropResponse,proto3" json:"heroGetPropResponse,omitempty"`
	PlayerHeartBeatResponse *PlayerHeartBeatResponse `protobuf:"bytes,5,opt,name=playerHeartBeatResponse,proto3" json:"playerHeartBeatResponse,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_response_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_response_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_response_proto_rawDescGZIP(), []int{0}
}

func (x *Response) GetResponseCode() ResponseCode {
	if x != nil {
		return x.ResponseCode
	}
	return ResponseCode_PlayerEnterGameResponseCode
}

func (x *Response) GetPlayerEnterGameResponse() *PlayerEnterGameResponse {
	if x != nil {
		return x.PlayerEnterGameResponse
	}
	return nil
}

func (x *Response) GetPlayerQuitGameResponse() *PlayerQuitGameResponse {
	if x != nil {
		return x.PlayerQuitGameResponse
	}
	return nil
}

func (x *Response) GetHeroGetPropResponse() *HeroGetPropResponse {
	if x != nil {
		return x.HeroGetPropResponse
	}
	return nil
}

func (x *Response) GetPlayerHeartBeatResponse() *PlayerHeartBeatResponse {
	if x != nil {
		return x.PlayerHeartBeatResponse
	}
	return nil
}

type PlayerEnterGameResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HeroId int32 `protobuf:"varint,1,opt,name=heroId,proto3" json:"heroId,omitempty"` //若为-1表示失败
}

func (x *PlayerEnterGameResponse) Reset() {
	*x = PlayerEnterGameResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_response_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PlayerEnterGameResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlayerEnterGameResponse) ProtoMessage() {}

func (x *PlayerEnterGameResponse) ProtoReflect() protoreflect.Message {
	mi := &file_response_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlayerEnterGameResponse.ProtoReflect.Descriptor instead.
func (*PlayerEnterGameResponse) Descriptor() ([]byte, []int) {
	return file_response_proto_rawDescGZIP(), []int{1}
}

func (x *PlayerEnterGameResponse) GetHeroId() int32 {
	if x != nil {
		return x.HeroId
	}
	return 0
}

type PlayerQuitGameResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	//todo 待构思
	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *PlayerQuitGameResponse) Reset() {
	*x = PlayerQuitGameResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_response_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PlayerQuitGameResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlayerQuitGameResponse) ProtoMessage() {}

func (x *PlayerQuitGameResponse) ProtoReflect() protoreflect.Message {
	mi := &file_response_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlayerQuitGameResponse.ProtoReflect.Descriptor instead.
func (*PlayerQuitGameResponse) Descriptor() ([]byte, []int) {
	return file_response_proto_rawDescGZIP(), []int{2}
}

func (x *PlayerQuitGameResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type HeroGetPropResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *HeroGetPropResponse) Reset() {
	*x = HeroGetPropResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_response_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HeroGetPropResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HeroGetPropResponse) ProtoMessage() {}

func (x *HeroGetPropResponse) ProtoReflect() protoreflect.Message {
	mi := &file_response_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HeroGetPropResponse.ProtoReflect.Descriptor instead.
func (*HeroGetPropResponse) Descriptor() ([]byte, []int) {
	return file_response_proto_rawDescGZIP(), []int{3}
}

func (x *HeroGetPropResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type PlayerHeartBeatResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HeartbeatId int32 `protobuf:"varint,1,opt,name=heartbeatId,proto3" json:"heartbeatId,omitempty"`
}

func (x *PlayerHeartBeatResponse) Reset() {
	*x = PlayerHeartBeatResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_response_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PlayerHeartBeatResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlayerHeartBeatResponse) ProtoMessage() {}

func (x *PlayerHeartBeatResponse) ProtoReflect() protoreflect.Message {
	mi := &file_response_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlayerHeartBeatResponse.ProtoReflect.Descriptor instead.
func (*PlayerHeartBeatResponse) Descriptor() ([]byte, []int) {
	return file_response_proto_rawDescGZIP(), []int{4}
}

func (x *PlayerHeartBeatResponse) GetHeartbeatId() int32 {
	if x != nil {
		return x.HeartbeatId
	}
	return 0
}

var File_response_proto protoreflect.FileDescriptor

var file_response_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0c, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9c, 0x03, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x37, 0x0a, 0x0c, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x43, 0x6f,
	0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x0c, 0x72,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x58, 0x0a, 0x17, 0x70,
	0x6c, 0x61, 0x79, 0x65, 0x72, 0x45, 0x6e, 0x74, 0x65, 0x72, 0x47, 0x61, 0x6d, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x45, 0x6e, 0x74, 0x65, 0x72,
	0x47, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52, 0x17, 0x70, 0x6c,
	0x61, 0x79, 0x65, 0x72, 0x45, 0x6e, 0x74, 0x65, 0x72, 0x47, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x55, 0x0a, 0x16, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x51,
	0x75, 0x69, 0x74, 0x47, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x6c,
	0x61, 0x79, 0x65, 0x72, 0x51, 0x75, 0x69, 0x74, 0x47, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x52, 0x16, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x51, 0x75, 0x69, 0x74,
	0x47, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4c, 0x0a, 0x13,
	0x68, 0x65, 0x72, 0x6f, 0x47, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x48, 0x65, 0x72, 0x6f, 0x47, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x70, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52, 0x13, 0x68, 0x65, 0x72, 0x6f, 0x47, 0x65, 0x74, 0x50, 0x72,
	0x6f, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x58, 0x0a, 0x17, 0x70, 0x6c,
	0x61, 0x79, 0x65, 0x72, 0x48, 0x65, 0x61, 0x72, 0x74, 0x42, 0x65, 0x61, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x48, 0x65, 0x61, 0x72, 0x74, 0x42,
	0x65, 0x61, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52, 0x17, 0x70, 0x6c, 0x61,
	0x79, 0x65, 0x72, 0x48, 0x65, 0x61, 0x72, 0x74, 0x42, 0x65, 0x61, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x31, 0x0a, 0x17, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x45, 0x6e,
	0x74, 0x65, 0x72, 0x47, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x16, 0x0a, 0x06, 0x68, 0x65, 0x72, 0x6f, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x06, 0x68, 0x65, 0x72, 0x6f, 0x49, 0x64, 0x22, 0x32, 0x0a, 0x16, 0x50, 0x6c, 0x61, 0x79, 0x65,
	0x72, 0x51, 0x75, 0x69, 0x74, 0x47, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x22, 0x2f, 0x0a, 0x13, 0x48,
	0x65, 0x72, 0x6f, 0x47, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x22, 0x3b, 0x0a, 0x17,
	0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x48, 0x65, 0x61, 0x72, 0x74, 0x42, 0x65, 0x61, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x68, 0x65, 0x61, 0x72, 0x74,
	0x62, 0x65, 0x61, 0x74, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x68, 0x65,
	0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x49, 0x64, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x3b, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_response_proto_rawDescOnce sync.Once
	file_response_proto_rawDescData = file_response_proto_rawDesc
)

func file_response_proto_rawDescGZIP() []byte {
	file_response_proto_rawDescOnce.Do(func() {
		file_response_proto_rawDescData = protoimpl.X.CompressGZIP(file_response_proto_rawDescData)
	})
	return file_response_proto_rawDescData
}

var file_response_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_response_proto_goTypes = []interface{}{
	(*Response)(nil),                // 0: proto.Response
	(*PlayerEnterGameResponse)(nil), // 1: proto.PlayerEnterGameResponse
	(*PlayerQuitGameResponse)(nil),  // 2: proto.PlayerQuitGameResponse
	(*HeroGetPropResponse)(nil),     // 3: proto.HeroGetPropResponse
	(*PlayerHeartBeatResponse)(nil), // 4: proto.PlayerHeartBeatResponse
	(ResponseCode)(0),               // 5: proto.ResponseCode
}
var file_response_proto_depIdxs = []int32{
	5, // 0: proto.Response.responseCode:type_name -> proto.ResponseCode
	1, // 1: proto.Response.playerEnterGameResponse:type_name -> proto.PlayerEnterGameResponse
	2, // 2: proto.Response.playerQuitGameResponse:type_name -> proto.PlayerQuitGameResponse
	3, // 3: proto.Response.heroGetPropResponse:type_name -> proto.HeroGetPropResponse
	4, // 4: proto.Response.playerHeartBeatResponse:type_name -> proto.PlayerHeartBeatResponse
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_response_proto_init() }
func file_response_proto_init() {
	if File_response_proto != nil {
		return
	}
	file_consts_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_response_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
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
		file_response_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PlayerEnterGameResponse); i {
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
		file_response_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PlayerQuitGameResponse); i {
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
		file_response_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HeroGetPropResponse); i {
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
		file_response_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PlayerHeartBeatResponse); i {
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
			RawDescriptor: file_response_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_response_proto_goTypes,
		DependencyIndexes: file_response_proto_depIdxs,
		MessageInfos:      file_response_proto_msgTypes,
	}.Build()
	File_response_proto = out.File
	file_response_proto_rawDesc = nil
	file_response_proto_goTypes = nil
	file_response_proto_depIdxs = nil
}
