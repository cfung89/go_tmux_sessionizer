package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

// TOML parser
func parser(filename string) ([]*Session, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var sessions []*Session
	var session *Session
	var window *Window

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || string(line[0]) == "#" {
			continue
		}
		var prev string
		switch line {
		case "[[sessions]]":
			session = &Session{}
			sessions = append(sessions, session)
			prev = "s"
		case "[[sessions.windows]]":
			window = &Window{}
			session.Windows = append(session.Windows, window)
			prev = "w"
		case "[[sessions.windows.panes]]":
			window.Panes = append(window.Panes, &Pane{})
			prev = "p"
		default:
			assert((string(line[0]) == "[" && string(line[len(line)-1]) == "]") ||
				(string(line[0]) != "[" && string(line[len(line)-1]) != "]"),
				invalidLine)
			assert((string(line[1]) == "[" && string(line[len(line)-2]) == "]") ||
				(string(line[1]) != "[" && string(line[len(line)-2]) != "]"),
				invalidLine)
			parts := strings.SplitN(line, "=", 3)
			assert(len(parts) == 2, invalidLine)
			parts[0] = strings.TrimSpace(parts[0])
			parts[1] = strings.TrimSpace(parts[1])
			switch prev {
			case "s":
				if parts[0] == "name" {
					assert(isString(parts[1]), fmt.Errorf("%w: value is not string.", invalidValue))
					session.Name = parts[1]
				} else if parts[1] == "root" {
					if exists, _ := dirExists(parts[1]); !exists {
						return nil, dirNotExist
					}
					assert(isString(parts[1]), fmt.Errorf("%w: value is not string.", invalidValue))
					session.Root = parts[1]
				} else {
					return nil, invalidKey
				}
			case "w":
				if parts[0] == "name" {
					assert(isString(parts[1]), fmt.Errorf("%w: value is not string.", invalidValue))
					window.Name = parts[1]
				} else if parts[1] == "command" {
					assert(isString(parts[1]), fmt.Errorf("%w: value is not string.", invalidValue))
					window.Command = parts[1]
				} else {
					return nil, invalidKey
				}
			case "p":
				if parts[1] == "command" {
					assert(isString(parts[1]), fmt.Errorf("%w: value is not string.", invalidValue))
					window.Command = parts[1]
				} else {
					return nil, invalidKey
				}
			default:
				return nil, errors.New("Internal error.")
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return sessions, nil
}

func isString(val string) bool {
	if string(val[0]) == "\"" && string(val[len(val)-1]) == "\"" {
		return true
	}
	return false
}
