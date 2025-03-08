package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func simpleSessionizer(root string) error {
	cmd := exec.Command("tmux new-session -c", root)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func confSessionizer(sessions []*Session) error {
	if sessions == nil {
		return errors.New("Sessions is nil.")
	}
	for i, session := range sessions {
		var name string
		if session.Name == "" {
			name = fmt.Sprintf("-s %d", i)
		} else {
			name = fmt.Sprintf("-s %s", session.Name)
		}

		var root string
		if session.Root != "" {
			root = fmt.Sprintf("-c %s", session.Root)
		}

		var defaultWinName string
		var defaultWinCmd string
		if session.DefaultWinInd != -1 {
			if session.Windows[session.DefaultWinInd].Name == "" {
				defaultWinName = "0"
			} else {
				defaultWinName = fmt.Sprintf("-n %s", session.Windows[session.DefaultWinInd].Name)
			}
			defaultWinCmd = session.Windows[session.DefaultWinInd].Command
			// add panes for default window
		}
		cmd := exec.Command("tmux new-session -d", root, name, defaultWinName)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return err
		}

		// fix name, since the flag is included
		defWinCmd := exec.Command("tmux send-keys -t", name, defaultWinCmd) // Fix
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

func listSessions() error {
	return nil
}
