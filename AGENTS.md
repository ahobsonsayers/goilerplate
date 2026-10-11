# AGENTS.md

Boilerplate repo for new Go projects. Single module, one package (`cmd/goilerplate`).

## Commands

- Lint (also fixes): `task lint` — runs `golangci-lint v2` via `go tool`; config in `.golangci.yaml`
- Format: `task format` — gofumpt via `go tool`
- Test: `task test` — `go test ./...`
- Build: `task build`; Docker: `task build:docker`; verify image user: `docker inspect --format '{{.Config.User}}'` (must be 65532 — distroless nonroot, keep it non-root)

## Toolchain quirks

- golangci-lint and gofumpt are `tool` directives in go.mod, run with `go tool <name>`.
- `go get -u ./...` does not bump the golangci-lint tool dependency; update it explicitly with `go get github.com/golangci/golangci-lint/v2@latest` and run `go mod tidy`.
- Go version is pinned in three places — bump together when it changes: go.mod `go` directive, Dockerfile `FROM golang:X`, `.golangci.yaml` `run.go`. CI reads `go-version-file: go.mod`.
- Running `go tool golangci-lint` may fail to link with cgo on new Xcode SDKs; `CGO_ENABLED=0 go tool golangci-lint run ./...` bypasses it. CI is unaffected.

## How this repo is used

- "goilerplate" is renamed throughout when a project is created from it (README First Steps, via `fd` + `sd`).
- Downstream repos pull updates via squash merge with `--allow-unrelated-histories` from a `boilerplate` remote (see README "Merge latest goilerplate changes").
- `.github/dependabot.yaml` covers three ecosystems (gomod, github-actions, docker), all weekly and grouped into one PR per ecosystem.