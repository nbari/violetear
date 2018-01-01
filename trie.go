package violetear

import (
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
	name        string
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
func (t *Trie) Set(path []string, handler http.Handler, method, version string) (*Trie, error) {
	if len(path) == 0 {
		return nil, errors.New("path cannot be empty")
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
		return node, nil
	}

	if key == "*" {
		return nil, errors.New("catch-all \"*\" must always be the final path element")
	}

	return node.Set(newpath, handler, method, version)
}

// Get returns a node
func (t *Trie) Get(path, version string) (*Trie, string, string, bool) {
	key, path := t.SplitPath(path)
	// search the key recursively on the tree
	if node, ok := t.contains(key, version); ok {
		if path == "" {
			return node, key, path, true
		}
		return node.Get(path, version)
	}
	// if not fount check for catchall or regex
	return t, key, path, false
}

// SplitPath returns first element of path and remaining path
func (t *Trie) SplitPath(path string) (string, string) {
	var key string
	if path == "" {
		return key, path
	} else if path == "/" {
		return path, ""
	}
	for i := 0; i < len(path); i++ {
		if path[i] == '/' {
			if i == 0 {
				return t.SplitPath(path[1:])
			}
			if i > 0 {
				key = path[:i]
				path = path[i:]
				if path == "/" {
					return key, ""
				}
				return key, path
			}
		}
	}
	return path, ""
}

// Name add custom name to node
func (t *Trie) Name(name string) *Trie {
	t.name = name
	return t
}
