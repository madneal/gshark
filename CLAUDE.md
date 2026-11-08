# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

GShark is a sensitive-information detection platform: a Go backend (`server/`) and a Vue 3 frontend (`web/`) in one repo. One Go binary provides three cobra subcommands: `serve` (web API), `scan` (scanner daemon), and `init` (DB init + admin seed; exit codes: 0 applied, 1 error, 2 already initialized).

## Commands

### Backend (from `server/`)

```bash
cp config-temp.yaml config.yaml   # then set the MySQL connection — required before anything runs
go build -o gshark .
go run main.go serve              # API server on :8888
go run main.go scan               # scanner (needs tokens + rules configured first)
go test ./...                     # all tests
go test ./search/githubsearch -run TestName   # single test
```

On macOS ARM, use `CGO_ENABLED=1` if CPU-usage percentages on the server-info page matter.

### Frontend (from `web/`)

```bash
npm install
npm run serve        # Vite dev server (proxies /api/ to backend :8888)
npm run build        # outputs web/dist
npm run test:login   # Playwright smoke tests; need Chrome and a running frontend
npm run test:pages   # env: FE_TEST_URL (default http://localhost:8081/), CHROME_PATH
```

The Playwright tests (`web/tests/`) forge a JWT locally in `tests/auth.mjs` — they don't log in through the API.

### Docker / release

```bash
./scripts/quick-docker.sh [--with-scan] [--admin-user U --admin-password P] [--skip-init]
docker compose up -d scan       # start scanner after tokens/rules are configured
./scripts/quick-release.sh      # release-package deployment (MySQL + Nginx)
./build.sh                      # cross-platform release zips into releases/
```

## Architecture

The backend follows the gin-vue-admin scaffold layering. A request flows: `router/` (route registration) → `api/v1/` (handlers) → `service/` (business logic) → `model/` (GORM models; `model/request`, `model/response` for DTOs).

- **Global singletons** in `global/`: `GVA_DB` (GORM/MySQL), `GVA_LOG` (zap), `GVA_CONFIG`, `GVA_VP` (viper). They're bootstrapped in `main.go`'s `init()` via `initialize/` (viper → zap → gorm). **`GVA_DB` may be nil** on first run (before DB init) — `serve` tolerates this so the web-based init flow works; `scan` refuses to start.
- **Auth**: routes split into a public group (login/captcha, DB init) and a private group behind `middleware.JWTAuth()` + `middleware.CasbinHandler()` (Casbin RBAC, roles/menus/APIs stored in DB).
- **Scanner** (`search/`): `ScanTask()` loops on a 15-minute interval, running provider tasks — `githubsearch` (code/issues/gist surfaces), `gitlabsearch` (global search with fallback to per-project crawl), `sourcegraphsearch` (Stream API; also picks up legacy `searchcode` rules), `postman`, `gobuster` (subdomains). Each provider run is wrapped in a scan-log lifecycle (start/heartbeat/finish, interrupted-run marking) recorded via `service/scan_log.go`. Results pass through filters (`service/filter.go`) and optional AI pre-ingest analysis (`service/ai.go`, OpenAI-compatible endpoints, fail-closed) before persisting to `search_result` page by page.
- **Frontend** (`web/src/`): Vue 3 + Element Plus + Vuex + vue-router; dynamic menus/routes come from the backend permission system. API clients in `src/api/`, pages in `src/view/`.

## Database migrations

GORM auto-migrates schema *additions* at startup (`initialize.MysqlTables`). Version-specific manual statements live in `sql.md` — run only those the release notes call for, never the whole file.

## Config notes

- `server/config.yaml` is the live config (gitignored content varies); `config-temp.yaml` is the template; `config.docker.yaml` is used in the Docker image.
- Scanner behavior knobs are under `search.` (e.g. `gitlab-discover-pages`, `gitlab-batch-size`, `sourcegraph-url`/`sourcegraph-token`); AI filtering under `system.ai_providers` (disabled by default).
