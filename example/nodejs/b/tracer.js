'use strict';

const api = require('@opentelemetry/api');
const grpc = require('@grpc/grpc-js');
const { BasicTracerProvider, BatchSpanProcessor } = require('@opentelemetry/sdk-trace-base');
const { OTLPTraceExporter } = require('@opentelemetry/exporter-trace-otlp-grpc');
const { Resource } = require('@opentelemetry/resources');
const { SemanticResourceAttributes } = require('@opentelemetry/semantic-conventions');
// const { registerInstrumentations } = require('@opentelemetry/instrumentation');
// const { HttpInstrumentation } = require('@opentelemetry/instrumentation-http');
const { W3CTraceContextPropagator } = require('@opentelemetry/core')

module.exports = (serviceName) => {
    const collectorOptions = {
        url: 'http://fission-collector-collector.observability.svc:4317',
        credentials: grpc.credentials.createInsecure(),
    };
    const provider = new BasicTracerProvider({
        resource: new Resource({
            [SemanticResourceAttributes.SERVICE_NAME]: serviceName,
        })
    });
    const exporter = new OTLPTraceExporter(collectorOptions);
    provider.addSpanProcessor(new BatchSpanProcessor(exporter));
    provider.register();

    api.propagation.setGlobalPropagator(new W3CTraceContextPropagator())
    // registerInstrumentations({
    //     instrumentations: [
    //         new HttpInstrumentation(),
    //     ]
    // })

    return provider.getTracer('');
}

