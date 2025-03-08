package main

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	sessions, err := parser("test/example.toml")
	check(err)
	for _, n := range sessions {
		fmt.Println(n.ToString())
	}
}
