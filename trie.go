package violetear

import (
	"fmt"
	"strings"
)

type Trie struct {
	Node    map[string]*Trie
	Handler map[string]string
	Level   int
}

func NewTrie() *Trie {
	t := &Trie{}
	t.Node = make(map[string]*Trie)
	t.Handler = make(map[string]string)
	return t
}

func (t *Trie) Set(path []string, handler string, method string, level ...bool) {
	key := path[0]
	newpath := path[1:]

	val, ok := t.Node[key]

	if !ok {
		val = NewTrie()
		t.Node[key] = val

		// increment level
		if len(level) > 0 {
			val.Level = t.Level + 1
		}
	}

	if len(newpath) == 0 {
		methods := strings.Split(method, ",")
		for _, v := range methods {
			val.Handler[strings.TrimSpace(v)] = handler
		}
		return
	}

	// recursive call with 4 argument set to true so that level can be
	// increased by 1 if newpath > than 1
	val.Set(newpath, handler, method, true)
}

func (t *Trie) Get(path []string) *Trie {

	key := path[0]
	newpath := path[1:]

	fmt.Println(key, newpath, t.Level)

	if val, ok := t.Node[key]; ok {
		if len(newpath) == 0 {
			return val
		}
		return val.Get(newpath)
	}

	return t
}
