package web

import (
	"fmt"
	"strings"
)

type router struct {
	// http method => tree
	trees map[string]*node
}

type node struct {
	path     string
	children map[string]*node // children path => children node

	// wild card child: /order/detail/*
	wildCardChild *node

	// parameter child: /order/detail/:id
	paramChild *node

	handler handleFunc
}

func newRouter() *router {
	return &router{
		trees: map[string]*node{},
	}
}

func (r *router) addRoute(method string, path string, handleFunc handleFunc) {
	if path == "" {
		panic("Route Check Error: [path] can not be empty!")
	}

	if path[:1] != "/" {
		panic("Route Check Error: [path] must be start with '/'!")
	}

	if len(path) > 1 && path[len(path)-1:] == "/" {
		panic("Route Check Error: [path] last character can not be '/'!")
	}

	currentNode, ok := r.trees[method]

	// Create tree if not exist
	if !ok {
		currentNode = &node{
			path: "/",
		}
		r.trees[method] = currentNode
	}

	if path == "/" {
		if currentNode.handler != nil {
			panic("Route Add More Than One Time: [/] Already added")
		}
		currentNode.handler = handleFunc
		return
	}

	// remove first "/"
	path2Split := path[1:]

	// Create child node if not exist
	segs := strings.Split(path2Split, "/")
	for _, seg := range segs {
		if seg == "" {
			panic("Route Check Error: [path] can not use continue '/', such as '//'!")
		}
		child := currentNode.childOrCreate(seg)
		currentNode = child
	}
	// currentNode now is the last seg's node
	if currentNode.handler != nil {
		panic(fmt.Sprintf("Route Add More Than One Time: [%s] Already added", path))
	}
	currentNode.handler = handleFunc
}

func (n *node) childOrCreate(seg string) *node {
	// go to wildCardChild
	if seg == "*" {
		if n.paramChild != nil {
			panic(fmt.Sprintf(
				`Parameter Child and Wild Card Child only can not exist at the same time: 
paramChild %s already exist!`, seg))
		}
		if n.wildCardChild == nil {
			n.wildCardChild = &node{path: seg}
		}
		return n.wildCardChild
	}

	// go to paramChild
	if seg[0] == ':' {
		if n.wildCardChild != nil {
			panic(fmt.Sprintf(
				`Parameter Child and Wild Card Child only can not exist at the same time: 
wildCardChild %s already exist!`, seg))
		}

		if n.paramChild == nil {
			n.paramChild = &node{path: seg}
		}
		return n.paramChild
	}

	// go to child
	// init node.child
	if n.children == nil {
		n.children = map[string]*node{}
	}

	child, ok := n.children[seg]
	if !ok {
		child = &node{
			path: seg,
		}
		n.children[seg] = child
	}
	return child
}

type matchInfo struct {
	n          *node
	pathParams map[string]string
}

func (r *router) findRoute(method string, path string) (*matchInfo, bool) {
	currentNode, ok := r.trees[method]
	if !ok {
		return nil, false
	}

	// root path
	if path == "/" {
		return &matchInfo{n: currentNode}, true
	}

	// remove first "/"
	path2Split := path[1:]
	segs := strings.Split(path2Split, "/")
	// used by paramChild node
	var pathParams map[string]string
	// used by tail wild card node, when regular match failed, try this
	var tailWildCardNode *node
	for _, seg := range segs {
		// regular match
		child, isParamChild, found := currentNode.childOf(seg)

		// try to cache tail wild card node
		wcNode, twcFound := currentNode.wildCardChildOf()
		if twcFound {
			tailWildCardNode = wcNode
		}

		// collect the path params in the request path
		if isParamChild {
			if pathParams == nil {
				pathParams = make(map[string]string)
			}
			pathParams[child.path[1:]] = seg
		}

		if !found {
			// regular match failed, then return cached tail wild card node
			if tailWildCardNode != nil {
				return &matchInfo{n: tailWildCardNode, pathParams: pathParams}, true
			}
			return nil, false
		}

		currentNode = child
	}

	return &matchInfo{n: currentNode, pathParams: pathParams}, true
}

func (n *node) childOf(path string) (node *node, isParamChild bool, found bool) {
	// 1. check children
	// 2. check paramChild
	// at last: check wildCardChild
	if n.children == nil {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.wildCardChild, false, n.wildCardChild != nil
	}
	child, ok := n.children[path]
	if !ok {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.wildCardChild, false, n.wildCardChild != nil
	}
	return child, false, ok
}

func (n *node) wildCardChildOf() (node *node, found bool) {
	if n.wildCardChild != nil {
		return n.wildCardChild, true
	}
	return n.wildCardChild, false
}
