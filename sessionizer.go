package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func simpleSessionizer(root string) error {
	// can also use os.Getenv("TERM_PROGRAM") == "tmux"
	root, err := filepath.Abs(root)
	base := strings.ReplaceAll(filepath.Base(root), " ", "")
	if err != nil {
		return err
	}
	if os.Getenv("TMUX") == "" {
		cmd := exec.Command("tmux", "new-session", "-c", root, "-s", base)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		return err
	}
	cmd := exec.Command("tmux", "new-session", "-d", "-c", root, "-s", base)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("tmux", "switch-client", "-t", base)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return err
}

func confSessionizer(sessions []*Session) error {
	if sessions == nil {
		return errors.New("Sessions is nil.")
	}
	for i, session := range sessions {
		name := fmt.Sprintf("%d", i)
		if session.Name != "" {
			name = session.Name
		}

		root := "."
		if session.Root != "" {
			root = session.Root
		}

		var cmd *exec.Cmd
		if i == 0 {
			if (session.Default == nil && session.Windows[0].Default) ||
				(session.Default != nil && !session.Windows[0].Default) {
				return internalErr
			}
			if session.Default.Name == "" {
				cmd = exec.Command("tmux", "new-session", "-d", "-c", root, "-s", name)
			} else {
				cmd = exec.Command("tmux", "new-session", "-d", "-c", root, "-s", name, "-n", session.Default.Name)
			}
		}
		// add panes for default window
		cmd = exec.Command("tmux", "new-session", "-d", "-c", root, "-s", name)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return err
		}

		// fix name, since the flag is included
		defWinCmd := exec.Command("tmux send-keys -t", name) // Fix
		defWinCmd.Stdin = os.Stdin
		defWinCmd.Stdout = os.Stdout
		defWinCmd.Stderr = os.Stderr
		err = defWinCmd.Run()
		if err != nil {
			return err
		}
		// for _, window := range session.Windows {

		// }
	}
	return nil
}

func listSessions() ([]string, error) {
	out, err := exec.Command("tmux", "list-sessions").Output()

	sessions := strings.Split(string(out), "\n")
	return sessions, err
}
