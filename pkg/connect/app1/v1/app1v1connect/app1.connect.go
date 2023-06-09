// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: app1/v1/app1.proto

package app1v1connect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	v1 "github.com/morning-night-dream/play-go-tracing/pkg/connect/app1/v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect_go.IsAtLeastVersion0_1_0

const (
	// HelloServiceName is the fully-qualified name of the HelloService service.
	HelloServiceName = "app1.v1.HelloService"
)

// HelloServiceClient is a client for the app1.v1.HelloService service.
type HelloServiceClient interface {
	Hello(context.Context, *connect_go.Request[v1.HelloRequest]) (*connect_go.Response[v1.HelloResponse], error)
}

// NewHelloServiceClient constructs a client for the app1.v1.HelloService service. By default, it
// uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewHelloServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) HelloServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &helloServiceClient{
		hello: connect_go.NewClient[v1.HelloRequest, v1.HelloResponse](
			httpClient,
			baseURL+"/app1.v1.HelloService/Hello",
			opts...,
		),
	}
}

// helloServiceClient implements HelloServiceClient.
type helloServiceClient struct {
	hello *connect_go.Client[v1.HelloRequest, v1.HelloResponse]
}

// Hello calls app1.v1.HelloService.Hello.
func (c *helloServiceClient) Hello(ctx context.Context, req *connect_go.Request[v1.HelloRequest]) (*connect_go.Response[v1.HelloResponse], error) {
	return c.hello.CallUnary(ctx, req)
}

// HelloServiceHandler is an implementation of the app1.v1.HelloService service.
type HelloServiceHandler interface {
	Hello(context.Context, *connect_go.Request[v1.HelloRequest]) (*connect_go.Response[v1.HelloResponse], error)
}

// NewHelloServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewHelloServiceHandler(svc HelloServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle("/app1.v1.HelloService/Hello", connect_go.NewUnaryHandler(
		"/app1.v1.HelloService/Hello",
		svc.Hello,
		opts...,
	))
	return "/app1.v1.HelloService/", mux
}

// UnimplementedHelloServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedHelloServiceHandler struct{}

func (UnimplementedHelloServiceHandler) Hello(context.Context, *connect_go.Request[v1.HelloRequest]) (*connect_go.Response[v1.HelloResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("app1.v1.HelloService.Hello is not implemented"))
}
