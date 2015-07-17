package main

import (
	"fmt"
	_ "reflect"
	"regexp"
)

type routes struct {
	dynamic map[string]regexp.Regexp
}

func main() {

	//	var routes Dynamic

	r := regexp.MustCompile("p([a-z]+)ch")

	route := routes{
		dynamic: map[string]regexp.Regexp{
			"test": *r,
		},
	}

	//	fmt.Println(set)
	fmt.Println(route.dynamic["test"])

	//fmt.Println(reflect.TypeOf(r))

	//fmt.Println(r.MatchString("peach"))
	//fmt.Println(r.MatchString("This is all Γςεεκ to me."))
	//fmt.Println(r.MatchString("This is all ⢓⢔⢕⢖⢗⢘⢙⢚⢛ to me."))
	//fmt.Println(r.MatchString("🌵 "))

}
