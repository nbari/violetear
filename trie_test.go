package violetear

/*
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

	_, _, _, err = trie.Get([]string{}, "")
	if err == nil {
		t.Error(err)
	}

	n, p, l, err := trie.Get([]string{"*"}, "")
	expect(t, err, nil)
	expectDeepEqual(t, p, []string{"*"})
	expect(t, l, true)
	n, p, l, err = trie.Get([]string{"*"}, "v5")
	expect(t, err, nil)
	n, p, l, err = trie.Get([]string{"*"}, "v4")
	expect(t, err, nil)
	n, p, l, err = trie.Get([]string{"*"}, "v3")
	expect(t, err, nil)
	expectDeepEqual(t, p, []string{"*"})
	expect(t, l, true)

	n, p, l, err = trie.Get([]string{"not_found"}, "")
	expect(t, err, nil)
	expectDeepEqual(t, p, []string{"not_found"})
	expect(t, l, false)
	expect(t, n.HasRegex, true)

	n, p, l, err = trie.Get([]string{":dynamic"}, "")
	expect(t, err, nil)
	expectDeepEqual(t, p, []string{":dynamic"})
	expect(t, l, true)
	n, p, l, err = trie.Get([]string{":dynamic"}, "v3")
	expect(t, err, nil)
	expectDeepEqual(t, p, []string{":dynamic"})
	expect(t, l, true)

	n, p, l, err = trie.Get([]string{"root"}, "")
	expect(t, err, nil)
	expectDeepEqual(t, p, []string{"root"})
	expect(t, l, true)
	expect(t, len(n.Node), 4)

	n, p, l, err = trie.Get([]string{"root", "v3"}, "v3")
	expect(t, err, nil)
	expectDeepEqual(t, p, []string{"v3"})
	expect(t, l, false)
	expect(t, len(n.Node), 1)

	n, p, l, err = trie.Get([]string{"root", "alpha"}, "")
	expect(t, err, nil)
	expectDeepEqual(t, p, []string{"alpha"})
	expect(t, l, true)
	expect(t, len(n.Node), 2)

	n, p, l, err = trie.Get([]string{"root", "not_found"}, "")
	expect(t, err, nil)
	expectDeepEqual(t, p, []string{"not_found"})
	expect(t, l, false)
	expect(t, len(n.Node), 4)

	n, p, l, err = trie.Get([]string{"root", "alpha1"}, "")
	expect(t, err, nil)
	expectDeepEqual(t, p, []string{"alpha1"})
	expect(t, l, true)
	expect(t, len(n.Node), 1)

	n, p, l, err = trie.Get([]string{"root", "alpha1", "any"}, "")
	expect(t, err, nil)
	expectDeepEqual(t, p, []string{"any"})
	expect(t, l, false)
	expect(t, len(n.Node), 1)
	expect(t, n.HasRegex, false)

	n, p, l, err = trie.Get([]string{"root", "alpha2", "any"}, "")
	expect(t, err, nil)
	expectDeepEqual(t, p, []string{"any"})
	expect(t, l, false)
	expect(t, len(n.Node), 1)
	expect(t, n.HasRegex, true)

	n, p, l, err = trie.Get([]string{"root", "alpha", "beta"}, "")
	expect(t, err, nil)
	expectDeepEqual(t, p, []string{"beta"})
	expect(t, l, true)
	expect(t, len(n.Node), 1)
	expect(t, n.HasRegex, false)

	n, p, l, err = trie.Get([]string{"root", "alpha", "beta", "gamma"}, "")
	expect(t, err, nil)
	expectDeepEqual(t, p, []string{"gamma"})
	expect(t, l, true)
	expect(t, len(n.Node), 0)
	expect(t, n.HasRegex, false)

	n, p, l, err = trie.Get([]string{"root", "alphaA", "betaB", "gammaC"}, "")
	expect(t, err, nil)
	expectDeepEqual(t, p, []string{"alphaA", "betaB", "gammaC"})
	expect(t, l, false)
	expect(t, len(n.Node), 4)
	expect(t, n.HasRegex, false)

	n, p, l, err = trie.Get([]string{"root", "alpha", "betaB", "gammaC"}, "")
	expect(t, err, nil)
	expectDeepEqual(t, p, []string{"betaB", "gammaC"})
	expect(t, l, false)
	expect(t, len(n.Node), 2)
	expect(t, n.HasRegex, false)

	n, p, l, err = trie.Get([]string{"root", "alpha", "betaB", "gamma", "delta"}, "")
	expect(t, err, nil)
	expectDeepEqual(t, p, []string{"betaB", "gamma", "delta"})
	expect(t, l, false)
	expect(t, len(n.Node), 2)
	expect(t, n.HasRegex, false)
}
*/
