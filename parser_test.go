package main

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	sessions, err := parser("test/.tms.toml")
	check(err)
	for _, n := range sessions {
		fmt.Println(n.ToString())
	}
}
