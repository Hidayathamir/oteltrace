package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Hidayathamir/oteltrace"
	"github.com/Hidayathamir/oteltrace/example/pbfoo"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func main() {
	os.Setenv("X_OTELTRACE_OTEL_OTLP_NEWRELIC_HOST", "otlp.nr-data.net:4317")
	os.Setenv("X_OTELTRACE_OTEL_OTLP_NEWRELIC_HEADER_API_KEY", "dummykey")
	os.Setenv("X_OTELTRACE_APP_SERVICE_NAME", "example_server")
	os.Setenv("X_OTELTRACE_APP_VERSION", "1.0.0")
	os.Setenv("X_OTELTRACE_APP_ENVIRONMENT", "dev")

	tp, err := oteltrace.NewTraceProvider(context.Background())
	fatalIfErr(err)
	defer func() {
		err := tp.Shutdown(context.Background())
		warnIfErr(err)
	}()

	go func() {
		grpcServer := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))
		pbfoo.RegisterExampleServer(grpcServer, &GRPCExampleServer{})
		lis, err := net.Listen("tcp", "localhost:4001")
		fatalIfErr(err)
		_ = grpcServer.Serve(lis)
	}()

	go func() {
		ginEngine := gin.Default()
		ginEngine.Use(otelgin.Middleware(""))
		ginEngine.GET("foo", controllerHTTPFoo)
		httpServer := &http.Server{Addr: "localhost:4000", Handler: ginEngine}
		_ = httpServer.ListenAndServe()
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	slog.Info("listens for the interrupt signal from the OS")
	<-ctx.Done()
	stop()
}

func warnIfErr(err error) {
	if err != nil {
		slog.Warn(err.Error())
	}
}

func fatalIfErr(err error) {
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func controllerHTTPFoo(c *gin.Context) {
	ctx, span := oteltrace.RecordSpan(c)
	defer span.End()

	time.Sleep(1 * time.Second)
	serviceFoo(ctx)

	c.JSON(http.StatusOK, gin.H{"trace_id": span.SpanContext().TraceID().String()})
}

type GRPCExampleServer struct {
	pbfoo.UnimplementedExampleServer
}

func (e *GRPCExampleServer) Foo(ctx context.Context, _ *pbfoo.ReqFoo) (*pbfoo.ResFoo, error) {
	ctx, span := oteltrace.RecordSpan(ctx)
	defer span.End()

	time.Sleep(1 * time.Second)
	serviceFoo(ctx)

	return &pbfoo.ResFoo{TraceId: span.SpanContext().TraceID().String()}, nil
}

func serviceFoo(ctx context.Context) {
	ctx, span := oteltrace.RecordSpan(ctx)
	defer span.End()

	time.Sleep(1 * time.Second)
	repoGetFoo(ctx)
}

func repoGetFoo(ctx context.Context) {
	ctx, span := oteltrace.RecordSpan(ctx)
	defer span.End()

	time.Sleep(1 * time.Second)
	slog.InfoContext(ctx, "dummyFoo")
}
