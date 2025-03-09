package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func simpleSessionizer(root string) (string, error) {
	base, err := getSName(root)
	if err != nil {
		return "", err
	}
	cmd := exec.Command("tmux", "new-session", "-d", "-c", root, "-s", base)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return base, err
}

func createSessions(path string, sessions []*Session) (string, error) {
	if sessions == nil {
		return "", errors.New("Sessions is nil.")
	}
	var attach string
	for i, session := range sessions {
		var err error
		if session.Root == "" {
			if path == "" {
				path, err = filepath.Abs(".")
			} else {
				path, err = filepath.Abs(path)
			}
			if err != nil {
				return "", err
			}
		} else {
			path = session.Root
		}

		name := session.Name
		if session.Name == "" {
			name, err = getSName(path)
			if err != nil {
				return "", err
			}
		}
		// Attach to the first session
		if i == 0 {
			attach = name
		}

		// Default window (with its panes)
		var cmd *exec.Cmd
		if session.Default.Name == "" {
			// no name for default window
			cmd = exec.Command("tmux", "new-session", "-d", "-c", path, "-s", name)
		} else {
			cmd = exec.Command("tmux", "new-session", "-d", "-c", path, "-s", name, "-n", session.Default.Name)
		}
		cmd.Stderr = os.Stderr
		if err = cmd.Run(); err != nil {
			return "", err
		}
		// Create panes for default window
		if err = createPanes(name, "1", path, session.Default.Panes); err != nil {
			return "", err
		}
		// Run command in default window, if given
		if session.Default.Command != "" {
			cmd = exec.Command("tmux", "send-keys", "-t", fmt.Sprintf("%s:%s.%d", name, "1", 1), session.Default.Command, "C-m")
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return "", err
			}
		}

		// Create remaining windows and their panes
		if err = createWindows(name, path, session.Windows); err != nil {
			return "", err
		}
	}
	return attach, nil
}

func createWindows(sName string, path string, windows []*Window) error {
	for i, n := range windows {
		add := exec.Command("tmux", "new-window", "-t", sName, "-n", n.Name, "-c", path)
		add.Stderr = os.Stderr
		if err := add.Run(); err != nil {
			return err
		}
		if err := createPanes(sName, fmt.Sprintf("%d", i+2), path, n.Panes); err != nil {
			return err
		}
		if n.Command != "" {
			cmd := exec.Command("tmux", "send-keys", "-t", fmt.Sprintf("%s:%d.1", sName, i+2), n.Command, "C-m")
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return err
			}
		}
	}
	return nil
}

func createPanes(sName string, wName string, path string, panes []*Pane) error {
	for i, n := range panes {
		if n.Orientation == "" {
			n.Orientation = "-h"
		}
		add := exec.Command("tmux", "split-window", "-t", fmt.Sprintf("%s:%s", sName, wName), n.Orientation, "-c", path)
		add.Stderr = os.Stderr
		if err := add.Run(); err != nil {
			return err
		}
		if n.Command != "" {
			cmd := exec.Command("tmux", "send-keys", "-t", fmt.Sprintf("%s:%s.%d", sName, wName, i+2), n.Command, "C-m")
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return err
			}
		}
	}
	return nil
}

func listSessions() ([]string, error) {
	out, err := exec.Command("tmux", "list-sessions").Output()
	sessions := strings.Split(string(out), "\n")
	return sessions, err
}

func switchClient(name string) error {
	if insideTmux() {
		cmd := exec.Command("tmux", "switch-client", "-t", name)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		return err
	}
	cmd := exec.Command("tmux", "attach", "-t", fmt.Sprintf("%s:1.1", name))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}
