# GShark 命令行工具

Release 压缩包内的 `gshark-cli` 用于管理已部署的 GShark。Windows 使用 `gshark-cli.exe`。

## 快速开始

```sh
./gshark-cli auth login --server https://gshark.example.com/api
./gshark-cli results list --output text
```

按提示输入用户名、密码和终端显示的验证码，登录信息会自动保存，过期后重新登录即可。直连后端时使用 `--server http://127.0.0.1:8888`。

## 常用命令

| 操作 | 命令 |
| --- | --- |
| 查看本地登录状态 | `auth status` |
| 清除本地登录信息 | `auth logout` |
| 查询未处理结果 | `results list --status 0` |
| 查看结果原文 | `results get 123 --include-sensitive` |
| 标记已处理 | `results set-status 123 124 --status processed` |
| 忽略结果 | `results set-status 123 --status ignored` |
| 导出 CSV | `results export --format csv --file results.csv` |
| 删除结果 | `results delete 123 124` |
| 查询规则 | `rules list` |
| 创建规则 | `rules create --type github --name example --content 'org:example password'` |
| 修改规则 | `rules update 7 --description '更新说明'` |
| 启用／停用规则 | `rules enable 7` / `rules disable 7` |
| 管理仓库 | `repos list` / `repos get 5` |
| 管理过滤规则 | `filters list` / `filters get 3` |
| 查看扫描概况 | `scans overview` |
| 查看扫描日志 | `scans logs --provider github` |

规则、仓库和过滤规则均支持 `list`、`get`、`create`、`update`、`delete`。完整参数通过 `./gshark-cli <命令> --help` 查看，例如 `./gshark-cli filters create --help`。

## 常用选项

- 默认输出 JSON，添加 `--output text` 可切换为文本。
- 列表支持 `--page 1 --limit 20`，每页最多 100 条。
- 结果查询和导出默认隐藏敏感字段，添加 `--include-sensitive` 可显示原文。
- 导出目标文件不能已存在。默认最多导出一万条，可用 `--max-results` 调整，并查看完成摘要的 `truncated` 字段确认是否截断。
- 脚本可设置 `GSHARK_SERVER` 和 `GSHARK_TOKEN`，代替交互登录。
