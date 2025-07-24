---
title: JustJump
---

# JustJump

A minimal, blazing-fast directory jumper for your shell. Effortlessly jump between frequently used directories with smart auto-completion and shell integration.

## Features
- Fast directory jumping
- Shell integration for Bash and Zsh
- Simple installation

## Installation

### Prerequisites
- Go compiler installed
- `curl` installed
- Sudo/root access to write to `/usr/local/bin`

### Install via script


Run the following command in your terminal:

```bash
curl -sSf https://rtech91.github.io/justjump/install.sh | bash
```

Or, clone the repo and run the script manually:

```bash
git clone https://github.com/rtech91/justjump.git
cd justjump
bash install.sh
```

This will:
- Install the JustJump binary to `/usr/local/bin`
- Download the appropriate shell integration file to `~/.justjumprc`
- Add the integration line to your shell rc file (`.bashrc` or `.zshrc`)

## Usage

After installation, reload your shell:

```bash
source ~/.bashrc   # or source ~/.zshrc
```

Jump to a directory:

```bash
jj <directory>
```

For more details, see the [README](https://github.com/rtech91/justjump).
