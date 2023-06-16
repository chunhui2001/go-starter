// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.14.0
// source: proto/txns.proto

package wallet

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

type Txns struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChainId              string `protobuf:"bytes,1,opt,name=chainId,proto3" json:"chainId,omitempty"`
	BlockHash            string `protobuf:"bytes,2,opt,name=blockHash,proto3" json:"blockHash,omitempty"`
	BlockNumber          uint32 `protobuf:"varint,3,opt,name=blockNumber,proto3" json:"blockNumber,omitempty"`
	From                 string `protobuf:"bytes,4,opt,name=from,proto3" json:"from,omitempty"`
	Gas                  uint64 `protobuf:"varint,5,opt,name=gas,proto3" json:"gas,omitempty"`
	GasPrice             uint64 `protobuf:"varint,6,opt,name=gasPrice,proto3" json:"gasPrice,omitempty"`
	Hash                 string `protobuf:"bytes,7,opt,name=hash,proto3" json:"hash,omitempty"`
	MethodId             string `protobuf:"bytes,8,opt,name=methodId,proto3" json:"methodId,omitempty"`
	Input                string `protobuf:"bytes,9,opt,name=input,proto3" json:"input,omitempty"`
	Nonce                uint32 `protobuf:"varint,10,opt,name=nonce,proto3" json:"nonce,omitempty"`
	To                   string `protobuf:"bytes,11,opt,name=to,proto3" json:"to,omitempty"`
	TransactionIndex     uint32 `protobuf:"varint,12,opt,name=transactionIndex,proto3" json:"transactionIndex,omitempty"`
	Value                uint64 `protobuf:"varint,13,opt,name=value,proto3" json:"value,omitempty"`
	Timestamp            uint32 `protobuf:"varint,14,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	MaxFeePerGas         uint64 `protobuf:"varint,15,opt,name=maxFeePerGas,proto3" json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas uint64 `protobuf:"varint,16,opt,name=maxPriorityFeePerGas,proto3" json:"maxPriorityFeePerGas,omitempty"`
}

func (x *Txns) Reset() {
	*x = Txns{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_txns_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Txns) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Txns) ProtoMessage() {}

func (x *Txns) ProtoReflect() protoreflect.Message {
	mi := &file_proto_txns_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Txns.ProtoReflect.Descriptor instead.
func (*Txns) Descriptor() ([]byte, []int) {
	return file_proto_txns_proto_rawDescGZIP(), []int{0}
}

func (x *Txns) GetChainId() string {
	if x != nil {
		return x.ChainId
	}
	return ""
}

func (x *Txns) GetBlockHash() string {
	if x != nil {
		return x.BlockHash
	}
	return ""
}

func (x *Txns) GetBlockNumber() uint32 {
	if x != nil {
		return x.BlockNumber
	}
	return 0
}

func (x *Txns) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *Txns) GetGas() uint64 {
	if x != nil {
		return x.Gas
	}
	return 0
}

func (x *Txns) GetGasPrice() uint64 {
	if x != nil {
		return x.GasPrice
	}
	return 0
}

func (x *Txns) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *Txns) GetMethodId() string {
	if x != nil {
		return x.MethodId
	}
	return ""
}

func (x *Txns) GetInput() string {
	if x != nil {
		return x.Input
	}
	return ""
}

func (x *Txns) GetNonce() uint32 {
	if x != nil {
		return x.Nonce
	}
	return 0
}

func (x *Txns) GetTo() string {
	if x != nil {
		return x.To
	}
	return ""
}

func (x *Txns) GetTransactionIndex() uint32 {
	if x != nil {
		return x.TransactionIndex
	}
	return 0
}

func (x *Txns) GetValue() uint64 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *Txns) GetTimestamp() uint32 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *Txns) GetMaxFeePerGas() uint64 {
	if x != nil {
		return x.MaxFeePerGas
	}
	return 0
}

func (x *Txns) GetMaxPriorityFeePerGas() uint64 {
	if x != nil {
		return x.MaxPriorityFeePerGas
	}
	return 0
}

type TxnsResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    int32   `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Data    []*Txns `protobuf:"bytes,2,rep,name=data,proto3" json:"data,omitempty"`
	Message string  `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *TxnsResult) Reset() {
	*x = TxnsResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_txns_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TxnsResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TxnsResult) ProtoMessage() {}

func (x *TxnsResult) ProtoReflect() protoreflect.Message {
	mi := &file_proto_txns_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TxnsResult.ProtoReflect.Descriptor instead.
func (*TxnsResult) Descriptor() ([]byte, []int) {
	return file_proto_txns_proto_rawDescGZIP(), []int{1}
}

func (x *TxnsResult) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *TxnsResult) GetData() []*Txns {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *TxnsResult) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_proto_txns_proto protoreflect.FileDescriptor

var file_proto_txns_proto_rawDesc = []byte{
	0x0a, 0x10, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x78, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xc6, 0x03, 0x0a, 0x04, 0x54, 0x78, 0x6e, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x63,
	0x68, 0x61, 0x69, 0x6e, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x68,
	0x61, 0x69, 0x6e, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61,
	0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48,
	0x61, 0x73, 0x68, 0x12, 0x20, 0x0a, 0x0b, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0b, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x4e,
	0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x10, 0x0a, 0x03, 0x67, 0x61, 0x73,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x67, 0x61, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x67,
	0x61, 0x73, 0x50, 0x72, 0x69, 0x63, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x67,
	0x61, 0x73, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12, 0x1a, 0x0a, 0x08, 0x6d,
	0x65, 0x74, 0x68, 0x6f, 0x64, 0x49, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6d,
	0x65, 0x74, 0x68, 0x6f, 0x64, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74,
	0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x12, 0x14, 0x0a,
	0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x6e, 0x6f,
	0x6e, 0x63, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x74, 0x6f, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x74, 0x6f, 0x12, 0x2a, 0x0a, 0x10, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x10, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12,
	0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x12, 0x22, 0x0a, 0x0c, 0x6d, 0x61, 0x78, 0x46, 0x65, 0x65, 0x50, 0x65, 0x72,
	0x47, 0x61, 0x73, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0c, 0x6d, 0x61, 0x78, 0x46, 0x65,
	0x65, 0x50, 0x65, 0x72, 0x47, 0x61, 0x73, 0x12, 0x32, 0x0a, 0x14, 0x6d, 0x61, 0x78, 0x50, 0x72,
	0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x46, 0x65, 0x65, 0x50, 0x65, 0x72, 0x47, 0x61, 0x73, 0x18,
	0x10, 0x20, 0x01, 0x28, 0x04, 0x52, 0x14, 0x6d, 0x61, 0x78, 0x50, 0x72, 0x69, 0x6f, 0x72, 0x69,
	0x74, 0x79, 0x46, 0x65, 0x65, 0x50, 0x65, 0x72, 0x47, 0x61, 0x73, 0x22, 0x55, 0x0a, 0x0a, 0x54,
	0x78, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x19, 0x0a,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x54, 0x78,
	0x6e, 0x73, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x42, 0x10, 0x5a, 0x0e, 0x70, 0x62, 0x2f, 0x64, 0x65, 0x6d, 0x6f, 0x2f, 0x77, 0x61,
	0x6c, 0x6c, 0x65, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_txns_proto_rawDescOnce sync.Once
	file_proto_txns_proto_rawDescData = file_proto_txns_proto_rawDesc
)

func file_proto_txns_proto_rawDescGZIP() []byte {
	file_proto_txns_proto_rawDescOnce.Do(func() {
		file_proto_txns_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_txns_proto_rawDescData)
	})
	return file_proto_txns_proto_rawDescData
}

var file_proto_txns_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_txns_proto_goTypes = []interface{}{
	(*Txns)(nil),       // 0: Txns
	(*TxnsResult)(nil), // 1: TxnsResult
}
var file_proto_txns_proto_depIdxs = []int32{
	0, // 0: TxnsResult.data:type_name -> Txns
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_txns_proto_init() }
func file_proto_txns_proto_init() {
	if File_proto_txns_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_txns_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Txns); i {
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
		file_proto_txns_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TxnsResult); i {
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
			RawDescriptor: file_proto_txns_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_txns_proto_goTypes,
		DependencyIndexes: file_proto_txns_proto_depIdxs,
		MessageInfos:      file_proto_txns_proto_msgTypes,
	}.Build()
	File_proto_txns_proto = out.File
	file_proto_txns_proto_rawDesc = nil
	file_proto_txns_proto_goTypes = nil
	file_proto_txns_proto_depIdxs = nil
}
