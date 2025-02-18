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
# 添加脚本（仅需添加一次）
cat <<'EOF' > /usr/local/bin/append-go-env
#!/bin/bash

function join() {
    local IFS="$1"
    shift
    echo "$*"
}

function append() {
  arr=()
  IFS=',' read -r -a arr <<< "$(go env "$1")"
  arr+=("$2")
  #arr=($(echo "${arr[@]}" | tr ' ' '\n' | sort -u | tr '\n' ' '))
  read -r -a arr <<< "$(echo "${arr[@]}" | tr ' ' '\n' | sort -u | tr '\n' ' ')"
  str=$(join , "${arr[@]}")
  go env -w "$1=$str"
  go env "$1"
}

append "$@"
EOF

chmod +x /usr/local/bin/append-go-env

# 使用ssh替换https（仅需设置一次）
git config --global url."git@gitlab.liasica.com:".insteadof "https://gitlab.liasica.com/"

# 设置环境变量（仅需设置一次）
append-go-env GOPRIVATE "nexis.run"
append-go-env GONOPROXY "nexis.run"
append-go-env GONOSUMDB "nexis.run"

# 安装依赖
go get -u -v nexis.run/nexa
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