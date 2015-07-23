package main

import (
	"encoding/json"
	"fmt"
)

type Host struct {
	Host  string
	Vroot string
}

type Config struct {
	Hosts []Host
}

var in = `{
    "hosts": [{
        "host": "*",
        "vroot": "default"
    }, {
        "host": "*.zunzun.io",
        "vroot": "default"
    }, {
        "host": "ejemplo.org",
        "vroot": "ejemplo"
    }, {
        "host": "api.ejemplo.org",
        "vroot": "ejemplo"
    }]
}`

func main() {
	var conf Config
	err := json.Unmarshal([]byte(in), &conf)

	if err != nil {
		fmt.Print("Error:", err)
	}

	fmt.Printf("%#v\n", conf)
}
