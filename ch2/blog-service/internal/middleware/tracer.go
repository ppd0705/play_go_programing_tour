package middleware

import (
	"block-service/global"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
)

func Tracing() func(c *gin.Context) {
	return func(c *gin.Context) {
		var ctx context.Context

		span := opentracing.SpanFromContext(c)
		if span != nil {
			span, ctx = opentracing.StartSpanFromContextWithTracer(c.Request.Context(), global.Tracer, c.Request.URL.Path, opentracing.ChildOf(span.Context()))
		} else {
			span, ctx = opentracing.StartSpanFromContextWithTracer(c.Request.Context(), global.Tracer, c.Request.URL.Path)
		}
		defer span.Finish()

		var tranceID string
		var spanID string
		spanCtx := span.Context()
		switch spanCtx.(type){
		case jaeger.SpanContext:
			tranceID = spanCtx.(jaeger.SpanContext).TraceID().String()
			spanID = spanCtx.(jaeger.SpanContext).SpanID().String()
		}
		c.Set("X-Trace-ID", tranceID)
		c.Set("X-Span-ID", spanID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}