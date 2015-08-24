package violetear

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os/user"
	"strings"
)

type Config struct {
	Versions []string
	Hosts    []Host
	Routes   map[string][]Route
}

type Host struct {
	Host, Vroot string
}

type Route struct {
	URL, Handler string
	Methods      []string
}

func GetConfig(file string) Config {

	if file[:2] == "~/" {
		usr, _ := user.Current()
		dir := usr.HomeDir
		file = strings.Replace(file, "~", dir, 1)
	}

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
