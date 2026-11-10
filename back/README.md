# 东东电竞 APP 后端

使用 Kratos v3、GORM 和 PostgreSQL。HTTP 契约采用 Protobuf API-first 方式维护，当前 180 项契约已拆分为业务服务和订单服务两个可独立运行的进程。

## 微服务划分

| 服务 | 入口 | 默认地址 | 职责 |
|---|---|---|---|
| Business API | `cmd/business/main.go` | `127.0.0.1:8000` | 登录注册、用户、陪练、目录、内容和运营管理 |
| Order API | `cmd/order/main.go` | `127.0.0.1:8001` | 下单、支付、履约、退款、售后、结算和资金管理 |

两个进程都有独立的 Kratos HTTP Server、路由注册、访问策略、Service 实现和健康检查。当前迁移阶段仍读取同一份 PostgreSQL 配置和模型迁移；后续拆库时只需分别注入各自 Repository，不需要再调整 HTTP 边界。

```text
cmd/business ─ internal/app.RunBusiness ─ server.NewBusinessHTTPServer
             └ api/business/v1 ─ service.BusinessService ─ biz ─ data

cmd/order    ─ internal/app.RunOrder    ─ server.NewOrderHTTPServer
             └ api/order/v1    ─ service.OrderService    ─ biz ─ data
```

## API 文件划分

- `api/common/v1/common.proto`：两个服务共享的消息类型。
- `api/business/v1/public.proto`：健康检查与公共目录。
- `api/business/v1/account.proto`：注册、登录和用户账号。
- `api/business/v1/coach.proto`：陪练资料和非订单能力。
- `api/business/v1/admin.proto`：业务后台接口。
- `api/business/v1/callback.proto`：身份和内容回调。
- `api/order/v1/user.proto`：用户订单、支付、退款和售后。
- `api/order/v1/coach.proto`：陪练接单、履约和收入结算。
- `api/order/v1/admin.proto`：订单、退款、支付、结算和对账后台。
- `api/order/v1/callback.proto`：支付、退款和结算回调。
- `api/order/v1/system.proto`：订单服务健康检查。

生成的 `*.pb.go` 和 `*_http.pb.go` 不手工修改。业务实现也按角色拆到 `internal/service/business_*.go` 和 `internal/service/order_*.go`，不再存在超大的 `platform.proto` 或 `platform.go`。

## 运行

Go 1.25 或以上。在 `back` 目录分别启动：

```sh
go run ./cmd/business
go run ./cmd/order
```

可用环境变量覆盖地址：

```sh
export DD_BUSINESS_HTTP_ADDR='127.0.0.1:8000'
export DD_ORDER_HTTP_ADDR='127.0.0.1:8001'
```

PostgreSQL 的 `host`、`database`、`user`、`password` 为必填项。项目启动时校验配置、连接并探活数据库，再对 `model.AllModels()` 执行 GORM AutoMigrate；项目不会自动创建 PostgreSQL 数据库，也不会插入默认管理员。完整配置见 `config/config.yaml` 和 `config/.env.example`。

## 实现范围

当前已完成两个进程的启动、180 项生成路由、认证入口、统一响应、数据库健康检查和目录查询。未完成的业务函数明确返回 501；渠道回调在验签、幂等和入账接入前固定返回 503，不返回假成功。

## 生成与验证

修改 Proto 后执行：

```sh
buf lint
buf generate --path api/common/v1 --path api/business/v1 --path api/order/v1
```

常用验证命令：

```sh
go test ./...
go vet ./...
go build ./cmd/business ./cmd/order
```

`api/v1/routes.json` 和 `api/v1/openapi.json` 仍是全部 180 项接口的汇总契约。项目启动后的真实 HTTP 用例见 `test/README.md`。

## 分层约定

- `server` 只组装 Kratos Server、中间件和生成路由。
- `service` 实现生成的 Service 接口并编排用例。
- `biz` 保存领域用例和 Repository 接口。
- `data` 实现 GORM Repository。
- `internal/app` 负责配置、日志、数据库和 Kratos 生命周期。

日志统一使用标准库 `slog`。数据库连接生命周期归 `internal/app` 所有；路由层不创建或关闭数据库连接。

## 微信登录配置

`POST /api/v1/auth/wechat/login` 已实现小程序与移动应用登录。服务端使用 `client_type` 选择受信任的 AppID/Secret，不接受客户端提交 OpenID。请配置 `auth.token_secret`（至少 32 字节）以及对应的 `wechat.mini_*` 或 `wechat.open_*` 凭据。登录会创建或复用用户与微信身份，并把刷新令牌 SHA-256 摘要及会话写入 Redis；PostgreSQL 不保存登录会话或刷新令牌。Redis 不可用时登录失败并返回业务状态码 `5000`。

## 统一业务状态码

Business 与 Order 共用 `internal/ecode` 中的四位业务状态码。所有响应固定包含 `code`、`error_msg`、`data`、`request_id`：`1000` 表示成功，`5000` 表示系统内部错误。完整清单见 `docs/业务状态码.md`。
