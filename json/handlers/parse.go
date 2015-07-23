package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Something struct {
	Handlers Handler
}

type Handler struct {
	Default []HandlerData
	Extra   []HandlerData
}

type HandlerData struct {
	URL     string
	Handler string
	Methods []string
}

func main() {

	// var data map[string]interface{}
	data := Something{}

	file, err := ioutil.ReadFile("handlers.json")
	if err != nil {
		fmt.Println(err)
	}

	if err := json.Unmarshal(file, &data); err != nil {
		panic(err)
	}

	//fmt.Println(data)
	fmt.Printf("%+v\n", data)
}
