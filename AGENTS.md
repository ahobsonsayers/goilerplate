# AGENTS.md

Boilerplate repo for new Go projects. Single module: `cmd/goilerplate` (cli), `worker` (example worker - delete if the project does not need one), `server` (chi + oapi-codegen API) and `frontend` (React app embedded into the binary).

## Commands

- Lint (also fixes): `task lint` — runs `golangci-lint v2` via `go tool`; config in `.golangci.yaml`
- Format: `task format` — gofumpt via `go tool`
- Test: `task test` — `go test ./...`
- Generate: `task generate` — runs oapi-codegen from the repo root and regenerates `server/server.gen.go` from `schema/openapi.yaml` (config: `oapi.config.yaml`). Do not edit the generated file; call `GetSpec()` rather than deprecated `GetSwagger()`.
- Build: `task build` (builds the frontend first); Docker: `task build:docker`; verify image user: `docker inspect --format '{{.Config.User}}'` (must be 65532 — distroless nonroot, keep it non-root)
- Frontend check (format, lint, types): `task check:fe` — Vite+ `vp check`, run through bun
- Frontend dev server: `task dev:fe` — Vite proxies `/api` to `localhost:8080`
- Run the cli with args: `task run -- greet` (or `go run ./cmd/goilerplate ...`)

## Toolchain quirks

- golangci-lint and gofumpt are `tool` directives in go.mod, run with `go tool <name>`.
- `go get -u ./...` does not bump the golangci-lint tool dependency; update it explicitly with `go get github.com/golangci/golangci-lint/v2@latest` and run `go mod tidy`.
- Go version is pinned in three places — bump together when it changes: go.mod `go` directive, Dockerfile `FROM golang:X`, `.golangci.yaml` `run.go`. CI reads `go-version-file: go.mod`.
- Running `go tool golangci-lint` may fail to link with cgo on new Xcode SDKs; `CGO_ENABLED=0 go tool golangci-lint run ./...` bypasses it. CI is unaffected.
- Config is strict: revive and gocritic run with enable-all. Common failures: unused method receivers must be nameless (`func (*Worker) Run()`), `//nolint` with no space before it, error assertions in tests must use `require` not `assert`.
- The frontend uses bun and Vite+ (`vite-plus` dev dependency, `bun run check`). `frontend/dist/.gitkeep` is committed because `//go:embed all:dist` requires the directory to exist; real build output is gitignored.

## Worker example

- `worker/` is a deletable boilerplate example, not a library: if the project needs no worker, delete `worker/` and the `worker` command in `cmd/goilerplate/main.go` (see README "Example worker").
- In cli wiring that awaits a worker, block on the Start result channel, not `Worker.WaitUntilStopped()`: the WaitGroup counter is 0 until `Start` runs, so `Wait` can return instantly and kill the process before the worker ticks (see comment in `runWorker`).

## How this repo is used

- "goilerplate" is renamed throughout when a project is created from it (README First Steps, via `fd` + `sd`).
- Downstream repos pull updates via squash merge with `--allow-unrelated-histories` from a `boilerplate` remote (see README "Merge latest goilerplate changes").
- `.github/dependabot.yaml` covers three ecosystems (gomod, github-actions, docker), all weekly and grouped into one PR per ecosystem.