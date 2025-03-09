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
	// can also use os.Getenv("TERM_PROGRAM") == "tmux"
	root, err := filepath.Abs(root)
	base := strings.ReplaceAll(filepath.Base(root), " ", "")
	if err != nil {
		return "", err
	}
	if os.Getenv("TMUX") == "" {
		cmd := exec.Command("tmux", "new-session", "-d", "-c", root, "-s", base)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		return base, err
	}
	cmd := exec.Command("tmux", "new-session", "-d", "-c", root, "-s", base)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return "", err
	}
	return base, err
}

func createSessions(sessions []*Session) error {
	if sessions == nil {
		return errors.New("Sessions is nil.")
	}
	for _, session := range sessions {
		var err error
		root := session.Root
		if session.Root == "" {
			root, err = filepath.Abs(".")
			if err != nil {
				return err
			}
		}

		name := session.Name
		if session.Name == "" {
			name = strings.ReplaceAll(filepath.Base(root), " ", "")
		}

		// Default window (with its panes)
		var cmd *exec.Cmd
		if session.Default.Name == "" {
			// no name for default window
			cmd = exec.Command("tmux", "new-session", "-d", "-c", root, "-s", name)
		} else {
			cmd = exec.Command("tmux", "new-session", "-d", "-c", root, "-s", name, "-n", session.Default.Name)
		}
		cmd.Stderr = os.Stderr
		if err = cmd.Run(); err != nil {
			return err
		}
		// Create panes for default window
		if err = createPanes(name, "1", session.Default.Panes); err != nil {
			return err
		}
		// Run command in default window, if given
		if session.Default.Command != "" {
			cmd = exec.Command("tmux", "send-keys", "-t", fmt.Sprintf("%s:%s.%d", name, "1", 1), session.Default.Command, "C-m")
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return err
			}
		}

		// Create remaining windows and their panes
		if err = createWindows(name, session.Windows); err != nil {
			return err
		}
	}
	return nil
}

func createWindows(sName string, windows []*Window) error {
	for i, n := range windows {
		add := exec.Command("tmux", "new-window", "-t", sName)
		add.Stderr = os.Stderr
		if err := add.Run(); err != nil {
			return err
		}
		if err := createPanes(sName, fmt.Sprintf("%d", i+2), n.Panes); err != nil {
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

func createPanes(sName string, wName string, panes []*Pane) error {
	for i, n := range panes {
		if n.Orientation == "" {
			n.Orientation = "-h"
		}
		add := exec.Command("tmux", "split-window", "-t", fmt.Sprintf("%s:%s", sName, wName), n.Orientation)
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
	cmd := exec.Command("tmux", "switch-client", "-t", name)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}
