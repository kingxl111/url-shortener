package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"runtime/debug"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServerOptions struct {
	logger        *slog.Logger
	panicHandler  func(p any) (err error)
	serverOptions []grpc.ServerOption
}

func (o *ServerOptions) WithLogger(logger *slog.Logger) {
	o.logger = logger
}

func (o *ServerOptions) WithPanicHandler(h func(p any) (err error)) {
	o.panicHandler = h
}

func (o *ServerOptions) WithServerOptions(v ...grpc.ServerOption) {
	o.serverOptions = append(o.serverOptions, v...)
}

func (o *ServerOptions) WithUnaryInterceptors(v ...grpc.UnaryServerInterceptor) {
	o.WithServerOptions(grpc.ChainUnaryInterceptor(v...))
}

func (o *ServerOptions) WithStreamInterceptors(v ...grpc.StreamServerInterceptor) {
	o.WithServerOptions(grpc.ChainStreamInterceptor(v...))
}

func (o *ServerOptions) NewServer() (*grpc.Server, error) {
	if o.logger == nil {
		o.logger = slog.Default()
	}
	if o.panicHandler == nil {
		o.panicHandler = func(p any) (err error) {
			o.logger.Error("recovered from panic", "panic", p, "stack", debug.Stack())
			return status.Errorf(codes.Internal, "%s", p)
		}
	}

	return grpc.NewServer(
		append(
			[]grpc.ServerOption{
				grpc.ChainUnaryInterceptor(
					logging.UnaryServerInterceptor(Interceptor(o.logger)),
					recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(o.panicHandler)),
				),
			}, o.serverOptions...,
		)...,
	), nil
}

func ListenAndServeContext(ctx context.Context, addr string, srv *grpc.Server) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s", addr))
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		srv.GracefulStop()
	}()
	return srv.Serve(lis)
}

func Interceptor(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}
