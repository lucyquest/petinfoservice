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
	PetInfoService_UpdateName_FullMethodName  = "/protocol.PetInfoService/UpdateName"
	PetInfoService_UpdateAge_FullMethodName   = "/protocol.PetInfoService/UpdateAge"
	PetInfoService_Add_FullMethodName         = "/protocol.PetInfoService/Add"
)

// PetInfoServiceClient is the client API for PetInfoService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PetInfoServiceClient interface {
	Get(ctx context.Context, in *PetGetRequest, opts ...grpc.CallOption) (*PetGetResponse, error)
	GetMultiple(ctx context.Context, in *PetGetMultipleRequest, opts ...grpc.CallOption) (*PetGetMultipleResponse, error)
	UpdateName(ctx context.Context, in *PetUpdateNameRequest, opts ...grpc.CallOption) (*PetUpdateNameResponse, error)
	UpdateAge(ctx context.Context, in *PetUpdateDateOfBirthRequest, opts ...grpc.CallOption) (*PetUpdateDateOfBirthResponse, error)
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

func (c *petInfoServiceClient) UpdateName(ctx context.Context, in *PetUpdateNameRequest, opts ...grpc.CallOption) (*PetUpdateNameResponse, error) {
	out := new(PetUpdateNameResponse)
	err := c.cc.Invoke(ctx, PetInfoService_UpdateName_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *petInfoServiceClient) UpdateAge(ctx context.Context, in *PetUpdateDateOfBirthRequest, opts ...grpc.CallOption) (*PetUpdateDateOfBirthResponse, error) {
	out := new(PetUpdateDateOfBirthResponse)
	err := c.cc.Invoke(ctx, PetInfoService_UpdateAge_FullMethodName, in, out, opts...)
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
	UpdateName(context.Context, *PetUpdateNameRequest) (*PetUpdateNameResponse, error)
	UpdateAge(context.Context, *PetUpdateDateOfBirthRequest) (*PetUpdateDateOfBirthResponse, error)
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
func (UnimplementedPetInfoServiceServer) UpdateName(context.Context, *PetUpdateNameRequest) (*PetUpdateNameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateName not implemented")
}
func (UnimplementedPetInfoServiceServer) UpdateAge(context.Context, *PetUpdateDateOfBirthRequest) (*PetUpdateDateOfBirthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAge not implemented")
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

func _PetInfoService_UpdateName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PetUpdateNameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PetInfoServiceServer).UpdateName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PetInfoService_UpdateName_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PetInfoServiceServer).UpdateName(ctx, req.(*PetUpdateNameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PetInfoService_UpdateAge_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PetUpdateDateOfBirthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PetInfoServiceServer).UpdateAge(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PetInfoService_UpdateAge_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PetInfoServiceServer).UpdateAge(ctx, req.(*PetUpdateDateOfBirthRequest))
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
			MethodName: "UpdateName",
			Handler:    _PetInfoService_UpdateName_Handler,
		},
		{
			MethodName: "UpdateAge",
			Handler:    _PetInfoService_UpdateAge_Handler,
		},
		{
			MethodName: "Add",
			Handler:    _PetInfoService_Add_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pet.proto",
}
