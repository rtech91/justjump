---
title: JustJump
---

# JustJump

A minimal, blazing-fast directory jumper for your shell. Effortlessly jump between frequently used directories with predefined suggestions and shell integration.

## Features

- Fast directory jumping with **Fuzzy Search**
- **Git Workspace Discovery** (`-W`) for worktree navigation
- **Back Jump** (`-`) for instant toggle between locations
- Shell integration for Bash and Zsh
- Zero-dependency, blazing-fast performance

## Smart Navigation

JustJump is more than just a `cd` list. It understands your workflow:

- **Fuzzy Guessing**: Type `jj project` to jump directly to your `some-project` directory. If it's a unique match, you jump instantly.
- **Git Worktrees**: Use `jj -W` to discover and jump between different branches/worktrees of the same repository.
- **Persistent Back Jump**: Use `jj -` to toggle back to your previous location across any terminal window.

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
- Add the integration line to your shell configuration file (`.bashrc` or `.zshrc`).

## Usage

After installation, reload your shell:

```bash
source ~/.bashrc   # or source ~/.zshrc
```

To jump to a directory:

```bash
jj
```

or

```bash
jj -G
```

For more details, see the [README](https://github.com/rtech91/justjump).
