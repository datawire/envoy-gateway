// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: internal/extension/proto/extension_service.proto

package proto

import (
	v32 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	v31 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	v33 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
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

type PostHTTPListenerTranslationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ExtensionContext *HTTPListenerExtensionContext `protobuf:"bytes,1,opt,name=extension_context,json=extensionContext,proto3" json:"extension_context,omitempty"`
	Listener         *v3.Listener                  `protobuf:"bytes,2,opt,name=listener,proto3" json:"listener,omitempty"`
	RouteTable       *v31.RouteConfiguration       `protobuf:"bytes,3,opt,name=route_table,json=routeTable,proto3" json:"route_table,omitempty"`
}

func (x *PostHTTPListenerTranslationRequest) Reset() {
	*x = PostHTTPListenerTranslationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_extension_proto_extension_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PostHTTPListenerTranslationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PostHTTPListenerTranslationRequest) ProtoMessage() {}

func (x *PostHTTPListenerTranslationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_extension_proto_extension_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PostHTTPListenerTranslationRequest.ProtoReflect.Descriptor instead.
func (*PostHTTPListenerTranslationRequest) Descriptor() ([]byte, []int) {
	return file_internal_extension_proto_extension_service_proto_rawDescGZIP(), []int{0}
}

func (x *PostHTTPListenerTranslationRequest) GetExtensionContext() *HTTPListenerExtensionContext {
	if x != nil {
		return x.ExtensionContext
	}
	return nil
}

func (x *PostHTTPListenerTranslationRequest) GetListener() *v3.Listener {
	if x != nil {
		return x.Listener
	}
	return nil
}

func (x *PostHTTPListenerTranslationRequest) GetRouteTable() *v31.RouteConfiguration {
	if x != nil {
		return x.RouteTable
	}
	return nil
}

type PostHTTPListenerTranslationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Clusters   []*v32.Cluster          `protobuf:"bytes,1,rep,name=clusters,proto3" json:"clusters,omitempty"`
	Listener   *v3.Listener            `protobuf:"bytes,2,opt,name=listener,proto3" json:"listener,omitempty"`
	RouteTable *v31.RouteConfiguration `protobuf:"bytes,3,opt,name=route_table,json=routeTable,proto3" json:"route_table,omitempty"`
	Secrets    []*v33.Secret           `protobuf:"bytes,4,rep,name=secrets,proto3" json:"secrets,omitempty"`
}

func (x *PostHTTPListenerTranslationResponse) Reset() {
	*x = PostHTTPListenerTranslationResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_extension_proto_extension_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PostHTTPListenerTranslationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PostHTTPListenerTranslationResponse) ProtoMessage() {}

func (x *PostHTTPListenerTranslationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_extension_proto_extension_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PostHTTPListenerTranslationResponse.ProtoReflect.Descriptor instead.
func (*PostHTTPListenerTranslationResponse) Descriptor() ([]byte, []int) {
	return file_internal_extension_proto_extension_service_proto_rawDescGZIP(), []int{1}
}

func (x *PostHTTPListenerTranslationResponse) GetClusters() []*v32.Cluster {
	if x != nil {
		return x.Clusters
	}
	return nil
}

func (x *PostHTTPListenerTranslationResponse) GetListener() *v3.Listener {
	if x != nil {
		return x.Listener
	}
	return nil
}

func (x *PostHTTPListenerTranslationResponse) GetRouteTable() *v31.RouteConfiguration {
	if x != nil {
		return x.RouteTable
	}
	return nil
}

func (x *PostHTTPListenerTranslationResponse) GetSecrets() []*v33.Secret {
	if x != nil {
		return x.Secrets
	}
	return nil
}

var File_internal_extension_proto_extension_service_proto protoreflect.FileDescriptor

var file_internal_extension_proto_extension_service_proto_rawDesc = []byte{
	0x0a, 0x30, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x65, 0x78, 0x74, 0x65, 0x6e,
	0x73, 0x69, 0x6f, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x78, 0x74, 0x65, 0x6e,
	0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x16, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79,
	0x2e, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x1a, 0x25, 0x65, 0x6e, 0x76, 0x6f,
	0x79, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72,
	0x2f, 0x76, 0x33, 0x2f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x27, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f,
	0x6c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x2f, 0x76, 0x33, 0x2f, 0x6c, 0x69, 0x73, 0x74,
	0x65, 0x6e, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x21, 0x65, 0x6e, 0x76, 0x6f,
	0x79, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x2f, 0x76,
	0x33, 0x2f, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x36, 0x65,
	0x6e, 0x76, 0x6f, 0x79, 0x2f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2f,
	0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x5f, 0x73, 0x6f, 0x63, 0x6b, 0x65, 0x74,
	0x73, 0x2f, 0x74, 0x6c, 0x73, 0x2f, 0x76, 0x33, 0x2f, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x30, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f,
	0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x93, 0x02, 0x0a, 0x22, 0x50, 0x6f, 0x73, 0x74,
	0x48, 0x54, 0x54, 0x50, 0x4c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x61,
	0x0a, 0x11, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x78, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x34, 0x2e, 0x65, 0x6e, 0x76, 0x6f,
	0x79, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69,
	0x6f, 0x6e, 0x2e, 0x48, 0x54, 0x54, 0x50, 0x4c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x45,
	0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x52,
	0x10, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78,
	0x74, 0x12, 0x3e, 0x0a, 0x08, 0x6c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2e, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x2e, 0x6c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x2e, 0x76, 0x33, 0x2e, 0x4c,
	0x69, 0x73, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x52, 0x08, 0x6c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x65,
	0x72, 0x12, 0x4a, 0x0a, 0x0b, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x5f, 0x74, 0x61, 0x62, 0x6c, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2e, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x2e, 0x76, 0x33, 0x2e, 0x52,
	0x6f, 0x75, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x0a, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x22, 0xbc, 0x02,
	0x0a, 0x23, 0x50, 0x6f, 0x73, 0x74, 0x48, 0x54, 0x54, 0x50, 0x4c, 0x69, 0x73, 0x74, 0x65, 0x6e,
	0x65, 0x72, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3c, 0x0a, 0x08, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2e,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x76,
	0x33, 0x2e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x08, 0x63, 0x6c, 0x75, 0x73, 0x74,
	0x65, 0x72, 0x73, 0x12, 0x3e, 0x0a, 0x08, 0x6c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2e, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x2e, 0x6c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x2e, 0x76, 0x33,
	0x2e, 0x4c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x52, 0x08, 0x6c, 0x69, 0x73, 0x74, 0x65,
	0x6e, 0x65, 0x72, 0x12, 0x4a, 0x0a, 0x0b, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x5f, 0x74, 0x61, 0x62,
	0x6c, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79,
	0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x2e, 0x76, 0x33,
	0x2e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x12,
	0x4b, 0x0a, 0x07, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x31, 0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x5f, 0x73, 0x6f,
	0x63, 0x6b, 0x65, 0x74, 0x73, 0x2e, 0x74, 0x6c, 0x73, 0x2e, 0x76, 0x33, 0x2e, 0x53, 0x65, 0x63,
	0x72, 0x65, 0x74, 0x52, 0x07, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x73, 0x32, 0xb2, 0x01, 0x0a,
	0x15, 0x45, 0x6e, 0x76, 0x6f, 0x79, 0x47, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x45, 0x78, 0x74,
	0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x98, 0x01, 0x0a, 0x1b, 0x50, 0x6f, 0x73, 0x74, 0x48,
	0x54, 0x54, 0x50, 0x4c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x3a, 0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x67, 0x61,
	0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x2e,
	0x50, 0x6f, 0x73, 0x74, 0x48, 0x54, 0x54, 0x50, 0x4c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x65, 0x72,
	0x54, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x3b, 0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61,
	0x79, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x50, 0x6f, 0x73, 0x74,
	0x48, 0x54, 0x54, 0x50, 0x4c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x42, 0x1a, 0x5a, 0x18, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x65, 0x78,
	0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_extension_proto_extension_service_proto_rawDescOnce sync.Once
	file_internal_extension_proto_extension_service_proto_rawDescData = file_internal_extension_proto_extension_service_proto_rawDesc
)

func file_internal_extension_proto_extension_service_proto_rawDescGZIP() []byte {
	file_internal_extension_proto_extension_service_proto_rawDescOnce.Do(func() {
		file_internal_extension_proto_extension_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_extension_proto_extension_service_proto_rawDescData)
	})
	return file_internal_extension_proto_extension_service_proto_rawDescData
}

var file_internal_extension_proto_extension_service_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_internal_extension_proto_extension_service_proto_goTypes = []interface{}{
	(*PostHTTPListenerTranslationRequest)(nil),  // 0: envoygateway.extension.PostHTTPListenerTranslationRequest
	(*PostHTTPListenerTranslationResponse)(nil), // 1: envoygateway.extension.PostHTTPListenerTranslationResponse
	(*HTTPListenerExtensionContext)(nil),        // 2: envoygateway.extension.HTTPListenerExtensionContext
	(*v3.Listener)(nil),                         // 3: envoy.config.listener.v3.Listener
	(*v31.RouteConfiguration)(nil),              // 4: envoy.config.route.v3.RouteConfiguration
	(*v32.Cluster)(nil),                         // 5: envoy.config.cluster.v3.Cluster
	(*v33.Secret)(nil),                          // 6: envoy.extensions.transport_sockets.tls.v3.Secret
}
var file_internal_extension_proto_extension_service_proto_depIdxs = []int32{
	2, // 0: envoygateway.extension.PostHTTPListenerTranslationRequest.extension_context:type_name -> envoygateway.extension.HTTPListenerExtensionContext
	3, // 1: envoygateway.extension.PostHTTPListenerTranslationRequest.listener:type_name -> envoy.config.listener.v3.Listener
	4, // 2: envoygateway.extension.PostHTTPListenerTranslationRequest.route_table:type_name -> envoy.config.route.v3.RouteConfiguration
	5, // 3: envoygateway.extension.PostHTTPListenerTranslationResponse.clusters:type_name -> envoy.config.cluster.v3.Cluster
	3, // 4: envoygateway.extension.PostHTTPListenerTranslationResponse.listener:type_name -> envoy.config.listener.v3.Listener
	4, // 5: envoygateway.extension.PostHTTPListenerTranslationResponse.route_table:type_name -> envoy.config.route.v3.RouteConfiguration
	6, // 6: envoygateway.extension.PostHTTPListenerTranslationResponse.secrets:type_name -> envoy.extensions.transport_sockets.tls.v3.Secret
	0, // 7: envoygateway.extension.EnvoyGatewayExtension.PostHTTPListenerTranslation:input_type -> envoygateway.extension.PostHTTPListenerTranslationRequest
	1, // 8: envoygateway.extension.EnvoyGatewayExtension.PostHTTPListenerTranslation:output_type -> envoygateway.extension.PostHTTPListenerTranslationResponse
	8, // [8:9] is the sub-list for method output_type
	7, // [7:8] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_internal_extension_proto_extension_service_proto_init() }
func file_internal_extension_proto_extension_service_proto_init() {
	if File_internal_extension_proto_extension_service_proto != nil {
		return
	}
	file_internal_extension_proto_extension_context_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_internal_extension_proto_extension_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PostHTTPListenerTranslationRequest); i {
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
		file_internal_extension_proto_extension_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PostHTTPListenerTranslationResponse); i {
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
			RawDescriptor: file_internal_extension_proto_extension_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_extension_proto_extension_service_proto_goTypes,
		DependencyIndexes: file_internal_extension_proto_extension_service_proto_depIdxs,
		MessageInfos:      file_internal_extension_proto_extension_service_proto_msgTypes,
	}.Build()
	File_internal_extension_proto_extension_service_proto = out.File
	file_internal_extension_proto_extension_service_proto_rawDesc = nil
	file_internal_extension_proto_extension_service_proto_goTypes = nil
	file_internal_extension_proto_extension_service_proto_depIdxs = nil
}
