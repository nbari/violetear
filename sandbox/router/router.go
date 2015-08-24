package main

import (
	"fmt"
)

type Route struct {
	handler string
	methods []string
}

type Router struct {
	routes map[string]map[string]Route
}

func main() {

	/**
	 *	Set hosts
	 *  (regex|host): vroot
	 */
	hosts := map[string]string{
		"*":      "vroot1",
		"*.test": "vroot1",
		".com":   "vroot1",
	}
	fmt.Println(hosts)
	fmt.Println(hosts[".com"])

	/**
	 * Routes
	 * vroot: routes{"regex", "handler", Methods}
	 */
	var r = Router{
		map[string]map[string]Route{
			"vroot1": {
				"r1": Route{"r2", []string{"post", "get"}},
			},
			"vroot2": {
				"r1": Route{"r2", []string{"post", "get"}},
			},
		},
	}

	fmt.Println(r, r.routes)
	fmt.Println(r.routes[hosts[".com"]])
	fmt.Println(r.routes[hosts[".com"]]["r1"])
	fmt.Println(r.routes[hosts[".com"]]["r1"].handler)
	fmt.Println(r.routes[hosts[".com"]]["r1"].methods)

}
