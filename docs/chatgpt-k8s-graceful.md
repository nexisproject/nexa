下面是几个在 Go & Kubernetes 环境中用于“平滑更新容器，等待进行中任务完成后才退出”推荐的方式和库：

---

### 1. 使用现成示例库：`integrii/go-k8s-graceful-termination`

这是一个小而实用的 GitHub 示例项目，实现了优雅关闭 Kubernetes Pod 的模式：

* 捕获 `SIGTERM`，并在 Pod 进入 Terminating 时继续服务正在接入的请求；
* 等待 `terminationGracePeriodSeconds` 内所有请求处理完毕，再退出；
* 配合 `preStop` hook 和 readiness probe，实现流量退下后再关闭进程 ([developer20.com][1], [linkerd.io][2], [reddit.com][3], [github.com][4])。

**优点**：直接借鉴其示例代码，并集成到你自己的服务中，快速实现 graceful shutdown。

**使用方式**：参考仓库的 README 和 Kubernetes YAML，设置 `terminationGracePeriodSeconds`, HTTP server shutdown, readiness probe 等。

---

### 2. 手写优雅退出方案（基于标准库）

如果你希望更灵活地控制“等待后台任务／Goroutine 完成”逻辑，可以采用 Go 标准实践：

```go
sigCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
defer stop()

var wg sync.WaitGroup
// 启动处理任务的 goroutine
wg.Add(1)
go func() {
  defer wg.Done()
  for {
    select {
    case task := <-taskCh:
      do(task)  // 你的任务处理
    case <-sigCtx.Done():
      return
    }
  }
}()

// HTTP server 启动...
go server.ListenAndServe()

<-sigCtx.Done()        // 收到终止信号
stopTraffic()         // 如关闭 readiness probe 或 stop new tasks
wg.Wait()             // 等待所有正在运行的任务完成
server.Shutdown(...)  // 等待 http handler 结束
```

* `sync.WaitGroup` 用于等待所有任务安全退出 ([developer20.com][1])；
* `signal.NotifyContext` 捕获 SIGTERM/SIGINT ([victoriametrics.com][5])；
* 配合 readiness probe 或 pre‑stop hook，使流量先退出再执行实际 shutdown。

---

### 3. Kubernetes PreStop + Readiness probe 配合

官方推荐在 Kubernetes Pod 中通过 `preStop` hook 加延迟，以及 readiness probe 配合，确保 Pod 排除流量后才关闭：

```yaml
lifecycle:
  preStop:
    exec:
      command: ["sleep", "10"]
readinessProbe:
  httpGet: ...
terminationGracePeriodSeconds: 30
```

结合 Go 中 `server.Shutdown()`，就能平滑退出正在处理的请求 ([victoriametrics.com][5], [kubernetes.io][6])。

---

### 总结建议

1. **快速起步**：直接 fork `integrii/go-k8s-graceful-termination`，它封装了完整流程；
2. **进阶需求（任务等待）**：使用 `sync.WaitGroup + signal.NotifyContext + server.Shutdown()` 模式；
3. **K8s 端配置必配项**：设置 `terminationGracePeriodSeconds`、preStop hook、以及 readiness probe。

---

如果你希望我帮你整理一个可直接使用的代码模版，或者推荐一个封装更完善的第三方库／框架，也可以再说，我可以深入挖掘！

[1]: https://developer20.com/golang-tips-and-trics-iii/?utm_source=chatgpt.com "Golang Tips & Tricks #3 - graceful shutdown - Developer 2.0"
[2]: https://linkerd.io/2-edge/tasks/graceful-shutdown/?utm_source=chatgpt.com "Graceful Pod Shutdown | Linkerd"
[3]: https://www.reddit.com/r/golang/comments/1e5d36v/terminating_elegantly_a_guide_to_graceful/?utm_source=chatgpt.com "Terminating Elegantly: A Guide to Graceful Shutdowns : r/golang"
[4]: https://github.com/integrii/go-k8s-graceful-termination?utm_source=chatgpt.com "integrii/go-k8s-graceful-termination: Kubernetes application ... - GitHub"
[5]: https://victoriametrics.com/blog/go-graceful-shutdown/?utm_source=chatgpt.com "Graceful Shutdown in Go: Practical Patterns - VictoriaMetrics"
[6]: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/?utm_source=chatgpt.com "Pod Lifecycle - Kubernetes"

---

### 参考
- [k8s~pod单副本的平滑部署](https://www.cnblogs.com/lori/p/18129433)
- [k8s~通过探针实现服务的平滑部署](https://www.cnblogs.com/lori/p/15623169.html)
- 