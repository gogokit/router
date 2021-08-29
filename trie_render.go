package router

import (
	"fmt"

	"github.com/58kg/tree_render"
)

type rNode struct {
	ids map[*node]string
	n   *node
}

func (n rNode) Id() string {
	return n.ids[n.n]
}

func (n rNode) Children() (ret []tree_render.Node) {
	for _, v := range n.n.children {
		ret = append(ret, rNode{
			ids: n.ids,
			n:   v,
		})
	}
	return ret
}

func (n rNode) String() string {
	return string(n.n.path) + func() string {
		if n.n.handler != nil {
			return " [#]"
		}
		return ""
	}()
}

func genIdByDFS(root *node, ids map[*node]string) {
	ids[root] = fmt.Sprintf("%d", len(ids))
	for _, v := range root.children {
		genIdByDFS(v, ids)
	}
}

func render(n *node) string {
	ids := make(map[*node]string)
	genIdByDFS(n, ids)
	return tree_render.Render(rNode{
		ids: ids,
		n:   n,
	}, 4)
}
