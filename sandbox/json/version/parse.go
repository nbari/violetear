package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	Versions []string
}

func main() {
	var conf Config
	file, err := ioutil.ReadFile("versions.json")
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(file, &conf)
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println(conf)

}
