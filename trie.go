package router

import (
	"bytes"
	"errors"
	"sort"
)

type node struct {
	path     []byte
	children []*node
	handler  interface{}
}

func (n *node) Register(path []byte, h interface{}) {
	if h == nil {
		panic("handler must not be nil")
	}

	if err := verify(path); err != nil {
		panic(err.Error())
	}

	treePath := bytes.Buffer{}
	fullPath := string(path)

	// 每次循环头检查之前len(path)>0且n为空间点或n.path和path公共前缀长度大于0
	// treePath为根节点到结点n父结点的路径
	for {
		if len(n.path) == 0 {
			// n是空结点
			n.genTree(path, h)
			return
		}

		// 至此，n.path和path的公共前缀的长度必大于0
		l := longestCommonPrefix(n.path, path)
		if n.path[0] == ':' {
			// 此时必须满足如下条件：
			// n.path和path的公共前缀的长度等于n.path的长度且如果path的长度大于公共前缀的长度则path中位于公共前缀之后的首字符必须为'/'
			if !(l == len(n.path) && (l == len(path) || path[l] == '/')) {
				treePath.Write(n.getToMostLeftNodePath())
				panic("'" + fullPath + "' conflict with the registered path '" + treePath.String() + "'")
			}
		}

		if l < len(n.path) {
			// 分裂结点n
			n.children = []*node{
				{
					path:     n.path[l:],
					children: n.children,
					handler:  n.handler,
				},
			}
			n.path = n.path[:l]
			n.handler = nil
		}

		treePath.Write(n.path)

		if l == len(path) {
			if n.handler != nil {
				panic("the current path '" + fullPath + "' handler has been registered")
			}

			if n.isWildcardParent() && n.children[0].handler != nil {
				treePath.Write(n.children[0].getToMostLeftNodePath())
				panic("'" + fullPath + "' conflict with the registered path '" + treePath.String() + "'")
			}

			n.handler = h
			return
		}

		path = path[l:]

		// 检查孩子是否已经存在
		if v := n.findChildren(path[0]); v != nil {
			if n.handler != nil && isWildcardSegment(path) {
				panic("'" + fullPath + "' conflict with the registered path '" + treePath.String() + "'")
			}
			if v.path[0] == '*' {
				treePath.Write(v.getToMostLeftNodePath())
				panic("'" + fullPath + "' conflict with the registered path '" + treePath.String() + "'")
			}
			n = v
			continue
		}

		// 至此，n不存在以path[0]为首字符的孩子结点，此时必须满足以下条件：
		// 1：n不是叶结点时，path不是通配符段且n不是通配符结点的父结点
		// 2：n是叶结点时，path不是通配符段且n不是通配符结点的父结点
		// 3：path[0]为通配符时n必须为叶结点
		// 上面1和2化简为：path不是通配符段且n不是通配符结点的父结点

		if isWildcardSegment(path) || n.isWildcardParent() || (isWildcard(path[0]) && !n.isLeaf()) {
			if !n.isLeaf() {
				treePath.Write(n.children[0].getToMostLeftNodePath())
			}
			panic("'" + fullPath + "' conflict with the registered path '" + treePath.String() + "'")
		}

		// 插入新的孩子
		child := &node{}
		n.children = append(n.children, child)
		// 对n.children重试排序，使满足按照结点路径首字母递增排序
		for i := len(n.children) - 1; i > 0 && n.children[i-1].path[0] > path[0]; i-- {
			n.children[i], n.children[i-1] = n.children[i-1], n.children[i]
		}
		n = child
	}
}

func (n *node) Lookup(path []byte) (h interface{}, p []UrlParam, redirect bool) {
	if !(len(path) > 0 && path[0] == '/') {
		return nil, nil, false
	}
	var np *node // 结点n的父节点
walk:
	for {
		switch n.path[0] {
		case '*':
			p = append(p, UrlParam{
				Key:   n.path[1:],
				Value: path,
			})
			return n.handler, p, false
		case ':':
			for i, c := range path {
				if c == '/' {
					p = append(p, UrlParam{
						Key:   n.path[1:],
						Value: path[:i],
					})
					path = path[i:]
					if v := n.findChildren(path[0]); v != nil {
						np = n
						n = v
						continue walk
					}
					// 没找到该节点
					return nil, nil, isSlash(path) && n.canHandle()
				}
			}

			if n.handler != nil {
				p = append(p, UrlParam{
					Key:   n.path[1:],
					Value: path,
				})
				return n.handler, p, false
			}

			v := n.findChildren('/')
			return nil, nil, v.canHandle() && isSlash(v.path)
		default:
			l := longestCommonPrefix(n.path, path)
			if l < len(n.path) {
				return nil, nil, (isSlash(path) && np.canHandle()) || (path[len(path)-1] != '/' && l+1 == len(n.path) && n.path[l] == '/' && n.canHandle())
			}

			// 至此l == len(n.path)

			if l == len(path) {
				if n.handler != nil {
					return n.handler, p, false
				}

				if n.isWildcardParent() && n.children[0].handler != nil {
					n = n.children[0]
					path = []byte("")
					continue walk
				}

				if path[len(path)-1] == '/' {
					return nil, nil, isSlash(path) && np.canHandle()
				}

				v := n.findChildren('/')
				if v.canHandle() && isSlash(v.path) {
					return nil, nil, true
				}

				v = nil
				if n.isWildcardParent() {
					v = n.children[0].findChildren('/')
				}

				return nil, nil, v.canHandle() && isSlash(v.path)
			}

			path = path[l:]

			if n.isWildcardParent() {
				np = n
				n = n.children[0]
				continue walk
			}

			if v := n.findChildren(path[0]); v != nil {
				np = n
				n = v
				continue walk
			}
			return nil, nil, n.canHandle() && isSlash(path)
		}
	}
}

// 将path插入到以为n为根的空树中，要求len(path)>0
func (n *node) genTree(path []byte, h interface{}) {
	for {
		wildcard, idx := findWildcard(path)
		if idx < 0 {
			n.path = path
			n.handler = h
			return
		}

		if idx > 0 {
			n.path = path[:idx]
			n.children = append(n.children, &node{})
			n = n.children[0]
		}

		n.path = wildcard
		if idx+len(wildcard) == len(path) {
			n.handler = h
			return
		}

		path = path[idx+len(wildcard):]
		n.children = append(n.children, &node{})
		n = n.children[0]
	}
}

func (n *node) findChildren(firstChar byte) *node {
	pos := sort.Search(len(n.children), func(i int) bool {
		return n.children[i].path[0] >= firstChar
	})
	if pos < len(n.children) && n.children[pos].path[0] == firstChar {
		return n.children[pos]
	}
	return nil
}

func (n *node) isLeaf() bool {
	return len(n.children) == 0
}

// 返回n是否是通配符结点的父结点
func (n *node) isWildcardParent() bool {
	return n != nil && len(n.children) == 1 && (n.children[0].path[0] == '*' || n.children[0].path[0] == ':')
}

// 返回恰好匹配到结点n的路径是否能找到对应handler
func (n *node) canHandle() bool {
	return n != nil && (n.handler != nil || (n.isWildcardParent() && n.children[0].handler != nil))
}

// 返回n到以n为根的子树中最左边结点的路径
func (n *node) getToMostLeftNodePath() []byte {
	buf := bytes.Buffer{}
	for {
		buf.Write(n.path)
		if n.isLeaf() {
			break
		}
		n = n.children[0]
	}
	return buf.Bytes()
}

func isWildcard(c byte) bool {
	return c == ':' || c == '*'
}

func isWildcardSegment(path []byte) bool {
	if !isWildcard(path[0]) {
		return false
	}
	for _, c := range path[1:] {
		if c == '/' {
			return false
		}
	}
	return true
}
func isSlash(path []byte) bool {
	return len(path) == 1 && path[0] == '/'
}

// path路径合法性检查, 路径首字符必须为'/'
func verify(path []byte) error {
	if !(len(path) > 0 && path[0] == '/') {
		return errors.New("first char must be '/'")
	}
	var lastWildcard byte // 当前路径段的最后一个通配符
	for i, c := range path[1:] {
		if c == '/' {
			if lastWildcard == '*' {
				return errors.New("there should be no '/' after the wildcard '*'")
			}
			if i > 0 && isWildcard(path[i]) {
				return errors.New("the name of wildcard segment must not empty")
			}
			lastWildcard = 0
			continue
		}

		if !isWildcard(c) {
			continue
		}

		if lastWildcard != 0 {
			return errors.New("the wildcard '*' and ':' should not exist in the same path segment")
		}

		if c == '*' && i > 0 && path[i] != '/' {
			return errors.New("the previous character of '*' must be '/'")
		}

		lastWildcard = c
	}
	if isWildcard(path[len(path)-1]) {
		return errors.New("the name of wildcard segment must not empty")
	}
	return nil
}

// 返回a,b最长公共前缀的长度
func longestCommonPrefix(a, b []byte) int {
	minLen := len(a)
	if len(b) < minLen {
		minLen = len(b)
	}

	l := 0
	for l < minLen && a[l] == b[l] {
		l++
	}
	return l
}

// 返回path中的通配符段及其第1个字符在path中的索引，未找到通配符段时返回的索引值小于0
func findWildcard(path []byte) ([]byte, int) {
	for i, c := range path {
		if c == '*' {
			return path[i:], i
		}

		if c != ':' {
			continue
		}

		for j := i + 1; j < len(path); j++ {
			if path[j] == '/' {
				return path[i:j], i
			}
		}

		return path[i:], i
	}
	return nil, -1
}
