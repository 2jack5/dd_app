# 测试目录

测试按模块维护：`system`、`common`、`access`、`catalog` 请求业务服务，`callback` 和 `order` 请求订单服务，`model` 是不请求 HTTP 的单元测试。

每个场景都是独立的顶层 `TestXxx`，可在 GoLand 或 VS Code 中直接点击运行，不依赖其他 case 的顺序。

## 运行

```bash
# 终端一：业务服务
go run ./cmd/business

# 终端二：订单服务
go run ./cmd/order

# 终端三：全部真实接口用例
go test -v ./test/...

# 单个模块或 case
go test -v ./test/order
go test -v ./test/system -run '^TestHealthEndpoint$'
```

默认业务地址为 `http://127.0.0.1:8000`，默认订单地址为 `http://127.0.0.1:8001`。可覆盖：

```bash
DD_TEST_BUSINESS_BASE_URL=http://127.0.0.1:9000 \
DD_TEST_ORDER_BASE_URL=http://127.0.0.1:9001 \
go test -v ./test/...
```

兼容旧的 `DD_TEST_BASE_URL`，它只覆盖业务服务地址。IDE 运行配置使用同样的环境变量。

## 编写 case

- 业务服务：`testkit.Request`、`RequestJSON`、`RequestWithHeaders`。
- 订单服务：`testkit.RequestOrder`、`RequestOrderJSON`、`RequestOrderWithHeaders`。
- `testkit.Require` 同时断言 HTTP 状态与业务码。
- 测试客户端统一校验 JSON、`request_id` 和 `X-Request-ID`。

连接失败时会明确提示需要启动哪个服务，不会静默改用内存 Router。
