# Nexa



## Tips



### Commit 格式规范

> 参考文章 [Commit message 和 Change log 编写指南](https://www.ruanyifeng.com/blog/2016/01/commit_message_change_log.html)

```
[<type>](<scope>) <subject> (#pr)
docs：                   文档变动
fix：                    bug 修复
feat：                   新增功能
feat-wip：               开发中的功能，比如某功能的部分代码。
improvement：            原有功能的优化和改进
style：                  代码风格调整
typo：                   代码或文档勘误
refactor：               代码重构（不涉及功能变动）
performance/optimize：   性能优化
test：                   单元测试的添加或修复
chore：                  构建工具的修改
revert：                 回滚
deps：                   第三方依赖库的修改
community：              社区相关的修改，如修改 Github Issue 模板等。
```

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