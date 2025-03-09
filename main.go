package main

import (
	"fmt"
	"os"
)

const configFile = ".tms.toml"

func main() {
	// flag for config file, and flag for projects config file, copy template files from .config

	// main
	var path string
	var err error
	if len(os.Args) == 1 {
		homeDir, err := os.UserHomeDir()
		check(err)
		path, err = fzf(homeDir)
		check(err)
	} else {
		if bool, _ := dirExists(os.Args[1]); bool {
			path, err = fzf(os.Args[1])
			check(err)
		}
	}

	filename := fmt.Sprintf("%s/%s", path, configFile)
	if exists, _ := fExists(filename); exists {
		toml, err := parser(filename)
		check(err)
		err = createSessions(toml)
		check(err)
	} else {
		name, err := simpleSessionizer(path)
		check(err)
		err = switchClient(name)
		check(err)
	}
}
