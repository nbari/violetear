package main

import (
	"fmt"
	"regexp"
)

type route struct {
	pattern *regexp.Regexp
	handler string
}

type RegexHandler struct {
	routes []*route
}

func main() {

	h := new(RegexHandler)

	my_route := route{regexp.MustCompile("p([a-z0-9]+)ch"), "/test"}
	h.routes = append(h.routes, &my_route)

	my_route2 := route{regexp.MustCompile("p([a-z]+)ch"), "/test2"}
	h.routes = append(h.routes, &my_route2)

	fmt.Println(my_route.handler)

	for k, v := range h.routes {
		fmt.Println(k, v)
	}

	//fmt.Println(reflect.TypeOf(r))

	//fmt.Println(r.MatchString("peach"))
	//fmt.Println(r.MatchString("This is all Γςεεκ to me."))
	//fmt.Println(r.MatchString("This is all ⢓⢔⢕⢖⢗⢘⢙⢚⢛ to me."))
	//fmt.Println(r.MatchString("🌵 "))

}
