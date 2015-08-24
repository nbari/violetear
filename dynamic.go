package violetear

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

type dynamicSet map[string]regexp.Regexp

func (d dynamicSet) Set(name string, regex string) error {
	if !strings.HasPrefix(name, ":") {
		log.Fatal("Dynamic route name must start with a colon ':'")
	}

	// fix regex
	if !strings.HasPrefix(regex, "^") {
		regex = fmt.Sprintf("^%s$", regex)
	}

	r := regexp.MustCompile(regex)
	d[name] = *r
	return nil
}
