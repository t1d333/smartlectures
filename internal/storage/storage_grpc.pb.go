// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.6
// source: storage/storage.proto

package service

import (
	context "context"
	status "google.golang.org/genproto/googleapis/rpc/status"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status1 "google.golang.org/grpc/status"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// NotesClient is the client API for Notes service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NotesClient interface {
	CreateNote(ctx context.Context, in *Note, opts ...grpc.CallOption) (*status.Status, error)
	UpdateNote(ctx context.Context, in *NoteUpdateRequest, opts ...grpc.CallOption) (*status.Status, error)
	CreateDir(ctx context.Context, in *Dir, opts ...grpc.CallOption) (*status.Status, error)
	DeleteDir(ctx context.Context, in *wrapperspb.Int32Value, opts ...grpc.CallOption) (*status.Status, error)
	DeleteNote(ctx context.Context, in *wrapperspb.Int32Value, opts ...grpc.CallOption) (*status.Status, error)
	Search(ctx context.Context, in *wrapperspb.StringValue, opts ...grpc.CallOption) (*SearchResponse, error)
}

type notesClient struct {
	cc grpc.ClientConnInterface
}

func NewNotesClient(cc grpc.ClientConnInterface) NotesClient {
	return &notesClient{cc}
}

func (c *notesClient) CreateNote(ctx context.Context, in *Note, opts ...grpc.CallOption) (*status.Status, error) {
	out := new(status.Status)
	err := c.cc.Invoke(ctx, "/Notes/CreateNote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notesClient) UpdateNote(ctx context.Context, in *NoteUpdateRequest, opts ...grpc.CallOption) (*status.Status, error) {
	out := new(status.Status)
	err := c.cc.Invoke(ctx, "/Notes/UpdateNote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notesClient) CreateDir(ctx context.Context, in *Dir, opts ...grpc.CallOption) (*status.Status, error) {
	out := new(status.Status)
	err := c.cc.Invoke(ctx, "/Notes/CreateDir", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notesClient) DeleteDir(ctx context.Context, in *wrapperspb.Int32Value, opts ...grpc.CallOption) (*status.Status, error) {
	out := new(status.Status)
	err := c.cc.Invoke(ctx, "/Notes/DeleteDir", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notesClient) DeleteNote(ctx context.Context, in *wrapperspb.Int32Value, opts ...grpc.CallOption) (*status.Status, error) {
	out := new(status.Status)
	err := c.cc.Invoke(ctx, "/Notes/DeleteNote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notesClient) Search(ctx context.Context, in *wrapperspb.StringValue, opts ...grpc.CallOption) (*SearchResponse, error) {
	out := new(SearchResponse)
	err := c.cc.Invoke(ctx, "/Notes/Search", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NotesServer is the server API for Notes service.
// All implementations must embed UnimplementedNotesServer
// for forward compatibility
type NotesServer interface {
	CreateNote(context.Context, *Note) (*status.Status, error)
	UpdateNote(context.Context, *NoteUpdateRequest) (*status.Status, error)
	CreateDir(context.Context, *Dir) (*status.Status, error)
	DeleteDir(context.Context, *wrapperspb.Int32Value) (*status.Status, error)
	DeleteNote(context.Context, *wrapperspb.Int32Value) (*status.Status, error)
	Search(context.Context, *wrapperspb.StringValue) (*SearchResponse, error)
	mustEmbedUnimplementedNotesServer()
}

// UnimplementedNotesServer must be embedded to have forward compatible implementations.
type UnimplementedNotesServer struct {
}

func (UnimplementedNotesServer) CreateNote(context.Context, *Note) (*status.Status, error) {
	return nil, status1.Errorf(codes.Unimplemented, "method CreateNote not implemented")
}
func (UnimplementedNotesServer) UpdateNote(context.Context, *NoteUpdateRequest) (*status.Status, error) {
	return nil, status1.Errorf(codes.Unimplemented, "method UpdateNote not implemented")
}
func (UnimplementedNotesServer) CreateDir(context.Context, *Dir) (*status.Status, error) {
	return nil, status1.Errorf(codes.Unimplemented, "method CreateDir not implemented")
}
func (UnimplementedNotesServer) DeleteDir(context.Context, *wrapperspb.Int32Value) (*status.Status, error) {
	return nil, status1.Errorf(codes.Unimplemented, "method DeleteDir not implemented")
}
func (UnimplementedNotesServer) DeleteNote(context.Context, *wrapperspb.Int32Value) (*status.Status, error) {
	return nil, status1.Errorf(codes.Unimplemented, "method DeleteNote not implemented")
}
func (UnimplementedNotesServer) Search(context.Context, *wrapperspb.StringValue) (*SearchResponse, error) {
	return nil, status1.Errorf(codes.Unimplemented, "method Search not implemented")
}
func (UnimplementedNotesServer) mustEmbedUnimplementedNotesServer() {}

// UnsafeNotesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NotesServer will
// result in compilation errors.
type UnsafeNotesServer interface {
	mustEmbedUnimplementedNotesServer()
}

func RegisterNotesServer(s grpc.ServiceRegistrar, srv NotesServer) {
	s.RegisterService(&Notes_ServiceDesc, srv)
}

func _Notes_CreateNote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Note)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotesServer).CreateNote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Notes/CreateNote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotesServer).CreateNote(ctx, req.(*Note))
	}
	return interceptor(ctx, in, info, handler)
}

func _Notes_UpdateNote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NoteUpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotesServer).UpdateNote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Notes/UpdateNote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotesServer).UpdateNote(ctx, req.(*NoteUpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Notes_CreateDir_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Dir)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotesServer).CreateDir(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Notes/CreateDir",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotesServer).CreateDir(ctx, req.(*Dir))
	}
	return interceptor(ctx, in, info, handler)
}

func _Notes_DeleteDir_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(wrapperspb.Int32Value)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotesServer).DeleteDir(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Notes/DeleteDir",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotesServer).DeleteDir(ctx, req.(*wrapperspb.Int32Value))
	}
	return interceptor(ctx, in, info, handler)
}

func _Notes_DeleteNote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(wrapperspb.Int32Value)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotesServer).DeleteNote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Notes/DeleteNote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotesServer).DeleteNote(ctx, req.(*wrapperspb.Int32Value))
	}
	return interceptor(ctx, in, info, handler)
}

func _Notes_Search_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(wrapperspb.StringValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotesServer).Search(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Notes/Search",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotesServer).Search(ctx, req.(*wrapperspb.StringValue))
	}
	return interceptor(ctx, in, info, handler)
}

// Notes_ServiceDesc is the grpc.ServiceDesc for Notes service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Notes_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Notes",
	HandlerType: (*NotesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateNote",
			Handler:    _Notes_CreateNote_Handler,
		},
		{
			MethodName: "UpdateNote",
			Handler:    _Notes_UpdateNote_Handler,
		},
		{
			MethodName: "CreateDir",
			Handler:    _Notes_CreateDir_Handler,
		},
		{
			MethodName: "DeleteDir",
			Handler:    _Notes_DeleteDir_Handler,
		},
		{
			MethodName: "DeleteNote",
			Handler:    _Notes_DeleteNote_Handler,
		},
		{
			MethodName: "Search",
			Handler:    _Notes_Search_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "storage/storage.proto",
}
