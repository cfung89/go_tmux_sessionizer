package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const configFile = ".tms.toml"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "kill" && insideTmux() {
		cmd := exec.Command("tmux", "kill-session")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		check(err)
	}

	helpPtr := flag.Bool("help", false, "Help")
	hPtr := flag.Bool("h", false, "Help")

	usageFile := "Sessionize according to the given configuration file path."
	configPtr := flag.String("file", "", usageFile)
	cPtr := flag.String("f", "", usageFile)

	usageCp := "Copy the given template file name from ~/.config/tms/templates/<filename>.toml to the given destination directory. The last two arguments are the arguments given to 'cp'"
	copyPtr := flag.Bool("copy", false, usageCp)
	cpPtr := flag.Bool("cp", false, usageCp)

	flag.Parse()
	assert((len(*configPtr) >= 0 && len(*cPtr) == 0) || (len(*configPtr) == 0 && len(*cPtr) >= 0), invalidArgument)

	var err error
	homeDir, err := os.UserHomeDir()
	check(err)
	if *hPtr || *helpPtr {
		flag.PrintDefaults()
		return
	} else if *copyPtr || *cpPtr {
		// copy file
		tail := flag.Args()
		assert(len(tail) >= 2, fmt.Errorf("%w: not enough arguments for copy (needs 2).", invalidArgument))
		from := tail[len(tail)-2]
		path, err := filepath.Abs(tail[len(tail)-1])
		check(err)
		exists, err := dirExists(path)
		assert(exists, fmt.Errorf("Last argument of copy must be a directory: %w", err))
		filename := fmt.Sprintf("%s/%s", path, configFile)
		cmd := exec.Command("cp", fmt.Sprintf("%s/.config/tms/templates/%s.toml", homeDir, from), filename)
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		check(err)
		parseAndCreate(path, filename)
		return
	} else if *configPtr != "" {
		parseAndCreate(filepath.Dir(*configPtr), *configPtr)
		return
	} else if *cPtr != "" {
		parseAndCreate(filepath.Dir(*cPtr), *cPtr)
		return
	}
	// main functionality
	var path string
	ignores, err := parseIgnoreFile(fmt.Sprintf("%s/.config/tms/tmsignore", homeDir))
	if err != nil {
		ignores = make([]string, 0)
	}
	if len(os.Args) == 1 {
		homeDir, err := os.UserHomeDir()
		check(err)
		path, err = fzf(homeDir, ignores)
		check(err)
	} else if bool, _ := dirExists(os.Args[1]); bool {
		path, err = fzf(os.Args[1], ignores)
		check(err)
	}

	if len(path) == 0 {
		return
	}
	filename := fmt.Sprintf("%s/%s", path, configFile)
	var name string
	if exists, _ := fileExists(filename); exists {
		parseAndCreate(path, filename)
	} else {
		name, err = simpleSessionizer(path)
		check(err)
		err = switchClient(name)
		check(err)
	}
}

func parseAndCreate(path string, filename string) {
	toml, err := parseToml(filename)
	check(err)
	name, err := createSessions(path, toml)
	check(err)
	err = switchClient(name)
	check(err)
}
