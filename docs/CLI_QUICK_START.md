# NEXA CLI 快速入门指南

> 快速上手 NEXA 框架实用工具的构建、安装和使用

## 📦 安装

### 自动安装（推荐）

使用一键安装脚本自动下载并安装最新版本：

```bash
curl -fsSL https://raw.githubusercontent.com/nexisproject/nexa/master/install.sh | bash
```

安装脚本会：
- 自动检测操作系统和架构
- 下载对应的二进制文件
- 安装到 `$GOPATH/bin` 或 `/usr/local/bin`
- 自动检查更新

### 手动安装

1. 访问 [GitHub Releases](https://github.com/nexisproject/nexa/releases)
2. 下载对应平台的二进制文件
3. 解压并移动到 PATH 目录

### 验证安装

```bash
nexa --version
# 输出: nexa version 0.1.0.c39a3be (built at 2026-01-20T03:54:55+00:00)
```

## 🔧 本地开发

### 前置要求

- Go 1.25.3+
- Git
- Make

### 克隆仓库

```bash
git clone https://github.com/nexisproject/nexa.git
cd nexa
```

### 构建命令

#### 构建当前平台

```bash
# 指定版本号
VERSION=0.1.0 make build-darwin-arm64   # macOS ARM64
VERSION=0.1.0 make build-darwin-amd64   # macOS Intel
VERSION=0.1.0 make build-linux-amd64    # Linux AMD64
VERSION=0.1.0 make build-linux-arm64    # Linux ARM64
VERSION=0.1.0 make build-windows-amd64  # Windows
```

#### 构建所有平台

```bash
VERSION=0.1.0 make all
```

构建产物位于 `bin/` 目录：

```
bin/
├── nexa-darwin-amd64       # macOS Intel
├── nexa-darwin-arm64       # macOS Apple Silicon
├── nexa-linux-amd64        # Linux AMD64
├── nexa-linux-arm64        # Linux ARM64
└── nexa-windows-amd64.exe  # Windows
```

#### 清理构建

```bash
make clean
```

### 自定义构建参数

```bash
# 完全自定义版本信息
VERSION=0.1.0 HASH=abc1234 BUILD_TIME=2026-01-20T10:00:00+00:00 make all
```

参数说明：
- `VERSION`: 版本号（必需）
- `HASH`: Git 提交哈希（默认自动获取）
- `BUILD_TIME`: 构建时间（默认自动获取）

## 📌 版本格式说明

### 版本号组成

```
{major}.{minor}.{patch}.{hash}
例如: 0.1.0.c39a3be
```

| 部分 | 说明 | 示例 |
|------|------|------|
| `major` | 主版本号 | `0` |
| `minor` | 次版本号 | `1` |
| `patch` | 修订版本 | `0` |
| `hash` | Git 提交短哈希 | `c39a3be` |

### 版本输出格式

```bash
$ nexa --version
nexa version 0.1.0.c39a3be (built at 2026-01-20T03:54:55+00:00)
```

### 版本比较规则

版本比较优先级：**基础版本号 > Git Hash（字典序）**

```
0.1.0.abc123 < 0.2.0.def456  # 基础版本不同，比较主版本号
0.1.0.abc123 < 0.1.0.def456  # 基础版本相同，比较 Git Hash
0.1.0.abc123 = 0.1.0.abc123  # 完全相同
```

## 🚀 CI/CD 工作流程

当推送代码到 GitLab 的 `master` 分支时，会自动触发 CI Pipeline：

### Stage 1: Sync 🔄

**触发条件**: 每次 `master` 分支推送

**功能**:
- 同步代码到 GitHub
- 强制推送 master 分支
- 同步所有 tags

### Stage 2: Build 🛠️

**触发条件**: 每次 `master` 分支推送

**功能**:
- 编译多平台二进制文件
  - Linux (AMD64, ARM64)
  - macOS (AMD64, ARM64)
  - Windows (AMD64)
- 自动生成版本信息
- 保存构建产物

### Stage 3: Release 📦

**触发条件**: `cmd/nexa/**/*` 文件变更时

**功能**:
- 创建 Git Tag（格式: `0.1.0.hash`）
- 推送 Tag 到 GitHub
- 创建 GitHub Release
- 上传所有平台的二进制文件

## 📁 项目结构

```
nexa/
├── cmd/nexa/              # 命令行工具入口
│   ├── main.go           # 主程序（版本号定义）
│   └── internal/         # 内部包
├── kit/                  # 工具包
├── pkg/                  # 公共包
├── docs/                 # 文档
│   ├── CLI_QUICK_START.md       # 本文档
│   ├── version-format-update.md # 版本格式说明
│   ├── gitlab-ci-setup.md       # CI 配置指南
│   └── time-format-explained.md # 时间格式说明
├── bin/                  # 构建产物（git ignored）
├── Makefile             # 构建配置
├── .gitlab-ci.yml       # GitLab CI 配置
└── install.sh           # 一键安装脚本
```

## 🧪 测试

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行指定包的测试
go test ./kit/...

# 详细输出
go test -v ./...

# 测试覆盖率
go test -cover ./...
```

### 验证构建

```bash
# 构建后验证版本
./bin/nexa-darwin-arm64 --version

# 验证帮助信息
./bin/nexa-darwin-arm64 --help
```

## 🔍 常用命令

### 查看版本

```bash
nexa --version
```

### 查看帮助

```bash
nexa --help
nexa [command] --help
```

### 配置管理

```bash
# 查看配置
nexa config show

# 设置配置
nexa config set <key> <value>
```

### 创建新项目

```bash
# 创建新项目
nexa new <project-name>
```

## 📚 相关文档

| 文档 | 说明 |
|------|------|
| [README.md](../README.md) | 项目介绍和总览 |
| [version-format-update.md](version-format-update.md) | 版本格式详细说明 |
| [gitlab-ci-setup.md](gitlab-ci-setup.md) | GitLab CI 配置指南 |
| [time-format-explained.md](time-format-explained.md) | 时间格式说明 |

## 💡 提示和技巧

### 自动更新

安装脚本会自动检测已安装的版本，如果有新版本会提示更新：

```bash
curl -fsSL https://raw.githubusercontent.com/nexisproject/nexa/master/install.sh | bash
```

### 开发环境配置

推荐在 `~/.bashrc` 或 `~/.zshrc` 中添加：

```bash
export GOPATH="$HOME/go"
export PATH="$PATH:$GOPATH/bin"
```

### CI/CD 触发策略

- **每次推送**: Sync + Build
- **文件变更**: Release（仅当 `cmd/nexa/**/*` 变更）
- **手动触发**: 在 GitLab CI/CD 页面手动运行

### 版本号管理

修改版本号需要更新 `cmd/nexa/main.go`:

```go
var (
    Version   = "0.1.0"  // 修改这里
    BuildTime = ""
    Hash      = ""
)
```

## ❓ 常见问题

### Q: 如何强制重新安装？

```bash
curl -fsSL https://raw.githubusercontent.com/nexisproject/nexa/master/install.sh | bash -s -- --force
```

### Q: 如何安装特定版本？

从 GitHub Releases 页面手动下载对应版本的二进制文件。

### Q: 构建失败怎么办？

1. 确保 Go 版本 >= 1.25.3
2. 检查 `go.mod` 依赖是否完整
3. 运行 `go mod tidy` 更新依赖
4. 清理后重新构建：`make clean && VERSION=0.1.0 make all`

### Q: CI Pipeline 失败？

检查以下项：
- GitLab CI/CD 变量是否正确配置（`GITHUB_TOKEN`, `SSH_PRIVATE_KEY`）
- GitHub SSH 公钥是否已添加
- 网络连接是否正常

## 📞 获取帮助

- **Issues**: [GitHub Issues](https://github.com/nexisproject/nexa/issues)
- **Discussions**: [GitHub Discussions](https://github.com/nexisproject/nexa/discussions)
- **文档**: [docs/](.)

---

**最后更新**: 2026-01-20  
**当前版本**: 0.1.0.c39a3be  
**状态**: ✅ 生产就绪
