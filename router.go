package router

type Router interface {
	Register(method, path string, handler interface{})
	Lookup(method, path string) (handler interface{}, param []UrlParam, redirect bool)
}

type UrlParam struct {
	Key   []byte
	Value []byte
}

func New() Router {
	return &trieRouter{
		trees: make(map[string]*node, 5),
	}
}

// trieRouter 通过预先配置的路由将请求分发到不同的处理程序
type trieRouter struct {
	trees map[string]*node // key为http method
}

func (r *trieRouter) Register(method, path string, handler interface{}) {
	if method == "" {
		panic("method must not be empty")
	}

	if handler == nil {
		panic("handler must not be nil")
	}

	root := r.trees[method]
	if root == nil {
		root = &node{}
		r.trees[method] = root
	}

	root.Register([]byte(path), handler)
}

// 返回method和path对应的handler和参数，如果未找到则在最后一个参数为true时表示存在path添加或删除尾部'/'后的路径对应的handler
func (r *trieRouter) Lookup(method, path string) (interface{}, []UrlParam, bool) {
	if root := r.trees[method]; root != nil {
		return root.Lookup([]byte(path))
	}
	return nil, nil, false
}
