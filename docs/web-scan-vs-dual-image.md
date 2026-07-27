# Web 扫描任务管理 vs 双镜像（serve / scan）

本文回答：在 GShark 里，**Web 能管理什么扫描相关工作**，以及与 **双进程/双镜像**（`gshark serve` + `gshark scan`）相比哪种部署更好。

## 1. 架构：三个「任务」不要混谈

| 概念 | 代码入口 | 谁执行 | Web 能否控制 |
|------|----------|--------|--------------|
| **主扫描循环** `ScanTask` | `server/search/scan.go` → 各平台 `RunTask` | **仅** `gshark scan` 进程 | **否**（无 API 启停） |
| **二次过滤** | `POST /searchResult/startSecFilterTask` | **serve** 进程内 `go func` | **是**（结果页按钮） |
| **AI 分析** | `POST /searchResult/startAITask` | **serve** 进程内 `go func` | **是**（结果页按钮） |
| **二次过滤状态** | `GET /searchResult/getTaskStatus` | serve 内存变量 `taskStatus` | **只读**（非 scan 循环状态） |
| **DB 种子 `task` 表** | `server/source/task.go` / `model.Task` | 初始化写入 | **Web 无管理页**；且当前 `search/*` **不读**该表 |

CLI 分叉（`server/main.go`）：

- `gshark serve` → `core.RunServer()`（HTTP / 管理面）
- `gshark scan` → `search.ScanTask()`（无限循环：GitLab → Searchcode → GitHub → gobuster → Postman）

Docker：

- `server/deploy/serve/Dockerfile` → `ENTRYPOINT ./gshark serve` → 容器 `gshark-server`
- `server/deploy/scan/Dockerfile` → `ENTRYPOINT ./gshark scan` → 容器 `gshark-scanner`
- `docker-compose.yaml`：`scan` 服务默认 **不**随 `quick-docker.sh` 启动；需 `./scripts/quick-docker.sh --with-scan` 或 `docker compose up -d scan`

```
                    ┌─────────────┐
   Browser ────────►│ gshark-web  │ nginx :8080
                    └──────┬──────┘
                           │ /api
                    ┌──────▼──────┐         MySQL
   管理面 ──────────►│gshark-server│◄──────── (rule/token/filter/
   serve            │ gshark serve│          search_result)
                    └─────────────┘
                           ▲ 读配置/规则/Token，写 search_result
                    ┌──────┴──────┐
   扫描面（可选）───►│gshark-scanner│
   scan             │ gshark scan │
                    └─────────────┘
```

共享状态：**同一 MySQL** + 同一 `config.docker.yaml`（挂载到 server 与 scan）。  
Web 改规则/Token → DB；**只有 scan 进程**会周期性 `GetValidRulesByType` / `ListTokenByType` 拉数据并写 `search_result`。

## 2. Web 实际能管理什么

### 2.1 配置与结果（管理面，不依赖 scan 容器是否在跑）

| 模块 | 前端 | 后端能力 |
|------|------|----------|
| 规则 | `web/src/view/rule/rule.vue` | CRUD / 开关；scan 用 `GetValidRulesByType`（`status=1`） |
| Token | `web/src/view/token/token.vue` | CRUD；scan 用 `ListTokenByType`（GitHub/GitLab 等） |
| 过滤 | `web/src/view/filter/filter.vue` | 扩展名/二次关键词等；scan 与二次过滤都会读 |
| 搜索结果 | `web/src/view/searchResult/searchResult.vue` | 列表/导出/按 ID 或 repo 改状态（确认/忽略） |
| 子域名等资产 | 菜单「搜索结果」下子域名页 | 列表管理（scan 中 gobuster 可产出） |

### 2.2 Web 内可点的「任务」——**不是**主扫描编排

结果页（`searchResult.vue`）：

- **启动二次过滤** → `startSecFilterTask`：在 **serve** 上对未处理 repo 再搜 `sec_keyword`，改结果状态；状态串 `taskStatus`（`stop`/`running`/…）仅在 serve 进程内存。
- **启动 AI 分析** → `startAITask`：在 **serve** 上对未处理结果调 AI，更新 status。

二者 **不**启动 `ScanTask`，也 **不**启停 `gshark-scanner` 容器。

### 2.3 Web **不能**做的

- 启动/停止多平台主扫描循环（`gshark scan` / `gshark-scanner`）
- 查看「scan 是否在跑、下一轮何时扫描」的编排仪表盘（无对应 API）
- 通过 UI 管理 `task` 表（种子有 github/gitlab/postman/searchcode 行，但前端无页面，且 `search/*` 未按 `task.task_status` 开关平台）

## 3. 双镜像如何分工

| | **serve**（`gshark-server`） | **scan**（`gshark-scanner`） |
|--|------------------------------|------------------------------|
| 入口 | `./gshark serve` | `./gshark scan` |
| 职责 | REST API、鉴权、规则/Token/结果 CRUD、二次过滤/AI | 周期调用各平台 `RunTask`，外网检索并 `SaveSearchResults` |
| 重启策略 | compose `restart: always` | compose `restart: "no"`（退出后不自动再起） |
| 默认部署 | `quick-docker.sh` **会**启动 | 仅 `--with-scan` 或手动 `compose up scan` |
| 失败形态 | 服务不可用 → 管理面挂 | 无 Token/规则可 panic 退出（见验证日志），**不影响** Web 管理 |

主扫描一轮间隔约 **900s**（`scan.go` 传入各 `RunTask` 的 sleep）。

## 4. 两种部署形态对比与推荐

| 形态 | 组成 | 适合 | 限制 |
|------|------|------|------|
| **A. 仅管理面** | mysql + server + web | 配置规则/Token、 dig 结果、二次过滤/AI、日常运营 | **无**持续外网敏感信息扫描；`search_result` 不会因 ScanTask 增长 |
| **B. 管理 + 扫描** | A + scan（`--with-scan`） | 生产/要持续发现敏感信息 | 需有效 Token 与启用规则；scan 崩溃需运维重启（`restart: no`） |

### 明确推荐

1. **只做配置与结果运营（本地调试 UI、导入规则、人工 triage）**  
   → 用形态 **A**（默认 `quick-docker.sh`）。更简单、资源更少；Web 任务按钮（二次过滤/AI）已够。

2. **要「系统自动扫 GitHub/GitLab/… 并进结果库」**  
   → 必须用形态 **B** 或等价地另起 `gshark scan` 进程。  
   **Web 不能替代 scan 容器。**

3. **更好的心智模型（优于「Web 里管扫描」）**  
   - **Web + serve** = 控制面（配置 + 结果 + 轻量后处理任务）  
   - **scan 镜像** = 数据面（持续采集）  
   两者通过 **MySQL** 解耦；先配规则/Token，再起 scan，比把扫描塞进 Web 进程更稳（长任务、外网限流、崩溃隔离）。

4. **何时不必起 scan**  
   - 已有结果数据只需 triage  
   - 仅用二次过滤/AI 加工现有未处理结果  
   - CI/文档环境只验证管理 API  

## 5. 验证证据

### 5.1 代码路径（门控）

- `server/main.go`：`serve` / `scan` / `init` 子命令  
- `server/search/scan.go`：`ScanTask` 循环  
- `server/search/githubsearch/search.go`：`RunTask` → `GetValidRulesByType("github")`  
- `server/service/rule.go`：`GetValidRulesByType`（Web 配置的启用规则）  
- `server/deploy/serve|scan/Dockerfile`：不同 ENTRYPOINT  
- `docker-compose.yaml`：`server` vs `scan`  
- `scripts/quick-docker.sh`：默认 `mysql server web`；`--with-scan` 另起 scan  
- `web/src/api/searchResult.js`：`startSecFilterTask` / `getTaskStatus` / `startAITask`  
- `server/api/search_result.go`：`taskStatus` 内存变量，与 ScanTask 无关  

### 5.2 运行观测（本机 Docker）

证据文件（会话 scratch，goal 验证用）：

- `{SCRATCH}/scan-deploy-obs.txt` — 默认栈 **无** `gshark-scanner`；server entrypoint 为 `./gshark serve`  
- `{SCRATCH}/task-api-obs.txt` — 已登录 `GET .../getTaskStatus` → `msg: "stop"`；`compose up -d scan` 后 scanner **跑 `gshark scan`**，因无 GitHub Token 在 `GetGithubClient` panic 退出（exit 2），server/web **仍运行**

### 5.3 仓库内回归测试

- `server/search/architecture_separation_test.go`：静态断言 serve/scan 分叉、ScanTask 平台列表、二次过滤 API 挂在 serve 路由、scan Dockerfile 入口为 `gshark scan`、Web 结果页按钮不指向 scan。

运行：

```bash
cd server && go test ./search -run TestArchitecture -count=1
```

## 6. 结论（一句话）

**Web 适合管理「规则 / Token / 过滤 / 结果」和「serve 上的二次过滤与 AI」；持续多平台扫描必须另起 `gshark scan`（或 scan 镜像）。默认双服务（serve+web）更好做日常管理；需要自动发现敏感信息时再加 scan——双镜像分离优于把扫描塞进 Web。**
