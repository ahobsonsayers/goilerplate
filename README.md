# goilerplate - boilerplate for my go projects

This repo contains the boilerplate for when I am starting a new go project, including task files, linting configuration and go workflows amongst other goodies

Linting configuration can also be [found in the gist here](https://gist.github.com/ahobsonsayers/5a6baccee157e5d5c1ac4c1ccd163348)

## First Steps

Once you have cloned `goilerplate`, replace all references to `goilerplate` with you project name:

> [!NOTE]
> You will need [fd](https://github.com/sharkdp/fd) and [sd](https://github.com/chmln/sd) installed

```bash
read -p "Enter project name: " PROJECT_NAME
fd -t f --hidden -E .git -x sd "goilerplate" "$PROJECT_NAME"
```

Remember to also delete your `.git` folder.

## Update project with latest goilerplate changes

Ask your coding agent:

> Add https://github.com/ahobsonsayers/goilerplate as a git remote named `boilerplate` (skip if it already exists), fetch it, then squash merge `boilerplate/main` into this repository with `--allow-unrelated-histories` and commit the result. Keep my project's own changes (project names, renamed files) and adopt any new boilerplate updates.
