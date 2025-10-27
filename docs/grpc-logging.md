# gRPC 日志记录

## 概述

Nexa 的 `kit/micro` 包提供了基于 zap 的 gRPC 请求日志记录功能，可以自动记录所有 gRPC 请求的详细信息。

## 功能特性

日志中间件会自动记录以下信息：

- **kind**: 传输类型（gRPC）
- **operation**: 调用的方法名
- **code**: gRPC 状态码
- **duration**: 请求耗时
- **reason**: 错误原因（仅在出错时）
- **error**: 详细错误信息（仅在出错时）

## 使用方法

### 1. 基本使用

`LoggingMiddleware` 已经默认集成到 `micro.Run` 函数中，无需额外配置：

```go
package main

import (
    "nexis.run/nexa/kit/micro"
    pb "your-project/api/proto"
)

func main() {
    server, errCh := micro.Run("my-service", ":9000", func(s *grpc.Server) {
        pb.RegisterYourServiceServer(s, &YourServiceImpl{})
    })
    
    // 等待错误或信号
    if err := <-errCh; err != nil {
        panic(err)
    }
}
```

### 2. 日志输出示例

#### 成功的请求

```
2025-10-27T10:15:30.123+0800    INFO    gRPC request completed    {
    "kind": "gRPC",
    "operation": "/api.v1.UserService/GetUser",
    "code": 0,
    "duration": "15.234ms"
}
```

#### 失败的请求

```
2025-10-27T10:15:35.456+0800    ERROR   gRPC request failed    {
    "kind": "gRPC",
    "operation": "/api.v1.UserService/CreateUser",
    "code": 3,
    "duration": "8.567ms",
    "reason": "invalid argument: user name cannot be empty",
    "error": "rpc error: code = InvalidArgument desc = invalid argument: user name cannot be empty"
}
```

### 3. 配置日志级别

可以通过配置 logger 来控制日志级别和输出目标：

```go
package main

import (
    "nexis.run/nexa/kit/configure"
    "nexis.run/nexa/kit/logger"
    "nexis.run/nexa/kit/micro"
)

func main() {
    // 配置日志
    logger.Setup(&configure.Logger{
        Name:   "my-grpc-service",
        Stdout: true,  // 输出到控制台
        Kafka: &configure.Kafka{
            Brokers: []string{"localhost:9092"},
            Topic:   "logs",
        },
    })
    
    // 启动 gRPC 服务
    server, errCh := micro.Run("my-service", ":9000", func(s *grpc.Server) {
        // 注册你的服务
    })
    
    if err := <-errCh; err != nil {
        panic(err)
    }
}
```

### 4. 自定义日志中间件

如果需要自定义日志行为，可以创建自己的中间件：

```go
package main

import (
    "context"
    "time"
    
    "github.com/go-kratos/kratos/v2/middleware"
    "github.com/go-kratos/kratos/v2/transport"
    "go.uber.org/zap"
)

func CustomLoggingMiddleware() middleware.Middleware {
    return func(handler middleware.Handler) middleware.Handler {
        return func(ctx context.Context, req interface{}) (interface{}, error) {
            startTime := time.Now()
            
            // 获取请求信息
            if info, ok := transport.FromServerContext(ctx); ok {
                zap.L().Debug("gRPC request started",
                    zap.String("operation", info.Operation()),
                    zap.Any("request", req),
                )
            }
            
            // 执行请求
            reply, err := handler(ctx, req)
            
            duration := time.Since(startTime)
            
            if err != nil {
                zap.L().Error("gRPC request error",
                    zap.Duration("duration", duration),
                    zap.Error(err),
                )
            } else {
                zap.L().Debug("gRPC request success",
                    zap.Duration("duration", duration),
                    zap.Any("reply", reply),
                )
            }
            
            return reply, err
        }
    }
}
```

### 5. gRPC 状态码说明

常见的 gRPC 状态码：

- `0` - OK: 成功
- `1` - CANCELLED: 操作被取消
- `2` - UNKNOWN: 未知错误
- `3` - INVALID_ARGUMENT: 无效参数
- `4` - DEADLINE_EXCEEDED: 超时
- `5` - NOT_FOUND: 未找到
- `7` - PERMISSION_DENIED: 权限拒绝
- `13` - INTERNAL: 内部错误
- `14` - UNAVAILABLE: 服务不可用
- `16` - UNAUTHENTICATED: 未认证

## 性能考虑

- 日志中间件的性能开销很小，通常在微秒级别
- 建议在生产环境使用 INFO 级别，开发环境使用 DEBUG 级别
- 大量请求时考虑使用异步日志写入（Kafka）

## 调试技巧

### 启用详细日志

如果需要查看请求和响应的详细内容，可以在开发环境中启用 DEBUG 级别：

```go
logger.Setup(&configure.Logger{
    Name:   "my-service",
    Stdout: true,
    Level:  "debug",  // 启用 DEBUG 级别
})
```

### 添加自定义字段

在 gRPC 方法中添加额外的日志信息：

```go
func (s *YourServiceImpl) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
    zap.L().Debug("GetUser called",
        zap.String("user_id", req.UserId),
        zap.String("client_ip", getClientIP(ctx)),
    )
    
    // 处理逻辑...
    
    return &pb.GetUserResponse{}, nil
}
```

## 最佳实践

1. **使用结构化日志**: 始终使用 zap 的字段而不是格式化字符串
2. **避免记录敏感信息**: 不要在日志中记录密码、令牌等敏感数据
3. **合理的日志级别**: 
   - ERROR: 错误和异常
   - INFO: 重要的业务事件
   - DEBUG: 详细的调试信息
4. **性能监控**: 关注 duration 字段，识别慢请求
5. **集中式日志**: 使用 Kafka 等工具收集日志，便于分析和监控

