package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Host struct {
	Host  string
	Vroot string
}

type Config struct {
	Hosts []Host
}

func main() {
	var conf Config
	file, err := ioutil.ReadFile("hosts.json")
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(file, &conf)
	if err != nil {
		fmt.Println("Error:", err)
	}

	//	fmt.Println(conf)
	//	fmt.Printf("%#v\n", conf)
	fmt.Printf("%+v\n", conf)

}
