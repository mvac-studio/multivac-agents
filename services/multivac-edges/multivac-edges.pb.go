// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v3.12.4
// source: multivac-edges.proto

package edges

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

type DeleteForwardEdgesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Source     *Vertex `protobuf:"bytes,1,opt,name=source,proto3" json:"source,omitempty"`
	TargetType string  `protobuf:"bytes,2,opt,name=targetType,proto3" json:"targetType,omitempty"`
}

func (x *DeleteForwardEdgesRequest) Reset() {
	*x = DeleteForwardEdgesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_multivac_edges_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteForwardEdgesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteForwardEdgesRequest) ProtoMessage() {}

func (x *DeleteForwardEdgesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_multivac_edges_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteForwardEdgesRequest.ProtoReflect.Descriptor instead.
func (*DeleteForwardEdgesRequest) Descriptor() ([]byte, []int) {
	return file_multivac_edges_proto_rawDescGZIP(), []int{0}
}

func (x *DeleteForwardEdgesRequest) GetSource() *Vertex {
	if x != nil {
		return x.Source
	}
	return nil
}

func (x *DeleteForwardEdgesRequest) GetTargetType() string {
	if x != nil {
		return x.TargetType
	}
	return ""
}

type GetForwardEdgesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Source     *Vertex `protobuf:"bytes,1,opt,name=source,proto3" json:"source,omitempty"`
	TargetType string  `protobuf:"bytes,2,opt,name=targetType,proto3" json:"targetType,omitempty"`
}

func (x *GetForwardEdgesRequest) Reset() {
	*x = GetForwardEdgesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_multivac_edges_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetForwardEdgesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetForwardEdgesRequest) ProtoMessage() {}

func (x *GetForwardEdgesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_multivac_edges_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetForwardEdgesRequest.ProtoReflect.Descriptor instead.
func (*GetForwardEdgesRequest) Descriptor() ([]byte, []int) {
	return file_multivac_edges_proto_rawDescGZIP(), []int{1}
}

func (x *GetForwardEdgesRequest) GetSource() *Vertex {
	if x != nil {
		return x.Source
	}
	return nil
}

func (x *GetForwardEdgesRequest) GetTargetType() string {
	if x != nil {
		return x.TargetType
	}
	return ""
}

type DeleteEdgeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DeleteEdgeRequest) Reset() {
	*x = DeleteEdgeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_multivac_edges_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteEdgeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteEdgeRequest) ProtoMessage() {}

func (x *DeleteEdgeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_multivac_edges_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteEdgeRequest.ProtoReflect.Descriptor instead.
func (*DeleteEdgeRequest) Descriptor() ([]byte, []int) {
	return file_multivac_edges_proto_rawDescGZIP(), []int{2}
}

func (x *DeleteEdgeRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type Vertex struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ref  string `protobuf:"bytes,1,opt,name=ref,proto3" json:"ref,omitempty"`
	Type string `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
}

func (x *Vertex) Reset() {
	*x = Vertex{}
	if protoimpl.UnsafeEnabled {
		mi := &file_multivac_edges_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Vertex) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Vertex) ProtoMessage() {}

func (x *Vertex) ProtoReflect() protoreflect.Message {
	mi := &file_multivac_edges_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Vertex.ProtoReflect.Descriptor instead.
func (*Vertex) Descriptor() ([]byte, []int) {
	return file_multivac_edges_proto_rawDescGZIP(), []int{3}
}

func (x *Vertex) GetRef() string {
	if x != nil {
		return x.Ref
	}
	return ""
}

func (x *Vertex) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

type EdgeCollection struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Edges []*Edge `protobuf:"bytes,1,rep,name=edges,proto3" json:"edges,omitempty"`
}

func (x *EdgeCollection) Reset() {
	*x = EdgeCollection{}
	if protoimpl.UnsafeEnabled {
		mi := &file_multivac_edges_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EdgeCollection) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EdgeCollection) ProtoMessage() {}

func (x *EdgeCollection) ProtoReflect() protoreflect.Message {
	mi := &file_multivac_edges_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EdgeCollection.ProtoReflect.Descriptor instead.
func (*EdgeCollection) Descriptor() ([]byte, []int) {
	return file_multivac_edges_proto_rawDescGZIP(), []int{4}
}

func (x *EdgeCollection) GetEdges() []*Edge {
	if x != nil {
		return x.Edges
	}
	return nil
}

type Edge struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Source  *Vertex `protobuf:"bytes,2,opt,name=source,proto3" json:"source,omitempty"`
	Target  *Vertex `protobuf:"bytes,3,opt,name=target,proto3" json:"target,omitempty"`
	Created int64   `protobuf:"varint,4,opt,name=created,proto3" json:"created,omitempty"`
	Updated int64   `protobuf:"varint,5,opt,name=updated,proto3" json:"updated,omitempty"`
}

func (x *Edge) Reset() {
	*x = Edge{}
	if protoimpl.UnsafeEnabled {
		mi := &file_multivac_edges_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Edge) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Edge) ProtoMessage() {}

func (x *Edge) ProtoReflect() protoreflect.Message {
	mi := &file_multivac_edges_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Edge.ProtoReflect.Descriptor instead.
func (*Edge) Descriptor() ([]byte, []int) {
	return file_multivac_edges_proto_rawDescGZIP(), []int{5}
}

func (x *Edge) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Edge) GetSource() *Vertex {
	if x != nil {
		return x.Source
	}
	return nil
}

func (x *Edge) GetTarget() *Vertex {
	if x != nil {
		return x.Target
	}
	return nil
}

func (x *Edge) GetCreated() int64 {
	if x != nil {
		return x.Created
	}
	return 0
}

func (x *Edge) GetUpdated() int64 {
	if x != nil {
		return x.Updated
	}
	return 0
}

var File_multivac_edges_proto protoreflect.FileDescriptor

var file_multivac_edges_proto_rawDesc = []byte{
	0x0a, 0x14, 0x6d, 0x75, 0x6c, 0x74, 0x69, 0x76, 0x61, 0x63, 0x2d, 0x65, 0x64, 0x67, 0x65, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x65, 0x64, 0x67, 0x65, 0x22, 0x61, 0x0a, 0x19,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x45, 0x64, 0x67,
	0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x24, 0x0a, 0x06, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x65, 0x64, 0x67, 0x65,
	0x2e, 0x56, 0x65, 0x72, 0x74, 0x65, 0x78, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12,
	0x1e, 0x0a, 0x0a, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x54, 0x79, 0x70, 0x65, 0x22,
	0x5e, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x46, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x45, 0x64, 0x67,
	0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x24, 0x0a, 0x06, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x65, 0x64, 0x67, 0x65,
	0x2e, 0x56, 0x65, 0x72, 0x74, 0x65, 0x78, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12,
	0x1e, 0x0a, 0x0a, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x54, 0x79, 0x70, 0x65, 0x22,
	0x23, 0x0a, 0x11, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x45, 0x64, 0x67, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x22, 0x2e, 0x0a, 0x06, 0x56, 0x65, 0x72, 0x74, 0x65, 0x78, 0x12, 0x10,
	0x0a, 0x03, 0x72, 0x65, 0x66, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x72, 0x65, 0x66,
	0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x22, 0x32, 0x0a, 0x0e, 0x45, 0x64, 0x67, 0x65, 0x43, 0x6f, 0x6c, 0x6c,
	0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x20, 0x0a, 0x05, 0x65, 0x64, 0x67, 0x65, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x65, 0x64, 0x67, 0x65, 0x2e, 0x45, 0x64, 0x67,
	0x65, 0x52, 0x05, 0x65, 0x64, 0x67, 0x65, 0x73, 0x22, 0x96, 0x01, 0x0a, 0x04, 0x45, 0x64, 0x67,
	0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x24, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0c, 0x2e, 0x65, 0x64, 0x67, 0x65, 0x2e, 0x56, 0x65, 0x72, 0x74, 0x65, 0x78, 0x52,
	0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x24, 0x0a, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x65, 0x64, 0x67, 0x65, 0x2e, 0x56,
	0x65, 0x72, 0x74, 0x65, 0x78, 0x52, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x12, 0x18, 0x0a,
	0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x75, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x64, 0x32, 0xf2, 0x02, 0x0a, 0x0b, 0x45, 0x64, 0x67, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x24, 0x0a, 0x08, 0x53, 0x61, 0x76, 0x65, 0x45, 0x64, 0x67, 0x65, 0x12, 0x0a, 0x2e,
	0x65, 0x64, 0x67, 0x65, 0x2e, 0x45, 0x64, 0x67, 0x65, 0x1a, 0x0a, 0x2e, 0x65, 0x64, 0x67, 0x65,
	0x2e, 0x45, 0x64, 0x67, 0x65, 0x22, 0x00, 0x12, 0x30, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x45, 0x64,
	0x67, 0x65, 0x73, 0x12, 0x0c, 0x2e, 0x65, 0x64, 0x67, 0x65, 0x2e, 0x56, 0x65, 0x72, 0x74, 0x65,
	0x78, 0x1a, 0x14, 0x2e, 0x65, 0x64, 0x67, 0x65, 0x2e, 0x45, 0x64, 0x67, 0x65, 0x43, 0x6f, 0x6c,
	0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x00, 0x12, 0x47, 0x0a, 0x0f, 0x47, 0x65, 0x74,
	0x46, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x45, 0x64, 0x67, 0x65, 0x73, 0x12, 0x1c, 0x2e, 0x65,
	0x64, 0x67, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x46, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x45, 0x64,
	0x67, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x65, 0x64, 0x67,
	0x65, 0x2e, 0x45, 0x64, 0x67, 0x65, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x22, 0x00, 0x12, 0x33, 0x0a, 0x0a, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x45, 0x64, 0x67, 0x65,
	0x12, 0x17, 0x2e, 0x65, 0x64, 0x67, 0x65, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x45, 0x64,
	0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0a, 0x2e, 0x65, 0x64, 0x67, 0x65,
	0x2e, 0x45, 0x64, 0x67, 0x65, 0x22, 0x00, 0x12, 0x3e, 0x0a, 0x16, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x41, 0x6c, 0x6c, 0x45, 0x64, 0x67, 0x65, 0x73, 0x42, 0x79, 0x53, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x12, 0x0c, 0x2e, 0x65, 0x64, 0x67, 0x65, 0x2e, 0x56, 0x65, 0x72, 0x74, 0x65, 0x78, 0x1a,
	0x14, 0x2e, 0x65, 0x64, 0x67, 0x65, 0x2e, 0x45, 0x64, 0x67, 0x65, 0x43, 0x6f, 0x6c, 0x6c, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x00, 0x12, 0x4d, 0x0a, 0x12, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x46, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x45, 0x64, 0x67, 0x65, 0x73, 0x12, 0x1f, 0x2e,
	0x65, 0x64, 0x67, 0x65, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x6f, 0x72, 0x77, 0x61,
	0x72, 0x64, 0x45, 0x64, 0x67, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14,
	0x2e, 0x65, 0x64, 0x67, 0x65, 0x2e, 0x45, 0x64, 0x67, 0x65, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x22, 0x00, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x2f, 0x65, 0x64, 0x67, 0x65,
	0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_multivac_edges_proto_rawDescOnce sync.Once
	file_multivac_edges_proto_rawDescData = file_multivac_edges_proto_rawDesc
)

func file_multivac_edges_proto_rawDescGZIP() []byte {
	file_multivac_edges_proto_rawDescOnce.Do(func() {
		file_multivac_edges_proto_rawDescData = protoimpl.X.CompressGZIP(file_multivac_edges_proto_rawDescData)
	})
	return file_multivac_edges_proto_rawDescData
}

var file_multivac_edges_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_multivac_edges_proto_goTypes = []interface{}{
	(*DeleteForwardEdgesRequest)(nil), // 0: edge.DeleteForwardEdgesRequest
	(*GetForwardEdgesRequest)(nil),    // 1: edge.GetForwardEdgesRequest
	(*DeleteEdgeRequest)(nil),         // 2: edge.DeleteEdgeRequest
	(*Vertex)(nil),                    // 3: edge.Vertex
	(*EdgeCollection)(nil),            // 4: edge.EdgeCollection
	(*Edge)(nil),                      // 5: edge.Edge
}
var file_multivac_edges_proto_depIdxs = []int32{
	3,  // 0: edge.DeleteForwardEdgesRequest.source:type_name -> edge.Vertex
	3,  // 1: edge.GetForwardEdgesRequest.source:type_name -> edge.Vertex
	5,  // 2: edge.EdgeCollection.edges:type_name -> edge.Edge
	3,  // 3: edge.Edge.source:type_name -> edge.Vertex
	3,  // 4: edge.Edge.target:type_name -> edge.Vertex
	5,  // 5: edge.EdgeService.SaveEdge:input_type -> edge.Edge
	3,  // 6: edge.EdgeService.GetEdges:input_type -> edge.Vertex
	1,  // 7: edge.EdgeService.GetForwardEdges:input_type -> edge.GetForwardEdgesRequest
	2,  // 8: edge.EdgeService.DeleteEdge:input_type -> edge.DeleteEdgeRequest
	3,  // 9: edge.EdgeService.DeleteAllEdgesBySource:input_type -> edge.Vertex
	0,  // 10: edge.EdgeService.DeleteForwardEdges:input_type -> edge.DeleteForwardEdgesRequest
	5,  // 11: edge.EdgeService.SaveEdge:output_type -> edge.Edge
	4,  // 12: edge.EdgeService.GetEdges:output_type -> edge.EdgeCollection
	4,  // 13: edge.EdgeService.GetForwardEdges:output_type -> edge.EdgeCollection
	5,  // 14: edge.EdgeService.DeleteEdge:output_type -> edge.Edge
	4,  // 15: edge.EdgeService.DeleteAllEdgesBySource:output_type -> edge.EdgeCollection
	4,  // 16: edge.EdgeService.DeleteForwardEdges:output_type -> edge.EdgeCollection
	11, // [11:17] is the sub-list for method output_type
	5,  // [5:11] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_multivac_edges_proto_init() }
func file_multivac_edges_proto_init() {
	if File_multivac_edges_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_multivac_edges_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteForwardEdgesRequest); i {
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
		file_multivac_edges_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetForwardEdgesRequest); i {
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
		file_multivac_edges_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteEdgeRequest); i {
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
		file_multivac_edges_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Vertex); i {
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
		file_multivac_edges_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EdgeCollection); i {
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
		file_multivac_edges_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Edge); i {
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
			RawDescriptor: file_multivac_edges_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_multivac_edges_proto_goTypes,
		DependencyIndexes: file_multivac_edges_proto_depIdxs,
		MessageInfos:      file_multivac_edges_proto_msgTypes,
	}.Build()
	File_multivac_edges_proto = out.File
	file_multivac_edges_proto_rawDesc = nil
	file_multivac_edges_proto_goTypes = nil
	file_multivac_edges_proto_depIdxs = nil
}
