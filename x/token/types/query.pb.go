// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: token/query.proto

package types

import (
	context "context"
	fmt "fmt"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type BalanceRequest struct {
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
}

func (m *BalanceRequest) Reset()         { *m = BalanceRequest{} }
func (m *BalanceRequest) String() string { return proto.CompactTextString(m) }
func (*BalanceRequest) ProtoMessage()    {}
func (*BalanceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ec043bcd18c4056e, []int{0}
}
func (m *BalanceRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BalanceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BalanceRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BalanceRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BalanceRequest.Merge(m, src)
}
func (m *BalanceRequest) XXX_Size() int {
	return m.Size()
}
func (m *BalanceRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_BalanceRequest.DiscardUnknown(m)
}

var xxx_messageInfo_BalanceRequest proto.InternalMessageInfo

type BalanceResponse struct {
	Balance      types.DecProto `protobuf:"bytes,1,opt,name=balance,proto3" json:"balance"`
	BalanceDelta types.DecProto `protobuf:"bytes,2,opt,name=balance_delta,json=balanceDelta,proto3" json:"balance_delta"`
	IsBanned     bool           `protobuf:"varint,3,opt,name=is_banned,json=isBanned,proto3" json:"is_banned,omitempty"`
}

func (m *BalanceResponse) Reset()         { *m = BalanceResponse{} }
func (m *BalanceResponse) String() string { return proto.CompactTextString(m) }
func (*BalanceResponse) ProtoMessage()    {}
func (*BalanceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ec043bcd18c4056e, []int{1}
}
func (m *BalanceResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BalanceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BalanceResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BalanceResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BalanceResponse.Merge(m, src)
}
func (m *BalanceResponse) XXX_Size() int {
	return m.Size()
}
func (m *BalanceResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_BalanceResponse.DiscardUnknown(m)
}

var xxx_messageInfo_BalanceResponse proto.InternalMessageInfo

func (m *BalanceResponse) GetBalance() types.DecProto {
	if m != nil {
		return m.Balance
	}
	return types.DecProto{}
}

func (m *BalanceResponse) GetBalanceDelta() types.DecProto {
	if m != nil {
		return m.BalanceDelta
	}
	return types.DecProto{}
}

func (m *BalanceResponse) GetIsBanned() bool {
	if m != nil {
		return m.IsBanned
	}
	return false
}

type PoolRequest struct {
}

func (m *PoolRequest) Reset()         { *m = PoolRequest{} }
func (m *PoolRequest) String() string { return proto.CompactTextString(m) }
func (*PoolRequest) ProtoMessage()    {}
func (*PoolRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ec043bcd18c4056e, []int{2}
}
func (m *PoolRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PoolRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PoolRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PoolRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PoolRequest.Merge(m, src)
}
func (m *PoolRequest) XXX_Size() int {
	return m.Size()
}
func (m *PoolRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PoolRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PoolRequest proto.InternalMessageInfo

type PoolResponse struct {
	Size_                  types.DecCoin  `protobuf:"bytes,1,opt,name=size,proto3" json:"size"`
	TotalDelta             types.DecProto `protobuf:"bytes,2,opt,name=total_delta,json=totalDelta,proto3" json:"total_delta"`
	NextDistributionHeight uint64         `protobuf:"varint,3,opt,name=next_distribution_height,json=nextDistributionHeight,proto3" json:"next_distribution_height,omitempty"`
}

func (m *PoolResponse) Reset()         { *m = PoolResponse{} }
func (m *PoolResponse) String() string { return proto.CompactTextString(m) }
func (*PoolResponse) ProtoMessage()    {}
func (*PoolResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ec043bcd18c4056e, []int{3}
}
func (m *PoolResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PoolResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PoolResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PoolResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PoolResponse.Merge(m, src)
}
func (m *PoolResponse) XXX_Size() int {
	return m.Size()
}
func (m *PoolResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PoolResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PoolResponse proto.InternalMessageInfo

func (m *PoolResponse) GetSize_() types.DecCoin {
	if m != nil {
		return m.Size_
	}
	return types.DecCoin{}
}

func (m *PoolResponse) GetTotalDelta() types.DecProto {
	if m != nil {
		return m.TotalDelta
	}
	return types.DecProto{}
}

func (m *PoolResponse) GetNextDistributionHeight() uint64 {
	if m != nil {
		return m.NextDistributionHeight
	}
	return 0
}

func init() {
	proto.RegisterType((*BalanceRequest)(nil), "token.BalanceRequest")
	proto.RegisterType((*BalanceResponse)(nil), "token.BalanceResponse")
	proto.RegisterType((*PoolRequest)(nil), "token.PoolRequest")
	proto.RegisterType((*PoolResponse)(nil), "token.PoolResponse")
}

func init() { proto.RegisterFile("token/query.proto", fileDescriptor_ec043bcd18c4056e) }

var fileDescriptor_ec043bcd18c4056e = []byte{
	// 493 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x52, 0x41, 0x6b, 0x13, 0x41,
	0x14, 0xde, 0xad, 0xa9, 0x49, 0x27, 0x55, 0x71, 0xd4, 0xb8, 0xc4, 0xba, 0x29, 0x11, 0xa1, 0xa0,
	0xee, 0x90, 0x2a, 0x22, 0x82, 0x97, 0xb8, 0x48, 0x6f, 0xd6, 0x3d, 0x7a, 0x09, 0xb3, 0xbb, 0x8f,
	0xcd, 0xe0, 0x76, 0xde, 0x76, 0x67, 0x52, 0x5a, 0xc5, 0x8b, 0x27, 0x8f, 0x82, 0x7f, 0xa0, 0xff,
	0xc2, 0xa3, 0x37, 0xe9, 0xb1, 0xe0, 0xc5, 0x93, 0x48, 0xe2, 0xc1, 0x9f, 0x21, 0x3b, 0x3b, 0x29,
	0xa9, 0x22, 0xd8, 0xdb, 0xcc, 0xfb, 0xe6, 0xfb, 0xde, 0x7c, 0xef, 0x7d, 0xe4, 0xb2, 0xc6, 0x57,
	0x20, 0xd9, 0xee, 0x04, 0xca, 0x83, 0xa0, 0x28, 0x51, 0x23, 0x5d, 0x36, 0xa5, 0xee, 0xd5, 0x0c,
	0x33, 0x34, 0x15, 0x56, 0x9d, 0x6a, 0xb0, 0xbb, 0x96, 0x21, 0x66, 0x39, 0x30, 0x5e, 0x08, 0xc6,
	0xa5, 0x44, 0xcd, 0xb5, 0x40, 0xa9, 0x2c, 0xea, 0x27, 0xa8, 0x76, 0x50, 0xb1, 0x98, 0x2b, 0x60,
	0x7b, 0x83, 0x18, 0x34, 0x1f, 0xb0, 0x04, 0x85, 0xac, 0xf1, 0xfe, 0x03, 0x72, 0x71, 0xc8, 0x73,
	0x2e, 0x13, 0x88, 0x60, 0x77, 0x02, 0x4a, 0x53, 0x8f, 0x34, 0x79, 0x9a, 0x96, 0xa0, 0x94, 0xe7,
	0xae, 0xbb, 0x1b, 0x2b, 0xd1, 0xfc, 0xfa, 0xb8, 0xf5, 0xfe, 0xb0, 0xe7, 0xfc, 0x3a, 0xec, 0x39,
	0xfd, 0x4f, 0x2e, 0xb9, 0x74, 0x42, 0x53, 0x05, 0x4a, 0x05, 0xf4, 0x09, 0x69, 0xc6, 0x75, 0xc9,
	0xf0, 0xda, 0x9b, 0x37, 0x83, 0xba, 0x77, 0x50, 0xf5, 0x0e, 0x6c, 0xef, 0x20, 0x84, 0x64, 0xbb,
	0xea, 0x3c, 0x6c, 0x1c, 0x7d, 0xef, 0x39, 0xd1, 0x9c, 0x43, 0xb7, 0xc8, 0x05, 0x7b, 0x1c, 0xa5,
	0x90, 0x6b, 0xee, 0x2d, 0xfd, 0xbf, 0xc8, 0xaa, 0x65, 0x86, 0x15, 0x91, 0xde, 0x20, 0x2b, 0x42,
	0x8d, 0x62, 0x2e, 0x25, 0xa4, 0xde, 0xb9, 0x75, 0x77, 0xa3, 0x15, 0xb5, 0x84, 0x1a, 0x9a, 0x7b,
	0xff, 0x3a, 0x69, 0x6f, 0x23, 0xe6, 0xd6, 0xec, 0x82, 0xa5, 0x2f, 0x2e, 0x59, 0xad, 0x11, 0xeb,
	0xe7, 0x21, 0x69, 0x28, 0xf1, 0x7a, 0x6e, 0x66, 0xed, 0x5f, 0xff, 0x78, 0x8a, 0x42, 0xda, 0x6f,
	0x98, 0xf7, 0x34, 0x24, 0x6d, 0x8d, 0x9a, 0xe7, 0x67, 0xb7, 0x41, 0x0c, 0xaf, 0x36, 0xf1, 0x88,
	0x78, 0x12, 0xf6, 0xf5, 0x28, 0x15, 0x4a, 0x97, 0x22, 0x9e, 0x54, 0x3b, 0x1d, 0x8d, 0x41, 0x64,
	0x63, 0x6d, 0x3c, 0x35, 0xa2, 0x4e, 0x85, 0x87, 0x0b, 0xf0, 0x96, 0x41, 0x37, 0x3f, 0xbb, 0x64,
	0xf9, 0x45, 0x15, 0x1e, 0x0a, 0xa4, 0x69, 0x97, 0x44, 0xaf, 0x05, 0x26, 0x42, 0xc1, 0xe9, 0x5d,
	0x77, 0x3b, 0x7f, 0x96, 0x6b, 0xef, 0xfd, 0x3b, 0xef, 0xbe, 0xfe, 0xfc, 0xb8, 0x74, 0x9b, 0xde,
	0x62, 0x29, 0x24, 0x20, 0x75, 0xc9, 0xea, 0x50, 0xee, 0x0d, 0x98, 0x1d, 0x35, 0x7b, 0x63, 0x53,
	0xf1, 0x96, 0x3e, 0x27, 0x8d, 0x6a, 0x70, 0x94, 0x5a, 0xb1, 0x85, 0xf9, 0x76, 0xaf, 0x9c, 0xaa,
	0x59, 0x75, 0xdf, 0xa8, 0x7b, 0xb4, 0xf3, 0xb7, 0x7a, 0x81, 0x98, 0x0f, 0x9f, 0x1d, 0x4d, 0x7d,
	0xf7, 0x78, 0xea, 0xbb, 0x3f, 0xa6, 0xbe, 0xfb, 0x61, 0xe6, 0x3b, 0xc7, 0x33, 0xdf, 0xf9, 0x36,
	0xf3, 0x9d, 0x97, 0x77, 0x33, 0xa1, 0xc7, 0x93, 0x38, 0x48, 0x70, 0x87, 0x85, 0x35, 0xf7, 0x9e,
	0x04, 0x7d, 0xa2, 0xb3, 0x6f, 0x95, 0xf4, 0x41, 0x01, 0x2a, 0x3e, 0x6f, 0x22, 0x7e, 0xff, 0x77,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x11, 0xef, 0xa5, 0x63, 0x52, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	Balance(ctx context.Context, in *BalanceRequest, opts ...grpc.CallOption) (*BalanceResponse, error)
	Pool(ctx context.Context, in *PoolRequest, opts ...grpc.CallOption) (*PoolResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Balance(ctx context.Context, in *BalanceRequest, opts ...grpc.CallOption) (*BalanceResponse, error) {
	out := new(BalanceResponse)
	err := c.cc.Invoke(ctx, "/token.Query/Balance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Pool(ctx context.Context, in *PoolRequest, opts ...grpc.CallOption) (*PoolResponse, error) {
	out := new(PoolResponse)
	err := c.cc.Invoke(ctx, "/token.Query/Pool", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	Balance(context.Context, *BalanceRequest) (*BalanceResponse, error)
	Pool(context.Context, *PoolRequest) (*PoolResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) Balance(ctx context.Context, req *BalanceRequest) (*BalanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Balance not implemented")
}
func (*UnimplementedQueryServer) Pool(ctx context.Context, req *PoolRequest) (*PoolResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Pool not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_Balance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BalanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Balance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/token.Query/Balance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Balance(ctx, req.(*BalanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Pool_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PoolRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Pool(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/token.Query/Pool",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Pool(ctx, req.(*PoolRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "token.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Balance",
			Handler:    _Query_Balance_Handler,
		},
		{
			MethodName: "Pool",
			Handler:    _Query_Pool_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "token/query.proto",
}

func (m *BalanceRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BalanceRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BalanceRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *BalanceResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BalanceResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BalanceResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.IsBanned {
		i--
		if m.IsBanned {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x18
	}
	{
		size, err := m.BalanceDelta.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size, err := m.Balance.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *PoolRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PoolRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PoolRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *PoolResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PoolResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PoolResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.NextDistributionHeight != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.NextDistributionHeight))
		i--
		dAtA[i] = 0x18
	}
	{
		size, err := m.TotalDelta.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size, err := m.Size_.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *BalanceRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *BalanceResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Balance.Size()
	n += 1 + l + sovQuery(uint64(l))
	l = m.BalanceDelta.Size()
	n += 1 + l + sovQuery(uint64(l))
	if m.IsBanned {
		n += 2
	}
	return n
}

func (m *PoolRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *PoolResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Size_.Size()
	n += 1 + l + sovQuery(uint64(l))
	l = m.TotalDelta.Size()
	n += 1 + l + sovQuery(uint64(l))
	if m.NextDistributionHeight != 0 {
		n += 1 + sovQuery(uint64(m.NextDistributionHeight))
	}
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *BalanceRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: BalanceRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BalanceRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *BalanceResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: BalanceResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BalanceResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Balance", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Balance.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BalanceDelta", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.BalanceDelta.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsBanned", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.IsBanned = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *PoolRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: PoolRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PoolRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *PoolResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: PoolResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PoolResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Size_", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Size_.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalDelta", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TotalDelta.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NextDistributionHeight", wireType)
			}
			m.NextDistributionHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NextDistributionHeight |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)
