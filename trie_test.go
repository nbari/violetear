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
	tt := []struct {
		name    string
		path    []string
		method  string
		version string
		err     bool
	}{
		{"root", []string{"/"}, "ALL", "", false},
		{"root v3", []string{"/"}, "ALL", "v3", false},
		{"root", []string{"/"}, "ALL", "", false},
		{"root", []string{"root"}, "ALL", "", false},
		{"dyn", []string{":dnyamic"}, "ALL", "", false},
		{"*", []string{"*"}, "ALL", "", false},
		{"root", []string{"root", "*"}, "ALL", "", false},
		{"root", []string{"root", "alpha"}, "ALL", "", false},
		{"root", []string{"root", ":dynamic"}, "ALL", "", false},
		{"alpha", []string{"alpha", "beta", "gamma"}, "ALL", "", false},
		{"* error", []string{"root", "*", "beta"}, "ALL", "", true},
		{"* error", []string{"*", ":dynamic"}, "ALL", "", true},
		{"root", []string{"root", "alpha", "beta"}, "ALL", "", false},
		{"root 4", []string{"root", "alpha", "beta", "gamma"}, "ALL", "", false},
		{"root 4", []string{"root", "alpha", "beta", "gamma"}, "ALL", "v3", false},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := trie.Set(tc.path, nil, tc.method, tc.version)
			expect(t, err != nil, tc.err)
		})
	}
}

func TestTrieGet(t *testing.T) {
	trie := &Trie{}
	tt := []struct {
		path    []string
		method  string
		version string
	}{
		{[]string{"*"}, "ALL", ""},
		{[]string{"*"}, "ALL", "v3"},
		{[]string{":dynamic"}, "ALL", ""},
		{[]string{":dynamic"}, "ALL", "v3"},
		{[]string{"root"}, "ALL", ""},
		{[]string{"root"}, "ALL", "v3"},
		{[]string{"root", "*"}, "ALL", ""},
		{[]string{"root", "*"}, "ALL", "v3"},
		{[]string{"root", "alpha", "*"}, "ALL", ""},
		{[]string{"root", "alpha", "*"}, "ALL", ""},
		{[]string{"root", "alpha1", "*"}, "ALL", ""},
		{[]string{"root", "alpha1", "*"}, "ALL", ""},
		{[]string{"root", "alpha2", ":dynamic"}, "ALL", ""},
		{[]string{"root", "alpha", "beta", "gamma"}, "ALL", ""},
		{[]string{"alpha", "*"}, "ALL", ""},
	}
	for _, tc := range tt {
		err := trie.Set(tc.path, nil, tc.method, tc.version)
		expect(t, err, nil)
	}

	// Get
	ttGet := []struct {
		path    string
		version string
		node    int
		key     string
		newPath string
		leaf    bool
	}{
		{"", "", 0, "", "", false},
		{"*", "", 0, "*", "", true},
		{"*", "v5", 0, "*", "", false},
		{"*", "v4", 0, "*", "", false},
		{"*", "v3", 0, "*", "", true},
		{":dynamic", "", 0, ":dynamic", "", true},
		{":dynamic", "v3", 0, ":dynamic", "", true},
		{"root", "", 4, "root", "", true},
		{"root/v3", "v3", 1, "v3", "", false},
		{"root/alpha", "", 2, "alpha", "", true},
		{"root/not_found", "", 4, "not_found", "", false},
		{"root/alpha1", "", 1, "alpha1", "", true},
	}
	for _, tc := range ttGet {
		n, k, p, l := trie.Get(tc.path, tc.version)
		if tc.node > 0 {
			expect(t, len(n.Node), tc.node)
		}
		expect(t, k, tc.key)
		expect(t, p, tc.newPath)
		expect(t, l, tc.leaf)
	}

	n, k, p, l := trie.Get("not_found", "")
	expect(t, k, "not_found")
	expect(t, p, "")
	expect(t, l, false)
	expect(t, n.HasRegex, true)

	n, k, p, l = trie.Get("root/alpha1/any", "")
	expect(t, k, "any")
	expect(t, p, "")
	expect(t, l, false)
	expect(t, len(n.Node), 1)
	expect(t, n.HasRegex, false)

	n, k, p, l = trie.Get("root/alpha2/any", "")
	expect(t, k, "any")
	expect(t, p, "")
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
