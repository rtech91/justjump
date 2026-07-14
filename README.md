# JustJump

JustJump is a simple tool to help you jump between directories quickly.

## Table of Contents

- [Requirements](#requirements)
- [Installation](#installation)
- [Shell integration](#shell-integration)
- [Usage](#usage)
- [Configuration](#configuration)
- [DDEV Integration](#ddev-integration)

## Requirements

JustJump requires the following dependencies to be installed on your system:
- [Go](https://golang.org/dl/)
- [Git](https://git-scm.com/downloads) (optional)
- [Make](https://www.gnu.org/software/make/)

## Installation

### DDEV Integration
If you are using DDEV, see [DDEV Integration instructions](./DDEV_INTEGRATION.md) for container setup.


### 1. Recommended: Install with the official script


For most users, run this command (no need to download the source):

```bash
curl -sSf https://rtech91.github.io/justjump/install.sh | bash
```

This will:
- Download and install JustJump to `~/.local/bin` (no sudo required)
- Set up shell integration for Bash or Zsh automatically

---

### 2. Advanced: Build from source (if you downloaded the source code)

If you have cloned or downloaded the source code, you can build and install manually:

```sh
make build-release
make install
```

This will build the binary and install it to `~/.local/bin`. You do not need sudo for this location.

To uninstall JustJump, run the following command:

```sh
make remove
```

## Shell integration

**`Important note: The shell integration is necessary to use JustJump in your shell properly.`**

To integrate JustJump with your shell, add the following line to your shell configuration file (e.g. `~/.bashrc`, `~/.zshrc`):

```sh
[ -f ~/.justjumprc ] && source ~/.justjumprc
```

You can automatically add this line to your shell config by running:

```sh
if [ -n "$ZSH_VERSION" ]; then
  echo '[ -f ~/.justjumprc ] && source ~/.justjumprc' >> ~/.zshrc
elif [ -n "$BASH_VERSION" ]; then
  echo '[ -f ~/.justjumprc ] && source ~/.justjumprc' >> ~/.bashrc
fi
```

Then copy the `.justjumprc` file to your home directory.
By default in the `misc/rc/` directory you can find the examples of the `.justjumprc` file for `bash` and `zsh` shells.

## Usage

To use JustJump, simply run `jj` in your terminal and select the directory you want to jump to.
```sh
jj
```

You can also pass arguments to `justjump` or `jj`.
- `--help`: Display help information
- `--global` or `-G`: Perform a global jump across registered projects
- `--workspaces` or `-W`: Discover nearby Git workspaces (worktrees)

### Smart Navigation & Discovery

JustJump includes several powerful features to speed up your navigation:

#### 1. Fuzzy Search & Guessing
You can pass a search pattern directly to `jj`. It uses **fuzzy matching** to find the best project.
- **Immediate Jump**: If only one project matches your pattern (e.g., `jj -G some-project`), it will jump there instantly without showing a menu.
- **Initial Filtering**: If multiple projects match, the menu will be pre-filtered to show only those options.
- **Shorthand**: You can use initials (e.g., `jj -G sp` for `some-project`).

#### 2. Workspace Discovery (`-W`)
Quickly jump between different **Git Worktrees** of the same project.
```sh
jj -W
```
This is perfect for ecosystems where you have multiple branches checked out into different folders. It uses the same fuzzy matching as the standard jump.

#### 3. Smart Fallback
If you are viewing a filtered list and realize you wanted something else, simply press **Escape**. Instead of exiting, JustJump will immediately switch to the **full list** of all available jump points.

#### 4. Back Jump (`-`)
Jump back to your previous successful destination across any terminal session.
```sh
jj -
```
Unlike `cd -`, this works across different terminal windows and persists even after you restart your computer.

### Add Command

The `add` command allows you to register a directory as a jump root.

To add a local jump root for the current directory:
```sh
jj add
```

To add a global jump root for the current directory:
```sh
jj add --global
```
or
```sh
jj add -G
```

### Remove Command

The `remove` command allows you to unregister a directory as a jump root.

To remove a local jump root for the current directory:
```sh
jj remove
```

To remove a global jump root for the current directory:
```sh
jj remove --global
```
or
```sh
jj remove -G
```

### Verify Command

The `verify` command checks whether registered jump roots (global) or jump points (local) still exist on disk. Non-existent entries are reported, and you can optionally clean them up.

To verify local folders:
```sh
jj verify
```

To verify global folders:
```sh
jj verify -G
```

When run in an interactive terminal, you will be **prompted to remove** each invalid entry individually.

To remove all invalid entries automatically without prompting, use the `--clean` / `-c` flag:
```sh
jj verify --clean
jj verify -G -c
```

## Configuration

JustJump uses a few configuration files to manage jump points and settings.

### Global Jump Root Registry

The global jump root registry is managed as a drop-in directory located at:
```sh
~/.config/justjump/jumproots.d/
```

Each global jump root is represented as a file with the `.jpath` extension. The file name corresponds to the jump root name, and the file content contains the absolute path to the directory.

To manually add a global jump root, create a file in the `jumproots.d` directory:
```sh
echo "/absolute/path/to/directory" > ~/.config/justjump/jumproots.d/yourproject.jpath
```

To remove a global jump root, delete the corresponding `.jpath` file:
```sh
rm ~/.config/justjump/jumproots.d/yourproject.jpath
```

### .justjump.yaml

This file contains predefined jump points that you can use to jump to your project directories quickly.

Example:
```yaml
jumppoints:
  - folder1/
  - subfolder/folder2/
  - folder3/
```

See and use example in the misc/config/local directory.
Don't forget to copy the `.justjump.yaml` file to your project directory and update the `jumppoints` list with your project directories.

#### Note:
**By default, the list of jump points will contain the project's root directory at the first position.**