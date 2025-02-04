// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package assessment

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AssessmentClient is the client API for Assessment service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AssessmentClient interface {
	// Triggers the assessment. Part of the private API. Not exposed as REST.
	TriggerAssessment(ctx context.Context, in *TriggerAssessmentRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// List all assessment results. Part of the public API, also exposed as REST.
	ListAssessmentResults(ctx context.Context, in *ListAssessmentResultsRequest, opts ...grpc.CallOption) (*ListAssessmentResultsResponse, error)
	// Assesses the evidence sent by the discovery. Part of the public API, also exposed as REST.
	AssessEvidence(ctx context.Context, in *AssessEvidenceRequest, opts ...grpc.CallOption) (*AssessEvidenceResponse, error)
	// Assesses stream of evidences sent by the discovery. Part of the public API. Not exposed as REST.
	AssessEvidences(ctx context.Context, opts ...grpc.CallOption) (Assessment_AssessEvidencesClient, error)
}

type assessmentClient struct {
	cc grpc.ClientConnInterface
}

func NewAssessmentClient(cc grpc.ClientConnInterface) AssessmentClient {
	return &assessmentClient{cc}
}

func (c *assessmentClient) TriggerAssessment(ctx context.Context, in *TriggerAssessmentRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/clouditor.Assessment/TriggerAssessment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *assessmentClient) ListAssessmentResults(ctx context.Context, in *ListAssessmentResultsRequest, opts ...grpc.CallOption) (*ListAssessmentResultsResponse, error) {
	out := new(ListAssessmentResultsResponse)
	err := c.cc.Invoke(ctx, "/clouditor.Assessment/ListAssessmentResults", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *assessmentClient) AssessEvidence(ctx context.Context, in *AssessEvidenceRequest, opts ...grpc.CallOption) (*AssessEvidenceResponse, error) {
	out := new(AssessEvidenceResponse)
	err := c.cc.Invoke(ctx, "/clouditor.Assessment/AssessEvidence", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *assessmentClient) AssessEvidences(ctx context.Context, opts ...grpc.CallOption) (Assessment_AssessEvidencesClient, error) {
	stream, err := c.cc.NewStream(ctx, &Assessment_ServiceDesc.Streams[0], "/clouditor.Assessment/AssessEvidences", opts...)
	if err != nil {
		return nil, err
	}
	x := &assessmentAssessEvidencesClient{stream}
	return x, nil
}

type Assessment_AssessEvidencesClient interface {
	Send(*AssessEvidenceRequest) error
	CloseAndRecv() (*emptypb.Empty, error)
	grpc.ClientStream
}

type assessmentAssessEvidencesClient struct {
	grpc.ClientStream
}

func (x *assessmentAssessEvidencesClient) Send(m *AssessEvidenceRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *assessmentAssessEvidencesClient) CloseAndRecv() (*emptypb.Empty, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(emptypb.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// AssessmentServer is the server API for Assessment service.
// All implementations must embed UnimplementedAssessmentServer
// for forward compatibility
type AssessmentServer interface {
	// Triggers the assessment. Part of the private API. Not exposed as REST.
	TriggerAssessment(context.Context, *TriggerAssessmentRequest) (*emptypb.Empty, error)
	// List all assessment results. Part of the public API, also exposed as REST.
	ListAssessmentResults(context.Context, *ListAssessmentResultsRequest) (*ListAssessmentResultsResponse, error)
	// Assesses the evidence sent by the discovery. Part of the public API, also exposed as REST.
	AssessEvidence(context.Context, *AssessEvidenceRequest) (*AssessEvidenceResponse, error)
	// Assesses stream of evidences sent by the discovery. Part of the public API. Not exposed as REST.
	AssessEvidences(Assessment_AssessEvidencesServer) error
	mustEmbedUnimplementedAssessmentServer()
}

// UnimplementedAssessmentServer must be embedded to have forward compatible implementations.
type UnimplementedAssessmentServer struct {
}

func (UnimplementedAssessmentServer) TriggerAssessment(context.Context, *TriggerAssessmentRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TriggerAssessment not implemented")
}
func (UnimplementedAssessmentServer) ListAssessmentResults(context.Context, *ListAssessmentResultsRequest) (*ListAssessmentResultsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAssessmentResults not implemented")
}
func (UnimplementedAssessmentServer) AssessEvidence(context.Context, *AssessEvidenceRequest) (*AssessEvidenceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AssessEvidence not implemented")
}
func (UnimplementedAssessmentServer) AssessEvidences(Assessment_AssessEvidencesServer) error {
	return status.Errorf(codes.Unimplemented, "method AssessEvidences not implemented")
}
func (UnimplementedAssessmentServer) mustEmbedUnimplementedAssessmentServer() {}

// UnsafeAssessmentServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AssessmentServer will
// result in compilation errors.
type UnsafeAssessmentServer interface {
	mustEmbedUnimplementedAssessmentServer()
}

func RegisterAssessmentServer(s grpc.ServiceRegistrar, srv AssessmentServer) {
	s.RegisterService(&Assessment_ServiceDesc, srv)
}

func _Assessment_TriggerAssessment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TriggerAssessmentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AssessmentServer).TriggerAssessment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clouditor.Assessment/TriggerAssessment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AssessmentServer).TriggerAssessment(ctx, req.(*TriggerAssessmentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Assessment_ListAssessmentResults_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListAssessmentResultsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AssessmentServer).ListAssessmentResults(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clouditor.Assessment/ListAssessmentResults",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AssessmentServer).ListAssessmentResults(ctx, req.(*ListAssessmentResultsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Assessment_AssessEvidence_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AssessEvidenceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AssessmentServer).AssessEvidence(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clouditor.Assessment/AssessEvidence",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AssessmentServer).AssessEvidence(ctx, req.(*AssessEvidenceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Assessment_AssessEvidences_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(AssessmentServer).AssessEvidences(&assessmentAssessEvidencesServer{stream})
}

type Assessment_AssessEvidencesServer interface {
	SendAndClose(*emptypb.Empty) error
	Recv() (*AssessEvidenceRequest, error)
	grpc.ServerStream
}

type assessmentAssessEvidencesServer struct {
	grpc.ServerStream
}

func (x *assessmentAssessEvidencesServer) SendAndClose(m *emptypb.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *assessmentAssessEvidencesServer) Recv() (*AssessEvidenceRequest, error) {
	m := new(AssessEvidenceRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Assessment_ServiceDesc is the grpc.ServiceDesc for Assessment service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Assessment_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "clouditor.Assessment",
	HandlerType: (*AssessmentServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "TriggerAssessment",
			Handler:    _Assessment_TriggerAssessment_Handler,
		},
		{
			MethodName: "ListAssessmentResults",
			Handler:    _Assessment_ListAssessmentResults_Handler,
		},
		{
			MethodName: "AssessEvidence",
			Handler:    _Assessment_AssessEvidence_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "AssessEvidences",
			Handler:       _Assessment_AssessEvidences_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "assessment.proto",
}
