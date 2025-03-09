package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func assert(truth bool, msg error) {
	if !truth {
		fmt.Println(msg)
		os.Exit(1)
	}
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func dirExists(path string) (bool, error) {
	dir, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, dirNotExist
		}
		return false, err
	}
	return dir.IsDir(), nil
}

func fileExists(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, fileNotExist
		}
		return false, err
	}
	return f.Mode().IsRegular(), nil
}

func insideTmux() bool {
	// can also use os.Getenv("TERM_PROGRAM") == "tmux"
	if os.Getenv("TMUX") == "" {
		return false
	}
	return true
}

func getSName(root string) (string, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	r := strings.NewReplacer(" ", "", ".", "")
	base := r.Replace(filepath.Base(root))
	return base, nil
}
