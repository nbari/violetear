// http://play.golang.org/p/1-xLIRwWoQ
package main

import (
	"encoding/json"
	"fmt"
)

var json_data = []byte(`
{
    "handlers": {
        "default": [{
            "url": "/test/.*",
            "handler": "my_handler",
            "methods": [
                "GET",
                "POST",
                "PUT"
            ]
        }, {
            "url": "/(md5|sha1|sha256|sha512)(/.*)?",
            "handler": "hash_handler",
            "methods": [
                "GET"
            ]
        }, {
            "url": "/.*",
            "handler": "default"
        }],
        "extra": [{
            "url": "/.*",
            "handler": "other_handler",
            "methods": [
                "ALL"
            ]
        }]
    }
}
`)

func main() {

	var data map[string]interface{}

	if err := json.Unmarshal(json_data, &data); err != nil {
		panic(err)
	}

	fmt.Println(data)
}
