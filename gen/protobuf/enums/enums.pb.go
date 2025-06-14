// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: enums.proto

package enums

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type DomainType int32

const (
	DomainType_DOMAIN_TYPE_UNSPECIFIED DomainType = 0
	DomainType_DOMAIN_TYPE_UNDEFINED   DomainType = 1
	DomainType_DOMAIN_TYPE_BLACKLIST   DomainType = 2
	DomainType_DOMAIN_TYPE_WHITELIST   DomainType = 3
)

// Enum value maps for DomainType.
var (
	DomainType_name = map[int32]string{
		0: "DOMAIN_TYPE_UNSPECIFIED",
		1: "DOMAIN_TYPE_UNDEFINED",
		2: "DOMAIN_TYPE_BLACKLIST",
		3: "DOMAIN_TYPE_WHITELIST",
	}
	DomainType_value = map[string]int32{
		"DOMAIN_TYPE_UNSPECIFIED": 0,
		"DOMAIN_TYPE_UNDEFINED":   1,
		"DOMAIN_TYPE_BLACKLIST":   2,
		"DOMAIN_TYPE_WHITELIST":   3,
	}
)

func (x DomainType) Enum() *DomainType {
	p := new(DomainType)
	*p = x
	return p
}

func (x DomainType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DomainType) Descriptor() protoreflect.EnumDescriptor {
	return file_enums_proto_enumTypes[0].Descriptor()
}

func (DomainType) Type() protoreflect.EnumType {
	return &file_enums_proto_enumTypes[0]
}

func (x DomainType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DomainType.Descriptor instead.
func (DomainType) EnumDescriptor() ([]byte, []int) {
	return file_enums_proto_rawDescGZIP(), []int{0}
}

var File_enums_proto protoreflect.FileDescriptor

const file_enums_proto_rawDesc = "" +
	"\n" +
	"\venums.proto*z\n" +
	"\n" +
	"DomainType\x12\x1b\n" +
	"\x17DOMAIN_TYPE_UNSPECIFIED\x10\x00\x12\x19\n" +
	"\x15DOMAIN_TYPE_UNDEFINED\x10\x01\x12\x19\n" +
	"\x15DOMAIN_TYPE_BLACKLIST\x10\x02\x12\x19\n" +
	"\x15DOMAIN_TYPE_WHITELIST\x10\x03B@Z>github.com/aerosystems/common-service/gen/protobuf/enums;enumsb\x06proto3"

var (
	file_enums_proto_rawDescOnce sync.Once
	file_enums_proto_rawDescData []byte
)

func file_enums_proto_rawDescGZIP() []byte {
	file_enums_proto_rawDescOnce.Do(func() {
		file_enums_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_enums_proto_rawDesc), len(file_enums_proto_rawDesc)))
	})
	return file_enums_proto_rawDescData
}

var file_enums_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_enums_proto_goTypes = []any{
	(DomainType)(0), // 0: DomainType
}
var file_enums_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_enums_proto_init() }
func file_enums_proto_init() {
	if File_enums_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_enums_proto_rawDesc), len(file_enums_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_enums_proto_goTypes,
		DependencyIndexes: file_enums_proto_depIdxs,
		EnumInfos:         file_enums_proto_enumTypes,
	}.Build()
	File_enums_proto = out.File
	file_enums_proto_goTypes = nil
	file_enums_proto_depIdxs = nil
}
