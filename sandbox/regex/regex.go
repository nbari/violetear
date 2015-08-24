package main

import (
	"fmt"
	"regexp"
)

type routes struct {
	dynamic map[string]regexp.Regexp
}

func main() {

	dynamic_route := make(map[string]string)

	dynamic_route["r_1"] = "p([a-z]+)ch"
	dynamic_route["r_2"] = "p([a-z0-9]+)ch"
	dynamic_route["r_3"] = "/simple"

	dynamic_set := make(map[string]regexp.Regexp)

	for k, v := range dynamic_route {
		r := regexp.MustCompile(v)
		dynamic_set[k] = *r
	}

	strings := []string{"peach", "peach2", "p3ch", "^/simple$"}

	for _, s := range strings {
		for k, v := range dynamic_set {
			if v.MatchString(s) {
				fmt.Printf("Match %s -> %v [%v]\n", dynamic_route[k], s, k)
			}
		}
	}

	//fmt.Println(reflect.TypeOf(r))

	//fmt.Println(r.MatchString("peach"))
	//fmt.Println(r.MatchString("This is all Γςεεκ to me."))
	//fmt.Println(r.MatchString("This is all ⢓⢔⢕⢖⢗⢘⢙⢚⢛ to me."))
	//fmt.Println(r.MatchString("🌵 "))

}
