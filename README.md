# goilerplate - boilerplate for my go projects

This repo contains the boilerplate for when I am starting a new go project, including task files, linting configuration and go workflows amongst other goodies

Linting configuration can also be [found in the gist here](https://gist.github.com/ahobsonsayers/5a6baccee157e5d5c1ac4c1ccd163348)

## Setup

Once you have cloned `goilerplate`, replace all references to `goilerplate` with you project name:

> [!NOTE]
> You will need [fd](https://github.com/sharkdp/fd) and [sd](https://github.com/chmln/sd) installed

```bash
read -p "Enter project name: " PROJECT_NAME
fd -t f -x sd "goilerplate" "$PROJECT_NAME"
```

Remember to also delete your `.git` folder.

## Merging latest goilerplate changes

To merge the latest goilerplate changes, run:

```bash
git fetch upstream
git merge -X ours --no-ff upstream/main
```
