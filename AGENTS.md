# AGENTS.md

## Repository overview

gostc-open is an intranet tunneling management platform (based on FRP/GOST). It consists of **four independent components**, each with its own Go module or frontend project. There is no shared root `go.mod`; always `cd` into the correct subdirectory before running Go or npm commands.

## Directory → component mapping

| Directory | Component | Module name | Language |
|-----------|-----------|-------------|----------|
| `server/` | Admin backend (API + embedded frontend) | `server` (Go 1.24) | Go |
| `web/` | Admin frontend (served embedded in server) | — | Vue 3 / Vite |
| `client/` | Client & node agent (CLI + GUI + webui) | `gostc-sub` (Go 1.24) | Go |
| `client/webui/frontend/` | Client local webui frontend | — | Vue 3 / Vite |
| `proxy/` | Gateway/proxy service (custom domains) | `proxy` (Go 1.25) | Go |

## Build commands

### Frontend → server embed (required before building server binary)

```sh
# 1. Build admin frontend
cd web && npm install && npm run build

# 2. Package into server embed
cd web && zip -r dist.zip dist
mv web/dist.zip server/web/
```

The server binary embeds the frontend via `//go:embed dist.zip` in `server/web/static.go`. If you skip this step, the server build fails with a missing embed file error.

The client webui uses the same pattern: `client/webui/frontend/` → build → zip → `client/webui/backend/web/dist.zip`.

### Go binaries

```sh
# Server (from server/)
go build -ldflags "-s -w" -o server main.go

# Client CLI (from client/)
go build -ldflags "-s -w" -o gostc ./cli/

# Client GUI - Windows only (from client/)
go build -ldflags="-H windowsgui -w -s" -o gostc-gui.exe ./gui/

# Proxy (from proxy/)
go build -ldflags "-s -w" -o gostc-proxy main.go
```

All builds use `CGO_ENABLED=0`. Cross-platform builds use goreleaser (`goreleaser release --snapshot --clean`), each directory has its own `.goreleaser.yaml`.

## GORM code generation

Server uses `gorm.io/gen`. Generated query code lives at `server/repository/query/`. After modifying models in `server/model/`, regenerate:

```sh
cd server && go run ./generate/main.go
```

Or: `cd server && bash generate.sh`

## Formatting check

CI enforces `gofmt` on the server module:

```sh
cd server && test -z "$(gofmt -l .)"
```

No golangci-lint or other linter is configured. Use `gofmt` to format Go code.

## CI workflows

- `build.yml`: triggered by non-beta version tags (`v*`, excluding `*beta*`). Builds all three components, publishes GitHub releases and Docker images to Docker Hub (`sianhh/gostc-admin`, `sianhh/gostc`, `sianhh/gostc-proxy`).
- `beta.yml`: triggered by beta tags (`*beta*`). Same flow, tags Docker images with beta version only (no `latest`).
- `frontend.yml`: triggered on pushes/PRs touching `web/**` or `server/**`. Builds frontend, embeds in server, runs `gofmt` check and `go build ./...`, publishes Docker image to GHCR.

## Key architectural notes

- **Separate Go modules**: each of `server/`, `client/`, `proxy/` is independent. They share no code at the Go level; don't try to import across them.
- **Database**: server uses SQLite via `ncruces/go-sqlite3` (pure Go, no CGO) with GORM. MySQL is also a supported driver (`gorm.io/driver/mysql`).
- **CLI framework**: all three Go components use `spf13/cobra` for command parsing.
- **RPC between components**: server ↔ client/node communication uses `lesismal/arpc` (not gRPC).
- **Frontend stack**: Vue 3 + Naive UI + Pinia + Vite. Both `web/` and `client/webui/frontend/` share the same stack and identical `package.json` structure.
- **Docker timezone**: all Dockerfiles hardcode `Asia/Shanghai` timezone.

## Code style conventions

### Naming

| Element | Convention | Example |
|---------|-----------|---------|
| Package names | `snake_case` (非标准 Go) | `package system_user`, `package gost_node` |
| Service packages | 统一用 `package service` | `server/service/admin/system_user/service.go` |
| Constants | `SCREAMING_SNAKE_CASE` | `ALLOW_EDIT`, `GOST_NODE_LIMIT_KIND_ALL` |
| File names | dot-separated by operation | `service.go`, `service.create.go`, `service.page.go` |
| Controller files | 固定命名 `api.go` | `server/controller/admin/system_user/api.go` |
| JSON tags | camelCase | `json:"inputBytes"`, `json:"createdAt"` |
| YAML config tags | kebab-case | `yaml:"auth-key"`, `yaml:"db-type"` |

### Import style

所有 import 放在**一个分组**内（stdlib、third-party、local 混在一起，不用空行分隔），仅 blank import (`_ "..."`) 用空行隔开：

```go
import (
    "errors"
    "github.com/google/uuid"
    "go.uber.org/zap"
    "server/model"
    "server/repository"
)
```

### Comments & error messages

- 注释和错误消息使用**中文**：`errors.New("该账号已被使用")`
- 日志消息：业务用中文 `"新增用户失败"`，基础设施用英文 `"server listen fail"`
- GORM `comment:` tag 和 binding `label:` tag 用中文
- 注释密度低，大部分函数无 doc comment

### Architecture pattern: Controller → Service → Repository

- **Controller** (`api.go`)：仅做 JSON 绑定 + 调用 service + 返回响应，无业务逻辑
- **Service**（`service.go` + `service.<op>.go`）：未导出 struct + 包级 nil 单例 `var Service *service`，每个 CRUD 操作一个文件
- **Repository**（`repository.Get("")`）：返回 `(*query.Query, memory.Interface, *zap.Logger)`，`domain` 参数始终为 `""`

典型 service 方法起手：
```go
db, _, log := repository.Get("")
```

### HTTP handler pattern

- **所有端点均为 POST**（包括读操作 `page`、`list`）
- 请求绑定：`c.ShouldBindJSON(&req)`
- 响应：HTTP 200 + `code` 字段区分状态，通过 `bean.Response.Ok/Fail/OkData/Param` 返回
- 分页：请求嵌入 `bean.PageParam`，响应用 `bean.NewPage(list, total)`
- Request/Response 类型定义在 **service 包**内，不在 controller 中

### Error handling

- `errors.New("中文用户消息")` 直接返回，**不使用 `%w` 包装**
- 内部错误 log 后返回通用消息：`log.Error("新增用户失败", zap.Error(err)); return errors.New("操作失败")`
- 非关键错误用 `_ =` 忽略

### Model pattern

- 所有 model 嵌入 `Base` struct（含 `Id` 自增主键 + `Code` UUID 业务主键 + `AllowEdit/AllowDel` + `Version` + timestamps）
- `Code` 通过 `BeforeCreate` hook 自动生成 UUID
- 自定义 GORM 类型：`ArrayStr`（JSON 数组）、`Map`（JSON 对象）

### Config pattern

- YAML 配置文件反序列化到 `global.Config`
- 不存在时自动生成默认配置
- 环境变量可覆盖（`GOSTC_` 前缀，`SCREAMING_SNAKE_CASE`）

### Bootstrap & init

Server 启动按固定顺序调用 `bootstrap.Init*()` 函数链。路由通过 `init()` 注册到 `bootstrap.Route` 函数变量。清理通过 `releaseFunc` 切片收集。

### Logging

- Server / Proxy：`go.uber.org/zap` + lumberjack 日志轮转，全局 `global.Logger`
- Client：自定义 `CircularLogger`（非 zap）
