# oteltrace

# go-opentelemetry oteltrace

OpenTelemetry Trace.

## Usage

<details close>
  <summary>Usage main.go</summary>

```go
package main

import (
	"context"

	"github.com/Hidayathamir/oteltrace"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func main() {
	tp, err := oteltrace.NewTraceProvider(context.Background())
	fatalIfErr(err)
	defer func() {
		err := tp.Shutdown(context.Background())
		warnIfErr(err)
	}()

	grpcServer = grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))

	ginEngine := gin.Default()
	ginEngine.Use(otelgin.Middleware(serviceName))
}
```

</details>

<details close>
  <summary>Usage controller-service-repository</summary>

```go
package main

import (
	"context"

	"github.com/Hidayathamir/oteltrace"
	"github.com/gin-gonic/gin"
)

func controllerHTTPFoo(c *gin.Context) {
	ctx, span := oteltrace.RecordSpan(c)
	defer span.End()
}

func controllerGRPCFoo(ctx context.Context, req *Req) (*Res, error) {
	ctx, span := oteltrace.RecordSpan(ctx)
	defer span.End()
}

func serviceFoo(ctx context.Context) {
	ctx, span := oteltrace.RecordSpan(ctx)
	defer span.End()
}

func repoGetFoo(ctx context.Context) {
	ctx, span := oteltrace.RecordSpan(ctx)
	defer span.End()
}
```

</details>

<details close>
  <summary>Usage grpc/http client</summary>

```go
package main

import (
	"context"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"google.golang.org/grpc"
)

func extapiHTTPFoo(ctx context.Context) {
	client = &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	res, err := client.Do(req)
}

func extapiGRPCFoo(ctx context.Context) {
	conn, err := grpc.NewClient(target, grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	client = pb.NewClient(conn)
	res, err := client.Foo(ctx, req)
}
```

</details>
