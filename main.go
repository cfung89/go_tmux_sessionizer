package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

const configFile = ".tms.toml"

func main() {
	// flag for config file, and flag for projects config file, copy template files from .config
	helpPtr := flag.Bool("help", false, "Help")
	hPtr := flag.Bool("h", false, "Help")

	configPtr := flag.String("config", "", "Configuration file")
	cPtr := flag.String("c", "", "Configuration file")

	var copyPtr StringList
	var cpPtr StringList
	flag.Var(&copyPtr, "copy", "Copy the given template file from ~/.config/tms to the given destination")
	flag.Var(&cpPtr, "cp", "Copy the given template file from ~/.config/tms to the given destination")

	flag.Parse()
	assert(*configPtr == *cPtr, invalidArgument)
	assert((len(copyPtr) >= 0 && len(cpPtr) == 0) || (len(copyPtr) == 0 && len(cpPtr) >= 0), invalidArgument)

	if *hPtr || *helpPtr {
		flag.PrintDefaults()
	} else if len(copyPtr) == 2 {
		// copy file
		cmd := exec.Command("cp", copyPtr[0], copyPtr[1])
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		check(err)
		parseAndCreate(*cPtr)
	} else if len(cpPtr) == 2 {
		// copy file
		cmd := exec.Command("cp", cpPtr[0], cpPtr[1])
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		check(err)
		parseAndCreate(*cPtr)
	} else if *configPtr != "" {
		parseAndCreate(*cPtr)
	} else if *cPtr != "" {
		parseAndCreate(*cPtr)
	} else {
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

		if path != "" {
			filename := fmt.Sprintf("%s/%s", path, configFile)
			var name string
			if exists, _ := fileExists(filename); exists {
				parseAndCreate(filename)
			} else {
				name, err = simpleSessionizer(path)
				check(err)
				err = switchClient(name)
				check(err)
			}
		}
	}
}

func parseAndCreate(filename string) {
	toml, err := parser(filename)
	check(err)
	name, err := createSessions(toml)
	check(err)
	err = switchClient(name)
	check(err)
}
