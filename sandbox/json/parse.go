package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	Versions []string
	Hosts    []Host
	Handlers Handler
}

type Host struct {
	Host  string
	Vroot string
}

type Handler map[string][]HandlerData

type HandlerData struct {
	URL     string
	Handler string
	Methods []string
}

func main() {

	// var data map[string]interface{}
	data := Config{}

	file, err := ioutil.ReadFile("router.json")
	if err != nil {
		fmt.Println(err)
	}

	if err := json.Unmarshal(file, &data); err != nil {
		panic(err)
	}

	//fmt.Println(data)
	fmt.Printf("%+v\n", data)
}
