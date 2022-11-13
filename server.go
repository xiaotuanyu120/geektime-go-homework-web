package web

import (
	"net"
	"net/http"
)

type handleFunc func(ctx *Context)

// ensure HTTPServer implement Server
var _ Server = &HTTPServer{}

type Server interface {
	http.Handler

	// Start a Server
	Start(addr string) error

	// AddRoute
	// register route logic here
	// - method, http request method
	// - path, http request path
	// - handleFunc, business logic func
	addRoute(method string, path string, handleFunc handleFunc)
}

type HTTPServer struct {
	*router
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{
		router: newRouter(),
	}
}

func (h *HTTPServer) Get(path string, handleFunc handleFunc) {
	h.addRoute(http.MethodGet, path, handleFunc)
}

func (h *HTTPServer) Post(path string, handleFunc handleFunc) {
	h.addRoute(http.MethodPost, path, handleFunc)
}

func (h *HTTPServer) Put(path string, handleFunc handleFunc) {
	h.addRoute(http.MethodPut, path, handleFunc)
}

func (h *HTTPServer) Delete(path string, handleFunc handleFunc) {
	h.addRoute(http.MethodDelete, path, handleFunc)
}

func (h *HTTPServer) Head(path string, handleFunc handleFunc) {
	h.addRoute(http.MethodHead, path, handleFunc)
}

func (h *HTTPServer) Options(path string, handleFunc handleFunc) {
	h.addRoute(http.MethodOptions, path, handleFunc)
}

func (h *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}

	h.Serve(ctx)
}

func (h *HTTPServer) Start(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	// start hook here

	return http.Serve(l, h)
}

func (h *HTTPServer) Serve(ctx *Context) {
	n, found := h.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !found || n.handler == nil {
		ctx.Resp.WriteHeader(404)
		ctx.Resp.Write([]byte("NOT FOUND"))
		return
	}
	n.handler(ctx)
}
