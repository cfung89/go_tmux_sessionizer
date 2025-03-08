package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func main() {
	out, err := search()
	check(err)
	fmt.Println(out)

	// filename := "test/example.toml"
	// assert(filename[len(filename)-5:] == ".toml", invalidFile)
	// toml, err := parser(filename)
	// check(err)
}

func search() (string, error) {
	fzfCmd := exec.Command("fzf")
	findCmd := exec.Command("find", "/home", "-type", "d")

	// show errors if fails
	findCmd.Stderr = os.Stderr
	fzfCmd.Stderr = os.Stderr

	// Piping
	r, w := io.Pipe()
	findCmd.Stdout = w
	fzfCmd.Stdin = r
	// defer w.Close()
	defer r.Close()

	// Fzf stdout
	var res strings.Builder
	fzfCmd.Stdout = &res

	// Can also use goroutine and Run, instead of Start and Wait for async use
	if err := findCmd.Start(); err != nil {
		return "", err
	}
	if err := fzfCmd.Start(); err != nil {
		return "", err
	}

	if err := findCmd.Wait(); err != nil {
		return "", err
	}
	w.Close()
	if err := fzfCmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 130 {
			os.Exit(0) // Cancelled
		}
		return "", err
	}
	return res.String(), nil
}
