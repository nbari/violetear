package main

import (
	"fmt"
)

type Hosts struct {
	host  string
	vroot string
}

type Route struct {
	regex   string
	handler string
	method  []string
}

type Router struct {
	host   string
	routes []Route
}

func main() {
	var route = Route{"regex", "handler", []string{"ALL"}}
	var r = Router{"default", []Route{route}}
	fmt.Println(r)
}
