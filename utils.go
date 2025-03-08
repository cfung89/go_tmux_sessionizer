package main

import (
	"fmt"
	"os"
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

func fExists(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, fileNotExist
		}
		return false, err
	}
	return f.Mode().IsRegular(), nil
}
