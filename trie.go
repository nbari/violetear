package violetear

import (
	"fmt"
	"strings"
)

type Trie struct {
	node    map[string]*Trie
	handler map[string]string
	level   int
}

func NewTrie() *Trie {
	t := &Trie{}
	t.node = make(map[string]*Trie)
	t.handler = make(map[string]string)
	return t
}

func (t *Trie) Set(path []string, handler string, method string, level ...bool) {
	key := path[0]
	newpath := path[1:]

	val, ok := t.node[key]

	if !ok {
		val = NewTrie()
		t.node[key] = val

		// increment level
		if len(level) > 0 {
			val.level = t.level + 1
		}
	}

	fmt.Println(val.level, key, newpath)

	if len(newpath) == 0 {
		methods := strings.Split(method, ",")
		for _, v := range methods {
			val.handler[strings.TrimSpace(v)] = handler
		}
		return
	}

	// recursive call with 4 argument set to true so that level can be
	// increased by 1 if newpath > than 1
	val.Set(newpath, handler, method, true)
}

func (t *Trie) Get(path []string) (level int, handler map[string]string) {

	key := path[0]
	newpath := path[1:]

	// check if the node on the trie exists and return current handler
	if val, ok := t.node[key]; ok {
		if len(newpath) == 0 {
			return val.level, val.handler
		}
		return val.Get(newpath)
	}

	///////
	fmt.Println("find the : regex")
	////

	return t.level, nil
}
