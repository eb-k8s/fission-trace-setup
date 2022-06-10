from .otel import *

from opentelemetry import trace
import opentelemetry.sdk.resources as resources
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter

resource = resources.Resource.create({"service.name": "Function-A-py"})
bsp = BatchSpanProcessor(OTLPSpanExporter(endpoint="fission-collector-collector.observability.svc:4317", insecure=True))
trace.set_tracer_provider(TracerProvider(active_span_processor=bsp, resource=resource))

tracer = trace.get_tracer("")
