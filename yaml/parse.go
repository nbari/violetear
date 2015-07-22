package main

import (
	"fmt"
	_ "gopkg.in/yaml.v2"
	"io/ioutil"
	_ "log"
)

type Version struct {
	Version string `json:"version" yaml:"version"`
}

type Versions []Version

type T struct {
	V []string
	H struct {
	}
}

func main() {
	data, err := ioutil.ReadFile("router.yml")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(data))

}
