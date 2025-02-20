// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: hyperlane/core/v1/tx.proto

package corev1

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
	Msg_CreateMailbox_FullMethodName           = "/hyperlane.core.v1.Msg/CreateMailbox"
	Msg_DispatchMessage_FullMethodName         = "/hyperlane.core.v1.Msg/DispatchMessage"
	Msg_ProcessMessage_FullMethodName          = "/hyperlane.core.v1.Msg/ProcessMessage"
	Msg_AnnounceValidator_FullMethodName       = "/hyperlane.core.v1.Msg/AnnounceValidator"
	Msg_CreateIgp_FullMethodName               = "/hyperlane.core.v1.Msg/CreateIgp"
	Msg_SetDestinationGasConfig_FullMethodName = "/hyperlane.core.v1.Msg/SetDestinationGasConfig"
	Msg_PayForGas_FullMethodName               = "/hyperlane.core.v1.Msg/PayForGas"
	Msg_Claim_FullMethodName                   = "/hyperlane.core.v1.Msg/Claim"
	Msg_UpdateParams_FullMethodName            = "/hyperlane.core.v1.Msg/UpdateParams"
)

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MsgClient interface {
	// CreateMailbox ...
	CreateMailbox(ctx context.Context, in *MsgCreateMailbox, opts ...grpc.CallOption) (*MsgCreateMailboxResponse, error)
	// DispatchMessage ...
	DispatchMessage(ctx context.Context, in *MsgDispatchMessage, opts ...grpc.CallOption) (*MsgDispatchMessageResponse, error)
	// ProcessMessage ...
	ProcessMessage(ctx context.Context, in *MsgProcessMessage, opts ...grpc.CallOption) (*MsgProcessMessageResponse, error)
	// AnnounceValidator ...
	AnnounceValidator(ctx context.Context, in *MsgAnnounceValidator, opts ...grpc.CallOption) (*MsgAnnounceValidatorResponse, error)
	// CreateIgp ...
	CreateIgp(ctx context.Context, in *MsgCreateIgp, opts ...grpc.CallOption) (*MsgCreateIgpResponse, error)
	// SetDestinationGasConfig ...
	SetDestinationGasConfig(ctx context.Context, in *MsgSetDestinationGasConfig, opts ...grpc.CallOption) (*MsgSetDestinationGasConfigResponse, error)
	// PayForGas ...
	PayForGas(ctx context.Context, in *MsgPayForGas, opts ...grpc.CallOption) (*MsgPayForGasResponse, error)
	// Claim ...
	Claim(ctx context.Context, in *MsgClaim, opts ...grpc.CallOption) (*MsgClaimResponse, error)
	// UpdateParams updates the module parameters.
	UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error)
}

type msgClient struct {
	cc grpc.ClientConnInterface
}

func NewMsgClient(cc grpc.ClientConnInterface) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) CreateMailbox(ctx context.Context, in *MsgCreateMailbox, opts ...grpc.CallOption) (*MsgCreateMailboxResponse, error) {
	out := new(MsgCreateMailboxResponse)
	err := c.cc.Invoke(ctx, Msg_CreateMailbox_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) DispatchMessage(ctx context.Context, in *MsgDispatchMessage, opts ...grpc.CallOption) (*MsgDispatchMessageResponse, error) {
	out := new(MsgDispatchMessageResponse)
	err := c.cc.Invoke(ctx, Msg_DispatchMessage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) ProcessMessage(ctx context.Context, in *MsgProcessMessage, opts ...grpc.CallOption) (*MsgProcessMessageResponse, error) {
	out := new(MsgProcessMessageResponse)
	err := c.cc.Invoke(ctx, Msg_ProcessMessage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) AnnounceValidator(ctx context.Context, in *MsgAnnounceValidator, opts ...grpc.CallOption) (*MsgAnnounceValidatorResponse, error) {
	out := new(MsgAnnounceValidatorResponse)
	err := c.cc.Invoke(ctx, Msg_AnnounceValidator_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) CreateIgp(ctx context.Context, in *MsgCreateIgp, opts ...grpc.CallOption) (*MsgCreateIgpResponse, error) {
	out := new(MsgCreateIgpResponse)
	err := c.cc.Invoke(ctx, Msg_CreateIgp_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) SetDestinationGasConfig(ctx context.Context, in *MsgSetDestinationGasConfig, opts ...grpc.CallOption) (*MsgSetDestinationGasConfigResponse, error) {
	out := new(MsgSetDestinationGasConfigResponse)
	err := c.cc.Invoke(ctx, Msg_SetDestinationGasConfig_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) PayForGas(ctx context.Context, in *MsgPayForGas, opts ...grpc.CallOption) (*MsgPayForGasResponse, error) {
	out := new(MsgPayForGasResponse)
	err := c.cc.Invoke(ctx, Msg_PayForGas_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) Claim(ctx context.Context, in *MsgClaim, opts ...grpc.CallOption) (*MsgClaimResponse, error) {
	out := new(MsgClaimResponse)
	err := c.cc.Invoke(ctx, Msg_Claim_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error) {
	out := new(MsgUpdateParamsResponse)
	err := c.cc.Invoke(ctx, Msg_UpdateParams_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
// All implementations must embed UnimplementedMsgServer
// for forward compatibility
type MsgServer interface {
	// CreateMailbox ...
	CreateMailbox(context.Context, *MsgCreateMailbox) (*MsgCreateMailboxResponse, error)
	// DispatchMessage ...
	DispatchMessage(context.Context, *MsgDispatchMessage) (*MsgDispatchMessageResponse, error)
	// ProcessMessage ...
	ProcessMessage(context.Context, *MsgProcessMessage) (*MsgProcessMessageResponse, error)
	// AnnounceValidator ...
	AnnounceValidator(context.Context, *MsgAnnounceValidator) (*MsgAnnounceValidatorResponse, error)
	// CreateIgp ...
	CreateIgp(context.Context, *MsgCreateIgp) (*MsgCreateIgpResponse, error)
	// SetDestinationGasConfig ...
	SetDestinationGasConfig(context.Context, *MsgSetDestinationGasConfig) (*MsgSetDestinationGasConfigResponse, error)
	// PayForGas ...
	PayForGas(context.Context, *MsgPayForGas) (*MsgPayForGasResponse, error)
	// Claim ...
	Claim(context.Context, *MsgClaim) (*MsgClaimResponse, error)
	// UpdateParams updates the module parameters.
	UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error)
	mustEmbedUnimplementedMsgServer()
}

// UnimplementedMsgServer must be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (UnimplementedMsgServer) CreateMailbox(context.Context, *MsgCreateMailbox) (*MsgCreateMailboxResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateMailbox not implemented")
}
func (UnimplementedMsgServer) DispatchMessage(context.Context, *MsgDispatchMessage) (*MsgDispatchMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DispatchMessage not implemented")
}
func (UnimplementedMsgServer) ProcessMessage(context.Context, *MsgProcessMessage) (*MsgProcessMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProcessMessage not implemented")
}
func (UnimplementedMsgServer) AnnounceValidator(context.Context, *MsgAnnounceValidator) (*MsgAnnounceValidatorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AnnounceValidator not implemented")
}
func (UnimplementedMsgServer) CreateIgp(context.Context, *MsgCreateIgp) (*MsgCreateIgpResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateIgp not implemented")
}
func (UnimplementedMsgServer) SetDestinationGasConfig(context.Context, *MsgSetDestinationGasConfig) (*MsgSetDestinationGasConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetDestinationGasConfig not implemented")
}
func (UnimplementedMsgServer) PayForGas(context.Context, *MsgPayForGas) (*MsgPayForGasResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PayForGas not implemented")
}
func (UnimplementedMsgServer) Claim(context.Context, *MsgClaim) (*MsgClaimResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Claim not implemented")
}
func (UnimplementedMsgServer) UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateParams not implemented")
}
func (UnimplementedMsgServer) mustEmbedUnimplementedMsgServer() {}

// UnsafeMsgServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MsgServer will
// result in compilation errors.
type UnsafeMsgServer interface {
	mustEmbedUnimplementedMsgServer()
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&Msg_ServiceDesc, srv)
}

func _Msg_CreateMailbox_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateMailbox)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreateMailbox(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_CreateMailbox_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreateMailbox(ctx, req.(*MsgCreateMailbox))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_DispatchMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgDispatchMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).DispatchMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_DispatchMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).DispatchMessage(ctx, req.(*MsgDispatchMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_ProcessMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgProcessMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ProcessMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_ProcessMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ProcessMessage(ctx, req.(*MsgProcessMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_AnnounceValidator_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgAnnounceValidator)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).AnnounceValidator(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_AnnounceValidator_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).AnnounceValidator(ctx, req.(*MsgAnnounceValidator))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_CreateIgp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateIgp)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreateIgp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_CreateIgp_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreateIgp(ctx, req.(*MsgCreateIgp))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_SetDestinationGasConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSetDestinationGasConfig)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SetDestinationGasConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_SetDestinationGasConfig_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SetDestinationGasConfig(ctx, req.(*MsgSetDestinationGasConfig))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_PayForGas_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgPayForGas)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).PayForGas(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_PayForGas_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).PayForGas(ctx, req.(*MsgPayForGas))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_Claim_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgClaim)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Claim(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_Claim_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Claim(ctx, req.(*MsgClaim))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_UpdateParams_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateParams(ctx, req.(*MsgUpdateParams))
	}
	return interceptor(ctx, in, info, handler)
}

// Msg_ServiceDesc is the grpc.ServiceDesc for Msg service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "hyperlane.core.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateMailbox",
			Handler:    _Msg_CreateMailbox_Handler,
		},
		{
			MethodName: "DispatchMessage",
			Handler:    _Msg_DispatchMessage_Handler,
		},
		{
			MethodName: "ProcessMessage",
			Handler:    _Msg_ProcessMessage_Handler,
		},
		{
			MethodName: "AnnounceValidator",
			Handler:    _Msg_AnnounceValidator_Handler,
		},
		{
			MethodName: "CreateIgp",
			Handler:    _Msg_CreateIgp_Handler,
		},
		{
			MethodName: "SetDestinationGasConfig",
			Handler:    _Msg_SetDestinationGasConfig_Handler,
		},
		{
			MethodName: "PayForGas",
			Handler:    _Msg_PayForGas_Handler,
		},
		{
			MethodName: "Claim",
			Handler:    _Msg_Claim_Handler,
		},
		{
			MethodName: "UpdateParams",
			Handler:    _Msg_UpdateParams_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "hyperlane/core/v1/tx.proto",
}
