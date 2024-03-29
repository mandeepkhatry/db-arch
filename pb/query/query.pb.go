// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pb/query/query.proto

package query

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Query struct {
	Query                string   `protobuf:"bytes,1,opt,name=query,proto3" json:"query,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Query) Reset()         { *m = Query{} }
func (m *Query) String() string { return proto.CompactTextString(m) }
func (*Query) ProtoMessage()    {}
func (*Query) Descriptor() ([]byte, []int) {
	return fileDescriptor_f02e40967707ef49, []int{0}
}

func (m *Query) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Query.Unmarshal(m, b)
}
func (m *Query) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Query.Marshal(b, m, deterministic)
}
func (m *Query) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Query.Merge(m, src)
}
func (m *Query) XXX_Size() int {
	return xxx_messageInfo_Query.Size(m)
}
func (m *Query) XXX_DiscardUnknown() {
	xxx_messageInfo_Query.DiscardUnknown(m)
}

var xxx_messageInfo_Query proto.InternalMessageInfo

func (m *Query) GetQuery() string {
	if m != nil {
		return m.Query
	}
	return ""
}

type Response struct {
	Result               map[string][]byte `protobuf:"bytes,1,rep,name=result,proto3" json:"result,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_f02e40967707ef49, []int{1}
}

func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (m *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(m, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetResult() map[string][]byte {
	if m != nil {
		return m.Result
	}
	return nil
}

type QueryTransferRequest struct {
	Request              *Query   `protobuf:"bytes,1,opt,name=request,proto3" json:"request,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryTransferRequest) Reset()         { *m = QueryTransferRequest{} }
func (m *QueryTransferRequest) String() string { return proto.CompactTextString(m) }
func (*QueryTransferRequest) ProtoMessage()    {}
func (*QueryTransferRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f02e40967707ef49, []int{2}
}

func (m *QueryTransferRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryTransferRequest.Unmarshal(m, b)
}
func (m *QueryTransferRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryTransferRequest.Marshal(b, m, deterministic)
}
func (m *QueryTransferRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryTransferRequest.Merge(m, src)
}
func (m *QueryTransferRequest) XXX_Size() int {
	return xxx_messageInfo_QueryTransferRequest.Size(m)
}
func (m *QueryTransferRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryTransferRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryTransferRequest proto.InternalMessageInfo

func (m *QueryTransferRequest) GetRequest() *Query {
	if m != nil {
		return m.Request
	}
	return nil
}

type QueryTransferResponse struct {
	Response             []*Response `protobuf:"bytes,1,rep,name=response,proto3" json:"response,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *QueryTransferResponse) Reset()         { *m = QueryTransferResponse{} }
func (m *QueryTransferResponse) String() string { return proto.CompactTextString(m) }
func (*QueryTransferResponse) ProtoMessage()    {}
func (*QueryTransferResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_f02e40967707ef49, []int{3}
}

func (m *QueryTransferResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryTransferResponse.Unmarshal(m, b)
}
func (m *QueryTransferResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryTransferResponse.Marshal(b, m, deterministic)
}
func (m *QueryTransferResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryTransferResponse.Merge(m, src)
}
func (m *QueryTransferResponse) XXX_Size() int {
	return xxx_messageInfo_QueryTransferResponse.Size(m)
}
func (m *QueryTransferResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryTransferResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryTransferResponse proto.InternalMessageInfo

func (m *QueryTransferResponse) GetResponse() []*Response {
	if m != nil {
		return m.Response
	}
	return nil
}

func init() {
	proto.RegisterType((*Query)(nil), "pb.Query")
	proto.RegisterType((*Response)(nil), "pb.Response")
	proto.RegisterMapType((map[string][]byte)(nil), "pb.Response.ResultEntry")
	proto.RegisterType((*QueryTransferRequest)(nil), "pb.QueryTransferRequest")
	proto.RegisterType((*QueryTransferResponse)(nil), "pb.QueryTransferResponse")
}

func init() { proto.RegisterFile("pb/query/query.proto", fileDescriptor_f02e40967707ef49) }

var fileDescriptor_f02e40967707ef49 = []byte{
	// 256 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x51, 0x4d, 0x4b, 0xc4, 0x30,
	0x10, 0xb5, 0x5d, 0xf6, 0x6b, 0x5a, 0x41, 0x42, 0x85, 0xb8, 0x20, 0x94, 0x78, 0xe9, 0xa9, 0x4a,
	0xbd, 0xf8, 0x71, 0x52, 0xd0, 0xbb, 0x51, 0x3c, 0x78, 0xdb, 0xca, 0x08, 0x62, 0x69, 0xb3, 0x93,
	0x64, 0x65, 0xff, 0xbd, 0x24, 0x69, 0x97, 0x75, 0xd9, 0x4b, 0x78, 0x33, 0x6f, 0xf2, 0xe6, 0xbd,
	0x04, 0x32, 0x55, 0x5f, 0xae, 0x2c, 0xd2, 0x26, 0x9c, 0xa5, 0xa2, 0xce, 0x74, 0x2c, 0x56, 0xb5,
	0x38, 0x87, 0xf1, 0x8b, 0x6b, 0xb1, 0x0c, 0xc6, 0x9e, 0xe3, 0x51, 0x1e, 0x15, 0x73, 0x19, 0x0a,
	0xf1, 0x0b, 0x33, 0x89, 0x5a, 0x75, 0xad, 0x46, 0x76, 0x05, 0x13, 0x42, 0x6d, 0x1b, 0xc3, 0xa3,
	0x7c, 0x54, 0x24, 0x15, 0x2f, 0x55, 0x5d, 0x0e, 0xac, 0x03, 0xb6, 0x31, 0x4f, 0xad, 0xa1, 0x8d,
	0xec, 0xe7, 0x16, 0xb7, 0x90, 0xec, 0xb4, 0xd9, 0x09, 0x8c, 0x7e, 0x70, 0x58, 0xe0, 0xa0, 0x5b,
	0xba, 0x5e, 0x36, 0x16, 0x79, 0x9c, 0x47, 0x45, 0x2a, 0x43, 0x71, 0x17, 0xdf, 0x44, 0xe2, 0x1e,
	0x32, 0xef, 0xeb, 0x8d, 0x96, 0xad, 0xfe, 0x42, 0x92, 0xb8, 0xb2, 0xa8, 0x0d, 0xbb, 0x80, 0x29,
	0x05, 0xe8, 0x75, 0x92, 0x6a, 0xee, 0x5c, 0xf8, 0x51, 0x39, 0x30, 0xe2, 0x01, 0x4e, 0xf7, 0x2e,
	0xf7, 0x11, 0x0a, 0x98, 0x51, 0x8f, 0xfb, 0x10, 0xe9, 0x6e, 0x08, 0xb9, 0x65, 0xab, 0x77, 0x48,
	0xbd, 0xc4, 0x2b, 0xd2, 0xfa, 0xfb, 0x13, 0xd9, 0x33, 0x1c, 0xff, 0x93, 0x64, 0x7c, 0xbb, 0x77,
	0xcf, 0xe2, 0xe2, 0xec, 0x00, 0x13, 0x54, 0xc5, 0xd1, 0xe3, 0xf4, 0x23, 0xbc, 0x6c, 0x3d, 0xf1,
	0x7f, 0x70, 0xfd, 0x17, 0x00, 0x00, 0xff, 0xff, 0xfc, 0xec, 0xb2, 0xa3, 0x9b, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryServiceClient is the client API for QueryService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryServiceClient interface {
	QueryTransfer(ctx context.Context, in *QueryTransferRequest, opts ...grpc.CallOption) (*QueryTransferResponse, error)
}

type queryServiceClient struct {
	cc *grpc.ClientConn
}

func NewQueryServiceClient(cc *grpc.ClientConn) QueryServiceClient {
	return &queryServiceClient{cc}
}

func (c *queryServiceClient) QueryTransfer(ctx context.Context, in *QueryTransferRequest, opts ...grpc.CallOption) (*QueryTransferResponse, error) {
	out := new(QueryTransferResponse)
	err := c.cc.Invoke(ctx, "/pb.QueryService/QueryTransfer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServiceServer is the server API for QueryService service.
type QueryServiceServer interface {
	QueryTransfer(context.Context, *QueryTransferRequest) (*QueryTransferResponse, error)
}

// UnimplementedQueryServiceServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServiceServer struct {
}

func (*UnimplementedQueryServiceServer) QueryTransfer(ctx context.Context, req *QueryTransferRequest) (*QueryTransferResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryTransfer not implemented")
}

func RegisterQueryServiceServer(s *grpc.Server, srv QueryServiceServer) {
	s.RegisterService(&_QueryService_serviceDesc, srv)
}

func _QueryService_QueryTransfer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryTransferRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServiceServer).QueryTransfer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.QueryService/QueryTransfer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServiceServer).QueryTransfer(ctx, req.(*QueryTransferRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _QueryService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.QueryService",
	HandlerType: (*QueryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "QueryTransfer",
			Handler:    _QueryService_QueryTransfer_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pb/query/query.proto",
}
