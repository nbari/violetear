package violetear

import (
	"testing"
)

func TestTrieSetEmpty(t *testing.T) {
	trie := &Trie{}
	err := trie.Set([]string{}, nil, "ALL", "")
	if err == nil {
		t.Error("path cannot be empty")
	}
}

func TestTrieSet(t *testing.T) {
	trie := &Trie{}

	err := trie.Set([]string{"/"}, nil, "ALL", "")
	expect(t, err, nil)

	err = trie.Set([]string{"/"}, nil, "ALL", "v3")
	expect(t, err, nil)

	err = trie.Set([]string{"/"}, nil, "ALL", "")
	expect(t, err, nil)

	err = trie.Set([]string{"root"}, nil, "ALL", "")
	expect(t, err, nil)

	err = trie.Set([]string{":dynamic"}, nil, "ALL", "")
	expect(t, err, nil)

	err = trie.Set([]string{"*"}, nil, "ALL", "")
	expect(t, err, nil)

	err = trie.Set([]string{"root", "*"}, nil, "ALL", "")
	expect(t, err, nil)

	err = trie.Set([]string{"root", "alpha"}, nil, "ALL", "")
	expect(t, err, nil)

	err = trie.Set([]string{"root", ":dynamic"}, nil, "ALL", "")
	expect(t, err, nil)

	// Catch-all must always be the final path element.
	err = trie.Set([]string{"root", "*", "beta"}, nil, "ALL", "")
	trie.Set([]string{"alpha", "beta", "gamma"}, nil, "ALL", "")
	if err == nil {
		t.Error(err)
	}

	err = trie.Set([]string{"*", ":dynamic"}, nil, "ALL", "")
	if err == nil {
		t.Error(err)
	}

	err = trie.Set([]string{"root", "alpha", "beta"}, nil, "ALL", "")
	expect(t, err, nil)

	err = trie.Set([]string{"root", "alpha", "beta", "gamma"}, nil, "ALL", "")
	expect(t, err, nil)

	err = trie.Set([]string{"root", "alpha", "beta", "gamma"}, nil, "ALL", "v3")
	expect(t, err, nil)
}

func TestTrieGet(t *testing.T) {
	trie := &Trie{}
	err := trie.Set([]string{"*"}, nil, "ALL", "")
	expect(t, err, nil)
	err = trie.Set([]string{"*"}, nil, "ALL", "v3")
	expect(t, err, nil)

	err = trie.Set([]string{":dynamic"}, nil, "ALL", "")
	expect(t, err, nil)
	err = trie.Set([]string{":dynamic"}, nil, "ALL", "v3")
	expect(t, err, nil)

	err = trie.Set([]string{"root"}, nil, "ALL", "")
	expect(t, err, nil)
	err = trie.Set([]string{"root"}, nil, "ALL", "v3")
	expect(t, err, nil)

	err = trie.Set([]string{"root", "*"}, nil, "ALL", "")
	expect(t, err, nil)
	err = trie.Set([]string{"root", "*"}, nil, "ALL", "v3")
	expect(t, err, nil)

	err = trie.Set([]string{"root", "alpha", "*"}, nil, "ALL", "")
	expect(t, err, nil)

	err = trie.Set([]string{"root", "alpha1", "*"}, nil, "ALL", "")
	expect(t, err, nil)

	err = trie.Set([]string{"root", "alpha2", ":dynamic"}, nil, "ALL", "")
	expect(t, err, nil)

	err = trie.Set([]string{"root", "alpha", "beta", "gamma"}, nil, "ALL", "")
	expect(t, err, nil)

	err = trie.Set([]string{"alpha", "*"}, nil, "ALL", "")
	expect(t, err, nil)

	_, k, p, l := trie.Get("", "")
	expect(t, k, "")
	expect(t, p, "")
	expect(t, l, false)

	_, k, p, l = trie.Get("*", "")
	expect(t, k, "*")
	expect(t, p, "")
	expect(t, l, true)

	_, k, p, l = trie.Get("*", "v5")
	expect(t, k, "*")
	expect(t, p, "")
	expect(t, l, false)

	_, k, p, l = trie.Get("*", "v4")
	expect(t, k, "*")
	expect(t, p, "")
	expect(t, l, false)

	_, k, p, l = trie.Get("*", "v3")
	expect(t, k, "*")
	expect(t, p, "")
	expect(t, l, true)

	n, k, p, l := trie.Get("not_found", "")
	expect(t, k, "not_found")
	expect(t, p, "")
	expect(t, l, false)
	expect(t, n.HasRegex, true)

	n, k, p, l = trie.Get(":dynamic", "")
	expect(t, k, ":dynamic")
	expect(t, p, "")
	expect(t, l, true)

	n, k, p, l = trie.Get(":dynamic", "v3")
	expect(t, k, ":dynamic")
	expect(t, p, "")
	expect(t, l, true)

	n, k, p, l = trie.Get("root", "")
	expect(t, k, "root")
	expect(t, p, "")
	expect(t, l, true)
	expect(t, len(n.Node), 4)

	n, k, p, l = trie.Get("root/v3", "v3")
	expect(t, k, "v3")
	expect(t, p, "")
	expect(t, l, false)
	expect(t, len(n.Node), 1)

	n, k, p, l = trie.Get("root/alpha", "")
	expect(t, k, "alpha")
	expect(t, p, "")
	expect(t, l, true)
	expect(t, len(n.Node), 2)

	n, k, p, l = trie.Get("root/not_found", "")
	expect(t, k, "not_found")
	expect(t, p, "")
	expect(t, l, false)
	expect(t, len(n.Node), 4)

	n, k, p, l = trie.Get("root/alpha1", "")
	expect(t, k, "alpha1")
	expect(t, p, "")
	expect(t, l, true)
	expect(t, len(n.Node), 1)

	n, k, p, l = trie.Get("root/alpha1/any", "")
	expect(t, k, "any")
	expect(t, p, "")
	expect(t, l, false)
	expect(t, len(n.Node), 1)
	expect(t, n.HasRegex, false)

	n, k, p, l = trie.Get("root/alpha2/any", "")
	expect(t, k, "any")
	expect(t, p, "")
	expect(t, err, nil)
	expect(t, l, false)
	expect(t, len(n.Node), 1)
	expect(t, n.HasRegex, true)

	n, k, p, l = trie.Get("root/alpha/beta", "")
	expect(t, k, "beta")
	expect(t, p, "")
	expect(t, l, true)
	expect(t, len(n.Node), 1)
	expect(t, n.HasRegex, false)

	n, k, p, l = trie.Get("root/alpha/beta/gamma", "")
	expect(t, k, "gamma")
	expect(t, p, "")
	expect(t, l, true)
	expect(t, len(n.Node), 0)
	expect(t, n.HasRegex, false)

	n, k, p, l = trie.Get("root/alphaA/betaB/gammaC", "")
	expect(t, k, "alphaA")
	expect(t, p, "/betaB/gammaC")
	expect(t, l, false)
	expect(t, len(n.Node), 4)
	expect(t, n.HasRegex, false)

	n, k, p, l = trie.Get("root/alpha/betaB/gammaC", "")
	expect(t, k, "betaB")
	expect(t, p, "/gammaC")
	expect(t, l, false)
	expect(t, len(n.Node), 2)
	expect(t, n.HasRegex, false)

	n, k, p, l = trie.Get("root/alpha/betaB/gamma/delta", "")
	expect(t, k, "betaB")
	expect(t, p, "/gamma/delta")
	expect(t, l, false)
	expect(t, len(n.Node), 2)
	expect(t, n.HasRegex, false)
}

func TestSplitPath(t *testing.T) {
	tt := []struct {
		in  string
		out []string
	}{
		{"/", []string{"/", ""}},
		{"//", []string{"/", ""}},
		{"///", []string{"/", ""}},
		{"////", []string{"/", ""}},
		{"/////", []string{"/", ""}},
		{"/hello", []string{"hello", ""}},
		{"/hello/world", []string{"hello", "/world"}},
		{"/hello/:world", []string{"hello", "/:world"}},
		{"*", []string{"*", ""}},
		{"/?foo=bar", []string{"?foo=bar", ""}},
	}

	trie := &Trie{}
	for _, tc := range tt {
		k, p := trie.SplitPath(tc.in)
		expect(t, k, tc.out[0])
		expect(t, p, tc.out[1])
	}
}
