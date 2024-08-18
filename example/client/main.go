package main

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/Hidayathamir/oteltrace"
	"github.com/Hidayathamir/oteltrace/example/pbfoo"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	os.Setenv("X_OTELTRACE_OTEL_OTLP_NEWRELIC_HOST", "otlp.nr-data.net:4317")
	os.Setenv("X_OTELTRACE_OTEL_OTLP_NEWRELIC_HEADER_API_KEY", "dummykey")
	os.Setenv("X_OTELTRACE_APP_SERVICE_NAME", "example_client")
	os.Setenv("X_OTELTRACE_APP_VERSION", "1.0.0")
	os.Setenv("X_OTELTRACE_APP_ENVIRONMENT", "dev")

	tp, err := oteltrace.NewTraceProvider(context.Background())
	fatalIfErr(err)
	defer func() {
		err := tp.Shutdown(context.Background())
		warnIfErr(err)
	}()

	extapiHTTPFoo()
	extapiGRPCFoo()
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

func extapiHTTPFoo() {
	ctx, span := oteltrace.RecordSpan(context.Background())
	defer span.End()

	client := &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:4000/foo", nil)
	fatalIfErr(err)
	res, err := client.Do(req)
	fatalIfErr(err)
	resBody, err := io.ReadAll(res.Body)
	fatalIfErr(err)
	body := map[string]any{}
	err = json.Unmarshal(resBody, &body)
	fatalIfErr(err)
	slog.Info("success", "body", body)
}

func extapiGRPCFoo() {
	ctx, span := oteltrace.RecordSpan(context.Background())
	defer span.End()

	conn, err := grpc.NewClient("localhost:4001", grpc.WithStatsHandler(otelgrpc.NewClientHandler()), grpc.WithTransportCredentials(insecure.NewCredentials()))
	fatalIfErr(err)
	client := pbfoo.NewExampleClient(conn)
	res, err := client.Foo(ctx, &pbfoo.ReqFoo{})
	fatalIfErr(err)
	slog.Info("success", "res", res)
}
