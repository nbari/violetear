package main

import (
	"fmt"
	"os/user"
	"strings"
)

func main() {

	path := "~/projects"
	if path[:2] == "~/" {
		usr, _ := user.Current()
		dir := usr.HomeDir
		path = strings.Replace(path, "~", dir, 1)
	}

	fmt.Println(path)
}
