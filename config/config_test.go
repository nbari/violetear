package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"testing"
)

func TestGet(t *testing.T) {

	yml_file, err := ioutil.ReadFile("config_test.yml")

	if err != nil {
		panic(err)
	}

	var data Config

	if err := yaml.Unmarshal(yml_file, &data); err != nil {
		panic(err)
	}

	/**
	 * output
	 *
	{[v0 v1 v2] [{* default} {*.zunzun.io default} {ejemplo.org ejemplo} {api.ejemplo.org ejemplo}] map[default:[{/test/.* test [GET POST PUT]} {/(md5|sha1|sha256|sha512)(/.*)? hash [GET]} {/.* default []}] ejemplo:[{/.* default [ALL]}]]}"
	*/
	versions := []string{"v0", "v1", "v2"}
	for i, v := range data.Versions {
		if v != versions[i] {
			t.Error(data.Versions, versions)
		}
	}

	hosts := []Host{
		{"*", "default"},
		{"*.zunzun.io", "default"},
		{"ejemplo.org", "ejemplo"},
		{"api.ejemplo.org", "ejemplo"},
	}

	for i, v := range data.Hosts {
		if v != hosts[i] {
			t.Error(i, v, hosts[i])
		}
	}

}
