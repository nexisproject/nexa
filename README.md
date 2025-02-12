# Nexa

## tips

### 拉取依赖仓库

```shell
# 使用ssh替换https
git config --global url."git@gitlab.liasica.com:".insteadof "https://gitlab.liasica.com/"

# 设置环境变量
go env -w GOPRIVATE="orba.plus"
go env -w GONOPROXY="orba.plus"
go env -w GONOSUMDB="orba.plus"

# 安装依赖
go get -u -v orba.plus/nexa
```

### 防止静态检查工具误报

```go
// 误报func
var _ = Setup

func Setup() {}

// 误报interface
var _ Hello = (*HelloImpl)(nil)

type Hello interface {
    World()
}

type HelloImpl struct{}
```

## 基本结构




### HAProxy
- [01 . HAProxy原理使用和配置](https://www.cnblogs.com/you-men/p/12979599.html)