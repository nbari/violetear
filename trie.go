package violetear

import (
	"strings"
)

type Trie struct {
	Node    map[string]*Trie
	handler map[string]string
}

func NewTrie() *Trie {
	t := &Trie{}
	t.Node = make(map[string]*Trie)
	t.handler = make(map[string]string)
	return t
}

func (t *Trie) Set(path []string, handler string, method ...string) {

	var methods string

	if len(method) > 0 {
		methods = method[0]
	}

	if len(path) == 0 {
		if len(methods) > 0 {
			methods := strings.Split(methods, ",")
			for _, v := range methods {
				t.handler[strings.TrimSpace(v)] = handler
			}
		} else {
			t.handler["ALL"] = handler
		}
		return
	}

	key := path[0]
	newpath := path[1:]

	res, ok := t.Node[key]

	if !ok {
		res = NewTrie()
		t.Node[key] = res
	}

	res.Set(newpath, handler, methods)
}

func (t *Trie) Get(path []string) (handler map[string]string, ok bool) {
	if len(path) == 0 {
		return t.handler, true
	}

	key := path[0]
	newpath := path[1:]

	res, ok := t.Node[key]

	if !ok {
		return nil, false
	}
	return res.Get(newpath)
}
