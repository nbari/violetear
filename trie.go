package violetear

type Trie struct {
	Node    map[string]*Trie
	handler map[string]string
}

func NewTrie() *Trie {
	t := &Trie{}
	t.Node = make(map[string]*Trie)
	return t
}

func (t *Trie) Set(path []string, handler map[string]string) {
	if len(path) == 0 {
		t.handler = handler
		return
	}

	key := path[0]
	newpath := path[1:]

	res, ok := t.Node[key]

	if !ok {
		res = NewTrie()
		t.Node[key] = res
	}

	res.Set(newpath, handler)

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
