package web

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func TestRouter_addRoute(t *testing.T) {
	// Construct testRouter AND Add Route
	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
		{
			method: http.MethodPost,
			path:   "/login",
		},
		{
			method: http.MethodGet,
			path:   "/*",
		},
		{
			method: http.MethodGet,
			path:   "/abc/*",
		},
		{
			method: http.MethodGet,
			path:   "/*/*",
		},
		{
			method: http.MethodGet,
			path:   "/*/abc/*",
		},
	}

	testRouter := newRouter()
	var fakeHandleFunc handleFunc = func(ctx *Context) {}
	for _, tr := range testRoutes {
		testRouter.addRoute(tr.method, tr.path, fakeHandleFunc)
	}

	// Construct mockRouter
	mockRouter := &router{
		trees: map[string]*node{
			http.MethodGet: {
				path:    "/",
				handler: fakeHandleFunc,
				children: map[string]*node{
					"user": {
						path:    "user",
						handler: fakeHandleFunc,
						children: map[string]*node{
							"home": {
								path:     "home",
								children: map[string]*node{},
								handler:  fakeHandleFunc,
							},
						},
					},
					"order": {
						path: "order",
						children: map[string]*node{
							"detail": {
								path:     "detail",
								children: map[string]*node{},
								handler:  fakeHandleFunc,
							},
						},
					},
					"abc": {
						path: "abc",
						wildChild: &node{
							path:    "*",
							handler: fakeHandleFunc,
						},
					},
				},
				wildChild: &node{
					path:    "*",
					handler: fakeHandleFunc,
					children: map[string]*node{
						"abc": {
							path: "abc",
							wildChild: &node{
								path:    "*",
								handler: fakeHandleFunc,
							},
						},
					},
					wildChild: &node{
						path:    "*",
						handler: fakeHandleFunc,
					},
				},
			},
			http.MethodPost: {
				path: "/",
				children: map[string]*node{
					"order": {
						path: "order",
						children: map[string]*node{
							"create": {
								path:     "create",
								children: map[string]*node{},
								handler:  fakeHandleFunc,
							},
						},
					},
					"login": {
						path:    "login",
						handler: fakeHandleFunc,
					},
				},
			},
		},
	}

	// TEST: Compare mockRouter and testRouter
	errMsg, ok := mockRouter.equal(testRouter)
	assert.True(t, ok, errMsg)

	// TEST: CHECK [path]
	r := newRouter()
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "", fakeHandleFunc)
	},
		"Route Check Error: [path] can not be empty!",
	)

	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "user/home", fakeHandleFunc)
	}, "Route Check Error: [path] must be start with '/'!")

	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/user/", fakeHandleFunc)
	}, "Route Check Error: [path] last character can not be '/'!")

	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/user///home", fakeHandleFunc)
	}, "Route Check Error: [path] can not use continue '/', such as '//'!")

	// TEST: REPEATED router
	r = newRouter()
	r.addRoute(http.MethodGet, "/", fakeHandleFunc)
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/", fakeHandleFunc)
	}, "Route Add More Than One Time: [/] Already added")

	r = newRouter()
	r.addRoute(http.MethodGet, "/a/b/c", fakeHandleFunc)
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/a/b/c", fakeHandleFunc)
	}, "Route Add More Than One Time: [/a/b/c] Already added")
}

func (r *router) equal(targetRouter *router) (errMsg string, equal bool) {
	for HTTPMethod, rootNode := range r.trees {
		// compare key(http method) exist or not
		targetRootNode, ok := targetRouter.trees[HTTPMethod]
		if !ok {
			return fmt.Sprintf("MATCH FAILED: [targetRouter's HTTPMETHOD: %s] not exist", HTTPMethod), false
		}

		// compare node
		if errMsg, equal := rootNode.equal(targetRootNode); !equal {
			return errMsg, false
		}
	}

	return "", true
}

func (n *node) equal(targetNode *node) (errMsg string, equal bool) {
	// compare path
	if n.path != targetNode.path {
		return fmt.Sprintf("MATCH FAILED: [path %s != %s]", n.path, targetNode.path), false
	}

	// compare wildChild node
	if n.wildChild != nil {
		if targetNode.wildChild == nil {
			return "MATCH FAILED: target wildChild is nil, but real wildChild is not nil", false
		}
		msg, ok := n.wildChild.equal(targetNode.wildChild)
		return msg, ok
	}

	// compare child node
	cl := len(n.children)
	tcl := len(targetNode.children)
	if cl != tcl {
		return fmt.Sprintf("MATCH FAILED: [child node len %d != %d]", cl, tcl), false
	}

	for cp, cn := range n.children {
		// compare key(child path) exist or not
		targetChildNode, ok := targetNode.children[cp]
		if !ok {
			return fmt.Sprintf("MATCH FAILED: [targetChildNode's path %s not exist]", cp), false
		}

		// compare child node recursion
		errMsg, equal := cn.equal(targetChildNode)
		if !equal {
			return errMsg, false
		}
	}

	// compare handleFunc
	nodeHandleFunc := reflect.ValueOf(n.handler)
	targetNodeHandleFunc := reflect.ValueOf(targetNode.handler)
	if nodeHandleFunc != targetNodeHandleFunc {
		return fmt.Sprintf("MATCH FAILED: [handleFunc not equal]"), false
	}

	return "", true
}

func TestRouter_findRoute(t *testing.T) {
	// Construct testRouter AND Add route
	testRoutes := []struct {
		name   string
		expect bool
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail",
		},
		{
			method: http.MethodGet,
			path:   "/order/*",
		},
		{
			method: http.MethodDelete,
			path:   "/",
		},
	}

	testRouter := newRouter()
	var fakeHandleFunc handleFunc = func(ctx *Context) {}
	for _, tr := range testRoutes {
		testRouter.addRoute(tr.method, tr.path, fakeHandleFunc)
	}

	// Construct testCases
	testCases := []struct {
		name   string
		expect bool
		method string
		path   string
		ret    *node
	}{
		{
			name:   "root",
			expect: true,
			method: http.MethodDelete,
			path:   "/",
			ret: &node{
				path:    "/",
				handler: fakeHandleFunc,
			},
		},
		{
			name:   "depth two",
			expect: true,
			method: http.MethodGet,
			path:   "/order/detail",
			ret: &node{
				path:    "detail",
				handler: fakeHandleFunc,
			},
		},
		{
			name:   "/order/abc > /order/*",
			expect: true,
			method: http.MethodGet,
			path:   "/order/abc",
			ret: &node{
				path:    "*",
				handler: fakeHandleFunc,
			},
		},
		{
			name:   "no handler",
			expect: true,
			method: http.MethodGet,
			path:   "/order",
			ret: &node{
				path: "order",
				children: map[string]*node{
					"detail": {
						path:    "detail",
						handler: fakeHandleFunc,
					},
				},
			},
		},
		{
			name:   "method not found",
			expect: false,
			method: http.MethodOptions,
			path:   "/order/detail",
		},
		{
			name:   "path not found",
			expect: false,
			method: http.MethodGet,
			path:   "/not/exist",
		},
	}

	// Test: find route
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			node, found := testRouter.findRoute(tc.method, tc.path)
			assert.Equal(t, tc.expect, found)
			if !found {
				return
			}
			msg, equal := tc.ret.equal(node)
			assert.True(t, equal, msg)
		})
	}
}
