package violetear

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type (
	dynRoutes []dynRoute

	dynRoute struct {
		name string
		rx   *regexp.Regexp
	}
)

func (d *dynRoutes) Set(name, regex string) error {
	if !strings.HasPrefix(name, ":") {
		return errors.New("Dynamic route name must start with a colon ':'")
	}

	// fix regex
	if !strings.HasPrefix(regex, "^") {
		regex = fmt.Sprintf("^%s$", regex)
	}

	r := regexp.MustCompile(regex)
	if d.Get(name) == nil {
		*d = append(*d, dynRoute{name, r})
	}
	return nil
}

func (d *dynRoutes) Get(name string) *regexp.Regexp {
	for _, r := range *d {
		if r.name == name {
			return r.rx
		}
	}
	return nil
}
