<p align="center">
  <img alt="GShark logo" src="https://s1.ax1x.com/2018/10/17/idhZvj.png" />
  <h3 align="center">GShark</h3>
  <p align="center">Scan public code platforms for exposed sensitive information.</p>
</p>

<div align="center">
  <a href="README_CN.md">🇨🇳 中文版</a> | <strong>🇺🇸 English</strong>
</div>

# GShark [![Go Report Card](https://goreportcard.com/badge/github.com/madneal/gshark)](https://goreportcard.com/report/github.com/madneal/gshark) [![Release](https://github.com/madneal/gshark/actions/workflows/release.yml/badge.svg)](https://github.com/madneal/gshark/actions/workflows/release.yml)

GShark is a Go/Gin backend with a Vue 3 management interface for monitoring sensitive information exposed on public code platforms. It supports GitHub, GitLab, Sourcegraph, Postman, and subdomain discovery through Gobuster.

GShark searches public APIs and indexes. Whether private repositories are visible depends on the platform API and token permissions. For the background and demos, see the [articles](https://mp.weixin.qq.com/mp/appmsgalbum?__biz=MzI3MjA3MTY3Mw==&action=getalbum&album_id=2376148333116850178#wechat_redirect), [videos](https://mp.weixin.qq.com/mp/appmsgalbum?__biz=MzI3MjA3MTY3Mw==&action=getalbum&album_id=1834365721464651778#wechat_redirect), and [wiki](https://github.com/madneal/gshark/wiki).

## Features

- GitHub code, Issue/PR, and public Gist searches
- GitLab global search with a project-crawl fallback
- Sourcegraph global code search (existing `searchcode` rules are also supported)
- Postman collection search and subdomain discovery
- Custom rules, keyword/extension filters, result review, and export
- Optional AI pre-ingest filtering with multiple OpenAI-compatible providers
- Docker deployment with scanner resource guardrails

## Quick start: Docker

Requirements: Docker Engine and Docker Compose.

```bash
git clone https://github.com/madneal/gshark.git
cd gshark

# Starts MySQL, backend, and web; the scanner is intentionally not started yet.
./scripts/quick-docker.sh \
  --admin-user myadmin \
  --admin-password 'change-this-password'
```

Open <http://localhost:8080>. After logging in, add tokens and rules, then start scanning:

```bash
docker compose up -d scan
```

The default account is `gshark / gshark` when no custom account is supplied. Change it immediately in any non-local deployment. The quick script initializes the database before starting the scanner; do not start `scan` manually before MySQL initialization is complete.

Useful Docker commands:

```bash
docker compose ps
docker compose logs -f server
docker compose logs -f scan
docker compose restart server
docker compose stop scan       # pause scanning without stopping the web UI
docker compose down            # keeps the ./mysql data directory
```

The scanner is limited to 512 MB RAM, 1 CPU, and `GOMEMLIMIT=384MiB` by default. Adjust the `scan` limits in [docker-compose.yaml](docker-compose.yaml) only after checking `docker stats gshark-scanner`. Never use `docker compose down -v` unless you intend to delete the database volume.

> **Security:** the sample Docker files contain local-development credentials. Change the MySQL and JWT secrets in `docker-compose.yaml` and `server/config.docker.yaml` before production use.

## First scan

1. Log in to the web UI at `http://localhost:8080`.
2. Add a token under the corresponding platform. Store tokens only in GShark; never commit them to the repository or paste them into an issue.
3. Add one rule per search expression. Examples:

   `text
   password in:file
   access_token org:example
   secret repo:owner/repository
   api_key extension:yaml
   `

   GitHub rule types are `github` (code), `github_issue` (Issue/PR), and `gist` (public Gists). A rule may enable more than one type, for example `github,github_issue,gist`.
4. Add keyword or extension filters when a rule produces too much noise. Filters are currently applied to GitHub code/Gist surfaces.
5. Optionally enable AI pre-ingest filtering. Configure one or more providers, use **Test AI configuration** in the system page, and enable the feature only after the test succeeds. It is disabled by default.
6. Start the scanner with `docker compose up -d scan` (or `./gshark scan` for a manual deployment). The scanner loops periodically while valid tokens and rules exist.
7. Review results in the result page: confirm genuine findings, ignore placeholders and false positives, and export results when needed.

## Platform configuration

- **GitHub:** create a personal access token from [GitHub token settings](https://github.com/settings/tokens). The token must be able to perform the searches required by your rules. GitHub search rate limits still apply; GShark retries rate-limited pages but cannot remove upstream limits.
- **GitLab:** configure the token and, for a self-hosted instance, the GitLab base URL. If global blob search is unavailable, GShark searches recently active public projects instead.
- **Sourcegraph:** configure `search.sourcegraph-url` (default `https://sourcegraph.com`) and `search.sourcegraph-token` or `SOURCEGRAPH_TOKEN` when required. Results depend on the instance index, timeout, and upstream result limits.
- **Postman/subdomain discovery:** configure the corresponding rules and wordlist before starting the scanner.

Rules can also be batch-imported from the rule CSV template. Existing rules with type `searchcode` are handled by the Sourcegraph provider and do not require a database migration.

## Optional AI filtering

AI filtering runs before a finding is written to `search_result`. Only a strict `real: true` verdict is stored; malformed responses, timeouts, HTTP errors, and transport errors fail closed. The built-in limits are a 30-second request timeout and 6,000 characters of evidence.

The feature is disabled by default. Configure providers in `server/config.yaml` (or `server/config.docker.yaml`):

```yaml
system:
  ai_analysis_enabled: true
  ai_providers:
    - name: primary
      server: https://api.openai.com/v1/chat/completions
      token: your-token
      model: gpt-4o-mini
    - name: backup
      server: https://example.com/v1/chat/completions
      token: backup-token
      model: backup-model
```

Providers are tried in order when a transport, authentication, or HTTP error occurs. Legacy `ai_server`, `ai_token`, and `model` fields remain supported as a single-provider fallback. The configuration test uses synthetic placeholder evidence and does not persist a result.

## Release package deployment

The release script supports macOS and Linux. It requires `curl`, `jq`, `unzip`, `nginx`, and `sudo`:

```bash
./scripts/quick-release.sh \
  --admin-user myadmin \
  --admin-password 'change-this-password'
```

It downloads the matching release zip, installs the frontend into the detected Nginx web root, proxies `/api/` to backend port `8888`, and serves the UI on port `8080`. To use a local package, add `--file ./gshark_linux_amd64.zip`. Use `--skip-init` when the database is already initialized.

## Manual deployment

Requirements: MySQL 8.0+, Go 1.25+, Node.js 20+, npm, and Nginx.

```bash
cd server
cp config-temp.yaml config.yaml
go build -o gshark
./gshark serve       # backend, default port 8888
./gshark scan        # scanner; run after the database is initialized
```

Build the frontend with `npm install && npm run build` in `web/`, serve `web/dist` from Nginx, and reverse proxy `/api/` to `127.0.0.1:8888`. The browser should use the Nginx port (8080 in the examples), not the backend port.

For an existing installation, back up MySQL, apply the version-specific statements in [sql.md](sql.md), restart the backend, and start the scanner after the migration completes. Do not assume that copying a new binary alone performs every historical migration.

## Configuration reference

The main configuration file is `server/config.yaml`. Start from [server/config-temp.yaml](server/config-temp.yaml); Docker uses [server/config.docker.yaml](server/config.docker.yaml). The `-c` flag and `GSHARK_CONFIG` environment variable can select another file:

```bash
./gshark -c /path/to/config.yaml serve
GSHARK_CONFIG=/path/to/config.yaml ./gshark scan
```

## Troubleshooting

- **401 from GitHub:** verify the token itself with the GitHub API, check its expiration and permissions, and make sure the token is assigned to the expected rule/platform.
- **No results:** confirm that `scan` is running, the rule syntax matches the platform, the token is valid, and the provider is reachable. Check `docker compose logs server scan`.
- **Scanner does not start:** wait for the MySQL health check and database initialization, then run `docker compose up -d scan`.
- **Slow page or high memory:** inspect `docker stats`, scanner logs, and MySQL health before changing resource limits.
- **GitHub rate limit:** reduce noisy rules and search scope. The scanner retries rate-limited pages, but upstream limits remain.
- **Upgrade issue:** check the release version, back up the database, apply the matching [sql.md](sql.md) section, and restart server before scanner.

## Development

Backend:

```bash
cd server
go mod tidy
cp config-temp.yaml config.yaml
go run main.go serve
```

Frontend:

```bash
cd web
npm install
npm run serve
```

Run the focused AI tests with `go test ./service -run 'Test(SearchResultContent|ParseSearchResultAnalysis|AnalyzeSearchResult|TestAIConfig)'`. Run the broader suite only in an isolated test database. On macOS ARM, set `CGO_ENABLED=1` when CPU percentages are needed in the server information page.

## Resources

### Articles

* [多平台的敏感信息监测工具-GShark](https://mp.weixin.qq.com/s?__biz=MzI3MjA3MTY3Mw==&mid=2247484283&idx=1&sn=3232df7d321c0f62ce61b7e6368204ad&chksm=eb396deddc4ee4fb0c825a378c085223b87fc45f05648d46e7bdc24a03fb83ad6c7ade414df7#rd)
* [GShark-监测你的 Github 敏感信息泄露](https://mp.weixin.qq.com/s?__biz=MzI3MjA3MTY3Mw==&mid=2247483770&idx=1&sn=9f02c2803e1c946e8c23b16ff3eba757&chksm=eb396fecdc4ee6fa2f378e846f354f45acf6e6f540cfd54190e9353df47c7707e3a2aadf714f&token=1578822041&lang=zh_CN#rd)

### Videos

* [GShark v1.5.0 版本及 Docker 使用指南](https://www.bilibili.com/video/BV1oUe3eBEMz/)
* [GShark v1.3.0 版本支持 Docker](https://www.bilibili.com/video/BV1BH4y1C7Ga/)
* [GShark 支持多种规则类型以及规则配置建议](https://www.bilibili.com/video/BV1uY4y177SX)
* [批量导入规则](https://mp.weixin.qq.com/s?__biz=MzI3MjA3MTY3Mw==&mid=2247484546&idx=1&sn=818915279c5199457340ade89d6cbd54&chksm=eb396a14dc4ee302039bcb1474380a6049dba84370345b7813049aa8feb49a98f89d47ec5d5b#rd)
* [GShark部署](https://mp.weixin.qq.com/s?__biz=MzI3MjA3MTY3Mw==&mid=2247484487&idx=1&sn=78f942ccf6861f433fc7f4a60564441c&chksm=eb396ad1dc4ee3c7505362da243433e54a2b558c96fbbb50f8b6cea87d1f9bc920b249b72705#rd)
* [windows 部署](https://mp.weixin.qq.com/s?__biz=MzI3MjA3MTY3Mw==&mid=2247484289&idx=1&sn=2b0f1c38b88c924ad514fb64b559b784&chksm=eb396d17dc4ee4018573dde6c3bfce83903c86034403539eaf1b87b89c4a4dd44f957a308818#rd)
* [GShark v1.0.2 版本发布](https://www.bilibili.com/video/BV1Zx4y1G7FX/)
* [GShark v1.1.0 更新内容介绍](https://www.bilibili.com/video/BV1aG4y1c72N/)

## License

[Apache License 2.0](https://github.com/madneal/gshark/blob/master/LICENSE)

## 404StarLink 2.0 - Galaxy

![](https://github.com/knownsec/404StarLink-Project/raw/master/logo.png)

GShark 是 404Team [星链计划2.0](https://github.com/knownsec/404StarLink2.0-Galaxy)中的一环，如果对 GShark 有任何疑问又或是想要找小伙伴交流，可以参考星链计划的加群方式。

- [https://github.com/knownsec/404StarLink2.0-Galaxy#community](https://github.com/knownsec/404StarLink2.0-Galaxy#community)
