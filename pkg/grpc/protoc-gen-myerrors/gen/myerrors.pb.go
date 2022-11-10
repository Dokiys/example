// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.21.9
// source: myerrors.proto

package gen

import (
	_ "github.com/go-kratos/kratos/v2/errors"
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

type ErrorReason int32

const (
	ErrorReason_InvalidParameter ErrorReason = 0
	ErrorReason_AccessForbidden  ErrorReason = 1
	ErrorReason_Unauthenticated  ErrorReason = 2
	ErrorReason_BusinessError    ErrorReason = 3
	ErrorReason_SystemError      ErrorReason = 4
	ErrorReason_NotFound         ErrorReason = 5
	ErrorReason_OrderNotFound    ErrorReason = 500
	ErrorReason_ItemNotFound     ErrorReason = 600
)

// Enum value maps for ErrorReason.
var (
	ErrorReason_name = map[int32]string{
		0:   "InvalidParameter",
		1:   "AccessForbidden",
		2:   "Unauthenticated",
		3:   "BusinessError",
		4:   "SystemError",
		5:   "NotFound",
		500: "OrderNotFound",
		600: "ItemNotFound",
	}
	ErrorReason_value = map[string]int32{
		"InvalidParameter": 0,
		"AccessForbidden":  1,
		"Unauthenticated":  2,
		"BusinessError":    3,
		"SystemError":      4,
		"NotFound":         5,
		"OrderNotFound":    500,
		"ItemNotFound":     600,
	}
)

func (x ErrorReason) Enum() *ErrorReason {
	p := new(ErrorReason)
	*p = x
	return p
}

func (x ErrorReason) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ErrorReason) Descriptor() protoreflect.EnumDescriptor {
	return file_myerrors_proto_enumTypes[0].Descriptor()
}

func (ErrorReason) Type() protoreflect.EnumType {
	return &file_myerrors_proto_enumTypes[0]
}

func (x ErrorReason) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ErrorReason.Descriptor instead.
func (ErrorReason) EnumDescriptor() ([]byte, []int) {
	return file_myerrors_proto_rawDescGZIP(), []int{0}
}

var File_myerrors_proto protoreflect.FileDescriptor

var file_myerrors_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x6d, 0x79, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x08, 0x6d, 0x79, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x1a, 0x1a, 0x6b, 0x72, 0x61, 0x74,
	0x6f, 0x73, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2a, 0xdc, 0x01, 0x0a, 0x0b, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x10, 0x49, 0x6e, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x10, 0x00, 0x1a, 0x04, 0xa8, 0x45,
	0x90, 0x03, 0x12, 0x19, 0x0a, 0x0f, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x46, 0x6f, 0x72, 0x62,
	0x69, 0x64, 0x64, 0x65, 0x6e, 0x10, 0x01, 0x1a, 0x04, 0xa8, 0x45, 0x93, 0x03, 0x12, 0x19, 0x0a,
	0x0f, 0x55, 0x6e, 0x61, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x65, 0x64,
	0x10, 0x02, 0x1a, 0x04, 0xa8, 0x45, 0x91, 0x03, 0x12, 0x17, 0x0a, 0x0d, 0x42, 0x75, 0x73, 0x69,
	0x6e, 0x65, 0x73, 0x73, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x10, 0x03, 0x1a, 0x04, 0xa8, 0x45, 0x90,
	0x03, 0x12, 0x15, 0x0a, 0x0b, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x10, 0x04, 0x1a, 0x04, 0xa8, 0x45, 0xf4, 0x03, 0x12, 0x12, 0x0a, 0x08, 0x4e, 0x6f, 0x74, 0x46,
	0x6f, 0x75, 0x6e, 0x64, 0x10, 0x05, 0x1a, 0x04, 0xa8, 0x45, 0x90, 0x03, 0x12, 0x18, 0x0a, 0x0d,
	0x4f, 0x72, 0x64, 0x65, 0x72, 0x4e, 0x6f, 0x74, 0x46, 0x6f, 0x75, 0x6e, 0x64, 0x10, 0xf4, 0x03,
	0x1a, 0x04, 0xa8, 0x45, 0x90, 0x03, 0x12, 0x17, 0x0a, 0x0c, 0x49, 0x74, 0x65, 0x6d, 0x4e, 0x6f,
	0x74, 0x46, 0x6f, 0x75, 0x6e, 0x64, 0x10, 0xd8, 0x04, 0x1a, 0x04, 0xa8, 0x45, 0x90, 0x03, 0x1a,
	0x04, 0xa0, 0x45, 0xf4, 0x03, 0x42, 0x1d, 0x5a, 0x1b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d,
	0x67, 0x65, 0x6e, 0x2d, 0x6d, 0x79, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2f, 0x67, 0x65, 0x6e,
	0x3b, 0x67, 0x65, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_myerrors_proto_rawDescOnce sync.Once
	file_myerrors_proto_rawDescData = file_myerrors_proto_rawDesc
)

func file_myerrors_proto_rawDescGZIP() []byte {
	file_myerrors_proto_rawDescOnce.Do(func() {
		file_myerrors_proto_rawDescData = protoimpl.X.CompressGZIP(file_myerrors_proto_rawDescData)
	})
	return file_myerrors_proto_rawDescData
}

var file_myerrors_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_myerrors_proto_goTypes = []interface{}{
	(ErrorReason)(0), // 0: myerrors.ErrorReason
}
var file_myerrors_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_myerrors_proto_init() }
func file_myerrors_proto_init() {
	if File_myerrors_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_myerrors_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_myerrors_proto_goTypes,
		DependencyIndexes: file_myerrors_proto_depIdxs,
		EnumInfos:         file_myerrors_proto_enumTypes,
	}.Build()
	File_myerrors_proto = out.File
	file_myerrors_proto_rawDesc = nil
	file_myerrors_proto_goTypes = nil
	file_myerrors_proto_depIdxs = nil
}
