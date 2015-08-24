package path

import (
	"fmt"
	"regexp"
	"strings"
)

func A(path string) []string {
	rex := regexp.MustCompile(`[^/ ]+`)
	out := rex.FindAllString(path, -1)
	fmt.Println(out)
	return out
}

func B(path string) []string {
	parts := strings.Split(path, "/")
	for i, v := range parts {
		if len(v) == 0 {
			parts = append(parts[:i], parts[i+1:]...)
		}
	}
	fmt.Println(parts)
	return parts
}
