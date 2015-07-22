package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type mytype []map[string]string

func main() {
	var data mytype
	file, err := ioutil.ReadFile("router.json")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(string(file))

	err = json.Unmarshal(file, &data)
	if err != nil {
		fmt.Println(data)
	}

	//	fmt.Println(data)

}
