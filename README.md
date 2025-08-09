# goilerplate - boilerplate for my go projects

This repo contains the boilerplate for when I am starting a new go project, including task files, linting configuration and go workflows amongst other goodies

Linting configuration can also be [found in the gist here](https://gist.github.com/ahobsonsayers/5a6baccee157e5d5c1ac4c1ccd163348)

## Setup

Once you have cloned goilerplate, replace all reference to goilerplate with you project name

```bash
read -p "Enter project name: " PROJECT_NAME
rm -rf .git
git init
find . -type f -exec sed -i 's|goileerplate|$PROJECT_NAME|g' {} \;
```

## Merging latest goilerplate changes

To merge the latest goilerplate changes, run:

```bash
git remote add -f upstream https://github.com/ahobsonsayers/goilerplate
git fetch upstream
git merge --no-commit --no-ff upstream/main
git reset README.md
```

This will not touch your README
