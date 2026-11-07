# GShark management CLI

Each release archive includes `gshark-cli` (`gshark-cli.exe` on Windows) alongside the server binary. It connects to an existing GShark instance over HTTP and needs no MySQL connection, config.yaml, frontend files, or local server. You can copy the CLI binary to another machine.

每个平台的 release 压缩包都包含独立的 `gshark-cli`（Windows 为 `gshark-cli.exe`）。它通过 HTTP 访问已部署的 GShark，无需本地数据库、服务端配置或前端文件。

## Connection and authentication

```sh
./gshark-cli auth login --server https://gshark.example.com/api --username alice
# Enter password (hidden) and the captcha displayed in the terminal.
./gshark-cli auth status
./gshark-cli results list
./gshark-cli auth logout
```

`auth login` uses the existing backend login and captcha endpoints. Enter your password at the hidden prompt, then read the captcha image rendered in a true-color terminal and enter its digits. Failed or expired captchas require a new login attempt. An interactive terminal is required. No browser token copying is necessary.

通过 `auth login` 输入用户名、密码和终端显示的验证码即可登录，无需复制网页登录 token。密码输入隐藏，验证码失败后重新运行登录命令。需要支持真彩色的交互终端。

The CLI saves the server address, username and JWT in `gshark-cli/credentials.json` under the OS user configuration directory (macOS: `~/Library/Application Support`; Linux: `$XDG_CONFIG_HOME` or `~/.config`; Windows: `%AppData%`). It saves no password. The file uses owner-only permissions on Unix. Only one login is saved at a time; logging in replaces it. Subsequent commands automatically use that server and token, persisting replacement tokens when the backend refreshes them.

`auth status` shows the saved account and local token expiry without printing the token; it does not verify revocation or permissions with the server. Expired tokens require `auth login` again. `auth logout` deletes local credentials; it does not revoke an already issued server token or clear environment variables.

凭据保存在操作系统的用户配置目录中，仅保存服务器、用户名和 JWT，不保存密码。后续命令自动读取并保存后端刷新返回的 token。`auth status` 显示本地凭据状态，不验证服务端权限；`auth logout` 删除本地凭据，不撤销服务端已签发的 token。

For automation, `--token` overrides `GSHARK_TOKEN`, which overrides saved credentials. `--server` overrides `GSHARK_SERVER`, which overrides the saved server. Saved tokens are never sent to a different server; log in again for that server. Prefer environment variables over token arguments to avoid shell history and process-argument exposure.

Direct backend: `http://127.0.0.1:8888`. Through the frontend Nginx proxy: `https://gshark.example.com/api`. Use HTTPS for remote connections. Redirects are rejected. Requests time out after 30 seconds; customize with `--timeout 10s`.

## Commands

```sh
./gshark-cli results list --status 0 --page 1 --limit 20
./gshark-cli results list --keyword 'example' --query 'repository-name'
./gshark-cli results get 123
./gshark-cli results get 123 --include-sensitive
./gshark-cli rules list --type github --name example
./gshark-cli scans overview
./gshark-cli scans logs --provider github --limit 20
./gshark-cli scans logs --status failed --cycle-id example
./gshark-cli results list --output text
```

`results list --status`: `-1` all (default), `0` unprocessed, `1` processed, `2` ignored. Lists return one page, with a default size of 20 and maximum of 100; use `--page` to continue. Filters map to existing backend semantics; result provider filtering is not supported by the current API.

Default output is JSON containing the backend `data` object, including list pagination metadata. Use `--output text` for interactive reading. Successful commands exit 0; validation, network, HTTP and API errors exit 1 and write diagnostics to stderr. Help works without credentials or a running server.

Result commands redact matches, source fragments, keywords and repository/file locations by default. Use `--include-sensitive` on result commands to return original fields. This is a client presentation feature, not server-side authorization. Rules and scan logs can also contain sensitive content; handle exports accordingly.

结果查询默认隐藏匹配内容、代码片段、关键词及仓库和文件位置；通过 `--include-sensitive` 显示原始字段。此开关仅控制客户端输出，不替代后端权限。规则和扫描日志也可能包含敏感信息。

## Build from source

```sh
cd server
go build -o gshark-cli ./cmd/gshark-cli
```

## Operations / 日常运营

```sh
# Mark selected results; status is required and IDs are explicit.
./gshark-cli results set-status 123 124 --status processed
./gshark-cli results set-status 125 --status ignored
./gshark-cli results set-status 125 --status unprocessed
./gshark-cli results delete 123 124

# Export across pages, with the same result filters as list.
./gshark-cli results export --status 0 --file unprocessed.jsonl
./gshark-cli results export --format csv --file results.csv --include-sensitive

# Manage scan rules.
./gshark-cli rules get 7
./gshark-cli rules create --type github --name example --content 'org:example password' --enabled=false
./gshark-cli rules update 7 --description 'Updated description' --match-pattern 'example'
./gshark-cli rules enable 7
./gshark-cli rules disable 7
./gshark-cli rules delete 7 8

# Manage repositories.
./gshark-cli repos list --page 1 --limit 20
./gshark-cli repos get 5
./gshark-cli repos create --type gitlab --url https://gitlab.example/team/project --project-id 42
./gshark-cli repos update 5 --description 'Tracked project'
./gshark-cli repos delete 5

# Manage GitHub whitelist/blacklist filters.
./gshark-cli filters list
./gshark-cli filters get 3
./gshark-cli filters create --type blacklist --class extension --content '.md'
./gshark-cli filters update 3 --content '.txt'
./gshark-cli filters delete 3
```

Result statuses are `unprocessed`, `processed`, and `ignored`. Status updates and deletes accept 1–100 explicit IDs separated by spaces; duplicate IDs are deduplicated. A successful mutation reports the backend response, not a verified affected-row count. Backend role permissions and operation auditing continue to apply.

Update commands first fetch the existing record and merge only supplied flags, preserving unspecified fields. An explicit empty string clears a text field. Rule state changes use `enable`/`disable`. This read-then-save follows the current backend API and has no optimistic locking: concurrent edits can overwrite each other.

Exports support JSONL and CSV. Files must not already exist; new files use owner-only permissions on Unix. Failed exports remove partial files. The default maximum is 10,000 results; set `--max-results` to change it and check the `truncated` field in the JSON completion summary. Result redaction applies unless `--include-sensitive` is supplied. CSV contains basic result columns; JSONL retains all returned fields. Formula-like CSV cells are prefixed with an apostrophe for spreadsheet safety. Paginated exports are not database snapshots; changes during export may cause rows to shift.

运营功能覆盖结果状态处理、显式 ID 批量删除、分页导出、规则增改及启停、仓库和过滤规则增删改查。更新时先读取详情，仅合并指定字段。过滤类型支持 `whitelist`/`blacklist`，种类支持 `extension`/`keyword`。导出默认脱敏、最多一万条，完成摘要中的 `truncated` 表明是否截断。并发修改可能影响更新与分页导出的一致性。

Continue using the server's `gshark serve`, `gshark scan`, and `gshark init` commands for service operation. The CLI does not submit remote scan jobs.

