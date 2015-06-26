package main

import (
	"fmt"
)

type Route struct {
	handler string
	method  []string
}

type Router struct {
	// vroute, regex
	routes map[string]map[string]Route
}

func main() {
	var hosts map[string]string
	hosts = make(map[string]string)
	hosts["teste1"] = "vroot1"
	hosts["teste2"] = "vroot2"
	fmt.Println(hosts)

	var r = Router{
		map[string]map[string]Route{
			"vroot1": map[string]Route{
				"regexA": Route{"handler", []string{"ALL"}},
				"regexB": Route{"handler", []string{"ALL"}},
				"regexC": Route{"handler", []string{"ALL"}},
			},
			"vroot2": map[string]Route{
				"regexA": Route{"handler", []string{"ALL"}},
				"regexB": Route{"handler", []string{"ALL"}},
				"regexC": Route{"handler", []string{"ALL"}},
			},
		}}
	fmt.Println(r)
}
