# GShark CLI

Manage a deployed GShark instance with `gshark-cli`, included in each release archive. On Windows, use `gshark-cli.exe`.

## Quick start

```sh
./gshark-cli auth login --server https://gshark.example.com/api
./gshark-cli results list --output text
```

The CLI opens the GShark website. Log in there (including the captcha), then click **Authorize CLI**. An existing browser session can authorize directly. The terminal never asks for your password or captcha. Login is saved automatically; repeat this flow when it expires.

For a direct backend connection, also specify its website: `auth login --server http://127.0.0.1:8888 --web-url http://localhost:8080`. The website must include the CLI authorization page from the same release and connect to that backend. With a frontend `/api` URL, the website URL is inferred; use `--web-url` to override it.

If the browser does not open, use the printed URL on the same computer. `--no-browser` only prints the URL; `--login-timeout 10m` extends the default five-minute wait. The callback listens only on `127.0.0.1`; for remote/headless scripts, use `GSHARK_TOKEN` instead.

## Common commands

| Task | Command |
| --- | --- |
| Saved login status | `auth status` |
| Remove saved login | `auth logout` |
| Unprocessed results | `results list --status 0` |
| Result details | `results get 123 --include-sensitive` |
| Mark processed | `results set-status 123 124 --status processed` |
| Ignore results | `results set-status 123 --status ignored` |
| Export CSV | `results export --format csv --file results.csv` |
| Delete results | `results delete 123 124` |
| List rules | `rules list` |
| Create a rule | `rules create --type github --name example --content 'org:example password'` |
| Update a rule | `rules update 7 --description 'Updated description'` |
| Enable / disable a rule | `rules enable 7` / `rules disable 7` |
| Manage repositories | `repos list` / `repos get 5` |
| Manage filters | `filters list` / `filters get 3` |
| Scan overview | `scans overview` |
| Scan logs | `scans logs --provider github` |

Rules, repositories and filters support `list`, `get`, `create`, `update` and `delete`. Run `./gshark-cli <command> --help` for parameters, for example `./gshark-cli filters create --help`.

## Options

- Output defaults to JSON; use `--output text` for reading.
- Lists support `--page 1 --limit 20` (maximum 100 per page).
- Result queries and exports hide sensitive fields by default; use `--include-sensitive` to reveal them.
- Export files must be new. Exports default to 10,000 records; adjust `--max-results` and check `truncated` in the completion summary.
- For scripts, set `GSHARK_SERVER` and `GSHARK_TOKEN` instead of interactive login.
