# Go HTTP Leak Guard - HTTP连接泄漏防护项目

**[English](./README.en.md)** | **[中文](./README.md)**

---

## 愿景

`go-http-leak-guard` 旨在通过一个具体的示例，展示在Go中常见的HTTP连接泄漏问题，并提供一种自动化的解决方案，以确保代码的稳定性和性能。

在高并发场景下，不正确的HTTP客户端使用会导致文件描述符耗尽，最终导致服务崩溃。本项目不仅重现了这个问题，还引入了一套基于AI的规则，可以在开发阶段自动检测和修复此类问题。

## 问题复现

### 漏洞客户端 (`vulnerable`)

`examples/vulnerable/main.go` 中的客户端故意忽略了对HTTP响应体的读取，直接调用 `resp.Body.Close()`。这导致HTTP连接无法被复用，最终耗尽系统资源。

#### 修复前性能表现

* **CPU使用率**：持续上升，直至服务无响应
* **文件描述符**：迅速增加，直至达到系统上限
* **内存占用**：平稳，但连接无法释放

![修复前性能图](https://user-images.githubusercontent.com/12345678/placeholder-before.png)
*(占位符：此处应为修复前性能监控截图)*

### 防护客户端 (`guarded`)

`examples/guarded/main.go` 中的客户端在关闭响应体之前，使用 `io.Copy(io.Discard, resp.Body)` 来排空响应体。这确保了连接可以被正确回收和复用，从而保证了服务的稳定性。

#### 修复后性能表现

* **CPU使用率**：保持平稳
* **文件描述符**：稳定在较低水平
* **内存占用**：平稳

![修复后性能图](https://user-images.githubusercontent.com/12345678/placeholder-after.png)
*(占位符：此处应为修复后性能监控截图)*

## 如何使用AI规则自动防护

本项目包含一套AI规则，可以自动检测和修复HTTP连接泄漏问题。这些规则定义在 `.cursor/rules/http-stability.mdc` 中，并通过 `.windsurfrules` 文件加载。

### 规则核心逻辑

- **检测**：扫描Go代码中所有的 `http.Client` 调用。
- **强制规范**：确保在 `resp.Body.Close()` 之前，响应体被完整读取或排空。
- **自动修复**：当检测到不规范的代码时，IDE（如Cursor）会提供一键修复建议，自动插入 `io.Copy(io.Discard, resp.Body)`。

### 错误与正确示范

#### ❌ Incorrect Code (will be flagged by AI rules)
```go
resp, err := client.Get(url)
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close() // Warning: Response body not drained
```

#### ✅ Correct Code (complies with AI rules)
```go
// BEST PRACTICE: Use an immediately-invoked function literal for atomic resource handling
resp, err := client.Get(url)
if err != nil {
    log.Fatal(err)
}
func() {
    _, _ = io.Copy(io.Discard, resp.Body)
    resp.Body.Close()
}()
```

## 如何运行

你可以根据需要选择运行不同的客户端来进行对比。

1. **运行漏洞客户端**
   ```bash
   # 启动服务器和漏洞客户端
   docker-compose up --build server vulnerable-client
   ```
   - 观察 `vulnerable-client` 容器的日志，会发现文件描述符持续增加。

2. **运行防护客户端**
   ```bash
   # 启动服务器和防护客户端
   docker-compose up --build server guarded-client
   ```
   - 观察 `guarded-client` 容器的日志，会发现文件描述符数量保持稳定。

3. **同时运行所有服务**
   ```bash
   # 启动所有定义的服务
   docker-compose up --build
   ```

---

# Go HTTP Leak Guard - Project for Preventing HTTP Connection Leaks

**[English](./README.en.md)** | **[中文](./README.md)**

---

## Vision

`go-http-leak-guard` aims to demonstrate a common HTTP connection leak issue in Go through a concrete example and provide an automated solution to ensure code stability and performance.

In high-concurrency scenarios, incorrect use of HTTP clients can lead to file descriptor exhaustion, eventually causing service crashes. This project not only reproduces the issue but also introduces a set of AI-based rules that can automatically detect and fix such problems during the development phase.

## Reproducing the Issue

### Vulnerable Client (`vulnerable`)

The client in `examples/vulnerable/main.go` intentionally omits reading the HTTP response body and directly calls `resp.Body.Close()`. This prevents HTTP connections from being reused, eventually exhausting system resources.

#### Performance Before the Fix

* **CPU Usage**: Continuously increases until the service becomes unresponsive.
* **File Descriptors**: Rapidly increases until it reaches the system limit.
* **Memory Usage**: Stable, but connections are not released.

![Performance Before Fix](https://user-images.githubusercontent.com/12345678/placeholder-before.png)
*(Placeholder: Performance monitoring screenshot before the fix)*

### Guarded Client (`guarded`)

The client in `examples/guarded/main.go` drains the response body using `io.Copy(io.Discard, resp.Body)` before closing it. This ensures that the connection can be properly recycled and reused, thus guaranteeing service stability.

#### Performance After the Fix

* **CPU Usage**: Remains stable.
* **File Descriptors**: Stays at a low and stable level.
* **Memory Usage**: Stable.

![Performance After Fix](https://user-images.githubusercontent.com/12345678/placeholder-after.png)
*(Placeholder: Performance monitoring screenshot after the fix)*

## How to Use AI Rules for Automatic Protection

This project includes a set of AI rules that can automatically detect and fix HTTP connection leak issues. These rules are defined in `.cursor/rules/http-stability.mdc` and loaded via the `.windsurfrules` file.

### Core Logic of the Rules

- **Detection**: Scans for all `http.Client` calls in Go code.
- **Enforcement**: Ensures that the response body is fully read or drained before `resp.Body.Close()` is called.
- **Auto-Fix**: When non-compliant code is detected, the IDE (like Cursor) will provide a one-click fix suggestion to automatically insert `io.Copy(io.Discard, resp.Body)`.

### Incorrect vs. Correct Examples

#### ❌ Incorrect Code (will be flagged by AI rules)
```go
resp, err := client.Get(url)
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close() // Warning: Response body not drained
```

#### ✅ Correct Code (complies with AI rules)
```go
resp, err := client.Get(url)
if err != nil {
    log.Fatal(err)
}
defer func() {
    _, _ = io.Copy(io.Discard, resp.Body)
    resp.Body.Close()
}()
```

## How to Run

You can choose to run different clients for comparison as needed.

1. **Run the Vulnerable Client**
   ```bash
   # Start the server and vulnerable client
   docker-compose up --build server vulnerable-client
   ```
   - Check the logs of the `vulnerable-client` container, and you will see the number of file descriptors continuously increasing.

2. **Run the Guarded Client**
   ```bash
   # Start the server and guarded client
   docker-compose up --build server guarded-client
   ```
   - Check the logs of the `guarded-client` container, and you will see the file descriptor count remains stable.

3. **Run All Services**
   ```bash
   # Start all defined services
   docker-compose up --build
   ```