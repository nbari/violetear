package violetear

import (
	"strings"
)

type Trie struct {
	Node     map[string]*Trie
	Handler  map[string]string
	HasRegex bool
}

func NewTrie() *Trie {
	t := &Trie{}
	t.Node = make(map[string]*Trie)
	t.Handler = make(map[string]string)
	return t
}

func (t *Trie) Set(path []string, handler string, method string) {
	key := path[0]
	newpath := path[1:]

	val, ok := t.Node[key]

	if !ok {
		val = NewTrie()
		t.Node[key] = val

		// check for regex ":"
		if strings.HasPrefix(key, ":") {
			t.HasRegex = true
		}
	}

	if len(newpath) == 0 {
		methods := strings.Split(method, ",")
		for _, v := range methods {
			val.Handler[strings.ToUpper(strings.TrimSpace(v))] = handler
		}
		return
	}

	val.Set(newpath, handler, method)
}

func (t *Trie) Get(path []string) (trie *Trie, leaf bool) {
	key := path[0]
	newpath := path[1:]

	if val, ok := t.Node[key]; ok {
		if len(newpath) == 0 {
			return val, true
		}
		return val.Get(newpath)
	}
	return t, false
}
