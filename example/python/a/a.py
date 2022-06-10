import http
from flask import request, Response
import requests

from opentelemetry.propagate import inject

from otel import tracer

def main():
    if request.method == "GET":
        with tracer.start_as_current_span("say hello to funcb-py"):

            # to start a child span:
            # with tracer.start_as_current_span("child span"):
            #     pass

            try:
                headers = {}
                inject(headers)
                resp = requests.get(
                    "http://router.fission.svc.cluster.local/funcb-py",
                    headers=headers,
                )
                
                print(resp.status_code)
                assert resp.status_code == 200

                return Response(status=http.HTTPStatus.OK)
            except Exception as e:
                print(e)
                return Response(status=http.HTTPStatus.INTERNAL_SERVER_ERROR)
    else:
        print("methods other than GET are not supported")
        return Response(status=http.HTTPStatus.BAD_REQUEST)
