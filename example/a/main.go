package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
)

// Initializes an OTLP exporter, and configures the corresponding trace and
// metric providers.
func initProvider() {
	ctx := context.Background()

	// Set up resource
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("Function-A"),
	)

	// Set up a trace exporter
	grpcOpts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint("fission-collector-collector.observability.svc:4317"),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
		otlptracegrpc.WithInsecure(),
	}
	traceExporter, err := otlptracegrpc.New(ctx, grpcOpts...)
	handleErr(err, "failed to create trace exporter")

	// Register the trace exporter with a TracerProvider, using
	// a batch span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(bsp),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tracerProvider)

	// Set tracecontext as global propagator.
	propagators := []propagation.TextMapPropagator{
		propagation.TraceContext{},
		propagation.Baggage{},
	}
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagators...))
}

func handleErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}

func init() {

	initProvider()
}

// HandlerA is the entry point for this fission function
func HandlerA(w http.ResponseWriter, r *http.Request) {

	bag, _ := baggage.Parse("funcname=HandlerA")
	ctx := baggage.ContextWithBaggage(context.Background(), bag)
	tracer := otel.Tracer("")

	// start a local span
	func(ctx context.Context) {
		_, span := tracer.Start(ctx, "doing some process")
		defer span.End()
		time.Sleep(time.Duration(1) * time.Second)
	}(ctx)

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	var respBody []byte


	err := func(ctx context.Context) error {
		// now ctx have all info about tracer and this span
		ctx, span := tracer.Start(ctx, "say hello to funcb", trace.WithSpanKind(trace.SpanKindClient))
		defer span.End()

		req, _ := http.NewRequestWithContext(
			ctx,
			"GET",
			"http://router.fission.svc.cluster.local/funcb",
			nil,
		)
		log.Println("Sending request...")
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalln(err)
		}
		w.Write([]byte("successfully invoked fn B"))
		respBody, err = ioutil.ReadAll(resp.Body)
		_ = resp.Body.Close()

		return err
	}(ctx)

	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Response Received: %s \n\n", respBody)
}
