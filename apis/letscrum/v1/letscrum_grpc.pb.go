// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: api/letscrum/v1/letscrum.proto

package v1

import (
	context "context"
	v1 "github.com/letscrum/letscrum/api/general/v1"
	v11 "github.com/letscrum/letscrum/api/project/v1"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// LetscrumClient is the client API for Letscrum service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LetscrumClient interface {
	GetVersion(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*v1.GetVersionResponse, error)
}

type letscrumClient struct {
	cc grpc.ClientConnInterface
}

func NewLetscrumClient(cc grpc.ClientConnInterface) LetscrumClient {
	return &letscrumClient{cc}
}

func (c *letscrumClient) GetVersion(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*v1.GetVersionResponse, error) {
	out := new(v1.GetVersionResponse)
	err := c.cc.Invoke(ctx, "/Letscrum/GetVersion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LetscrumServer is the server API for Letscrum service.
// All implementations must embed UnimplementedLetscrumServer
// for forward compatibility
type LetscrumServer interface {
	GetVersion(context.Context, *emptypb.Empty) (*v1.GetVersionResponse, error)
	mustEmbedUnimplementedLetscrumServer()
}

// UnimplementedLetscrumServer must be embedded to have forward compatible implementations.
type UnimplementedLetscrumServer struct {
}

func (UnimplementedLetscrumServer) GetVersion(context.Context, *emptypb.Empty) (*v1.GetVersionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVersion not implemented")
}
func (UnimplementedLetscrumServer) mustEmbedUnimplementedLetscrumServer() {}

// UnsafeLetscrumServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LetscrumServer will
// result in compilation errors.
type UnsafeLetscrumServer interface {
	mustEmbedUnimplementedLetscrumServer()
}

func RegisterLetscrumServer(s grpc.ServiceRegistrar, srv LetscrumServer) {
	s.RegisterService(&Letscrum_ServiceDesc, srv)
}

func _Letscrum_GetVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LetscrumServer).GetVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Letscrum/GetVersion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LetscrumServer).GetVersion(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Letscrum_ServiceDesc is the grpc.ServiceDesc for Letscrum service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Letscrum_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Letscrum",
	HandlerType: (*LetscrumServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetVersion",
			Handler:    _Letscrum_GetVersion_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/letscrum/v1/letscrum.proto",
}

// ProjectClient is the client API for Project service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProjectClient interface {
	CreateProject(ctx context.Context, in *v11.CreateProjectRequest, opts ...grpc.CallOption) (*v11.CreateProjectResponse, error)
	UpdateProject(ctx context.Context, in *v11.UpdateProjectRequest, opts ...grpc.CallOption) (*v11.UpdateProjectResponse, error)
	DeleteProject(ctx context.Context, in *v11.DeleteProjectRequest, opts ...grpc.CallOption) (*v11.DeleteProjectResponse, error)
	ListProject(ctx context.Context, in *v11.ListProjectRequest, opts ...grpc.CallOption) (*v11.ListProjectResponse, error)
	GetProject(ctx context.Context, in *v11.GetProjectRequest, opts ...grpc.CallOption) (*v11.GetProjectResponse, error)
	GetSprint(ctx context.Context, in *v11.CreateSprintRequest, opts ...grpc.CallOption) (*v11.CreateSprintResponse, error)
}

type projectClient struct {
	cc grpc.ClientConnInterface
}

func NewProjectClient(cc grpc.ClientConnInterface) ProjectClient {
	return &projectClient{cc}
}

func (c *projectClient) CreateProject(ctx context.Context, in *v11.CreateProjectRequest, opts ...grpc.CallOption) (*v11.CreateProjectResponse, error) {
	out := new(v11.CreateProjectResponse)
	err := c.cc.Invoke(ctx, "/Project/CreateProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectClient) UpdateProject(ctx context.Context, in *v11.UpdateProjectRequest, opts ...grpc.CallOption) (*v11.UpdateProjectResponse, error) {
	out := new(v11.UpdateProjectResponse)
	err := c.cc.Invoke(ctx, "/Project/UpdateProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectClient) DeleteProject(ctx context.Context, in *v11.DeleteProjectRequest, opts ...grpc.CallOption) (*v11.DeleteProjectResponse, error) {
	out := new(v11.DeleteProjectResponse)
	err := c.cc.Invoke(ctx, "/Project/DeleteProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectClient) ListProject(ctx context.Context, in *v11.ListProjectRequest, opts ...grpc.CallOption) (*v11.ListProjectResponse, error) {
	out := new(v11.ListProjectResponse)
	err := c.cc.Invoke(ctx, "/Project/ListProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectClient) GetProject(ctx context.Context, in *v11.GetProjectRequest, opts ...grpc.CallOption) (*v11.GetProjectResponse, error) {
	out := new(v11.GetProjectResponse)
	err := c.cc.Invoke(ctx, "/Project/GetProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectClient) GetSprint(ctx context.Context, in *v11.CreateSprintRequest, opts ...grpc.CallOption) (*v11.CreateSprintResponse, error) {
	out := new(v11.CreateSprintResponse)
	err := c.cc.Invoke(ctx, "/Project/GetSprint", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProjectServer is the server API for Project service.
// All implementations must embed UnimplementedProjectServer
// for forward compatibility
type ProjectServer interface {
	CreateProject(context.Context, *v11.CreateProjectRequest) (*v11.CreateProjectResponse, error)
	UpdateProject(context.Context, *v11.UpdateProjectRequest) (*v11.UpdateProjectResponse, error)
	DeleteProject(context.Context, *v11.DeleteProjectRequest) (*v11.DeleteProjectResponse, error)
	ListProject(context.Context, *v11.ListProjectRequest) (*v11.ListProjectResponse, error)
	GetProject(context.Context, *v11.GetProjectRequest) (*v11.GetProjectResponse, error)
	GetSprint(context.Context, *v11.CreateSprintRequest) (*v11.CreateSprintResponse, error)
	mustEmbedUnimplementedProjectServer()
}

// UnimplementedProjectServer must be embedded to have forward compatible implementations.
type UnimplementedProjectServer struct {
}

func (UnimplementedProjectServer) CreateProject(context.Context, *v11.CreateProjectRequest) (*v11.CreateProjectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateProject not implemented")
}
func (UnimplementedProjectServer) UpdateProject(context.Context, *v11.UpdateProjectRequest) (*v11.UpdateProjectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateProject not implemented")
}
func (UnimplementedProjectServer) DeleteProject(context.Context, *v11.DeleteProjectRequest) (*v11.DeleteProjectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteProject not implemented")
}
func (UnimplementedProjectServer) ListProject(context.Context, *v11.ListProjectRequest) (*v11.ListProjectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListProject not implemented")
}
func (UnimplementedProjectServer) GetProject(context.Context, *v11.GetProjectRequest) (*v11.GetProjectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProject not implemented")
}
func (UnimplementedProjectServer) GetSprint(context.Context, *v11.CreateSprintRequest) (*v11.CreateSprintResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSprint not implemented")
}
func (UnimplementedProjectServer) mustEmbedUnimplementedProjectServer() {}

// UnsafeProjectServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProjectServer will
// result in compilation errors.
type UnsafeProjectServer interface {
	mustEmbedUnimplementedProjectServer()
}

func RegisterProjectServer(s grpc.ServiceRegistrar, srv ProjectServer) {
	s.RegisterService(&Project_ServiceDesc, srv)
}

func _Project_CreateProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v11.CreateProjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServer).CreateProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Project/CreateProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServer).CreateProject(ctx, req.(*v11.CreateProjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Project_UpdateProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v11.UpdateProjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServer).UpdateProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Project/UpdateProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServer).UpdateProject(ctx, req.(*v11.UpdateProjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Project_DeleteProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v11.DeleteProjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServer).DeleteProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Project/DeleteProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServer).DeleteProject(ctx, req.(*v11.DeleteProjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Project_ListProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v11.ListProjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServer).ListProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Project/ListProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServer).ListProject(ctx, req.(*v11.ListProjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Project_GetProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v11.GetProjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServer).GetProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Project/GetProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServer).GetProject(ctx, req.(*v11.GetProjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Project_GetSprint_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v11.CreateSprintRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServer).GetSprint(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Project/GetSprint",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServer).GetSprint(ctx, req.(*v11.CreateSprintRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Project_ServiceDesc is the grpc.ServiceDesc for Project service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Project_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Project",
	HandlerType: (*ProjectServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateProject",
			Handler:    _Project_CreateProject_Handler,
		},
		{
			MethodName: "UpdateProject",
			Handler:    _Project_UpdateProject_Handler,
		},
		{
			MethodName: "DeleteProject",
			Handler:    _Project_DeleteProject_Handler,
		},
		{
			MethodName: "ListProject",
			Handler:    _Project_ListProject_Handler,
		},
		{
			MethodName: "GetProject",
			Handler:    _Project_GetProject_Handler,
		},
		{
			MethodName: "GetSprint",
			Handler:    _Project_GetSprint_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/letscrum/v1/letscrum.proto",
}