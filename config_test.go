// Test
//
// go test config.go config_test.go

package violetear

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestGetConfig(t *testing.T) {

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

	routes := map[string][]Route{
		"default": {
			{
				"/test/.*", "test", []string{"GET", "POST", "PUT"},
			},
			{
				"/(md5|sha1|sha256|sha512)(/.*)?", "hash", []string{"GET"},
			},
			{
				"/.*", "default", nil,
			},
		},
		"ejemplo": {
			{
				"/.*", "default", []string{"ALL"},
			},
		},
	}

	for k, v := range routes {
		for i, v := range v {
			if !reflect.DeepEqual(data.Routes[k][i], v) {
				t.Error(data.Routes[k][i], v)
			}
		}
	}

	if !reflect.DeepEqual(data.Routes, routes) {
		t.Error(data.Routes, routes)
	}

}
