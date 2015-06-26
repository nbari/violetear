package main

import (
	"fmt"
)

type Hosts struct {
	host, vroot string
}

type Route struct {
	regex, handler string
}

type Methods struct {
	methods map[string]struct{}
}

type Router struct {
	routes map[Hosts]map[Route]Methods
}

func main() {

	var r = Router{
		map[Hosts]map[Route]Methods{
			Hosts{"*", "vroot"}: {
				Route{"r2", "h"}: Methods{map[string]struct{}{
					"post": {},
					"get":  {},
				}},
				Route{"r3", "h"}: Methods{map[string]struct{}{
					"post": {},
					"get":  {},
				}},
			},
			Hosts{"*", "vroot2"}: {
				Route{"r2", "h"}: Methods{map[string]struct{}{
					"PUT":  {},
					"GET":  {},
					"POST": {},
				}},
			},
		},
	}

	fmt.Println(r)
}
