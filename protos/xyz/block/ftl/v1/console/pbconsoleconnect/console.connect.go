// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: xyz/block/ftl/v1/console/console.proto

package pbconsoleconnect

import (
	context "context"
	errors "errors"
	v1 "github.com/TBD54566975/ftl/protos/xyz/block/ftl/v1"
	console "github.com/TBD54566975/ftl/protos/xyz/block/ftl/v1/console"
	connect_go "github.com/bufbuild/connect-go"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect_go.IsAtLeastVersion1_7_0

const (
	// ConsoleServiceName is the fully-qualified name of the ConsoleService service.
	ConsoleServiceName = "xyz.block.ftl.v1.console.ConsoleService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// ConsoleServicePingProcedure is the fully-qualified name of the ConsoleService's Ping RPC.
	ConsoleServicePingProcedure = "/xyz.block.ftl.v1.console.ConsoleService/Ping"
	// ConsoleServiceGetModulesProcedure is the fully-qualified name of the ConsoleService's GetModules
	// RPC.
	ConsoleServiceGetModulesProcedure = "/xyz.block.ftl.v1.console.ConsoleService/GetModules"
	// ConsoleServiceGetCallsProcedure is the fully-qualified name of the ConsoleService's GetCalls RPC.
	ConsoleServiceGetCallsProcedure = "/xyz.block.ftl.v1.console.ConsoleService/GetCalls"
	// ConsoleServiceGetRequestCallsProcedure is the fully-qualified name of the ConsoleService's
	// GetRequestCalls RPC.
	ConsoleServiceGetRequestCallsProcedure = "/xyz.block.ftl.v1.console.ConsoleService/GetRequestCalls"
	// ConsoleServiceStreamTimelineProcedure is the fully-qualified name of the ConsoleService's
	// StreamTimeline RPC.
	ConsoleServiceStreamTimelineProcedure = "/xyz.block.ftl.v1.console.ConsoleService/StreamTimeline"
	// ConsoleServiceGetTimelineProcedure is the fully-qualified name of the ConsoleService's
	// GetTimeline RPC.
	ConsoleServiceGetTimelineProcedure = "/xyz.block.ftl.v1.console.ConsoleService/GetTimeline"
)

// ConsoleServiceClient is a client for the xyz.block.ftl.v1.console.ConsoleService service.
type ConsoleServiceClient interface {
	// Ping service for readiness.
	Ping(context.Context, *connect_go.Request[v1.PingRequest]) (*connect_go.Response[v1.PingResponse], error)
	GetModules(context.Context, *connect_go.Request[console.GetModulesRequest]) (*connect_go.Response[console.GetModulesResponse], error)
	GetCalls(context.Context, *connect_go.Request[console.GetCallsRequest]) (*connect_go.Response[console.GetCallsResponse], error)
	GetRequestCalls(context.Context, *connect_go.Request[console.GetRequestCallsRequest]) (*connect_go.Response[console.GetRequestCallsResponse], error)
	StreamTimeline(context.Context, *connect_go.Request[console.StreamTimelineRequest]) (*connect_go.ServerStreamForClient[console.StreamTimelineResponse], error)
	GetTimeline(context.Context, *connect_go.Request[console.TimelineQuery]) (*connect_go.Response[console.GetTimelineResponse], error)
}

// NewConsoleServiceClient constructs a client for the xyz.block.ftl.v1.console.ConsoleService
// service. By default, it uses the Connect protocol with the binary Protobuf Codec, asks for
// gzipped responses, and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply
// the connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewConsoleServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) ConsoleServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &consoleServiceClient{
		ping: connect_go.NewClient[v1.PingRequest, v1.PingResponse](
			httpClient,
			baseURL+ConsoleServicePingProcedure,
			connect_go.WithIdempotency(connect_go.IdempotencyNoSideEffects),
			connect_go.WithClientOptions(opts...),
		),
		getModules: connect_go.NewClient[console.GetModulesRequest, console.GetModulesResponse](
			httpClient,
			baseURL+ConsoleServiceGetModulesProcedure,
			opts...,
		),
		getCalls: connect_go.NewClient[console.GetCallsRequest, console.GetCallsResponse](
			httpClient,
			baseURL+ConsoleServiceGetCallsProcedure,
			opts...,
		),
		getRequestCalls: connect_go.NewClient[console.GetRequestCallsRequest, console.GetRequestCallsResponse](
			httpClient,
			baseURL+ConsoleServiceGetRequestCallsProcedure,
			opts...,
		),
		streamTimeline: connect_go.NewClient[console.StreamTimelineRequest, console.StreamTimelineResponse](
			httpClient,
			baseURL+ConsoleServiceStreamTimelineProcedure,
			opts...,
		),
		getTimeline: connect_go.NewClient[console.TimelineQuery, console.GetTimelineResponse](
			httpClient,
			baseURL+ConsoleServiceGetTimelineProcedure,
			opts...,
		),
	}
}

// consoleServiceClient implements ConsoleServiceClient.
type consoleServiceClient struct {
	ping            *connect_go.Client[v1.PingRequest, v1.PingResponse]
	getModules      *connect_go.Client[console.GetModulesRequest, console.GetModulesResponse]
	getCalls        *connect_go.Client[console.GetCallsRequest, console.GetCallsResponse]
	getRequestCalls *connect_go.Client[console.GetRequestCallsRequest, console.GetRequestCallsResponse]
	streamTimeline  *connect_go.Client[console.StreamTimelineRequest, console.StreamTimelineResponse]
	getTimeline     *connect_go.Client[console.TimelineQuery, console.GetTimelineResponse]
}

// Ping calls xyz.block.ftl.v1.console.ConsoleService.Ping.
func (c *consoleServiceClient) Ping(ctx context.Context, req *connect_go.Request[v1.PingRequest]) (*connect_go.Response[v1.PingResponse], error) {
	return c.ping.CallUnary(ctx, req)
}

// GetModules calls xyz.block.ftl.v1.console.ConsoleService.GetModules.
func (c *consoleServiceClient) GetModules(ctx context.Context, req *connect_go.Request[console.GetModulesRequest]) (*connect_go.Response[console.GetModulesResponse], error) {
	return c.getModules.CallUnary(ctx, req)
}

// GetCalls calls xyz.block.ftl.v1.console.ConsoleService.GetCalls.
func (c *consoleServiceClient) GetCalls(ctx context.Context, req *connect_go.Request[console.GetCallsRequest]) (*connect_go.Response[console.GetCallsResponse], error) {
	return c.getCalls.CallUnary(ctx, req)
}

// GetRequestCalls calls xyz.block.ftl.v1.console.ConsoleService.GetRequestCalls.
func (c *consoleServiceClient) GetRequestCalls(ctx context.Context, req *connect_go.Request[console.GetRequestCallsRequest]) (*connect_go.Response[console.GetRequestCallsResponse], error) {
	return c.getRequestCalls.CallUnary(ctx, req)
}

// StreamTimeline calls xyz.block.ftl.v1.console.ConsoleService.StreamTimeline.
func (c *consoleServiceClient) StreamTimeline(ctx context.Context, req *connect_go.Request[console.StreamTimelineRequest]) (*connect_go.ServerStreamForClient[console.StreamTimelineResponse], error) {
	return c.streamTimeline.CallServerStream(ctx, req)
}

// GetTimeline calls xyz.block.ftl.v1.console.ConsoleService.GetTimeline.
func (c *consoleServiceClient) GetTimeline(ctx context.Context, req *connect_go.Request[console.TimelineQuery]) (*connect_go.Response[console.GetTimelineResponse], error) {
	return c.getTimeline.CallUnary(ctx, req)
}

// ConsoleServiceHandler is an implementation of the xyz.block.ftl.v1.console.ConsoleService
// service.
type ConsoleServiceHandler interface {
	// Ping service for readiness.
	Ping(context.Context, *connect_go.Request[v1.PingRequest]) (*connect_go.Response[v1.PingResponse], error)
	GetModules(context.Context, *connect_go.Request[console.GetModulesRequest]) (*connect_go.Response[console.GetModulesResponse], error)
	GetCalls(context.Context, *connect_go.Request[console.GetCallsRequest]) (*connect_go.Response[console.GetCallsResponse], error)
	GetRequestCalls(context.Context, *connect_go.Request[console.GetRequestCallsRequest]) (*connect_go.Response[console.GetRequestCallsResponse], error)
	StreamTimeline(context.Context, *connect_go.Request[console.StreamTimelineRequest], *connect_go.ServerStream[console.StreamTimelineResponse]) error
	GetTimeline(context.Context, *connect_go.Request[console.TimelineQuery]) (*connect_go.Response[console.GetTimelineResponse], error)
}

// NewConsoleServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewConsoleServiceHandler(svc ConsoleServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle(ConsoleServicePingProcedure, connect_go.NewUnaryHandler(
		ConsoleServicePingProcedure,
		svc.Ping,
		connect_go.WithIdempotency(connect_go.IdempotencyNoSideEffects),
		connect_go.WithHandlerOptions(opts...),
	))
	mux.Handle(ConsoleServiceGetModulesProcedure, connect_go.NewUnaryHandler(
		ConsoleServiceGetModulesProcedure,
		svc.GetModules,
		opts...,
	))
	mux.Handle(ConsoleServiceGetCallsProcedure, connect_go.NewUnaryHandler(
		ConsoleServiceGetCallsProcedure,
		svc.GetCalls,
		opts...,
	))
	mux.Handle(ConsoleServiceGetRequestCallsProcedure, connect_go.NewUnaryHandler(
		ConsoleServiceGetRequestCallsProcedure,
		svc.GetRequestCalls,
		opts...,
	))
	mux.Handle(ConsoleServiceStreamTimelineProcedure, connect_go.NewServerStreamHandler(
		ConsoleServiceStreamTimelineProcedure,
		svc.StreamTimeline,
		opts...,
	))
	mux.Handle(ConsoleServiceGetTimelineProcedure, connect_go.NewUnaryHandler(
		ConsoleServiceGetTimelineProcedure,
		svc.GetTimeline,
		opts...,
	))
	return "/xyz.block.ftl.v1.console.ConsoleService/", mux
}

// UnimplementedConsoleServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedConsoleServiceHandler struct{}

func (UnimplementedConsoleServiceHandler) Ping(context.Context, *connect_go.Request[v1.PingRequest]) (*connect_go.Response[v1.PingResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("xyz.block.ftl.v1.console.ConsoleService.Ping is not implemented"))
}

func (UnimplementedConsoleServiceHandler) GetModules(context.Context, *connect_go.Request[console.GetModulesRequest]) (*connect_go.Response[console.GetModulesResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("xyz.block.ftl.v1.console.ConsoleService.GetModules is not implemented"))
}

func (UnimplementedConsoleServiceHandler) GetCalls(context.Context, *connect_go.Request[console.GetCallsRequest]) (*connect_go.Response[console.GetCallsResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("xyz.block.ftl.v1.console.ConsoleService.GetCalls is not implemented"))
}

func (UnimplementedConsoleServiceHandler) GetRequestCalls(context.Context, *connect_go.Request[console.GetRequestCallsRequest]) (*connect_go.Response[console.GetRequestCallsResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("xyz.block.ftl.v1.console.ConsoleService.GetRequestCalls is not implemented"))
}

func (UnimplementedConsoleServiceHandler) StreamTimeline(context.Context, *connect_go.Request[console.StreamTimelineRequest], *connect_go.ServerStream[console.StreamTimelineResponse]) error {
	return connect_go.NewError(connect_go.CodeUnimplemented, errors.New("xyz.block.ftl.v1.console.ConsoleService.StreamTimeline is not implemented"))
}

func (UnimplementedConsoleServiceHandler) GetTimeline(context.Context, *connect_go.Request[console.TimelineQuery]) (*connect_go.Response[console.GetTimelineResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("xyz.block.ftl.v1.console.ConsoleService.GetTimeline is not implemented"))
}
