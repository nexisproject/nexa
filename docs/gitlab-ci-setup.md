# GitLab CI 配置指南

本文档说明如何配置 GitLab CI 以自动发布到 GitHub。

## 前置要求

1. GitLab 项目（私有仓库）
2. GitHub 项目（公开仓库）
3. GitHub Personal Access Token
4. SSH 私钥（用于 Git 推送）

## 配置步骤

### 1. 创建 GitHub Personal Access Token

1. 登录 GitHub
2. 进入 Settings → Developer settings → Personal access tokens → Tokens (classic)
3. 点击 "Generate new token (classic)"
4. 填写信息：
   - **Note**: `GitLab CI for nexa` (或其他描述性名称)
   - **Expiration**: 建议选择 "No expiration" 或设置较长的过期时间
   - **Select scopes**: 勾选以下权限
     - ✅ `repo` (完整的仓库访问权限)
       - ✅ `repo:status`
       - ✅ `repo_deployment`
       - ✅ `public_repo`
       - ✅ `repo:invite`
       - ✅ `security_events`
5. 点击 "Generate token"
6. **重要**: 立即复制生成的 token（格式类似 `ghp_xxxxxxxxxxxxxxxxxxxx`），离开页面后将无法再次查看

### 2. 在 GitLab 配置环境变量

1. 进入 GitLab 项目：`https://gitlab.liasica.com/nexis/nexa`
2. 进入 Settings → CI/CD → Variables
3. 点击 "Add variable"
4. 添加以下变量：

#### GITHUB_TOKEN

- **Key**: `GITHUB_TOKEN`
- **Value**: 粘贴刚才创建的 GitHub Personal Access Token（如：`ghp_xxxxxxxxxxxxxxxxxxxx`）
- **Type**: Variable
- **Environment scope**: All (default)
- **Protect variable**: ✅ (推荐勾选，仅在受保护的分支/标签上可用)
- **Mask variable**: ✅ (推荐勾选，在日志中隐藏该值)
- **Expand variable reference**: ❌ (不勾选)

#### SSH_PRIVATE_KEY

- **Key**: `SSH_PRIVATE_KEY`
- **Value**: 粘贴 SSH 私钥内容（包括 `-----BEGIN OPENSSH PRIVATE KEY-----` 和 `-----END OPENSSH PRIVATE KEY-----`）
- **Type**: File (如果是文件路径) 或 Variable (如果是完整内容)
- **Environment scope**: All (default)
- **Protect variable**: ✅ (推荐勾选)
- **Mask variable**: ❌ (SSH 私钥太长无法 mask)
- **Expand variable reference**: ❌ (不勾选)

### 3. SSH 密钥配置

如果还没有为 GitHub 配置 SSH 密钥：

1. 生成 SSH 密钥对：
   ```bash
   ssh-keygen -t ed25519 -C "ci@nexis.run" -f ~/.ssh/github_ci
   ```

2. 将公钥添加到 GitHub：
   - 复制公钥内容：`cat ~/.ssh/github_ci.pub`
   - 进入 GitHub Settings → SSH and GPG keys → New SSH key
   - 粘贴公钥内容并保存

3. 将私钥内容复制到 GitLab CI/CD Variables：
   ```bash
   cat ~/.ssh/github_ci
   ```

### 4. 验证配置

提交代码到 `master` 分支后，GitLab CI 会自动：

1. **Build 阶段**: 编译多平台二进制文件
2. **Release 阶段**:
   - 验证 GitHub Token
   - 获取完整 Git 历史
   - 推送代码到 GitHub
   - 创建 Git Tag
   - 创建 GitHub Release
   - 上传编译产物

## 常见问题

### 401 Bad credentials

**原因**: GitHub Token 无效、过期或未设置

**解决方案**:
1. 检查 GitLab CI/CD Variables 中是否正确设置了 `GITHUB_TOKEN`
2. 确认 Token 没有过期
3. 确认 Token 具有 `repo` 权限
4. 尝试重新生成 Token

### remote rejected: failed

**原因**: 浅克隆导致 Git 对象缺失

**解决方案**: 已在 CI 配置中添加 `GIT_DEPTH: 0` 和 `git fetch --unshallow`

### SSH 连接失败

**原因**: SSH 私钥未正确配置

**解决方案**:
1. 确认 SSH 公钥已添加到 GitHub
2. 确认 GitLab CI/CD Variables 中的 SSH_PRIVATE_KEY 包含完整的私钥内容
3. 检查私钥格式是否正确（包括首尾的注释行）

### 找不到编译产物

**原因**: build 阶段的 artifacts 未正确传递

**解决方案**: 已在配置中添加 `artifacts` 和 `needs: [build]`

## CI 配置说明

### 关键配置项

```yaml
variables:
  GIT_DEPTH: 0  # 获取完整历史，避免推送失败

artifacts:
  paths:
    - bin/  # 保存编译产物
  expire_in: 1 hour  # 1小时后自动清理

needs: [build]  # release 阶段依赖 build 阶段的产物
```

### 工作流程

```
master 分支 push
    ↓
Build 阶段
    ├─ 编译多平台二进制文件
    └─ 保存到 artifacts
    ↓
Release 阶段
    ├─ 验证 GitHub Token
    ├─ 获取完整 Git 历史
    ├─ 推送代码到 GitHub
    ├─ 创建并推送 Tag
    ├─ 创建 GitHub Release
    └─ 上传 artifacts 中的二进制文件
```

## 安全建议

1. ✅ 使用 "Protect variable" 保护敏感变量
2. ✅ 使用 "Mask variable" 在日志中隐藏 Token
3. ✅ 定期更新 GitHub Token
4. ✅ 使用专用的 SSH 密钥对，不要使用个人密钥
5. ✅ 限制 Token 的最小权限范围
6. ✅ 在受保护的分支（如 master）上运行 CI

## 参考链接

- [GitHub Personal Access Tokens](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens)
- [GitHub REST API - Releases](https://docs.github.com/en/rest/releases/releases)
- [GitLab CI/CD Variables](https://docs.gitlab.com/ee/ci/variables/)
- [GitLab CI/CD SSH Keys](https://docs.gitlab.com/ee/ci/ssh_keys/)
