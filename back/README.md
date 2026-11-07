# 东东电竞APP后端接口框架

使用Kratos v3 + GORM + PostgreSQL。按已修订的接口文档组织180个接口，保留55表模型及17组SQL迁移。

## 启动关系

```text
cmd/server/main.go
  └─ api.NewRouter(cfg, pg)                     internal/api/router.go
       ├─ 仅注册HTTP接口，接收main初始化的数据库连接
       ├─ 集中注册180个接口，按五个板块分组
       └─ 指向各板块文件中的独立处理函数
  └─ kratos.New(..., kratos.Server(httpServer)).Run()
```

`main.go`启动Kratos服务；所有GET/POST/PATCH/PUT等接口注册语句都集中在`internal/api/router.go`。`public.go`、`user.go`、`coach.go`、`admin.go`、`callback.go`仅保留各板块的接口函数、请求结构和相关辅助函数，不包含路由注册。

## 文件划分

|文件|内容|
|---|---|
|`cmd/server/main.go`|程序入口，调用api.NewRouter并启动Kratos|
|`internal/api/router.go`|读取启动配置、创建HTTP服务、集中注册全部180项接口，按五个板块分组|
|`internal/api/public.go`|公共基础、登录、首页、游戏目录、找陪练（22项）|
|`internal/api/user.go`|用户账户、通知、上传、订单、支付退款、评价售后（31项）|
|`internal/api/coach.go`|陪练入驻、认证、服务排期、履约、收益结算（34项）|
|`internal/api/admin.go`|后台权限、审核、商品运营、交易售后、财务统计（88项）|
|`internal/api/callback.go`|支付、退款、结算、实名、内容审核回调（5项）|
|`internal/infra/pg.go`|PostgreSQL/GORM初始化、连接池、超时探活与关闭|
|`internal/api/middleware.go`|请求号、恢复、身份/权限/幂等键入口检查|
|`internal/api/response.go`|统一响应和错误返回|
|`internal/api/router_test.go`|文档与路由一致、独立处理函数匹配、入口检查测试|
|`internal/model/`|已有55表模型|
|`migrations/`|已有分类SQL迁移，启动不自动执行|

上一轮生成的internal/biz、data、service、server、auth、conf框架文件已移除，不再通过这些层转发接口。

## 运行

Go版本1.25或以上。在back目录执行：

```sh
go run ./cmd/server
```

默认地址`127.0.0.1:8000`。`config/.env.example`列出环境变量；程序读取系统环境，不自动加载.env文件。未配置`DD_DATABASE_URL`时允许启动路由框架，`/healthz`返回200，数据库就绪检查和数据查询返回503。

连接自己的PostgreSQL时先完成migrations，再设置环境变量：

```sh
export DD_HTTP_ADDR='127.0.0.1:8000'
export DD_DATABASE_URL='postgres://dd_api:replace-me@127.0.0.1:5432/dongdong?sslmode=disable&TimeZone=UTC'
go run ./cmd/server
```

示例为占位凭据；生产TLS按实际环境配置。指定数据库却连接失败会启动失败，不自动建库、不插入默认管理员、不执行AutoMigrate。

## 接口实现范围

当前完成启动、180项显式路由、各接口独立处理函数和文档中的请求结构。健康/数据库检查与游戏、区服、段位、分类共6项有基础实现；目录查询直接在public.go中调用GORM并转换为公开DTO。

其余业务函数明确返回501。请求结构注明必填/可选，但规划接口尚未执行JSON绑定、字段校验和业务事务，不能把请求结构存在当作业务完成。后续就在对应板块函数中逐项实现。

需要登录的接口无凭据返回401；尚未配置可信身份适配器时带凭据返回503。后台动作检查权限；资源归属、数据范围、账号/陪练状态由实际业务继续校验。禁止通过Header直接指定可信用户ID或权限。幂等Header检查不等于数据库事务幂等已经完成。

渠道回调在真实验签/幂等/入账接入前固定503，不返回假成功，不修改资金。

## 文档与测试

- `docs/接口设计.md`：V0.2按板块业务设计及当前启动结构。
- `api/v1/openapi.json`：可导入Apifox的请求/响应契约，180项操作。
- `api/v1/routes.json`：与文档同步的操作清单，请求schema名称已更新。
- `docs/接口设计复查意见_V0.2.md`：历史复查问题；此前路由空文件问题已在此次重建解决。

```sh
go test ./...
go vet ./...
go build ./cmd/server
```

测试核对实际注册路径与OpenAPI、180条路由各自的处理函数、未实现接口不假成功、身份权限/幂等入口及分页边界。未连接真实PostgreSQL或支付渠道，尚未验证真实订单、退款、出款事务与并发。

数据库实现统一位于 `internal/infra/pg.go`；`main.go`的“2.基础设施”调用 `infra.NewPG`，同处注册defer调用 `infra.ClosePG`。路由只接收连接。`api/database.go` 已删除；API中只保留请求级连接获取和HTTP错误映射。


## 统一日志与配置

日志统一使用Go标准库slog，初始化函数仍在 `internal/infra/logger.go`。`cmd/server/main.go`先读取 `config/config.yaml`，再执行 `logger.Init(cfg.Log.Level, cfg.Log.Format)`，把 `logger.L()` 传给Kratos。HTTP请求日志、infra辅助日志及GORM日志使用同一个默认实例，不再使用Zap或Sugar接口。

`config/config.yaml`的参数均留空，只保留字段与说明供自行填写。空日志配置默认info/json；空数据库地址只运行接口框架；HTTP留空使用默认本地地址。数值与时长请按注释类型填写。程序现优先查找config目录，兼容旧configs目录。

main将加载的配置传给 `api.NewRouter(cfg, pg)`，HTTP端口来自这份配置；PostgreSQL参数由main.go使用并初始化连接；保留原DD环境变量覆盖以兼容既有运行方式。NX环境变量经配置加载器覆盖YAML。Redis仍未接入初始化，保留配置项供后续使用。

数据库生命周期归main.go所有：初始化成功即defer关闭，路由构建失败或Kratos停止后都释放连接池；router.go不初始化、不关闭PostgreSQL。连接串留空时允许以无数据库模式运行接口框架。
