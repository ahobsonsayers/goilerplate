# AGENTS.md

Boilerplate repo for new Go projects. Single module, two packages: `cmd/goilerplate` (cli) and `worker` (example worker - delete if the project does not need one).

## Commands

- Lint (also fixes): `task lint` — runs `golangci-lint v2` via `go tool`; config in `.golangci.yaml`
- Format: `task format` — gofumpt via `go tool`
- Test: `task test` — `go test ./...`
- Build: `task build`; Docker: `task build:docker`; verify image user: `docker inspect --format '{{.Config.User}}'` (must be 65532 — distroless nonroot, keep it non-root)
- Run the cli with args: `task run -- greet` (or `go run ./cmd/goilerplate ...`)

## Toolchain quirks

- golangci-lint and gofumpt are `tool` directives in go.mod, run with `go tool <name>`.
- `go get -u ./...` does not bump the golangci-lint tool dependency; update it explicitly with `go get github.com/golangci/golangci-lint/v2@latest` and run `go mod tidy`.
- Go version is pinned in three places — bump together when it changes: go.mod `go` directive, Dockerfile `FROM golang:X`, `.golangci.yaml` `run.go`. CI reads `go-version-file: go.mod`.
- Running `go tool golangci-lint` may fail to link with cgo on new Xcode SDKs; `CGO_ENABLED=0 go tool golangci-lint run ./...` bypasses it. CI is unaffected.
- Config is strict: revive and gocritic run with enable-all. Common failures: unused method receivers must be nameless (`func (*Worker) Run()`), `//nolint` with no space before it, error assertions in tests must use `require` not `assert`.

## Worker example

- `worker/` is a deletable boilerplate example, not a library: if the project needs no worker, delete `worker/` and the `worker` command in `cmd/goilerplate/main.go` (see README "Example worker").
- In cli wiring that awaits a worker, block on the Start result channel, not `Worker.WaitUntilStopped()`: the WaitGroup counter is 0 until `Start` runs, so `Wait` can return instantly and kill the process before the worker ticks (see comment in `runWorker`).

## How this repo is used

- "goilerplate" is renamed throughout when a project is created from it (README First Steps, via `fd` + `sd`).
- Downstream repos pull updates via squash merge with `--allow-unrelated-histories` from a `boilerplate` remote (see README "Merge latest goilerplate changes").
- `.github/dependabot.yaml` covers three ecosystems (gomod, github-actions, docker), all weekly and grouped into one PR per ecosystem.