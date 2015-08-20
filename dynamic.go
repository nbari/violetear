package violetear

import (
	"log"
	"regexp"
	"strings"
)

type dynamicSet map[string]regexp.Regexp

func NewDynamicSet() dynamicSet {
	return make(dynamicSet)
}

func (d dynamicSet) Set(name string, regex string) error {
	if !strings.HasPrefix(name, ":") {
		log.Fatal("Dynamic route name must start with a colon ':'")
	}

	r := regexp.MustCompile(regex)
	d[name] = *r
	return nil
}
