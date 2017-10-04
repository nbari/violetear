package violetear

import (
	"bytes"
	"errors"
	"net/http"
	"strings"
)

// MethodHandler keeps HTTP Method and http.handler
type MethodHandler struct {
	Method  string
	Handler http.Handler
}

// Trie data structure
type Trie struct {
	Handler     []MethodHandler
	HasCatchall bool
	HasRegex    bool
	Node        []*Trie
	path        string
	version     string
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

	node, ok := t.contains(key, version)

	if !ok {
		node = &Trie{
			path:    key,
			version: version,
		}
		t.Node = append(t.Node, node)

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
		methods := strings.FieldsFunc(method, func(c rune) bool {
			return c == ','
		})
		for _, v := range methods {
			node.Handler = append(node.Handler, MethodHandler{strings.ToUpper(strings.TrimSpace(v)), handler})
		}
		return nil
	}

	if key == "*" {
		return errors.New("Catch-all \"*\" must always be the final path element")
	}

	return node.Set(newpath, handler, method, version)
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

// Split path by "/"
func (t *Trie) Split(path string) []string {
	if path == "" {
		return []string{"/"}
	}
	var key bytes.Buffer
	for i, rune := range path {
		if rune == '/' && i > 0 {
			return []string{key.String(), path[i:]}
		} else if rune == '*' {
			return []string{"*"}
		} else if rune != '/' {
			key.WriteRune(rune)
		}
	}
	if key.Len() > 0 {
		return []string{key.String()}
	}
	return nil
}
