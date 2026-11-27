# goilerplate - boilerplate for my go projects

This repo contains the boilerplate for when I am starting a new go project, including task files, linting configuration and go workflows amongst other goodies

Linting configuration can also be [found in the gist here](https://gist.github.com/ahobsonsayers/5a6baccee157e5d5c1ac4c1ccd163348)

## Commands

- `greet` - a simple example cli command you can replace
- `serve` - run the API server (see [API server](#api-server))
- `worker` - example worker (see [Example worker](#example-worker)) - delete if your project does not need a worker

## API server

The `server` package ([server/serve.go](server/serve.go)) starts a chi router with request-time OpenAPI validation. Routes are generated from [schema/openapi.yaml](schema/openapi.yaml) by oapi-codegen (config in [oapi.config.yaml](oapi.config.yaml)).

```bash
task generate
task run -- serve --port 8080
curl localhost:8080/api/health
```

To add an endpoint, add the path to the schema, run `task generate`, then implement the method on `Server` in [server/server.go](server/server.go). Do not edit `server/server.gen.go`.

If your project does not need a server, delete `server/`, `schema/` and `oapi.config.yaml`, then remove the `serve` command in [cmd/goilerplate/main.go](cmd/goilerplate/main.go), the `generate` task in [Taskfile.yaml](Taskfile.yaml) and the `server` dependencies from `go.mod`.

## Frontend

The `frontend` package ([frontend/embed.go](frontend/embed.go)) embeds the built React app and serves it at `/`; the API stays under `/api`. It uses React, Tailwind and shadcn components, with Vite+ (Vite, Oxlint, Oxfmt, TypeScript) as the toolchain.

```bash
task build        # builds the frontend, then the binary
task check:fe     # format, lint and type-check
```

For development, run the API and the Vite dev server separately. Vite proxies `/api` to `localhost:8080`.

```bash
task run -- serve --port 8080
task dev:fe       # http://localhost:5173
```

If your project does not need a frontend, delete `frontend/`, then remove the frontend FS argument from `server.Start` and the static route in [server/serve.go](server/serve.go). In [server/server_test.go](server/server_test.go), delete `TestFrontend` and the `testing/fstest` import, and call `newRouter()` without an argument. Then remove the `frontend` import in [cmd/goilerplate/main.go](cmd/goilerplate/main.go), the `build:fe`, `check:fe` and `dev:fe` tasks in [Taskfile.yaml](Taskfile.yaml), the frontend stage in [Dockerfile](Dockerfile), and the `frontend` job in [.github/workflows/lint-test.yaml](.github/workflows/lint-test.yaml).

## Example worker

The `worker` package ([worker/worker.go](worker/worker.go)) contains an example of a worker that runs in a goroutine, with proper synchronisation. Run it with:

```bash
go run ./cmd/goilerplate worker --period 200ms
```

Replace `(*Worker).Run` in `worker/worker.go` with the work your project needs.

If your project does not need a worker, delete `worker/` and the `worker` command in [cmd/goilerplate/main.go](cmd/goilerplate/main.go).

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
