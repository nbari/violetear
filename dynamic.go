package violetear

import (
	"errors"
	"regexp"
	"strings"
)

type dynamicSet map[string]*regexp.Regexp

func (d dynamicSet) Set(name, regex string) error {
	if !strings.HasPrefix(name, ":") {
		return errors.New("dynamic route name must start with a colon ':'")
	}

	r := regexp.MustCompile(regex)
	d[name] = r

	return nil
}
