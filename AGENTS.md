# AGENTS.md

## Cursor Cloud specific instructions

This is a Go project (`prometheus-nginxlog-exporter`) that parses NGINX log files and exposes Prometheus metrics on port 4040.

### Build & Test

- **Lint**: `go vet ./...`
- **Unit tests**: `go test ./...` (covers `pkg/config`, `pkg/parser/jsonparser`, `pkg/parser/textparser`, `pkg/relabeling`)
- **Build**: `CGO_ENABLED=0 go build -a -installsuffix cgo -o prometheus-nginxlog-exporter .`
- **Acceptance tests**: Require the built binary at `./prometheus-nginxlog-exporter` in the repo root, then `behave` (Python BDD tests in `features/`). Install test deps with `pip install behave requests`.

### Running the exporter

The binary requires a config file (HCL or YAML). Example:

```sh
./prometheus-nginxlog-exporter -config-file <path-to-config>
```

Metrics are served at `http://localhost:4040/metrics`. The exporter tails log files — only new lines appended after startup are counted (it does not re-read existing content on startup).

### Key caveats

- The acceptance tests start and stop the exporter binary automatically. They create temporary files in `.behave-sandbox/`. No external services are needed.
- The built binary (`prometheus-nginxlog-exporter`) is `.gitignore`d — rebuild after pulling changes.
- `go fmt` is required for all code changes (CI enforces it).
- No database or external services needed for development or testing.
