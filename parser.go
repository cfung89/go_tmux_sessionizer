package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// TOML parser
func parser(filename string) ([]*Session, error) {
	if _, err := fileExists(filename); err != nil {
		return nil, err
	}
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var sessions []*Session
	var session *Session
	var window *Window

	scanner := bufio.NewScanner(f)
	var prev string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || string(line[0]) == "#" {
			continue
		}
		switch line {
		case "[[sessions]]":
			session = &Session{Default: nil}
			sessions = append(sessions, session)
			prev = "s"
		case "[[sessions.windows]]":
			window = &Window{}
			session.Windows = append(session.Windows, window)
			prev = "w"
		case "[[sessions.windows.panes]]":
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
				switch parts[0] {
				case "name":
					if !isString(parts[1]) {
						return nil, fmt.Errorf("%w: value is not string.", invalidValue)
					}
					session.Name = strings.Trim(parts[1], "\"")
				case "root":
					if !isString(parts[1]) {
						return nil, fmt.Errorf("%w: value is not string.", invalidValue)
					}
					parts[1] = strings.Trim(parts[1], "\"")
					if exists, _ := dirExists(parts[1]); !exists {
						return nil, dirNotExist
					}
					session.Root = parts[1]
				default:
					return nil, fmt.Errorf("%w: key of session is invalid", invalidKey)
				}
			case "w":
				switch parts[0] {
				case "name":
					if !isString(parts[1]) {
						return nil, fmt.Errorf("%w: value is not string.", invalidValue)
					}
					window.Name = strings.Trim(parts[1], "\"")
				case "command":
					if !isString(parts[1]) {
						return nil, fmt.Errorf("%w: value is not string.", invalidValue)
					}
					window.Command = strings.Trim(parts[1], "\"")
				case "default":
					val, err := parseBool(parts[1])
					if err != nil {
						return nil, fmt.Errorf("%w: value is not string.", invalidValue)
					}
					if val && session.Default != nil {
						return nil, fmt.Errorf("%w: More than default window in session.", invalidValue)
					}
					window.Default = val
					if val {
						session.Windows = session.Windows[:len(session.Windows)-1]
						session.Default = window
					}
				default:
					return nil, fmt.Errorf("%w: key of window is invalid", invalidKey)
				}
			case "p":
				switch parts[0] {
				case "command":
					if !isString(parts[1]) {
						return nil, fmt.Errorf("%w: value is not string.", invalidValue)
					}
					window.Panes = append(window.Panes, &Pane{Command: strings.Trim(parts[1], "\"")})
				case "orientation":
					if !isString(parts[1]) {
						return nil, fmt.Errorf("%w: value is not string.", invalidValue)
					}
					str := strings.Trim(parts[1], "\"")
					if str == "-h" || str == "-v" {
						window.Panes = append(window.Panes, &Pane{Orientation: str})
					}
				default:
					return nil, fmt.Errorf("%w: key of pane is invalid", invalidKey)
				}
			default:
				return nil, internalErr
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return sessions, nil
}

func isString(val string) bool {
	if (string(val[0]) == "\"" && string(val[len(val)-1]) == "\"") ||
		(string(val[0]) == "'" && string(val[len(val)-1]) == "'") {
		return true
	}
	return false
}

func parseBool(val string) (bool, error) {
	if val == "true" {
		return true, nil
	} else if val == "false" {
		return false, nil
	}
	return false, invalidValue
}
