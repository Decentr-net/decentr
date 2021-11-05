// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: community/community.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type Category int32

const (
	Category_UNDEFINED              Category = 0
	Category_WORLD_NEWS             Category = 1
	Category_TRAVEL_AND_TOURISM     Category = 2
	Category_SCIENCE_AND_TECHNOLOGY Category = 3
	Category_STRANGE_WORLD          Category = 4
	Category_ARTS_AND_ENTERTAINMENT Category = 5
	Category_WRITERS_AND_WRITING    Category = 6
	Category_HEALTH_AND_FITNESS     Category = 7
	Category_CRYPTO_AND_BLOCKCHAIN  Category = 8
	Category_SPORTS                 Category = 9
)

var Category_name = map[int32]string{
	0: "UNDEFINED",
	1: "WORLD_NEWS",
	2: "TRAVEL_AND_TOURISM",
	3: "SCIENCE_AND_TECHNOLOGY",
	4: "STRANGE_WORLD",
	5: "ARTS_AND_ENTERTAINMENT",
	6: "WRITERS_AND_WRITING",
	7: "HEALTH_AND_FITNESS",
	8: "CRYPTO_AND_BLOCKCHAIN",
	9: "SPORTS",
}

var Category_value = map[string]int32{
	"UNDEFINED":              0,
	"WORLD_NEWS":             1,
	"TRAVEL_AND_TOURISM":     2,
	"SCIENCE_AND_TECHNOLOGY": 3,
	"STRANGE_WORLD":          4,
	"ARTS_AND_ENTERTAINMENT": 5,
	"WRITERS_AND_WRITING":    6,
	"HEALTH_AND_FITNESS":     7,
	"CRYPTO_AND_BLOCKCHAIN":  8,
	"SPORTS":                 9,
}

func (x Category) String() string {
	return proto.EnumName(Category_name, int32(x))
}

func (Category) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_663aacd61135a87b, []int{0}
}

type LikeWeight int32

const (
	LikeWeight_LikeWeightZero LikeWeight = 0
	LikeWeight_LikeWeightUp   LikeWeight = 1
	LikeWeight_LikeWeightDown LikeWeight = -1
)

var LikeWeight_name = map[int32]string{
	0:  "LikeWeightZero",
	1:  "LikeWeightUp",
	-1: "LikeWeightDown",
}

var LikeWeight_value = map[string]int32{
	"LikeWeightZero": 0,
	"LikeWeightUp":   1,
	"LikeWeightDown": -1,
}

func (x LikeWeight) String() string {
	return proto.EnumName(LikeWeight_name, int32(x))
}

func (LikeWeight) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_663aacd61135a87b, []int{1}
}

type Params struct {
	Moderators []string       `protobuf:"bytes,1,rep,name=moderators,proto3" json:"moderators,omitempty"`
	FixedGas   FixedGasParams `protobuf:"bytes,2,opt,name=fixed_gas,json=fixedGas,proto3" json:"fixed_gas"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_663aacd61135a87b, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetModerators() []string {
	if m != nil {
		return m.Moderators
	}
	return nil
}

func (m *Params) GetFixedGas() FixedGasParams {
	if m != nil {
		return m.FixedGas
	}
	return FixedGasParams{}
}

type FixedGasParams struct {
	CretePost uint64 `protobuf:"varint,1,opt,name=crete_post,json=cretePost,proto3" json:"crete_post,omitempty"`
	SetLike   uint64 `protobuf:"varint,2,opt,name=set_like,json=setLike,proto3" json:"set_like,omitempty"`
	Follow    uint64 `protobuf:"varint,3,opt,name=follow,proto3" json:"follow,omitempty"`
	Unfollow  uint64 `protobuf:"varint,4,opt,name=unfollow,proto3" json:"unfollow,omitempty"`
}

func (m *FixedGasParams) Reset()         { *m = FixedGasParams{} }
func (m *FixedGasParams) String() string { return proto.CompactTextString(m) }
func (*FixedGasParams) ProtoMessage()    {}
func (*FixedGasParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_663aacd61135a87b, []int{1}
}
func (m *FixedGasParams) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *FixedGasParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_FixedGasParams.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *FixedGasParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FixedGasParams.Merge(m, src)
}
func (m *FixedGasParams) XXX_Size() int {
	return m.Size()
}
func (m *FixedGasParams) XXX_DiscardUnknown() {
	xxx_messageInfo_FixedGasParams.DiscardUnknown(m)
}

var xxx_messageInfo_FixedGasParams proto.InternalMessageInfo

func (m *FixedGasParams) GetCretePost() uint64 {
	if m != nil {
		return m.CretePost
	}
	return 0
}

func (m *FixedGasParams) GetSetLike() uint64 {
	if m != nil {
		return m.SetLike
	}
	return 0
}

func (m *FixedGasParams) GetFollow() uint64 {
	if m != nil {
		return m.Follow
	}
	return 0
}

func (m *FixedGasParams) GetUnfollow() uint64 {
	if m != nil {
		return m.Unfollow
	}
	return 0
}

type Post struct {
	Uuid         string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Owner        string   `protobuf:"bytes,2,opt,name=owner,proto3" json:"owner,omitempty"`
	Title        string   `protobuf:"bytes,3,opt,name=title,proto3" json:"title,omitempty"`
	PreviewImage string   `protobuf:"bytes,4,opt,name=preview_image,json=previewImage,proto3" json:"preview_image,omitempty"`
	Category     Category `protobuf:"varint,5,opt,name=category,proto3,enum=community.Category" json:"category,omitempty"`
	Text         string   `protobuf:"bytes,6,opt,name=text,proto3" json:"text,omitempty"`
}

func (m *Post) Reset()         { *m = Post{} }
func (m *Post) String() string { return proto.CompactTextString(m) }
func (*Post) ProtoMessage()    {}
func (*Post) Descriptor() ([]byte, []int) {
	return fileDescriptor_663aacd61135a87b, []int{2}
}
func (m *Post) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Post) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Post.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Post) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Post.Merge(m, src)
}
func (m *Post) XXX_Size() int {
	return m.Size()
}
func (m *Post) XXX_DiscardUnknown() {
	xxx_messageInfo_Post.DiscardUnknown(m)
}

var xxx_messageInfo_Post proto.InternalMessageInfo

func (m *Post) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *Post) GetOwner() string {
	if m != nil {
		return m.Owner
	}
	return ""
}

func (m *Post) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Post) GetPreviewImage() string {
	if m != nil {
		return m.PreviewImage
	}
	return ""
}

func (m *Post) GetCategory() Category {
	if m != nil {
		return m.Category
	}
	return Category_UNDEFINED
}

func (m *Post) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

type Like struct {
	Owner     string     `protobuf:"bytes,1,opt,name=owner,proto3" json:"owner,omitempty"`
	PostOwner string     `protobuf:"bytes,2,opt,name=post_owner,json=postOwner,proto3" json:"post_owner,omitempty"`
	PostUuid  string     `protobuf:"bytes,3,opt,name=post_uuid,json=postUuid,proto3" json:"post_uuid,omitempty"`
	Weight    LikeWeight `protobuf:"varint,4,opt,name=weight,proto3,enum=community.LikeWeight" json:"weight,omitempty"`
}

func (m *Like) Reset()         { *m = Like{} }
func (m *Like) String() string { return proto.CompactTextString(m) }
func (*Like) ProtoMessage()    {}
func (*Like) Descriptor() ([]byte, []int) {
	return fileDescriptor_663aacd61135a87b, []int{3}
}
func (m *Like) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Like) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Like.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Like) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Like.Merge(m, src)
}
func (m *Like) XXX_Size() int {
	return m.Size()
}
func (m *Like) XXX_DiscardUnknown() {
	xxx_messageInfo_Like.DiscardUnknown(m)
}

var xxx_messageInfo_Like proto.InternalMessageInfo

func (m *Like) GetOwner() string {
	if m != nil {
		return m.Owner
	}
	return ""
}

func (m *Like) GetPostOwner() string {
	if m != nil {
		return m.PostOwner
	}
	return ""
}

func (m *Like) GetPostUuid() string {
	if m != nil {
		return m.PostUuid
	}
	return ""
}

func (m *Like) GetWeight() LikeWeight {
	if m != nil {
		return m.Weight
	}
	return LikeWeight_LikeWeightZero
}

func init() {
	proto.RegisterEnum("community.Category", Category_name, Category_value)
	proto.RegisterEnum("community.LikeWeight", LikeWeight_name, LikeWeight_value)
	proto.RegisterType((*Params)(nil), "community.Params")
	proto.RegisterType((*FixedGasParams)(nil), "community.FixedGasParams")
	proto.RegisterType((*Post)(nil), "community.Post")
	proto.RegisterType((*Like)(nil), "community.Like")
}

func init() { proto.RegisterFile("community/community.proto", fileDescriptor_663aacd61135a87b) }

var fileDescriptor_663aacd61135a87b = []byte{
	// 674 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x93, 0xcf, 0x6e, 0xe2, 0x48,
	0x10, 0xc6, 0xe9, 0xc4, 0x21, 0xb8, 0x36, 0x41, 0xde, 0xce, 0x9f, 0x05, 0xa2, 0x78, 0x11, 0x7b,
	0x41, 0x91, 0x12, 0x6f, 0xb2, 0xd7, 0xbd, 0x10, 0x70, 0xc0, 0xbb, 0xc4, 0x46, 0x6d, 0xb3, 0x28,
	0xb9, 0x58, 0x06, 0x1a, 0x62, 0x05, 0x68, 0x64, 0x37, 0x21, 0xb9, 0xec, 0x75, 0xae, 0xf3, 0x28,
	0xf3, 0x18, 0x39, 0xe6, 0x38, 0xa7, 0xd1, 0x28, 0x99, 0xf7, 0x98, 0x51, 0xb7, 0x09, 0x30, 0x3e,
	0x55, 0xfd, 0xbe, 0x52, 0x7f, 0x5f, 0x95, 0x64, 0xc8, 0xf7, 0xd8, 0x78, 0x3c, 0x9b, 0x84, 0xfc,
	0xc9, 0x58, 0x56, 0x67, 0xd3, 0x88, 0x71, 0x86, 0xd5, 0x25, 0x28, 0xec, 0x0f, 0xd9, 0x90, 0x49,
	0x6a, 0x88, 0x2a, 0x19, 0x28, 0xe8, 0x3d, 0x16, 0x8f, 0x59, 0x6c, 0x74, 0x83, 0x98, 0x1a, 0x0f,
	0xe7, 0x5d, 0xca, 0x83, 0x73, 0xa3, 0xc7, 0xc2, 0x49, 0xa2, 0x97, 0x06, 0x90, 0x6e, 0x05, 0x51,
	0x30, 0x8e, 0xb1, 0x0e, 0x30, 0x66, 0x7d, 0x1a, 0x05, 0x9c, 0x45, 0x71, 0x0e, 0x15, 0x37, 0xcb,
	0x2a, 0x59, 0x23, 0xf8, 0x6f, 0x50, 0x07, 0xe1, 0x23, 0xed, 0xfb, 0xc3, 0x20, 0xce, 0x6d, 0x14,
	0x51, 0xf9, 0x97, 0x8b, 0xfc, 0xd9, 0x2a, 0xcf, 0x95, 0xd0, 0xea, 0x41, 0x9c, 0xbc, 0x76, 0xa9,
	0x3c, 0x7f, 0xf9, 0x3d, 0x45, 0x32, 0x83, 0x05, 0x2d, 0xfd, 0x0f, 0xd9, 0x9f, 0x27, 0xf0, 0x31,
	0x40, 0x2f, 0xa2, 0x9c, 0xfa, 0x53, 0x16, 0xf3, 0x1c, 0x2a, 0xa2, 0xb2, 0x42, 0x54, 0x49, 0x5a,
	0x2c, 0xe6, 0x38, 0x0f, 0x99, 0x98, 0x72, 0x7f, 0x14, 0xde, 0x53, 0xe9, 0xa6, 0x90, 0xed, 0x98,
	0xf2, 0x66, 0x78, 0x4f, 0xf1, 0x21, 0xa4, 0x07, 0x6c, 0x34, 0x62, 0xf3, 0xdc, 0xa6, 0x14, 0x16,
	0x1d, 0x2e, 0x40, 0x66, 0x36, 0x59, 0x28, 0x8a, 0x54, 0x96, 0x7d, 0xe9, 0x13, 0x02, 0x45, 0xbe,
	0x8b, 0x41, 0x99, 0xcd, 0xc2, 0xbe, 0x34, 0x54, 0x89, 0xac, 0xf1, 0x3e, 0x6c, 0xb1, 0xf9, 0x84,
	0x46, 0xd2, 0x48, 0x25, 0x49, 0x23, 0x28, 0x0f, 0xf9, 0x88, 0x4a, 0x17, 0x95, 0x24, 0x0d, 0xfe,
	0x03, 0x76, 0xa7, 0x11, 0x7d, 0x08, 0xe9, 0xdc, 0x0f, 0xc7, 0xc1, 0x90, 0x4a, 0x27, 0x95, 0xec,
	0x2c, 0xa0, 0x25, 0x18, 0x36, 0x20, 0xd3, 0x0b, 0x38, 0x1d, 0xb2, 0xe8, 0x29, 0xb7, 0x55, 0x44,
	0xe5, 0xec, 0xc5, 0xde, 0xda, 0xa9, 0xaa, 0x0b, 0x89, 0x2c, 0x87, 0x44, 0x2a, 0x4e, 0x1f, 0x79,
	0x2e, 0x9d, 0xa4, 0x12, 0x75, 0xe9, 0x03, 0x02, 0x45, 0xee, 0xbb, 0x8c, 0x87, 0xd6, 0xe3, 0x1d,
	0x03, 0x88, 0xcb, 0xf9, 0xeb, 0xc9, 0x55, 0x41, 0x1c, 0x29, 0x1f, 0x81, 0x6c, 0x7c, 0xb9, 0x6c,
	0xb2, 0x41, 0x46, 0x80, 0xb6, 0x58, 0xf8, 0x14, 0xd2, 0x73, 0x1a, 0x0e, 0xef, 0xb8, 0x4c, 0x9f,
	0xbd, 0x38, 0x58, 0x4b, 0x27, 0x2c, 0x3b, 0x52, 0x24, 0x8b, 0xa1, 0x93, 0x6f, 0x08, 0x32, 0xef,
	0xa1, 0xf1, 0x2e, 0xa8, 0x6d, 0xbb, 0x66, 0x5e, 0x59, 0xb6, 0x59, 0xd3, 0x52, 0x38, 0x0b, 0xd0,
	0x71, 0x48, 0xb3, 0xe6, 0xdb, 0x66, 0xc7, 0xd5, 0x10, 0x3e, 0x04, 0xec, 0x91, 0xca, 0x7f, 0x66,
	0xd3, 0xaf, 0xd8, 0x35, 0xdf, 0x73, 0xda, 0xc4, 0x72, 0xaf, 0xb5, 0x0d, 0x5c, 0x80, 0x43, 0xb7,
	0x6a, 0x99, 0x76, 0xd5, 0x4c, 0x04, 0xb3, 0xda, 0xb0, 0x9d, 0xa6, 0x53, 0xbf, 0xd1, 0x36, 0xf1,
	0xaf, 0xb0, 0xeb, 0x7a, 0xa4, 0x62, 0xd7, 0x4d, 0x5f, 0xbe, 0xa5, 0x29, 0x62, 0xbc, 0x42, 0x3c,
	0x57, 0xce, 0x9a, 0xb6, 0x67, 0x12, 0xaf, 0x62, 0xd9, 0xd7, 0xa6, 0xed, 0x69, 0x5b, 0xf8, 0x37,
	0xd8, 0xeb, 0x10, 0xcb, 0x33, 0x49, 0x22, 0x8b, 0xda, 0xb2, 0xeb, 0x5a, 0x5a, 0x78, 0x37, 0xcc,
	0x4a, 0xd3, 0x6b, 0x48, 0x7e, 0x65, 0x79, 0xb6, 0xe9, 0xba, 0xda, 0x36, 0xce, 0xc3, 0x41, 0x95,
	0xdc, 0xb4, 0x3c, 0x47, 0xf2, 0xcb, 0xa6, 0x53, 0xfd, 0xb7, 0xda, 0xa8, 0x58, 0xb6, 0x96, 0xc1,
	0x00, 0x69, 0xb7, 0xe5, 0x10, 0xcf, 0xd5, 0xd4, 0x13, 0x07, 0x60, 0xb5, 0x3c, 0xc6, 0x90, 0x5d,
	0x75, 0xb7, 0x34, 0x62, 0x5a, 0x0a, 0x6b, 0xb0, 0xb3, 0x62, 0xed, 0xa9, 0x86, 0xf0, 0xd1, 0xfa,
	0x54, 0x8d, 0xcd, 0x27, 0xda, 0xf7, 0xf7, 0x0f, 0x5d, 0xfe, 0xf3, 0xfc, 0xaa, 0xa3, 0x97, 0x57,
	0x1d, 0x7d, 0x7d, 0xd5, 0xd1, 0xc7, 0x37, 0x3d, 0xf5, 0xf2, 0xa6, 0xa7, 0x3e, 0xbf, 0xe9, 0xa9,
	0xdb, 0x3f, 0x87, 0x21, 0xbf, 0x9b, 0x75, 0xc5, 0xd9, 0x8d, 0x1a, 0xed, 0xd1, 0x09, 0x8f, 0x4e,
	0x27, 0x94, 0x1b, 0xfd, 0xa4, 0x36, 0x1e, 0x57, 0x7f, 0xba, 0xc1, 0x9f, 0xa6, 0x34, 0xee, 0xa6,
	0xe5, 0xff, 0xfa, 0xd7, 0x8f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x3f, 0x30, 0x77, 0x60, 0x0d, 0x04,
	0x00, 0x00,
}

func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.FixedGas.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintCommunity(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Moderators) > 0 {
		for iNdEx := len(m.Moderators) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Moderators[iNdEx])
			copy(dAtA[i:], m.Moderators[iNdEx])
			i = encodeVarintCommunity(dAtA, i, uint64(len(m.Moderators[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *FixedGasParams) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *FixedGasParams) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *FixedGasParams) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Unfollow != 0 {
		i = encodeVarintCommunity(dAtA, i, uint64(m.Unfollow))
		i--
		dAtA[i] = 0x20
	}
	if m.Follow != 0 {
		i = encodeVarintCommunity(dAtA, i, uint64(m.Follow))
		i--
		dAtA[i] = 0x18
	}
	if m.SetLike != 0 {
		i = encodeVarintCommunity(dAtA, i, uint64(m.SetLike))
		i--
		dAtA[i] = 0x10
	}
	if m.CretePost != 0 {
		i = encodeVarintCommunity(dAtA, i, uint64(m.CretePost))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Post) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Post) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Post) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Text) > 0 {
		i -= len(m.Text)
		copy(dAtA[i:], m.Text)
		i = encodeVarintCommunity(dAtA, i, uint64(len(m.Text)))
		i--
		dAtA[i] = 0x32
	}
	if m.Category != 0 {
		i = encodeVarintCommunity(dAtA, i, uint64(m.Category))
		i--
		dAtA[i] = 0x28
	}
	if len(m.PreviewImage) > 0 {
		i -= len(m.PreviewImage)
		copy(dAtA[i:], m.PreviewImage)
		i = encodeVarintCommunity(dAtA, i, uint64(len(m.PreviewImage)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Title) > 0 {
		i -= len(m.Title)
		copy(dAtA[i:], m.Title)
		i = encodeVarintCommunity(dAtA, i, uint64(len(m.Title)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Owner) > 0 {
		i -= len(m.Owner)
		copy(dAtA[i:], m.Owner)
		i = encodeVarintCommunity(dAtA, i, uint64(len(m.Owner)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Uuid) > 0 {
		i -= len(m.Uuid)
		copy(dAtA[i:], m.Uuid)
		i = encodeVarintCommunity(dAtA, i, uint64(len(m.Uuid)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Like) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Like) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Like) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Weight != 0 {
		i = encodeVarintCommunity(dAtA, i, uint64(m.Weight))
		i--
		dAtA[i] = 0x20
	}
	if len(m.PostUuid) > 0 {
		i -= len(m.PostUuid)
		copy(dAtA[i:], m.PostUuid)
		i = encodeVarintCommunity(dAtA, i, uint64(len(m.PostUuid)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.PostOwner) > 0 {
		i -= len(m.PostOwner)
		copy(dAtA[i:], m.PostOwner)
		i = encodeVarintCommunity(dAtA, i, uint64(len(m.PostOwner)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Owner) > 0 {
		i -= len(m.Owner)
		copy(dAtA[i:], m.Owner)
		i = encodeVarintCommunity(dAtA, i, uint64(len(m.Owner)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintCommunity(dAtA []byte, offset int, v uint64) int {
	offset -= sovCommunity(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Moderators) > 0 {
		for _, s := range m.Moderators {
			l = len(s)
			n += 1 + l + sovCommunity(uint64(l))
		}
	}
	l = m.FixedGas.Size()
	n += 1 + l + sovCommunity(uint64(l))
	return n
}

func (m *FixedGasParams) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.CretePost != 0 {
		n += 1 + sovCommunity(uint64(m.CretePost))
	}
	if m.SetLike != 0 {
		n += 1 + sovCommunity(uint64(m.SetLike))
	}
	if m.Follow != 0 {
		n += 1 + sovCommunity(uint64(m.Follow))
	}
	if m.Unfollow != 0 {
		n += 1 + sovCommunity(uint64(m.Unfollow))
	}
	return n
}

func (m *Post) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Uuid)
	if l > 0 {
		n += 1 + l + sovCommunity(uint64(l))
	}
	l = len(m.Owner)
	if l > 0 {
		n += 1 + l + sovCommunity(uint64(l))
	}
	l = len(m.Title)
	if l > 0 {
		n += 1 + l + sovCommunity(uint64(l))
	}
	l = len(m.PreviewImage)
	if l > 0 {
		n += 1 + l + sovCommunity(uint64(l))
	}
	if m.Category != 0 {
		n += 1 + sovCommunity(uint64(m.Category))
	}
	l = len(m.Text)
	if l > 0 {
		n += 1 + l + sovCommunity(uint64(l))
	}
	return n
}

func (m *Like) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Owner)
	if l > 0 {
		n += 1 + l + sovCommunity(uint64(l))
	}
	l = len(m.PostOwner)
	if l > 0 {
		n += 1 + l + sovCommunity(uint64(l))
	}
	l = len(m.PostUuid)
	if l > 0 {
		n += 1 + l + sovCommunity(uint64(l))
	}
	if m.Weight != 0 {
		n += 1 + sovCommunity(uint64(m.Weight))
	}
	return n
}

func sovCommunity(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozCommunity(x uint64) (n int) {
	return sovCommunity(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommunity
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Moderators", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommunity
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthCommunity
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCommunity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Moderators = append(m.Moderators, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FixedGas", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommunity
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthCommunity
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthCommunity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.FixedGas.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCommunity(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCommunity
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *FixedGasParams) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommunity
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: FixedGasParams: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: FixedGasParams: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CretePost", wireType)
			}
			m.CretePost = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommunity
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CretePost |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SetLike", wireType)
			}
			m.SetLike = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommunity
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SetLike |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Follow", wireType)
			}
			m.Follow = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommunity
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Follow |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Unfollow", wireType)
			}
			m.Unfollow = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommunity
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Unfollow |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipCommunity(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCommunity
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Post) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommunity
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Post: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Post: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Uuid", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommunity
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthCommunity
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCommunity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Uuid = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Owner", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommunity
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthCommunity
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCommunity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Owner = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Title", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommunity
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthCommunity
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCommunity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Title = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PreviewImage", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommunity
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthCommunity
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCommunity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PreviewImage = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Category", wireType)
			}
			m.Category = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommunity
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Category |= Category(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Text", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommunity
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthCommunity
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCommunity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Text = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCommunity(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCommunity
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Like) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommunity
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Like: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Like: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Owner", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommunity
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthCommunity
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCommunity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Owner = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PostOwner", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommunity
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthCommunity
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCommunity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PostOwner = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PostUuid", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommunity
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthCommunity
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCommunity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PostUuid = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Weight", wireType)
			}
			m.Weight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommunity
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Weight |= LikeWeight(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipCommunity(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCommunity
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipCommunity(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowCommunity
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowCommunity
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowCommunity
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthCommunity
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupCommunity
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthCommunity
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthCommunity        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowCommunity          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupCommunity = fmt.Errorf("proto: unexpected end of group")
)
