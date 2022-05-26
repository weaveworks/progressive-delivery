// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: api/prog/prog.proto

package api

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
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

type GetVersionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetVersionRequest) Reset() {
	*x = GetVersionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_prog_prog_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetVersionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetVersionRequest) ProtoMessage() {}

func (x *GetVersionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_prog_prog_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetVersionRequest.ProtoReflect.Descriptor instead.
func (*GetVersionRequest) Descriptor() ([]byte, []int) {
	return file_api_prog_prog_proto_rawDescGZIP(), []int{0}
}

type GetVersionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version string `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *GetVersionResponse) Reset() {
	*x = GetVersionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_prog_prog_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetVersionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetVersionResponse) ProtoMessage() {}

func (x *GetVersionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_prog_prog_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetVersionResponse.ProtoReflect.Descriptor instead.
func (*GetVersionResponse) Descriptor() ([]byte, []int) {
	return file_api_prog_prog_proto_rawDescGZIP(), []int{1}
}

func (x *GetVersionResponse) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

type Gate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name      string            `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Namespace string            `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	Phase     string            `protobuf:"bytes,3,opt,name=phase,proto3" json:"phase,omitempty"`
	Metadata  map[string]string `protobuf:"bytes,4,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	HookType  string            `protobuf:"bytes,5,opt,name=hook_type,json=hookType,proto3" json:"hook_type,omitempty"`
}

func (x *Gate) Reset() {
	*x = Gate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_prog_prog_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Gate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Gate) ProtoMessage() {}

func (x *Gate) ProtoReflect() protoreflect.Message {
	mi := &file_api_prog_prog_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Gate.ProtoReflect.Descriptor instead.
func (*Gate) Descriptor() ([]byte, []int) {
	return file_api_prog_prog_proto_rawDescGZIP(), []int{2}
}

func (x *Gate) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Gate) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *Gate) GetPhase() string {
	if x != nil {
		return x.Phase
	}
	return ""
}

func (x *Gate) GetMetadata() map[string]string {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *Gate) GetHookType() string {
	if x != nil {
		return x.HookType
	}
	return ""
}

type GateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name      string            `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Namespace string            `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	Phase     string            `protobuf:"bytes,3,opt,name=phase,proto3" json:"phase,omitempty"`
	Metadata  map[string]string `protobuf:"bytes,4,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	HookType  string            `protobuf:"bytes,5,opt,name=hook_type,json=hookType,proto3" json:"hook_type,omitempty"`
}

func (x *GateRequest) Reset() {
	*x = GateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_prog_prog_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GateRequest) ProtoMessage() {}

func (x *GateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_prog_prog_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GateRequest.ProtoReflect.Descriptor instead.
func (*GateRequest) Descriptor() ([]byte, []int) {
	return file_api_prog_prog_proto_rawDescGZIP(), []int{3}
}

func (x *GateRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *GateRequest) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *GateRequest) GetPhase() string {
	if x != nil {
		return x.Phase
	}
	return ""
}

func (x *GateRequest) GetMetadata() map[string]string {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *GateRequest) GetHookType() string {
	if x != nil {
		return x.HookType
	}
	return ""
}

type GateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GateResponse) Reset() {
	*x = GateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_prog_prog_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GateResponse) ProtoMessage() {}

func (x *GateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_prog_prog_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GateResponse.ProtoReflect.Descriptor instead.
func (*GateResponse) Descriptor() ([]byte, []int) {
	return file_api_prog_prog_proto_rawDescGZIP(), []int{4}
}

type GateListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GateListRequest) Reset() {
	*x = GateListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_prog_prog_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GateListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GateListRequest) ProtoMessage() {}

func (x *GateListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_prog_prog_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GateListRequest.ProtoReflect.Descriptor instead.
func (*GateListRequest) Descriptor() ([]byte, []int) {
	return file_api_prog_prog_proto_rawDescGZIP(), []int{5}
}

type GateListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Gates []*Gate `protobuf:"bytes,1,rep,name=gates,proto3" json:"gates,omitempty"`
}

func (x *GateListResponse) Reset() {
	*x = GateListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_prog_prog_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GateListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GateListResponse) ProtoMessage() {}

func (x *GateListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_prog_prog_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GateListResponse.ProtoReflect.Descriptor instead.
func (*GateListResponse) Descriptor() ([]byte, []int) {
	return file_api_prog_prog_proto_rawDescGZIP(), []int{6}
}

func (x *GateListResponse) GetGates() []*Gate {
	if x != nil {
		return x.Gates
	}
	return nil
}

type GateProceedRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *GateProceedRequest) Reset() {
	*x = GateProceedRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_prog_prog_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GateProceedRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GateProceedRequest) ProtoMessage() {}

func (x *GateProceedRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_prog_prog_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GateProceedRequest.ProtoReflect.Descriptor instead.
func (*GateProceedRequest) Descriptor() ([]byte, []int) {
	return file_api_prog_prog_proto_rawDescGZIP(), []int{7}
}

func (x *GateProceedRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type GateProceedResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GateProceedResponse) Reset() {
	*x = GateProceedResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_prog_prog_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GateProceedResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GateProceedResponse) ProtoMessage() {}

func (x *GateProceedResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_prog_prog_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GateProceedResponse.ProtoReflect.Descriptor instead.
func (*GateProceedResponse) Descriptor() ([]byte, []int) {
	return file_api_prog_prog_proto_rawDescGZIP(), []int{8}
}

var File_api_prog_prog_proto protoreflect.FileDescriptor

var file_api_prog_prog_proto_rawDesc = []byte{
	0x0a, 0x13, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x67, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d,
	0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x13, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x2e, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0xd9, 0x01, 0x0a, 0x04, 0x47, 0x61, 0x74,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61,
	0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70,
	0x61, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x68, 0x61, 0x73, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x70, 0x68, 0x61, 0x73, 0x65, 0x12, 0x2f, 0x0a, 0x08, 0x6d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x47, 0x61,
	0x74, 0x65, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x1b, 0x0a, 0x09, 0x68, 0x6f,
	0x6f, 0x6b, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x68,
	0x6f, 0x6f, 0x6b, 0x54, 0x79, 0x70, 0x65, 0x1a, 0x3b, 0x0a, 0x0d, 0x4d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x3a, 0x02, 0x38, 0x01, 0x22, 0xe7, 0x01, 0x0a, 0x0b, 0x47, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65,
	0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6e, 0x61, 0x6d,
	0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x68, 0x61, 0x73, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x70, 0x68, 0x61, 0x73, 0x65, 0x12, 0x36, 0x0a, 0x08,
	0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x47, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x12, 0x1b, 0x0a, 0x09, 0x68, 0x6f, 0x6f, 0x6b, 0x5f, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x68, 0x6f, 0x6f, 0x6b, 0x54, 0x79, 0x70,
	0x65, 0x1a, 0x3b, 0x0a, 0x0d, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x0e,
	0x0a, 0x0c, 0x47, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x11,
	0x0a, 0x0f, 0x47, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x22, 0x2f, 0x0a, 0x10, 0x47, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1b, 0x0a, 0x05, 0x67, 0x61, 0x74, 0x65, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x47, 0x61, 0x74, 0x65, 0x52, 0x05, 0x67, 0x61, 0x74,
	0x65, 0x73, 0x22, 0x28, 0x0a, 0x12, 0x47, 0x61, 0x74, 0x65, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x65,
	0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x15, 0x0a, 0x13,
	0x47, 0x61, 0x74, 0x65, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x65, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x32, 0xd5, 0x02, 0x0a, 0x1a, 0x50, 0x72, 0x6f, 0x67, 0x72, 0x65, 0x73, 0x73,
	0x69, 0x76, 0x65, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x79, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x4a, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x12, 0x12, 0x2e, 0x47, 0x65, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x47, 0x65, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x13, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x0d, 0x12, 0x0b, 0x2f, 0x76, 0x31, 0x2f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x41,
	0x0a, 0x04, 0x47, 0x61, 0x74, 0x65, 0x12, 0x0c, 0x2e, 0x47, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x0d, 0x2e, 0x47, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x1c, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x16, 0x22, 0x11, 0x2f, 0x76, 0x31,
	0x2f, 0x63, 0x61, 0x6e, 0x61, 0x72, 0x69, 0x65, 0x73, 0x2f, 0x67, 0x61, 0x74, 0x65, 0x3a, 0x01,
	0x2a, 0x12, 0x4b, 0x0a, 0x08, 0x47, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x10, 0x2e,
	0x47, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x11, 0x2e, 0x47, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x1a, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x14, 0x12, 0x12, 0x2f, 0x76, 0x31, 0x2f,
	0x63, 0x61, 0x6e, 0x61, 0x72, 0x69, 0x65, 0x73, 0x2f, 0x67, 0x61, 0x74, 0x65, 0x73, 0x12, 0x5b,
	0x0a, 0x0b, 0x47, 0x61, 0x74, 0x65, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x65, 0x64, 0x12, 0x13, 0x2e,
	0x47, 0x61, 0x74, 0x65, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x65, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x11, 0x2e, 0x47, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x24, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1e, 0x22, 0x19, 0x2f,
	0x76, 0x31, 0x2f, 0x63, 0x61, 0x6e, 0x61, 0x72, 0x69, 0x65, 0x73, 0x2f, 0x67, 0x61, 0x74, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x65, 0x64, 0x3a, 0x01, 0x2a, 0x42, 0x30, 0x5a, 0x2e, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x77, 0x65, 0x61, 0x76, 0x65, 0x77,
	0x6f, 0x72, 0x6b, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x67, 0x72, 0x65, 0x73, 0x73, 0x69, 0x76, 0x65,
	0x2d, 0x64, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x79, 0x2f, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_prog_prog_proto_rawDescOnce sync.Once
	file_api_prog_prog_proto_rawDescData = file_api_prog_prog_proto_rawDesc
)

func file_api_prog_prog_proto_rawDescGZIP() []byte {
	file_api_prog_prog_proto_rawDescOnce.Do(func() {
		file_api_prog_prog_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_prog_prog_proto_rawDescData)
	})
	return file_api_prog_prog_proto_rawDescData
}

var file_api_prog_prog_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_api_prog_prog_proto_goTypes = []interface{}{
	(*GetVersionRequest)(nil),   // 0: GetVersionRequest
	(*GetVersionResponse)(nil),  // 1: GetVersionResponse
	(*Gate)(nil),                // 2: Gate
	(*GateRequest)(nil),         // 3: GateRequest
	(*GateResponse)(nil),        // 4: GateResponse
	(*GateListRequest)(nil),     // 5: GateListRequest
	(*GateListResponse)(nil),    // 6: GateListResponse
	(*GateProceedRequest)(nil),  // 7: GateProceedRequest
	(*GateProceedResponse)(nil), // 8: GateProceedResponse
	nil,                         // 9: Gate.MetadataEntry
	nil,                         // 10: GateRequest.MetadataEntry
}
var file_api_prog_prog_proto_depIdxs = []int32{
	9,  // 0: Gate.metadata:type_name -> Gate.MetadataEntry
	10, // 1: GateRequest.metadata:type_name -> GateRequest.MetadataEntry
	2,  // 2: GateListResponse.gates:type_name -> Gate
	0,  // 3: ProgressiveDeliveryService.GetVersion:input_type -> GetVersionRequest
	3,  // 4: ProgressiveDeliveryService.Gate:input_type -> GateRequest
	5,  // 5: ProgressiveDeliveryService.GateList:input_type -> GateListRequest
	7,  // 6: ProgressiveDeliveryService.GateProceed:input_type -> GateProceedRequest
	1,  // 7: ProgressiveDeliveryService.GetVersion:output_type -> GetVersionResponse
	4,  // 8: ProgressiveDeliveryService.Gate:output_type -> GateResponse
	6,  // 9: ProgressiveDeliveryService.GateList:output_type -> GateListResponse
	6,  // 10: ProgressiveDeliveryService.GateProceed:output_type -> GateListResponse
	7,  // [7:11] is the sub-list for method output_type
	3,  // [3:7] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_api_prog_prog_proto_init() }
func file_api_prog_prog_proto_init() {
	if File_api_prog_prog_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_prog_prog_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetVersionRequest); i {
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
		file_api_prog_prog_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetVersionResponse); i {
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
		file_api_prog_prog_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Gate); i {
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
		file_api_prog_prog_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GateRequest); i {
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
		file_api_prog_prog_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GateResponse); i {
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
		file_api_prog_prog_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GateListRequest); i {
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
		file_api_prog_prog_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GateListResponse); i {
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
		file_api_prog_prog_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GateProceedRequest); i {
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
		file_api_prog_prog_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GateProceedResponse); i {
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
			RawDescriptor: file_api_prog_prog_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_prog_prog_proto_goTypes,
		DependencyIndexes: file_api_prog_prog_proto_depIdxs,
		MessageInfos:      file_api_prog_prog_proto_msgTypes,
	}.Build()
	File_api_prog_prog_proto = out.File
	file_api_prog_prog_proto_rawDesc = nil
	file_api_prog_prog_proto_goTypes = nil
	file_api_prog_prog_proto_depIdxs = nil
}
