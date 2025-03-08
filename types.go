package main

import (
	"errors"
	"fmt"
	"strings"
)

var (
	invalidLine  = errors.New("Invalid line.")
	invalidInput = errors.New("Invalid input.")
	invalidKey   = errors.New("Invalid key.")
	invalidValue = errors.New("Invalid value.")
	invalidFile  = errors.New("Invalid file type, must be TOML.")
	dirNotExist  = errors.New("Directory does not exist.")
	fileNotExist = errors.New("File does not exist.")
)

type Session struct {
	Name          string
	Root          string
	DefaultWinInd int
	Windows       []*Window
}

func (s *Session) ToString() string {
	str := fmt.Sprintf("Session:\n\t- Name: %s\n\t- Root: %s\n\t- DefaultWinInd: %d\n\t- Windows:\n",
		s.Name, s.Root, s.DefaultWinInd)
	for _, n := range s.Windows {
		str += fmt.Sprintf("\t%s\n", n.ToString(2))
	}
	return str
}

type Window struct {
	Name    string
	Default bool
	Command string // default pane
	Panes   []*Pane
}

func (w *Window) ToString(count int) string {
	if count < 0 {
		count = 1
	}
	tabs := strings.Repeat("\t", count)
	str := fmt.Sprintf("Window:\n%s- Name: %s\n%s- Default: %t\n%s- Command: %s\n%s- Panes:\n",
		tabs, w.Name, tabs, w.Default, tabs, w.Command, tabs)
	for _, n := range w.Panes {
		str += fmt.Sprintf("%s\t%s\n", tabs, n.ToString())
	}
	return str
}

type Pane struct {
	Command string
}

func (p *Pane) ToString() string {
	return fmt.Sprintf("Pane - Command: %s\n", p.Command)
}
