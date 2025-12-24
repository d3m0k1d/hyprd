# Hyprd

Event-driven daemon for Hyprland window manager. Listens to IPC events and executes user-defined rules via TOML configuration.

# Table of contents
1. [Overview](#overview)
2. [Requirements](#requirements)
3. [Installation](#installation)
4. [Usage](#usage)
5. [License](#license)

## Overview
Hyprd is a lightweight, event-driven daemon designed to work seamlessly with the Hyprland window manager.
It monitors Hyprland's IPC events in real time and executes custom rules defined in a user-configurable TOML file.
This allows users to automate window management tasks such as workspace assignments, window class handling, tiling behavior, and more—without modifying the core window manager.
## Requirements

- Go 1.25 or higher

## Installation

Check the realases if you need the latest version: https://github.com/d3m0k1d/hyprd/releases

also you can compile it from sources:
``` shell
git clone https://github.com/d3m0k1d/hyprd.git
cd hyprd
go mod tidy
go build -o hyprd ./cmd
```

## Usage
After first execution, Hyprd will generate a default configuration file at ~/.config/hyprd/config.toml. Edit this file to define your desired rules. The daemon should be started automatically with your session or manually via hyprd.
## License
BSD-3-Clause
