package main

import (
	"errors"
)

var (
	invalidLine  = errors.New("Invalid line.")
	invalidInput = errors.New("Invalid input.")
	invalidKey   = errors.New("Invalid key.")
	invalidValue = errors.New("Invalid value.")
	invalidFile  = errors.New("Invalid file type, must be TOML.")
	dirNotExist  = errors.New("Directory does not exist.")
)

type Session struct {
	Name    string
	Root    string
	Windows []*Window
}

type Window struct {
	Name    string
	Command string // default pane
	Panes   []*Pane
}

type Pane struct {
	Command string
}
