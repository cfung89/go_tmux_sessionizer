package main

import (
	"fmt"
	"os"
)

const configFile = ".tms.toml"

func main() {
	var path string
	var err error
	if len(os.Args) == 1 {
		path, err = fzf(".")
		check(err)
	} else {
		if bool, _ := dirExists(os.Args[1]); bool {
			path, err = fzf(os.Args[1])
			check(err)
		}
	}

	if exists, _ := fExists(fmt.Sprintf("%s/%s", path, configFile)); exists {
		filename := "test/example.toml"
		toml, err := parser(filename)
		check(err)
		err = confSessionizer(toml)
		check(err)
	}

	err = simpleSessionizer(path)
	check(err)
}
