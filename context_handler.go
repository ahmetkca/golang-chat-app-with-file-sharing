package main

import (
	"net/http"

	"golang.org/x/net/context"
)

type requestKey int

const requestIDKey requestKey = 0

// Helper test function to test new Handler with context
func newContextWithRequestID(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, requestIDKey, req.Header.Get("X-Request-ID"))
}

func requestIDFromContext(ctx context.Context) string {
	return ctx.Value(requestIDKey).(string)
}

// Modified Handler that takes context along with ResponseWriter interaface and Request struct
type ContextHandler interface {
	ServeHTTPContext(context.Context, http.ResponseWriter, *http.Request)
}

// Modified HandlerFunc which additionally takes context as well
type ContextHandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

// ContextHandlerFunc can call ServeHTTPContext and can be used as a ServeHTTPContext func as well as ContextHandler
func (h ContextHandlerFunc) ServeHTTPContext(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	h(ctx, rw, req)
}

func middleware(h ContextHandler) ContextHandler {
	return ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
		ctx = newContextWithRequestID(ctx, req)
		h.ServeHTTPContext(ctx, rw, req)
	})
}

// Adapter which will allow developer to use vanilla http.Handler with the Context
type ContextAdapter struct {
	ctx     context.Context
	handler ContextHandler
}

// ContextAdapter can be used as a Handler (with context)
func (ca *ContextAdapter) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ca.handler.ServeHTTPContext(ca.ctx, rw, req)
}
