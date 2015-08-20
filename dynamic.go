package violetear

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type dynamicSet map[string]regexp.Regexp

func NewDynamicSet() dynamicSet {
	return make(dynamicSet)
}

func (d dynamicSet) Set(name string, regex string) error {
	if !strings.HasPrefix(name, ":") {
		fmt.Fprintf(os.Stderr, "Dynamic route name must start with a colon ':'\n")
		os.Exit(1)
	}

	r := regexp.MustCompile(regex)
	d[name] = *r
	return nil
}
