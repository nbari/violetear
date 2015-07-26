package main

import (
	"fmt"
	"path/filepath"
)

func main() {

	paths := []string{
		"/tmp/",
		"~/projects",
	}

	base := "/"

	for _, p := range paths {
		rel, err := filepath.Rel(base, p)
		fmt.Printf("%q: %q %v\n", p, rel, err)
	}

}
