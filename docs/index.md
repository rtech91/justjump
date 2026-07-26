---
title: JustJump
---

# JustJump

A minimal, blazing-fast directory jumper for your shell. Effortlessly jump between frequently used directories with predefined suggestions and shell integration.

## Features

- Fast directory jumping with **Fuzzy Search & Guessing**
- **Global & Local Jump Roots** (`-G`) across your registered projects
- **Project Root Jump** (`-R`) to jump directly to the current project's global jump root
- **Git Workspace Discovery** (`-W`) for worktree navigation
- **Persistent Back Jump** (`-`) for instant toggle between locations across terminal sessions
- **Management Commands**: `add`, `remove`, and `verify` (with auto-clean support)
- Shell integration for Bash and Zsh
- Zero-dependency, blazing-fast performance

## Smart Navigation

JustJump is more than just a `cd` list. It understands your workflow:

- **Fuzzy Guessing**: Type `jj project` to jump directly to your `some-project` directory. If it's a unique match, you jump instantly.
- **Project Root Jump**: Use `jj -R` from any subfolder inside a registered global jump root to return instantly to that root.
- **Git Worktrees**: Use `jj -W` to discover and jump between different branches/worktrees of the same repository.
- **Persistent Back Jump**: Use `jj -` to toggle back to your previous location across any terminal window.
- **Smart Fallback**: Pressing **Escape** while viewing a filtered list resets to the full list of available jump points instead of exiting.

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

### Basic Navigation

- **Interactive Local Jump**:
  ```bash
  jj
  ```
- **Global Jump**:
  ```bash
  jj -G
  ```
- **Project Root Jump**:
  ```bash
  jj -R
  ```
- **Git Workspace Jump**:
  ```bash
  jj -W
  ```
- **Back Jump**:
  ```bash
  jj -
  ```
- **Fuzzy Search Jump**:
  ```bash
  jj -G myproject
  ```

### Managing Jump Roots & Points

- **Add Jump Root**:
  ```bash
  jj add         # Add local jump root
  jj add -G      # Add global jump root
  ```
- **Remove Jump Root**:
  ```bash
  jj remove      # Remove local jump root
  jj remove -G   # Remove global jump root
  ```
- **Verify & Clean Stale Roots**:
  ```bash
  jj verify           # Interactively verify local jump points
  jj verify -G        # Interactively verify global jump roots
  jj verify -G -c     # Clean non-existent global jump roots automatically
  ```

For more details, see the [README](https://github.com/rtech91/justjump).
