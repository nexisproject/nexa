在 Kubernetes 中进行容器平滑更新（如滚动更新或蓝绿部署）时，若需确保容器**等待所有进行中的任务完成后再销毁**，核心在于**优雅终止（Graceful Shutdown）与任务生命周期管理**。以下是结合 Golang 库和 Kubernetes 特性的完整方案：

---

### ⚙️ 一、核心机制：Kubernetes 优雅终止流程
当 Pod 被终止时（如更新触发的删除），Kubernetes 会：
1. **发送 `SIGTERM` 信号**：通知容器开始关闭。
2. **等待 `terminationGracePeriodSeconds`**（默认 30 秒）。
3. **强制终止（`SIGKILL`）**：若超时未退出。

要等待任务完成，需在 Golang 中捕获 `SIGTERM` 并延迟退出，直到任务结束。

---

### 💻 二、Golang 实现优雅关闭（无需额外库）
通过标准库即可实现任务跟踪与信号处理：
```go
package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// 示例：启动后台任务
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return // 收到终止信号时退出
			default:
				// 执行任务逻辑（如处理 HTTP 请求、队列消息等）
			}
		}
	}()

	// 捕获 SIGTERM 信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM)
	<-sigCh // 阻塞直到收到信号

	// 启动优雅关闭
	cancel()          // 通知所有任务停止接收新工作
	wg.Wait()         // 等待进行中任务完成
	time.Sleep(1 * time.Second) // 可选：额外清理时间
}
```

---

### 🔧 三、增强：Kubernetes 生命周期钩子
在 Pod 配置中添加 `preStop` 钩子，延长等待时间：
```yaml
containers:
- name: my-app
  image: my-app:v1
  lifecycle:
    preStop:
      exec:
        command: ["/bin/sh", "-c", "sleep 30"] # 留出更多时间等待任务完成
  terminationGracePeriodSeconds: 60 # 需大于 preStop 耗时
```

---

### 📚 四、推荐 Golang 库与工具
若需更高级的任务管理，可结合以下库：

1. **Kubernetes Client-GO**
    - **作用**：监听 Pod 终止事件，动态更新任务状态。
    - **场景**：需将任务状态上报给 Kubernetes（如通过 Annotation）。
    - 示例：
      ```go
      // 监听自身 Pod 的 DELETE 事件
      watcher, _ := clientset.CoreV1().Pods(namespace).Watch(ctx, metav1.ListOptions{FieldSelector: "metadata.name=" + podName})
      for event := range watcher.ResultChan() {
          if event.Type == watch.Deleted {
              // 触发优雅关闭逻辑
          }
      }
      ```

2. **Reloader**
    - **作用**：监控 ConfigMap/Secret 变更并触发滚动更新，但**不直接管理任务**。
    - **整合建议**：在更新前通过 API 通知应用“排空任务”（Draining），再执行 `kubectl rollout restart`。

3. **自定义 Operator 框架（如 Kubebuilder/OperatorSDK）**
    - **场景**：需为有状态任务（如批处理作业）设计自定义资源（CRD）。
    - **能力**：定义 `Task` CRD，在更新前检查 `.status.activeTasks == 0` 才允许删除 Pod。

---

### 🚀 五、完整方案建议
| **步骤**                | **工具/库**                  | **关键配置**                                  |
|-------------------------|------------------------------|----------------------------------------------|
| 1. 应用监听 `SIGTERM`   | Golang `os/signal`           | 结合 `sync.WaitGroup` 跟踪任务               |
| 2. 延长优雅关闭窗口     | Kubernetes `preStop` Hook    | `terminationGracePeriodSeconds` 需 > `preStop` 耗时 |
| 3. 状态上报（可选）     | client-go                    | Pod Annotations 标记任务数        |
| 4. 自动触发更新         | Reloader 或 Argo Rollouts    | 监控配置变化并触发更新 |

---

### 💎 结论
- **基础需求**：仅用 **Golang 标准库 + `preStop` 钩子** 即可实现等待任务完成。
- **高级需求**（如集群级任务调度）：通过 **client-go 监听事件** 或 **自定义 Operator** 扩展。
- **生产建议**：
    - 为关键任务添加超时机制（避免 `WaitGroup` 永久阻塞）。
    - 结合 `Readiness Probe` 在关闭时返回失败，确保流量不再进入。

参考实现示例：[Graceful Shutdown in Kubernetes](https://github.com/kubernetes/sample-controller/blob/master/docs/controller-client-go.md)。