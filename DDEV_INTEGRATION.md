# DDEV Integration for JustJump

To make JustJump available inside your DDEV container, add the following to your project's `config.yaml`:

```yaml
webimage_extra_packages: [golang]
hooks:
  post-start:
    - exec: curl -sSf https://rtech91.github.io/justjump/install.sh | bash
    - exec: echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
```

## Steps
1. Open your DDEV project's `.ddev/config.yaml` file.
2. Add or merge the above configuration under the appropriate sections.
3. Run `ddev start` or `ddev restart` to apply the changes and install JustJump inside the container.

## Notes
- This will install Go and JustJump in the container and update the PATH for Bash users.
- For Zsh users, you may also want to update `~/.zshrc` similarly.
- For more information about DDEV, see the [DDEV Documentation](https://ddev.readthedocs.io/en/latest/).

---

This integration ensures JustJump is available for use in your DDEV development environment.
