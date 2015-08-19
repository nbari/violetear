package violetear

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type dynamic_set map[string]regexp.Regexp

func NewDynamic() dynamic_set {
	return make(dynamic_set)
}

func (d dynamic_set) Set(name string, regex string) error {
	if !strings.HasPrefix(name, ":") {
		fmt.Fprintf(os.Stderr, "Dynamic route name must start with colon ':'")
		os.Exit(1)
	}

	r := regexp.MustCompile(regex)
	d[name] = *r
	return nil
}
