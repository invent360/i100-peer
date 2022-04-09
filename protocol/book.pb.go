// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: protocol/book.protocol

package protocol

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

type Book struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Isbn int32  `protobuf:"varint,2,opt,name=isbn,proto3" json:"isbn,omitempty"`
}

func (x *Book) Reset() {
	*x = Book{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_book_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Book) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Book) ProtoMessage() {}

func (x *Book) ProtoReflect() protoreflect.Message {
	mi := &file_proto_book_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Book.ProtoReflect.Descriptor instead.
func (*Book) Descriptor() ([]byte, []int) {
	return file_proto_book_proto_rawDescGZIP(), []int{0}
}

func (x *Book) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Book) GetIsbn() int32 {
	if x != nil {
		return x.Isbn
	}
	return 0
}

var File_proto_book_proto protoreflect.FileDescriptor

var file_proto_book_proto_rawDesc = []byte{
	0x0a, 0x10, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x62, 0x6f, 0x6f, 0x6b, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x2e, 0x0a, 0x04, 0x42, 0x6f, 0x6f, 0x6b, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x69, 0x73, 0x62, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x69, 0x73,
	0x62, 0x6e, 0x42, 0x08, 0x5a, 0x06, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_book_proto_rawDescOnce sync.Once
	file_proto_book_proto_rawDescData = file_proto_book_proto_rawDesc
)

func file_proto_book_proto_rawDescGZIP() []byte {
	file_proto_book_proto_rawDescOnce.Do(func() {
		file_proto_book_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_book_proto_rawDescData)
	})
	return file_proto_book_proto_rawDescData
}

var file_proto_book_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_proto_book_proto_goTypes = []interface{}{
	(*Book)(nil), // 0: Book
}
var file_proto_book_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_book_proto_init() }
func file_proto_book_proto_init() {
	if File_proto_book_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_book_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Book); i {
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
			RawDescriptor: file_proto_book_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_book_proto_goTypes,
		DependencyIndexes: file_proto_book_proto_depIdxs,
		MessageInfos:      file_proto_book_proto_msgTypes,
	}.Build()
	File_proto_book_proto = out.File
	file_proto_book_proto_rawDesc = nil
	file_proto_book_proto_goTypes = nil
	file_proto_book_proto_depIdxs = nil
}
