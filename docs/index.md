---
title: JustJump
---

# JustJump

A minimal, blazing-fast directory jumper for your shell. Effortlessly jump between frequently used directories with directory suggestions and shell integration.

## Features
- Fast directory jumping
- Shell integration for Bash and Zsh
- Simple installation

## Installation

### Prerequisites
- Go compiler (required)
- `curl` (required)

### Install via script


Run the following command in your terminal:

```bash
curl -sSf https://rtech91.github.io/justjump/install.sh | bash
```

This will:
- Install the JustJump binary to `~/.local/bin` folder
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
