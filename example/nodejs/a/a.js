'use strict';

const api = require('@opentelemetry/api')
const { defaultTextMapSetter, ROOT_CONTEXT, trace } = require('@opentelemetry/api');
const tracer = require('./tracer')('Function-A-js');
const http = require('http');

module.exports = async function(context) {
    console.log(api.propagation.fields())
    if (context.request.method == "GET") {
        const span = tracer.startSpan('TriggerB');
        console.log('start a span')

        // http headers
        let carrier = {}
        api.propagation.inject(trace.setSpanContext(ROOT_CONTEXT, span.spanContext()), carrier, defaultTextMapSetter)

        http.get({
            host: 'router.fission.svc.cluster.local',
            port: '80',
            path: '/funcb-js',
            headers: carrier,
        }, (response) => {
            const body = [];
            response.on('data', (chunk) => body.push(chunk));
            response.on('end', () => {
                console.log(body.toString());
                span.end();
            });
        });
        return {
            status: 200,
        }
    }
}