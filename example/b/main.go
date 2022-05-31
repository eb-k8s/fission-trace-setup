package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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
	// Get Resource
	res := resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String("Function-B"))

	// Get Exporter
	grpcOpts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint("fission-collector-collector.observability.svc:4317"),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
		otlptracegrpc.WithInsecure(),
	}

	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, grpcOpts...)
	handleErr(err, "failed to create trace exporter")

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// set global propagator to tracecontext (the default is no-op).
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

// var HandlerB func(http.ResponseWriter, *http.Request)

func init() {
	initProvider()
}

// HandlerB is the entry point for this fission function
func HandlerB(w http.ResponseWriter, r *http.Request) {
	var tracer trace.Tracer
	if tracer == nil {
		if span := trace.SpanFromContext(r.Context()); span.SpanContext().IsValid() {
			tracer = span.TracerProvider().Tracer("")
		} else {
			tracer = otel.GetTracerProvider().Tracer("")
		}
	}
	// Extract context from carrier
	propagators := otel.GetTextMapPropagator()
	ctx := propagators.Extract(r.Context(), propagation.HeaderCarrier(r.Header))

	ctx, span := tracer.Start(
		ctx,
		"handle req from funca",
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(semconv.PeerServiceKey.String("TestService-A")),
	)
	defer span.End()

	bag := baggage.FromContext(ctx)
	funcname := attribute.Key("funcname")
	span.AddEvent("handling request", trace.WithAttributes(funcname.String(bag.Member("funcname").Value())))
	time.Sleep(time.Duration(2) * time.Second)

	log.Println("fn B is invoked")
	_, _ = io.WriteString(w, "I am fn B\n")

}

