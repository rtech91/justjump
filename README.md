# JustJump

JustJump is a simple tool to help you jump between directories quickly.

## Table of Contents

- [Requirements](#requirements)
- [Installation](#installation)
- [Shell integration](#shell-integration)
- [Usage](#usage)
- [Configuration](#configuration)
- [Contributing](#contributing)
- [License](#license)

## Requirements

JustJump requires the following dependencies to be installed on your system:
- [Go](https://golang.org/dl/)
- [Git](https://git-scm.com/downloads) (optional)
- [Make](https://www.gnu.org/software/make/)

## Installation

To install JustJump, you can use the provided Makefile. Run the following commands:

```sh
make build-release
make install
```
This will build the binary and install it to `/usr/local/bin/` directory.
Installation may require root privileges, so you may need to use `sudo` if you encounter permission issues.

**Alternatively, you can install JustJump directly using Go:**

```sh
GOBIN=/usr/local/bin go install github.com/rtech91/justjump@latest
```

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

The `verify` command checks if the local or global folders exist.

To verify local folders, run:
```sh
jj verify
```

To verify global folders, run:
```sh
jj verify --global
```

or 
```sh
jj verify -G
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