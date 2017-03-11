package violetear

import (
	"errors"
	"net/http"
	"strings"
)

// Trie data structure
type Trie struct {
	Handler     map[string]http.Handler
	HasCatchall bool
	HasRegex    bool
	Node        []*Trie
	path        string
	version     string
}

// NewTrie returns a new Trie
func NewTrie() *Trie {
	return &Trie{
		Node:    make([]*Trie, 0),
		Handler: map[string]http.Handler{},
	}
}

// contains check if path exists on node
func (t *Trie) contains(path, version string) (*Trie, bool) {
	for _, n := range t.Node {
		if n.path == path && n.version == version {
			return n, true
		}
	}
	return nil, false
}

// Set adds a node (url part) to the Trie
func (t *Trie) Set(path []string, handler http.Handler, method, version string) error {
	if len(path) == 0 {
		return errors.New("path cannot be empty")
	}

	key := path[0]
	newpath := path[1:]

	val, ok := t.contains(key, version)

	if !ok {
		val = NewTrie()
		val.path = key
		val.version = version
		t.Node = append(t.Node, val)

		// check for regex ":"
		if strings.HasPrefix(key, ":") {
			t.HasRegex = true
		}

		// check for Catch-all "*"
		if key == "*" {
			t.HasCatchall = true
		}
	}

	if len(newpath) == 0 {
		methods := strings.Split(method, ",")
		for _, v := range methods {
			val.Handler[strings.ToUpper(strings.TrimSpace(v))] = handler
		}
		return nil
	}

	if key == "*" {
		return errors.New("Catch-all \"*\" must always be the final path element.")
	}

	return val.Set(newpath, handler, method, version)
}

// Get returns a node
func (t *Trie) Get(path []string, version string) (trie *Trie, p []string, leaf bool, err error) {
	if len(path) == 0 {
		err = errors.New("path cannot be empty")
		return
	}

	key := path[0]
	newpath := path[1:]

	if val, ok := t.contains(key, version); ok {
		if len(newpath) == 0 {
			return val, path, true, nil
		}
		return val.Get(newpath, version)
	}

	return t, path, false, nil
}
