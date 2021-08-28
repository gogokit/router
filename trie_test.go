package router

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gogokit/treeprint"

	. "github.com/smartystreets/goconvey/convey"
)

// 测试普通结构体
func TestTrie(t *testing.T) {
	Convey("Register", t, func() {
		Convey("no_conflict", func() {
			paths := []string{
				"/",
				"//",
				"///",
				"/a",
				"/a/",
				"/aa",
				"/aa/",
				"/aa/:version1/",
				"/aa/:version1/:version2/",
				"/aa/:version1/:version2/a/*all",
				"/bbb/*all",
				"/bbc/",
				"/bbc",
			}
			/*
			                    / [#]
			                    |
			  --------------------------------------
			  |            |                       |
			  / [#]        a [#]                   bb
			  |            |                       |
			  |        ----------            -------------
			  |        |        |            |           |
			  / [#]    / [#]    a [#]        b/          c [#]
			                    |            |           |
			                    / [#]        *all [#]    / [#]
			                    |
			                    :version1
			                    |
			                    / [#]
			                    |
			                    :version2
			                    |
			                    / [#]
			                    |
			                    a/
			                    |
			                    *all [#]
			*/
			expect := []byte{32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 47, 32, 91, 35, 93, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 124, 10, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 10, 124, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 124, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 124, 10, 47, 32, 91, 35, 93, 32, 32, 32, 32, 32, 32, 32, 32, 97, 32, 91, 35, 93, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 98, 98, 10, 124, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 124, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 124, 10, 124, 32, 32, 32, 32, 32, 32, 32, 32, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 10, 124, 32, 32, 32, 32, 32, 32, 32, 32, 124, 32, 32, 32, 32, 32, 32, 32, 32, 124, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 124, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 124, 10, 47, 32, 91, 35, 93, 32, 32, 32, 32, 47, 32, 91, 35, 93, 32, 32, 32, 32, 97, 32, 91, 35, 93, 32, 32, 32, 32, 32, 32, 32, 32, 98, 47, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 99, 32, 91, 35, 93, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 124, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 124, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 124, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 47, 32, 91, 35, 93, 32, 32, 32, 32, 32, 32, 32, 32, 42, 97, 108, 108, 32, 91, 35, 93, 32, 32, 32, 32, 47, 32, 91, 35, 93, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 124, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 58, 118, 101, 114, 115, 105, 111, 110, 49, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 124, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 47, 32, 91, 35, 93, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 124, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 58, 118, 101, 114, 115, 105, 111, 110, 50, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 124, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 47, 32, 91, 35, 93, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 124, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 97, 47, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 124, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 42, 97, 108, 108, 32, 91, 35, 93}
			root := &node{}
			for _, v := range paths {
				root.Register([]byte(v), func(rw http.ResponseWriter, r *http.Request, up []UrlParam) {})
			}
			So(render(root), ShouldEqual, string(expect))
		})

		Convey("path_illegal_1", func() {
			paths := []string{
				"/a",
				"/a:/",
			}
			expect := "the name of wildcard segment must not empty"
			root := &node{}
			var errMsg string
			for _, v := range paths {
				func() {
					defer func() {
						if err := recover(); err != nil {
							errMsg = err.(string)
						}
					}()
					root.Register([]byte(v), func(rw http.ResponseWriter, r *http.Request, up []UrlParam) {})
				}()
			}
			So(errMsg, ShouldEqual, expect)
		})

		Convey("path_illegal_2", func() {
			paths := []string{
				"/a",
				"/a*/",
			}
			expect := "the previous character of '*' must be '/'"
			root := &node{}
			var errMsg string
			for _, v := range paths {
				func() {
					defer func() {
						if err := recover(); err != nil {
							errMsg = err.(string)
						}
					}()
					root.Register([]byte(v), func(rw http.ResponseWriter, r *http.Request, up []UrlParam) {})
				}()
			}
			So(errMsg, ShouldEqual, expect)
		})

		Convey("path_illegal_3", func() {
			paths := []string{
				"/a/",
				"/a/*",
			}
			expect := "the name of wildcard segment must not empty"
			root := &node{}
			var errMsg string
			for _, v := range paths {
				func() {
					defer func() {
						if err := recover(); err != nil {
							errMsg = err.(string)
						}
					}()
					root.Register([]byte(v), func(rw http.ResponseWriter, r *http.Request, up []UrlParam) {})
				}()
			}
			So(errMsg, ShouldEqual, expect)
		})

		Convey("conflict_because_same_path", func() {
			paths := []string{
				"/a",
				"/a",
			}
			expect := "the current path '/a' handler has been registered"
			root := &node{}
			var errMsg string
			for _, v := range paths {
				func() {
					defer func() {
						if err := recover(); err != nil {
							errMsg = err.(string)
						}
					}()
					root.Register([]byte(v), func(rw http.ResponseWriter, r *http.Request, up []UrlParam) {})
				}()
			}
			So(errMsg, ShouldEqual, expect)
		})

		Convey("conflict_because_path_tail_wildcard_1", func() {
			paths := []string{
				"/a",
				"/a:version",
			}
			expect := "'/a:version' conflict with the registered path '/a'"
			root := &node{}
			var errMsg string
			for _, v := range paths {
				func() {
					defer func() {
						if err := recover(); err != nil {
							errMsg = err.(string)
						}
					}()
					root.Register([]byte(v), func(rw http.ResponseWriter, r *http.Request, up []UrlParam) {})
				}()
			}
			So(errMsg, ShouldEqual, expect)
		})

		Convey("conflict_because_path_tail_wildcard_2", func() {
			paths := []string{
				"/a/",
				"/a/:version",
			}
			expect := "'/a/:version' conflict with the registered path '/a/'"
			root := &node{}
			var errMsg string
			for _, v := range paths {
				func() {
					defer func() {
						if err := recover(); err != nil {
							errMsg = err.(string)
						}
					}()
					root.Register([]byte(v), func(rw http.ResponseWriter, r *http.Request, up []UrlParam) {})
				}()
			}
			So(errMsg, ShouldEqual, expect)
		})

		Convey("conflict_because_path_wildcard_*_1", func() {
			paths := []string{
				"/a/",
				"/a/*all",
			}
			expect := "'/a/*all' conflict with the registered path '/a/'"
			root := &node{}
			var errMsg string
			for _, v := range paths {
				func() {
					defer func() {
						if err := recover(); err != nil {
							errMsg = err.(string)
						}
					}()
					root.Register([]byte(v), func(rw http.ResponseWriter, r *http.Request, up []UrlParam) {})
				}()
			}
			So(errMsg, ShouldEqual, expect)
		})

		Convey("conflict_because_path_wildcard_*_2", func() {
			paths := []string{
				"/a/*al",
				"/a/*all",
			}
			expect := "'/a/*all' conflict with the registered path '/a/*al'"
			root := &node{}
			var errMsg string
			for _, v := range paths {
				func() {
					defer func() {
						if err := recover(); err != nil {
							errMsg = err.(string)
						}
					}()
					root.Register([]byte(v), func(rw http.ResponseWriter, r *http.Request, up []UrlParam) {})
				}()
			}
			So(errMsg, ShouldEqual, expect)
		})

		Convey("conflict_because_path_wildcard_:_1", func() {
			paths := []string{
				"/a/",
				"/a/:all",
			}
			expect := "'/a/:all' conflict with the registered path '/a/'"
			root := &node{}
			var errMsg string
			for _, v := range paths {
				func() {
					defer func() {
						if err := recover(); err != nil {
							errMsg = err.(string)
						}
					}()
					root.Register([]byte(v), func(rw http.ResponseWriter, r *http.Request, up []UrlParam) {})
				}()
			}
			So(errMsg, ShouldEqual, expect)
		})

		Convey("conflict_because_path_wildcard_:_2", func() {
			paths := []string{
				"/a/:al",
				"/a/:all",
			}
			expect := "'/a/:all' conflict with the registered path '/a/:al'"
			root := &node{}
			var errMsg string
			for _, v := range paths {
				func() {
					defer func() {
						if err := recover(); err != nil {
							errMsg = err.(string)
						}
					}()
					root.Register([]byte(v), func(rw http.ResponseWriter, r *http.Request, up []UrlParam) {})
				}()
			}
			So(errMsg, ShouldEqual, expect)
		})

		Convey("conflict_because_path_wildcard_:_*", func() {
			paths := []string{
				"/a/:all",
				"/a/*all",
			}
			expect := "'/a/*all' conflict with the registered path '/a/:all'"
			root := &node{}
			var errMsg string
			for _, v := range paths {
				func() {
					defer func() {
						if err := recover(); err != nil {
							errMsg = err.(string)
						}
					}()
					root.Register([]byte(v), func(rw http.ResponseWriter, r *http.Request, up []UrlParam) {})
				}()
			}
			So(errMsg, ShouldEqual, expect)
		})

		Convey("conflict_because_path_wildcard_/:/", func() {
			paths := []string{
				"/a//",
				"/a/:version/",
			}
			expect := "'/a/:version/' conflict with the registered path '/a//'"
			root := &node{}
			var errMsg string
			for _, v := range paths {
				func() {
					defer func() {
						if err := recover(); err != nil {
							errMsg = err.(string)
						}
					}()
					root.Register([]byte(v), func(rw http.ResponseWriter, r *http.Request, up []UrlParam) {})
				}()
			}
			So(errMsg, ShouldEqual, expect)
		})
	})

	Convey("Lookup", t, func() {
		paths := []string{
			"/",
			"//",
			"///",
			"/a",
			"/a/",
			"/aa",
			"/aa/",
			"/aa/:version1/",
			"/aa/:version1/:version2/",
			"/aa/:version1/:version2/a/*all",
			"/aa/:version1/:version2/b:version3/*all",
			"/c/:version1/:version2/:version3//",
			"/d/:version1/*all",
			"/bbb/*all",
			"/bbc/",
			"/bbc",
		}

		root := &node{}
		var s string
		for _, v := range paths {
			vv := v
			root.Register([]byte(v), func(rw http.ResponseWriter, r *http.Request, up []UrlParam) {
				s = vv
			})
		}

		/*
		                    / [#]
		                    |
		  --------------------------------------
		  |            |                       |
		  / [#]        a [#]                   bb
		  |            |                       |
		  |        ----------            -------------
		  |        |        |            |           |
		  / [#]    / [#]    a [#]        b/          c [#]
		                    |            |           |
		                    / [#]        *all [#]    / [#]
		                    |
		                    :version1
		                    |
		                    / [#]
		                    |
		                    :version2
		                    |
		                    / [#]
		                    |
		                    a/
		                    |
		                    *all [#]
		*/

		Convey("match_absolute_path", func() {
			ppaths := []string{
				"/",
				"//",
				"///",
				"/a",
				"/a/",
				"/aa",
				"/aa/",
				"/bbc/",
				"/bbc",
			}

			for _, v := range ppaths {
				h, param, tsr := root.Lookup([]byte(v))
				So(h, ShouldNotEqual, nil)
				So(len(param), ShouldEqual, 0)
				So(tsr, ShouldEqual, false)
				h.(func(rw http.ResponseWriter, r *http.Request, up []UrlParam))(nil, nil, nil)
				So(s, ShouldEqual, v)
			}
		})

		Convey("match_param", func() {
			h, param, tsr := root.Lookup([]byte("/aa/param1/"))
			So(h, ShouldNotEqual, nil)
			So(len(param), ShouldEqual, 1)
			So(string(param[0].Key), ShouldEqual, "version1")
			So(string(param[0].Value), ShouldEqual, "param1")
			So(tsr, ShouldEqual, false)
			h.(func(rw http.ResponseWriter, r *http.Request, up []UrlParam))(nil, nil, nil)
			So(s, ShouldEqual, "/aa/:version1/")

			h, param, tsr = root.Lookup([]byte("/aa/param1/param2/"))
			So(h, ShouldNotEqual, nil)
			So(len(param), ShouldEqual, 2)
			So(string(param[0].Key), ShouldEqual, "version1")
			So(string(param[0].Value), ShouldEqual, "param1")
			So(string(param[1].Key), ShouldEqual, "version2")
			So(string(param[1].Value), ShouldEqual, "param2")
			So(tsr, ShouldEqual, false)
			h.(func(rw http.ResponseWriter, r *http.Request, up []UrlParam))(nil, nil, nil)
			So(s, ShouldEqual, "/aa/:version1/:version2/")

			h, param, tsr = root.Lookup([]byte("/aa/param1/param2/a/param3"))
			So(h, ShouldNotEqual, nil)
			So(len(param), ShouldEqual, 3)
			So(string(param[0].Key), ShouldEqual, "version1")
			So(string(param[0].Value), ShouldEqual, "param1")
			So(string(param[1].Key), ShouldEqual, "version2")
			So(string(param[1].Value), ShouldEqual, "param2")
			So(string(param[2].Key), ShouldEqual, "all")
			So(string(param[2].Value), ShouldEqual, "param3")
			So(tsr, ShouldEqual, false)
			h.(func(rw http.ResponseWriter, r *http.Request, up []UrlParam))(nil, nil, nil)
			So(s, ShouldEqual, "/aa/:version1/:version2/a/*all")

			h, param, tsr = root.Lookup([]byte("/bbb/param1"))
			So(h, ShouldNotEqual, nil)
			So(len(param), ShouldEqual, 1)
			So(string(param[0].Key), ShouldEqual, "all")
			So(string(param[0].Value), ShouldEqual, "param1")
			So(tsr, ShouldEqual, false)
			h.(func(rw http.ResponseWriter, r *http.Request, up []UrlParam))(nil, nil, nil)
			So(s, ShouldEqual, "/bbb/*all")
		})

		Convey("redirect", func() {
			h, param, tsr := root.Lookup([]byte("/aa/param1"))
			So(h, ShouldEqual, nil)
			So(len(param), ShouldEqual, 0)
			So(tsr, ShouldEqual, true)

			h, param, tsr = root.Lookup([]byte("/aa/param1/param2//"))
			So(h, ShouldEqual, nil)
			So(len(param), ShouldEqual, 0)
			So(tsr, ShouldEqual, true)

			h, param, tsr = root.Lookup([]byte("/aa/param1/param2/a"))
			So(h, ShouldEqual, nil)
			So(len(param), ShouldEqual, 0)
			So(tsr, ShouldEqual, true)

			h, param, tsr = root.Lookup([]byte("/aa/param1/param2"))
			So(h, ShouldEqual, nil)
			So(len(param), ShouldEqual, 0)
			So(tsr, ShouldEqual, true)

			h, param, tsr = root.Lookup([]byte("/bbb"))
			So(h, ShouldEqual, nil)
			So(len(param), ShouldEqual, 0)
			So(tsr, ShouldEqual, true)

			// "/aa/:version1/:version2/b:version3/*all",
			h, param, tsr = root.Lookup([]byte("/aa/p1/p2/b"))
			So(h, ShouldEqual, nil)
			So(len(param), ShouldEqual, 0)
			So(tsr, ShouldEqual, true)

			h, param, tsr = root.Lookup([]byte("/bb"))
			So(h, ShouldEqual, nil)
			So(len(param), ShouldEqual, 0)
			So(tsr, ShouldEqual, false)

			h, param, tsr = root.Lookup([]byte("/c/p1/p2/p3"))
			So(h, ShouldEqual, nil)
			So(len(param), ShouldEqual, 0)
			So(tsr, ShouldEqual, false)

			h, param, tsr = root.Lookup([]byte("/c/p1/p2/p3/"))
			So(h, ShouldEqual, nil)
			So(len(param), ShouldEqual, 0)
			So(tsr, ShouldEqual, false)

			h, param, tsr = root.Lookup([]byte("/d/p"))
			So(h, ShouldEqual, nil)
			So(len(param), ShouldEqual, 0)
			So(tsr, ShouldEqual, true)
		})
	})
}

type rNode struct {
	ids map[*node]string
	n   *node
}

func (n rNode) Id() string {
	return n.ids[n.n]
}

func (n rNode) Children() (ret []treeprint.Node) {
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
	return treeprint.Print(rNode{
		ids: ids,
		n:   n,
	}, 4)
}
