package main

import (
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

	// confSessionizer will be fixed
	// filename := fmt.Sprintf("%s/%s", path, configFile)
	// if exists, _ := fExists(filename); exists {
	// 	toml, err := parser(filename)
	// 	check(err)
	// 	err = confSessionizer(toml)
	// 	check(err)
	// }

	err = simpleSessionizer(path)
	check(err)
}
