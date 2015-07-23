package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

<<<<<<< HEAD
type H struct {
	host, vroot string
}

type T struct {
	//	Versions []string
	//	Hosts    []map[string]string
	Handlers []map[string]map[string]string
}

// map[versions:[v0 v1 v2] hosts:[map[host:* vroot:default] map[host:*.zunzun.io vroot:default] map[host:ejemplo.org vroot:ejemplo] map[host:api.ejemplo.org vroot:ejemplo]] handlers:map[default:[map[url:/test/.* handler:my_handler methods:[GET POST PUT]] map[url:/(md5|sha1|sha256|sha512)(/.*)? handler:hash_handler methods:[GET]] map[url:/.* handler:default]] ejemplo:[map[url:/.* handler:other_handler methods:[ALL]]]]]

func main() {
	//	var data map[string]interface{}
	var data T
	file, err := ioutil.ReadFile("router.json")
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(file, &data)
	if err != nil {
		fmt.Println(data)
	}

	fmt.Println(data)

}
