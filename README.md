# Go tmux-sessionizer

Tmux sessionizer written in Go.
Inspired by the [Primeagen's tmux-sessionizer](https://github.com/ThePrimeagen/tmux-sessionizer).

Dependencies: `fzf` and `tmux`

## Features

- Fuzzyfind through directories and create a tmux session in the chosen directory.
- Ignoring directories from the search space with a generic `tmsignore` file.
- Custom TOML parser.
- Parse TOML file and create tmux sessions, windows, panes based on the given configuration.
- Copy template configuration files into project directory.

## Install from package

Pre-built packages are found on the [Releases](https://github.com/cfung89/go_tmux_sessionizer/releases) page.

## Install from source

In order to build from source, run the following command.
```
git clone git@github.com:cfung89/go_tmux_sessionizer.git
cd go_tmux_sessionizer
go build .
```

## Add to system packages
```
sudo cp go_tmux_sessionizer /usr/local/bin/tms
```

To use it with a keybind, add to `.bashrc`:
```
bind '"\C-f":"tms\n"'
```

## Usage

Arguments are given as follows:
```
tms [ starting-point ] [ -h | -help ] ([ -f | -file ] { <config_file_path>} ) ( [ -cp | -copy ] { <template_name> <destination_dir> } )
```

`tms` searches the directory tree and opens the chosen directory through `fzf`. If no starting-point is specified, `$HOME` is assumed.
It also checks for a `.tms.toml` configuration file in the chosen directory and will create a tmux session based on that configuration file, if found.

`tms kill` is the equivalent short form of `tmux kill-session`.

### Ignore file

The ignore file is stored in `$HOME/.config/tms/tmsignore`.
It works similarly to `.gitignore` files, with certain minor differences:
- Trailing forward slashes (*'/'*) are ignored.
- Octothorpes (*'#'*) are used for comments.
- If the pattern contains no slashes, it acts as a global filter: it matches any directory whose name matches the pattern.
- If the pattern contains a slash (at the start or middle), it acts as a path filter: it matches only directories whose absolute path **contains** the pattern.
- Wildcards may be used.
- If a tilde character (*'~'*) is found at the beginning of the pattern, it is replaced with `$HOME`.

See an example at `./test/tmsignore`.

### Configuration file

Using the `-f` or `-file` flags with the path to a configuration file starts the session(s), window(s), and pane(s) outlined in the file.

<table>
  <thead>
    <tr>
      <td>Field</td>
      <td>Type</td>
      <td>Description</td>
    <tr>
  </thead>
  <tbody>
    <tr>
      <td><code>[[sessions]]</code></td>
      <td>Array</td>
      <td>Defines a new tmux session.</td>
    <tr>
    <tr>
      <td><code>name</code></td>
      <td>String</td>
      <td>(Optional) Name of the session. Defaults to the project directory name if omitted.</td>
    <tr>
    <tr>
      <td><code>root</code></td>
      <td>String</td>
      <td>(Optional) Root directory of the session. Defaults to current directory or project directory if omitted.</td>
    <tr>
    <tr>
      <td><code>[[sessions.windows]]</code></td>
      <td>Array</td>
      <td>Defines a new tmux window.</td>
    <tr>
    <tr>
      <td><code>name</code></td>
      <td>String</td>
      <td>(Optional) Name of the window. Defaults to the default tmux window name if omitted.</td>
    <tr>
    <tr>
      <td><code>default</code></td>
      <td>Boolean</td>
      <td>(Optional) Defines whether the window is the default first window pane of the session. Defaults to false if omitted.</td>
    <tr>
    <tr>
      <td><code>command</code></td>
      <td>String</td>
      <td>(Optional) Runs the given command in the default window pane.</td>
    <tr>
    <tr>
      <td><code>[[sessions.windows.panes]]</code></td>
      <td>Array</td>
      <td>Defines a new tmux pane.</td>
    <tr>
    <tr>
      <td><code>command</code></td>
      <td>String</td>
      <td>(Optional) Runs the given command in the new pane.</td>
    <tr>
    <tr>
      <td><code>orientation</code></td>
      <td>String</td>
      <td>(Optional) Splits the pane horizontally or vertically. Either "-h" or "-v". Defaults to "-h" if omitted.</td>
    <tr>
  </tbody>
</table>

### Template files

Template files are stored in `$HOME/.config/tms/templates`.
With the `-cp` or `-copy` flags, a template file is copied from the `<template_name>` to the `<destination_dir>` with these two arguments being the last two arguments in the command.
The template file is then parsed and the tmux session is started.
Template files can be named anything.

## License

This software is distributed under the MIT License. See LICENSE for details.
