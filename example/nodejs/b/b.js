'use strict';

const api = require('@opentelemetry/api')
const { defaultTextMapGetter, ROOT_CONTEXT, trace } = require('@opentelemetry/api');
const tracer = require('./tracer')('Function-B-js');

module.exports = async function(context) {
    console.log(context.request.headers);
    if (context.request.method == "GET") {
        const parentCtx = api.propagation.extract(ROOT_CONTEXT, context.request.headers, defaultTextMapGetter)
        // to get currentSpan
        // const currentSpan = trace.getSpan(parentCtx);
        // console.log(currentSpan.spanContext().traceId);
        const span = tracer.startSpan(
            'heandleRequest',
            {
                kind: 1, // server
            },
            parentCtx,
        )
        // do some process
        span.end();
        return {
            status: 200,
        }
    }
}