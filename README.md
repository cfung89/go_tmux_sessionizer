# Go Tmux sessionizer

Tmux sessionizer written in Go.
Inspired by the [Primeagen's tmux-sessionizer](https://github.com/ThePrimeagen/tmux-sessionizer).

Dependencies: `fzf` and `tmux`

## Usage

Install:
```
git clone git@github.com:cfung89/go_tmux_sessionizer.git
cd go_tmux_sessionizer
go build .
```

Add to system packages:
```
sudo cp go_tmux_sessionizer /usr/local/bin/tms
```

To use it with a keybind, add to `.bashrc`:
```
bind '"\C-f":"tms\n"'
```

## TODO

- [ ] Finish sessionizer for configuration files.
