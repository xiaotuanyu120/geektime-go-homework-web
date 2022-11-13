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
	path      string
	children  map[string]*node // children path => children node
	wildChild *node
	handler   handleFunc
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
	// go to wildChild
	if seg == "*" {
		if n.wildChild == nil {
			n.wildChild = &node{path: seg}
		}
		return n.wildChild
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

func (r *router) findRoute(method string, path string) (*node, bool) {
	currentNode, ok := r.trees[method]
	if !ok {
		return nil, false
	}

	// root path
	if path == "/" {
		return currentNode, true
	}

	// remove first "/"
	path2Split := path[1:]
	segs := strings.Split(path2Split, "/")
	for _, seg := range segs {
		child, found := currentNode.childOf(seg)
		if !found {
			return nil, false
		}
		currentNode = child
	}

	return currentNode, true
}

func (n *node) childOf(path string) (*node, bool) {
	if n.children == nil {
		return n.wildChild, n.wildChild != nil
	}
	child, ok := n.children[path]
	if !ok {
		return n.wildChild, n.wildChild != nil
	}
	return child, ok
}
