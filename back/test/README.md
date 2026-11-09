# 测试目录

测试按业务模块放在独立目录中：

- `test/system/`：健康检查、就绪检查。
- `test/common/`：HTTP公共行为与错误响应。
- `test/access/`：用户、陪练和管理端的未登录访问控制。
- `test/catalog/`：游戏和服务目录接口。
- `test/callback/`：外部渠道未接入时保持关闭，不返回假成功。
- `test/model/`：模型值对象的单元测试，不请求HTTP接口。
- `test/testkit/`：各模块复用的本地HTTP测试客户端。

每个场景都是一个顶层 `TestXxx` 函数，可以在 GoLand 或 VS Code 中点击函数旁的运行按钮单独执行。每个 case 都自行构造完整请求，不依赖其他 case 的执行顺序。接口用例不创建内存 Router，而是直接请求已经启动的 `http://127.0.0.1:8000`，运行日志会显示请求地址、HTTP 状态、业务码和 `request_id`。

## 运行

```bash
# 启动项目（终端一）
go run ./cmd/server

# 运行全部本地接口用例（终端二）
go test -v ./test/...

# 按模块运行
go test -v ./test/system
go test -v ./test/catalog

# 单独运行一个case
go test -v ./test/system -run '^TestHealthEndpoint$'

# 只运行不需要启动服务的模型单元测试
go test -v ./test/model
```

默认地址可以通过环境变量覆盖：

```bash
DD_TEST_BASE_URL=http://127.0.0.1:9000 go test -v ./test/system
```

IDE 单独运行时也会读取 `DD_TEST_BASE_URL`。可在 GoLand 的 Go Test 运行配置或 VS Code 的测试环境变量中设置它。连接失败时，case 会直接提示先在 `back` 目录执行 `go run ./cmd/server`，不会静默改用内存服务。

## 编写 case

- 无请求体：`testkit.Request(t, method, path, nil)`。
- JSON 请求体：`testkit.RequestJSON(t, method, path, payload)`。
- 自定义 Header：`testkit.RequestWithHeaders(t, method, path, body, headers)`。
- 使用 `testkit.Require` 同时断言 HTTP 状态和统一响应业务码；测试客户端还会统一校验 JSON、`request_id` 及 `X-Request-ID`。

新增用例时放入对应模块目录，并保持“一个场景一个顶层测试函数”，不要用表驱动子测试合并互不相关的接口场景。
