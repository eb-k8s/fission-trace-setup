import http
from time import sleep
from flask import request, Response

from opentelemetry.propagate import extract
from opentelemetry.trace import SpanKind

from otel import tracer

def main():
    if request.method == "GET":
        with tracer.start_as_current_span("handle req from funca", context=extract(request.headers), kind=SpanKind.SERVER) as span:
            span.add_event("handling request")
            sleep(1)
            return Response(status=http.HTTPStatus.OK)