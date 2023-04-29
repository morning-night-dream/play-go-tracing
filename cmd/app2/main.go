package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	app1v1 "github.com/morning-night-dream/play-go-tracing/pkg/connect/app1/v1"
	"github.com/morning-night-dream/play-go-tracing/pkg/connect/app1/v1/app1v1connect"
	"github.com/morning-night-dream/play-go-tracing/trace"

	"github.com/bufbuild/connect-go"
	otelconnect "github.com/bufbuild/connect-opentelemetry-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var tracer = otel.Tracer("app2/app2-service")

func main() {
	exporter, err := trace.NewExporter()
	if err != nil {
		log.Fatal(err)
	}

	reource := trace.NewResource("app2-service", "1.0.0")
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(reource),
	)
	otel.SetTracerProvider(tracerProvider)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Printf("Failed to shutdown tracer provider: %v", err)
		}
	}()

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	url := os.Getenv("API_SERVER_URL")
	if url == "" {
		log.Fatal("API_SERVER_URL is not set")
	}

	client := app1v1connect.NewHelloServiceClient(http.DefaultClient, url, connect.WithInterceptors(otelconnect.NewInterceptor(
		otelconnect.WithTracerProvider(tracerProvider),
		otelconnect.WithPropagator(propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		)),
		otelconnect.WithTrustRemote(),
	)))

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "app2/hello")
		defer span.End()
		log.Print("handle hello")

		res, err := client.Hello(ctx, connect.NewRequest(&app1v1.HelloRequest{}))
		if err != nil {
			log.Printf("failed to call hello: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write([]byte(fmt.Sprintf("%s, and Goodbye", res.Msg.Message)))
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
