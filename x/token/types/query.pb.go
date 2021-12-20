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
	Balance types.DecProto `protobuf:"bytes,1,opt,name=balance,proto3" json:"balance"`
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

func init() {
	proto.RegisterType((*BalanceRequest)(nil), "token.BalanceRequest")
	proto.RegisterType((*BalanceResponse)(nil), "token.BalanceResponse")
}

func init() { proto.RegisterFile("token/query.proto", fileDescriptor_ec043bcd18c4056e) }

var fileDescriptor_ec043bcd18c4056e = []byte{
	// 328 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x90, 0xb1, 0x4a, 0x33, 0x41,
	0x14, 0x85, 0x77, 0x7f, 0xfe, 0x18, 0x1d, 0x41, 0x71, 0x51, 0x09, 0x41, 0x27, 0x61, 0xab, 0x14,
	0x3a, 0x43, 0xa2, 0x95, 0x60, 0x13, 0x82, 0x75, 0x4c, 0x69, 0x37, 0xbb, 0xb9, 0xac, 0x4b, 0x92,
	0xb9, 0x9b, 0x9d, 0x89, 0x18, 0xc4, 0xc6, 0xca, 0x52, 0xf0, 0x05, 0xf2, 0x38, 0x29, 0x03, 0x36,
	0x56, 0x22, 0x59, 0x0b, 0x1f, 0x43, 0x76, 0x67, 0x12, 0xd0, 0x6e, 0xe6, 0x5c, 0xce, 0xf9, 0xee,
	0xb9, 0x64, 0x4f, 0xe3, 0x00, 0x24, 0x1f, 0x4f, 0x20, 0x9d, 0xb2, 0x24, 0x45, 0x8d, 0x5e, 0xa9,
	0x90, 0xaa, 0xfb, 0x11, 0x46, 0x58, 0x28, 0x3c, 0x7f, 0x99, 0x61, 0xf5, 0x28, 0x42, 0x8c, 0x86,
	0xc0, 0x45, 0x12, 0x73, 0x21, 0x25, 0x6a, 0xa1, 0x63, 0x94, 0xca, 0x4e, 0x69, 0x88, 0x6a, 0x84,
	0x8a, 0x07, 0x42, 0x01, 0xbf, 0x6b, 0x06, 0xa0, 0x45, 0x93, 0x87, 0x18, 0x4b, 0x33, 0xf7, 0xcf,
	0xc9, 0x4e, 0x5b, 0x0c, 0x85, 0x0c, 0xa1, 0x07, 0xe3, 0x09, 0x28, 0xed, 0x55, 0x48, 0x59, 0xf4,
	0xfb, 0x29, 0x28, 0x55, 0x71, 0xeb, 0x6e, 0x63, 0xab, 0xb7, 0xfa, 0x5e, 0x6c, 0x3e, 0xcf, 0x6a,
	0xce, 0xf7, 0xac, 0xe6, 0xf8, 0x5d, 0xb2, 0xbb, 0x76, 0xa9, 0x04, 0xa5, 0x02, 0xef, 0x92, 0x94,
	0x03, 0x23, 0x15, 0xb6, 0xed, 0xd6, 0x31, 0x33, 0x68, 0x96, 0xa3, 0x99, 0x45, 0xb3, 0x0e, 0x84,
	0xdd, 0x1c, 0xdc, 0xfe, 0x3f, 0xff, 0xa8, 0x39, 0xbd, 0x95, 0xa7, 0x35, 0x20, 0xa5, 0xeb, 0xbc,
	0xb1, 0x17, 0x90, 0xb2, 0x8d, 0xf6, 0x0e, 0x58, 0xd1, 0x9b, 0xfd, 0x5e, 0xb0, 0x7a, 0xf8, 0x57,
	0x36, 0x1b, 0xf8, 0x8d, 0xa7, 0xb7, 0xaf, 0xd7, 0x7f, 0xbe, 0x57, 0xe7, 0x7d, 0x08, 0x41, 0xea,
	0x94, 0x9b, 0x4b, 0x5a, 0x04, 0x7f, 0xb0, 0x3d, 0x1e, 0xdb, 0x57, 0xf3, 0x25, 0x75, 0x17, 0x4b,
	0xea, 0x7e, 0x2e, 0xa9, 0xfb, 0x92, 0x51, 0x67, 0x91, 0x51, 0xe7, 0x3d, 0xa3, 0xce, 0xcd, 0x49,
	0x14, 0xeb, 0xdb, 0x49, 0xc0, 0x42, 0x1c, 0xf1, 0x8e, 0x49, 0x39, 0x95, 0xa0, 0xd7, 0x89, 0xf7,
	0x36, 0x53, 0x4f, 0x13, 0x50, 0xc1, 0x46, 0x71, 0xc3, 0xb3, 0x9f, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x0f, 0xea, 0x6d, 0x85, 0xb3, 0x01, 0x00, 0x00,
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

// QueryServer is the server API for Query service.
type QueryServer interface {
	Balance(context.Context, *BalanceRequest) (*BalanceResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) Balance(ctx context.Context, req *BalanceRequest) (*BalanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Balance not implemented")
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

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "token.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Balance",
			Handler:    _Query_Balance_Handler,
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
