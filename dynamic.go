package violetear

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type dynamicRoute struct {
	name string
	rx   *regexp.Regexp
}

func (v *Router) dynamicRoutesSet(name, regex string) error {
	if !strings.HasPrefix(name, ":") {
		return errors.New("Dynamic route name must start with a colon ':'")
	}

	// fix regex
	if !strings.HasPrefix(regex, "^") {
		regex = fmt.Sprintf("^%s$", regex)
	}

	r := regexp.MustCompile(regex)
	if v.dynamicRoutesGet(name) == nil {
		v.dynamicRoutes = append(v.dynamicRoutes, dynamicRoute{name, r})
	}
	return nil
}

func (v *Router) dynamicRoutesGet(name string) *regexp.Regexp {
	for _, r := range v.dynamicRoutes {
		if r.name == name {
			return r.rx
		}
	}
	return nil
}
