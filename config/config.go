package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Versions []string
	Hosts    []Host
	Routes   Route
}

type Host struct {
	Host  string
	Vroot string
}

type Route map[string][]HandlerData

type HandlerData struct {
	URL     string
	Handler string
	Methods []string
}

func Get(file string) Config {

	yml_file, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	var data Config

	if err := yaml.Unmarshal(yml_file, &data); err != nil {
		panic(err)
	}

	return data
}
