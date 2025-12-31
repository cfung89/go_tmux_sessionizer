package main

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	sessions, err := parseToml("./test/.tms.toml")
	check(err)
	for _, n := range sessions {
		fmt.Println(n.ToString())
	}
}

func TestIgnoreFile(t *testing.T) {
	ignores, err := parseIgnoreFile("./test/tmsignore")
	check(err)
	for _, n := range ignores {
		fmt.Println(n)
	}
}
