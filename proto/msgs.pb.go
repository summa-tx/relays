// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.15.3
// source: msgs.proto

package proto

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// MsgIngestHeaderChain defines a IngestHeaderChain message
type MsgIngestHeaderChain struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Signer  string           `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer,omitempty"`
	Headers []*BitcoinHeader `protobuf:"bytes,2,rep,name=headers,proto3" json:"headers,omitempty"`
}

func (x *MsgIngestHeaderChain) Reset() {
	*x = MsgIngestHeaderChain{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msgs_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MsgIngestHeaderChain) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MsgIngestHeaderChain) ProtoMessage() {}

func (x *MsgIngestHeaderChain) ProtoReflect() protoreflect.Message {
	mi := &file_msgs_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MsgIngestHeaderChain.ProtoReflect.Descriptor instead.
func (*MsgIngestHeaderChain) Descriptor() ([]byte, []int) {
	return file_msgs_proto_rawDescGZIP(), []int{0}
}

func (x *MsgIngestHeaderChain) GetSigner() string {
	if x != nil {
		return x.Signer
	}
	return ""
}

func (x *MsgIngestHeaderChain) GetHeaders() []*BitcoinHeader {
	if x != nil {
		return x.Headers
	}
	return nil
}

// MsgIngestDifficultyChange defines a IngestDifficultyChange message
type MsgIngestDifficultyChange struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Signer  string           `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer,omitempty"`
	Start   []byte           `protobuf:"bytes,2,opt,name=start,proto3" json:"start,omitempty"`
	Headers []*BitcoinHeader `protobuf:"bytes,3,rep,name=headers,proto3" json:"headers,omitempty"`
}

func (x *MsgIngestDifficultyChange) Reset() {
	*x = MsgIngestDifficultyChange{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msgs_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MsgIngestDifficultyChange) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MsgIngestDifficultyChange) ProtoMessage() {}

func (x *MsgIngestDifficultyChange) ProtoReflect() protoreflect.Message {
	mi := &file_msgs_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MsgIngestDifficultyChange.ProtoReflect.Descriptor instead.
func (*MsgIngestDifficultyChange) Descriptor() ([]byte, []int) {
	return file_msgs_proto_rawDescGZIP(), []int{1}
}

func (x *MsgIngestDifficultyChange) GetSigner() string {
	if x != nil {
		return x.Signer
	}
	return ""
}

func (x *MsgIngestDifficultyChange) GetStart() []byte {
	if x != nil {
		return x.Start
	}
	return nil
}

func (x *MsgIngestDifficultyChange) GetHeaders() []*BitcoinHeader {
	if x != nil {
		return x.Headers
	}
	return nil
}

// MsgMarkNewHeaviest defines a MarkNewHeaviest message
type MsgMarkNewHeaviest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Signer      string `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer,omitempty"`
	Ancestor    []byte `protobuf:"bytes,2,opt,name=ancestor,proto3" json:"ancestor,omitempty"`
	CurrentBest []byte `protobuf:"bytes,3,opt,name=current_best,json=currentBest,proto3" json:"current_best,omitempty"`
	NewBest     []byte `protobuf:"bytes,4,opt,name=new_best,json=newBest,proto3" json:"new_best,omitempty"`
	Limit       uint32 `protobuf:"varint,5,opt,name=limit,proto3" json:"limit,omitempty"`
}

func (x *MsgMarkNewHeaviest) Reset() {
	*x = MsgMarkNewHeaviest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msgs_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MsgMarkNewHeaviest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MsgMarkNewHeaviest) ProtoMessage() {}

func (x *MsgMarkNewHeaviest) ProtoReflect() protoreflect.Message {
	mi := &file_msgs_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MsgMarkNewHeaviest.ProtoReflect.Descriptor instead.
func (*MsgMarkNewHeaviest) Descriptor() ([]byte, []int) {
	return file_msgs_proto_rawDescGZIP(), []int{2}
}

func (x *MsgMarkNewHeaviest) GetSigner() string {
	if x != nil {
		return x.Signer
	}
	return ""
}

func (x *MsgMarkNewHeaviest) GetAncestor() []byte {
	if x != nil {
		return x.Ancestor
	}
	return nil
}

func (x *MsgMarkNewHeaviest) GetCurrentBest() []byte {
	if x != nil {
		return x.CurrentBest
	}
	return nil
}

func (x *MsgMarkNewHeaviest) GetNewBest() []byte {
	if x != nil {
		return x.NewBest
	}
	return nil
}

func (x *MsgMarkNewHeaviest) GetLimit() uint32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

// MsgNewRequest defines a NewRequest message
type MsgNewRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Signer    string `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer,omitempty"`
	Spends    []byte `protobuf:"bytes,2,opt,name=spends,proto3" json:"spends,omitempty"`
	Pays      []byte `protobuf:"bytes,3,opt,name=pays,proto3" json:"pays,omitempty"`
	PaysValue uint64 `protobuf:"varint,4,opt,name=pays_value,json=paysValue,proto3" json:"pays_value,omitempty"`
	NumConfs  uint32 `protobuf:"varint,5,opt,name=num_confs,json=numConfs,proto3" json:"num_confs,omitempty"`
	Origin    Origin `protobuf:"varint,6,opt,name=origin,proto3,enum=relayproto.Origin" json:"origin,omitempty"`
	Action    []byte `protobuf:"bytes,7,opt,name=action,proto3" json:"action,omitempty"`
}

func (x *MsgNewRequest) Reset() {
	*x = MsgNewRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msgs_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MsgNewRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MsgNewRequest) ProtoMessage() {}

func (x *MsgNewRequest) ProtoReflect() protoreflect.Message {
	mi := &file_msgs_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MsgNewRequest.ProtoReflect.Descriptor instead.
func (*MsgNewRequest) Descriptor() ([]byte, []int) {
	return file_msgs_proto_rawDescGZIP(), []int{3}
}

func (x *MsgNewRequest) GetSigner() string {
	if x != nil {
		return x.Signer
	}
	return ""
}

func (x *MsgNewRequest) GetSpends() []byte {
	if x != nil {
		return x.Spends
	}
	return nil
}

func (x *MsgNewRequest) GetPays() []byte {
	if x != nil {
		return x.Pays
	}
	return nil
}

func (x *MsgNewRequest) GetPaysValue() uint64 {
	if x != nil {
		return x.PaysValue
	}
	return 0
}

func (x *MsgNewRequest) GetNumConfs() uint32 {
	if x != nil {
		return x.NumConfs
	}
	return 0
}

func (x *MsgNewRequest) GetOrigin() Origin {
	if x != nil {
		return x.Origin
	}
	return Origin_LOCAL
}

func (x *MsgNewRequest) GetAction() []byte {
	if x != nil {
		return x.Action
	}
	return nil
}

// MsgProvideProof defines a NewRequest message
type MsgProvideProof struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Signer string          `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer,omitempty"`
	Filled *FilledRequests `protobuf:"bytes,2,opt,name=filled,proto3" json:"filled,omitempty"`
}

func (x *MsgProvideProof) Reset() {
	*x = MsgProvideProof{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msgs_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MsgProvideProof) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MsgProvideProof) ProtoMessage() {}

func (x *MsgProvideProof) ProtoReflect() protoreflect.Message {
	mi := &file_msgs_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MsgProvideProof.ProtoReflect.Descriptor instead.
func (*MsgProvideProof) Descriptor() ([]byte, []int) {
	return file_msgs_proto_rawDescGZIP(), []int{4}
}

func (x *MsgProvideProof) GetSigner() string {
	if x != nil {
		return x.Signer
	}
	return ""
}

func (x *MsgProvideProof) GetFilled() *FilledRequests {
	if x != nil {
		return x.Filled
	}
	return nil
}

var File_msgs_proto protoreflect.FileDescriptor

var file_msgs_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x6d, 0x73, 0x67, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x72, 0x65,
	0x6c, 0x61, 0x79, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0c, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x63, 0x0a, 0x14, 0x4d, 0x73, 0x67, 0x49, 0x6e, 0x67,
	0x65, 0x73, 0x74, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x12, 0x16,
	0x0a, 0x06, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x73, 0x69, 0x67, 0x6e, 0x65, 0x72, 0x12, 0x33, 0x0a, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x72, 0x65, 0x6c, 0x61, 0x79, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x42, 0x69, 0x74, 0x63, 0x6f, 0x69, 0x6e, 0x48, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x52, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x22, 0x7e, 0x0a, 0x19, 0x4d,
	0x73, 0x67, 0x49, 0x6e, 0x67, 0x65, 0x73, 0x74, 0x44, 0x69, 0x66, 0x66, 0x69, 0x63, 0x75, 0x6c,
	0x74, 0x79, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x69, 0x67, 0x6e,
	0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x72,
	0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x12, 0x33, 0x0a, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x72, 0x65, 0x6c, 0x61, 0x79, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x42, 0x69, 0x74, 0x63, 0x6f, 0x69, 0x6e, 0x48, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x52, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x22, 0x9c, 0x01, 0x0a, 0x12,
	0x4d, 0x73, 0x67, 0x4d, 0x61, 0x72, 0x6b, 0x4e, 0x65, 0x77, 0x48, 0x65, 0x61, 0x76, 0x69, 0x65,
	0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x61, 0x6e,
	0x63, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x61, 0x6e,
	0x63, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e,
	0x74, 0x5f, 0x62, 0x65, 0x73, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b, 0x63, 0x75,
	0x72, 0x72, 0x65, 0x6e, 0x74, 0x42, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x6e, 0x65, 0x77,
	0x5f, 0x62, 0x65, 0x73, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x6e, 0x65, 0x77,
	0x42, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x22, 0xd3, 0x01, 0x0a, 0x0d, 0x4d,
	0x73, 0x67, 0x4e, 0x65, 0x77, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06,
	0x73, 0x69, 0x67, 0x6e, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x69,
	0x67, 0x6e, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x70, 0x65, 0x6e, 0x64, 0x73, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x73, 0x70, 0x65, 0x6e, 0x64, 0x73, 0x12, 0x12, 0x0a, 0x04,
	0x70, 0x61, 0x79, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x70, 0x61, 0x79, 0x73,
	0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x61, 0x79, 0x73, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x70, 0x61, 0x79, 0x73, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12,
	0x1b, 0x0a, 0x09, 0x6e, 0x75, 0x6d, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x73, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x08, 0x6e, 0x75, 0x6d, 0x43, 0x6f, 0x6e, 0x66, 0x73, 0x12, 0x2a, 0x0a, 0x06,
	0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x72,
	0x65, 0x6c, 0x61, 0x79, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e,
	0x52, 0x06, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x22, 0x5d, 0x0a, 0x0f, 0x4d, 0x73, 0x67, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x50, 0x72,
	0x6f, 0x6f, 0x66, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x72, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x72, 0x12, 0x32, 0x0a, 0x06, 0x66,
	0x69, 0x6c, 0x6c, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x72, 0x65,
	0x6c, 0x61, 0x79, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x46, 0x69, 0x6c, 0x6c, 0x65, 0x64, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x6c, 0x65, 0x64, 0x42,
	0x22, 0x5a, 0x20, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x75,
	0x6d, 0x6d, 0x61, 0x2d, 0x74, 0x78, 0x2f, 0x72, 0x65, 0x6c, 0x61, 0x79, 0x73, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_msgs_proto_rawDescOnce sync.Once
	file_msgs_proto_rawDescData = file_msgs_proto_rawDesc
)

func file_msgs_proto_rawDescGZIP() []byte {
	file_msgs_proto_rawDescOnce.Do(func() {
		file_msgs_proto_rawDescData = protoimpl.X.CompressGZIP(file_msgs_proto_rawDescData)
	})
	return file_msgs_proto_rawDescData
}

var file_msgs_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_msgs_proto_goTypes = []interface{}{
	(*MsgIngestHeaderChain)(nil),      // 0: relayproto.MsgIngestHeaderChain
	(*MsgIngestDifficultyChange)(nil), // 1: relayproto.MsgIngestDifficultyChange
	(*MsgMarkNewHeaviest)(nil),        // 2: relayproto.MsgMarkNewHeaviest
	(*MsgNewRequest)(nil),             // 3: relayproto.MsgNewRequest
	(*MsgProvideProof)(nil),           // 4: relayproto.MsgProvideProof
	(*BitcoinHeader)(nil),             // 5: relayproto.BitcoinHeader
	(Origin)(0),                       // 6: relayproto.Origin
	(*FilledRequests)(nil),            // 7: relayproto.FilledRequests
}
var file_msgs_proto_depIdxs = []int32{
	5, // 0: relayproto.MsgIngestHeaderChain.headers:type_name -> relayproto.BitcoinHeader
	5, // 1: relayproto.MsgIngestDifficultyChange.headers:type_name -> relayproto.BitcoinHeader
	6, // 2: relayproto.MsgNewRequest.origin:type_name -> relayproto.Origin
	7, // 3: relayproto.MsgProvideProof.filled:type_name -> relayproto.FilledRequests
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_msgs_proto_init() }
func file_msgs_proto_init() {
	if File_msgs_proto != nil {
		return
	}
	file_shared_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_msgs_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MsgIngestHeaderChain); i {
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
		file_msgs_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MsgIngestDifficultyChange); i {
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
		file_msgs_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MsgMarkNewHeaviest); i {
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
		file_msgs_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MsgNewRequest); i {
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
		file_msgs_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MsgProvideProof); i {
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
			RawDescriptor: file_msgs_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_msgs_proto_goTypes,
		DependencyIndexes: file_msgs_proto_depIdxs,
		MessageInfos:      file_msgs_proto_msgTypes,
	}.Build()
	File_msgs_proto = out.File
	file_msgs_proto_rawDesc = nil
	file_msgs_proto_goTypes = nil
	file_msgs_proto_depIdxs = nil
}
