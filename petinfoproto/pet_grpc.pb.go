// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.19.6
// source: pet.proto

package petinfoproto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	PetInfoService_Get_FullMethodName         = "/protocol.PetInfoService/Get"
	PetInfoService_GetMultiple_FullMethodName = "/protocol.PetInfoService/GetMultiple"
	PetInfoService_EditName_FullMethodName    = "/protocol.PetInfoService/EditName"
	PetInfoService_EditAge_FullMethodName     = "/protocol.PetInfoService/EditAge"
	PetInfoService_Add_FullMethodName         = "/protocol.PetInfoService/Add"
)

// PetInfoServiceClient is the client API for PetInfoService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PetInfoServiceClient interface {
	Get(ctx context.Context, in *PetGetRequest, opts ...grpc.CallOption) (*PetGetResponse, error)
	GetMultiple(ctx context.Context, in *PetGetMultipleRequest, opts ...grpc.CallOption) (*PetGetMultipleResponse, error)
	EditName(ctx context.Context, in *PetEditNameRequest, opts ...grpc.CallOption) (*PetEditNameResponse, error)
	EditAge(ctx context.Context, in *PetEditDateOfBirthRequest, opts ...grpc.CallOption) (*PetEditDateOfBirthResponse, error)
	Add(ctx context.Context, in *PetAddRequest, opts ...grpc.CallOption) (*PetAddResponse, error)
}

type petInfoServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPetInfoServiceClient(cc grpc.ClientConnInterface) PetInfoServiceClient {
	return &petInfoServiceClient{cc}
}

func (c *petInfoServiceClient) Get(ctx context.Context, in *PetGetRequest, opts ...grpc.CallOption) (*PetGetResponse, error) {
	out := new(PetGetResponse)
	err := c.cc.Invoke(ctx, PetInfoService_Get_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *petInfoServiceClient) GetMultiple(ctx context.Context, in *PetGetMultipleRequest, opts ...grpc.CallOption) (*PetGetMultipleResponse, error) {
	out := new(PetGetMultipleResponse)
	err := c.cc.Invoke(ctx, PetInfoService_GetMultiple_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *petInfoServiceClient) EditName(ctx context.Context, in *PetEditNameRequest, opts ...grpc.CallOption) (*PetEditNameResponse, error) {
	out := new(PetEditNameResponse)
	err := c.cc.Invoke(ctx, PetInfoService_EditName_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *petInfoServiceClient) EditAge(ctx context.Context, in *PetEditDateOfBirthRequest, opts ...grpc.CallOption) (*PetEditDateOfBirthResponse, error) {
	out := new(PetEditDateOfBirthResponse)
	err := c.cc.Invoke(ctx, PetInfoService_EditAge_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *petInfoServiceClient) Add(ctx context.Context, in *PetAddRequest, opts ...grpc.CallOption) (*PetAddResponse, error) {
	out := new(PetAddResponse)
	err := c.cc.Invoke(ctx, PetInfoService_Add_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PetInfoServiceServer is the server API for PetInfoService service.
// All implementations must embed UnimplementedPetInfoServiceServer
// for forward compatibility
type PetInfoServiceServer interface {
	Get(context.Context, *PetGetRequest) (*PetGetResponse, error)
	GetMultiple(context.Context, *PetGetMultipleRequest) (*PetGetMultipleResponse, error)
	EditName(context.Context, *PetEditNameRequest) (*PetEditNameResponse, error)
	EditAge(context.Context, *PetEditDateOfBirthRequest) (*PetEditDateOfBirthResponse, error)
	Add(context.Context, *PetAddRequest) (*PetAddResponse, error)
	mustEmbedUnimplementedPetInfoServiceServer()
}

// UnimplementedPetInfoServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPetInfoServiceServer struct {
}

func (UnimplementedPetInfoServiceServer) Get(context.Context, *PetGetRequest) (*PetGetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedPetInfoServiceServer) GetMultiple(context.Context, *PetGetMultipleRequest) (*PetGetMultipleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMultiple not implemented")
}
func (UnimplementedPetInfoServiceServer) EditName(context.Context, *PetEditNameRequest) (*PetEditNameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditName not implemented")
}
func (UnimplementedPetInfoServiceServer) EditAge(context.Context, *PetEditDateOfBirthRequest) (*PetEditDateOfBirthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditAge not implemented")
}
func (UnimplementedPetInfoServiceServer) Add(context.Context, *PetAddRequest) (*PetAddResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Add not implemented")
}
func (UnimplementedPetInfoServiceServer) mustEmbedUnimplementedPetInfoServiceServer() {}

// UnsafePetInfoServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PetInfoServiceServer will
// result in compilation errors.
type UnsafePetInfoServiceServer interface {
	mustEmbedUnimplementedPetInfoServiceServer()
}

func RegisterPetInfoServiceServer(s grpc.ServiceRegistrar, srv PetInfoServiceServer) {
	s.RegisterService(&PetInfoService_ServiceDesc, srv)
}

func _PetInfoService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PetGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PetInfoServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PetInfoService_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PetInfoServiceServer).Get(ctx, req.(*PetGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PetInfoService_GetMultiple_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PetGetMultipleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PetInfoServiceServer).GetMultiple(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PetInfoService_GetMultiple_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PetInfoServiceServer).GetMultiple(ctx, req.(*PetGetMultipleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PetInfoService_EditName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PetEditNameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PetInfoServiceServer).EditName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PetInfoService_EditName_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PetInfoServiceServer).EditName(ctx, req.(*PetEditNameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PetInfoService_EditAge_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PetEditDateOfBirthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PetInfoServiceServer).EditAge(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PetInfoService_EditAge_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PetInfoServiceServer).EditAge(ctx, req.(*PetEditDateOfBirthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PetInfoService_Add_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PetAddRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PetInfoServiceServer).Add(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PetInfoService_Add_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PetInfoServiceServer).Add(ctx, req.(*PetAddRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PetInfoService_ServiceDesc is the grpc.ServiceDesc for PetInfoService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PetInfoService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protocol.PetInfoService",
	HandlerType: (*PetInfoServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _PetInfoService_Get_Handler,
		},
		{
			MethodName: "GetMultiple",
			Handler:    _PetInfoService_GetMultiple_Handler,
		},
		{
			MethodName: "EditName",
			Handler:    _PetInfoService_EditName_Handler,
		},
		{
			MethodName: "EditAge",
			Handler:    _PetInfoService_EditAge_Handler,
		},
		{
			MethodName: "Add",
			Handler:    _PetInfoService_Add_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pet.proto",
}
