package router

import (
	"context"
	"net/http"
)

type Router interface {
	Register(method, path string, handler Handler)
	GetHandler(method, path string) (Handler, []UrlParam, bool)
}

type Handler func(http.ResponseWriter, *http.Request, []UrlParam)

type UrlParam struct {
	Key   []byte
	Value []byte
}

const MatchedRoutePathKey = "$matched_router_path"

// trieRouter 通过预先配置的路由将请求分发到不同的处理程序
type trieRouter struct {
	trees map[string]*node // key为http method
}

func New() Router {
	return &trieRouter{
		trees: make(map[string]*node, 5),
	}
}

func (r *trieRouter) Register(method, path string, handler Handler) {
	if method == "" {
		panic("method must not be empty")
	}

	if handler == nil {
		panic("handler must not be nil")
	}

	root := r.trees[method]
	if root == nil {
		r.trees[method] = &node{}
	}

	root.AddHandler([]byte(path), func(resp http.ResponseWriter, req *http.Request, p []UrlParam) {
		req.WithContext(context.WithValue(req.Context(), MatchedRoutePathKey, path))
		handler(resp, req, p)
	})
}

// 返回method和path对应的handler和参数，如果未找到则在最后一个参数为true时表示存在path添加或删除尾部'/'后的路径对应的handler
func (r *trieRouter) GetHandler(method, path string) (Handler, []UrlParam, bool) {
	if root := r.trees[method]; root != nil {
		return root.GetHandler([]byte(path))
	}
	return nil, nil, false
}
