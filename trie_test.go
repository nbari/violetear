package violetear

import (
	"fmt"
	"net/http"
	"testing"
)

func TestTrieNew(t *testing.T) {
	trie := NewTrie()
	my_trie := &Trie{
		Node:    make(map[string]*Trie),
		Handler: make(map[string]http.Handler),
	}
	expectDeepEqual(t, trie, my_trie)
}

func TestTrieSetEmpty(t *testing.T) {
	trie := NewTrie()
	err := trie.Set([]string{}, nil, "ALL")
	if err == nil {
		t.Error("path cannot be empty")
	}
}

func TestTrieSet(t *testing.T) {
	trie := NewTrie()

	err := trie.Set([]string{"/"}, nil, "ALL")
	expect(t, err, nil)

	err = trie.Set([]string{"/"}, nil, "ALL")
	expect(t, err, nil)

	err = trie.Set([]string{"root"}, nil, "ALL")
	expect(t, err, nil)

	err = trie.Set([]string{":dynamic"}, nil, "ALL")
	expect(t, err, nil)

	err = trie.Set([]string{"*"}, nil, "ALL")
	expect(t, err, nil)

	err = trie.Set([]string{"root", "*"}, nil, "ALL")
	expect(t, err, nil)

	err = trie.Set([]string{"root", "alpha"}, nil, "ALL")
	expect(t, err, nil)

	err = trie.Set([]string{"root", ":dynamic"}, nil, "ALL")
	expect(t, err, nil)

	// Catch-all must always be the final path element.
	err = trie.Set([]string{"root", "*", "beta"}, nil, "ALL")
	trie.Set([]string{"alpha", "beta", "gamma"}, nil, "ALL")
	if err == nil {
		t.Error(err)
	}

	err = trie.Set([]string{"*", ":dynamic"}, nil, "ALL")
	if err == nil {
		t.Error(err)
	}

	err = trie.Set([]string{"root", "alpha", "beta"}, nil, "ALL")
	expect(t, err, nil)

	err = trie.Set([]string{"root", "alpha", "beta", "gamma"}, nil, "ALL")
	expect(t, err, nil)
}

func TestTrieGet(t *testing.T) {
	trie := NewTrie()
	err := trie.Set([]string{"*"}, nil, "ALL")
	expect(t, err, nil)

	err = trie.Set([]string{":dynamic"}, nil, "ALL")
	expect(t, err, nil)

	err = trie.Set([]string{"root"}, nil, "ALL")
	expect(t, err, nil)

	err = trie.Set([]string{"root", "*"}, nil, "ALL")
	expect(t, err, nil)

	err = trie.Set([]string{"root", "alpha", "*"}, nil, "ALL")
	expect(t, err, nil)

	err = trie.Set([]string{"root", "alpha1", "*"}, nil, "ALL")
	expect(t, err, nil)

	err = trie.Set([]string{"root", "alpha2", "*"}, nil, "ALL")
	expect(t, err, nil)

	err = trie.Set([]string{"root", "alpha", "beta", "gamma"}, nil, "ALL")
	expect(t, err, nil)

	err = trie.Set([]string{"alpha", "*"}, nil, "ALL")
	expect(t, err, nil)

	_, _, _, err = trie.Get([]string{})
	if err == nil {
		t.Error(err)
	}

	n, p, l, err := trie.Get([]string{"*"})
	expect(t, err, nil)
	fmt.Println(n, p, l)

}
