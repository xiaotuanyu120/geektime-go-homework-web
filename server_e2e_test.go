//go:build e2e

package web

import (
	"fmt"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	s := NewHTTPServer()
	// method 1: use http package
	//http.ListenAndServe(":8080", s)

	s.addRoute(http.MethodGet, "/", func(ctx *Context) {
		_, err := ctx.Resp.Write([]byte("homepage"))
		if err != nil {
			return
		}
	})

	s.Get("/user", func(ctx *Context) {
		fmt.Println("something")
	})

	s.Get("/order/detail", func(ctx *Context) {
		_, err := ctx.Resp.Write([]byte("hello, /order/detail"))
		if err != nil {
			return
		}
	})

	s.Get("/order/detail/:id", func(ctx *Context) {
		_, err := ctx.Resp.Write([]byte(fmt.Sprintf("hello, /order/detail/%s", ctx.PathParams["id"])))
		if err != nil {
			return
		}
	})

	s.Get("/order/detail/3", func(ctx *Context) {
		_, err := ctx.Resp.Write([]byte("hello, /order/detail/3 (static route)"))
		if err != nil {
			return
		}
	})

	s.Get("/order/*", func(ctx *Context) {
		_, err := ctx.Resp.Write([]byte(fmt.Sprintf("hello, %s", ctx.Req.URL.Path)))
		if err != nil {
			return
		}
	})

	err := s.Start(":8080")
	if err != nil {
		return
	}
}
