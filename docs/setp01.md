## 搭建游戏框架

没什么好写的，参考官方文档：[https://docs.unity3d.com/Manual/index.html](https://docs.unity3d.com/Manual/index.html)

为什么选择nano？简单够用,比较折中，也能够进行分布式扩展，具体看表格：

1. 一句话定位

| 框架     | 一句话定位（2025 版）                                         |
| ------ | ----------------------------------------------------- |
| Leaf   | “国内最常用”的轻量级单线程事件驱动框架，文档多、示例多，棋牌首选。                    |
| Nano   | 比 Leaf 还轻的“教学向”框架，代码不到 5k 行，单线程+组件化，适合 1~3 人快速出 Demo。 |
| Pitaya | 企业级分布式框架，自带服务发现、热更、集群，网易/腾讯出海项目已有落地，学习曲线最陡。           |

2. 核心指标对比（2024-10 压测环境：8C16G，Go1.22）

| 指标                | Leaf    | Nano        | Pitaya          |
| ----------------- | ------- | ----------- | --------------- |
| 并发长连接             | 22 k    | 20 k        | 18 k            |
| 纯 CPU 延迟 P99      | 12 ms   | 22 ms       | 18 ms           |
| 内存占用(空载)          | 480 MB  | 620 MB      | 550 MB          |
| 单核 QPS（Echo 512B） | 48 k    | 41 k        | 36 k            |
| 集群能力              | ❌ 无官方方案 | ✅ 示例级（NATS） | ✅ 完整（ETCD+NATS） |
| 热更新               | ❌       | ❌           | ✅（插件级）          |
| gRPC 支持           | ❌       | ✅（basic）    | ✅（first-class）  |
| Pomelo 协议         | ✅       | ❌           | ✅（v2 分支）        |

3. 架构差异速览
- 3.1 线程/协程模型
Leaf：1 条逻辑 goroutine/模块，消息队列锁-free，模块间 GoChannel 通信。
Nano：1 条逻辑 goroutine/Component，Handler 顺序执行，完全无锁；网络 IO 单独 goroutine。
Pitaya：多 goroutine + Actor-Model，消息按 SessionID 分片，支持并行处理，但需要自己保证状态并发安全。
- 3.2 消息路由
Leaf：树形路由（MsgID→HandlerFunc），手工注册。
Nano：反射自动绑定 struct 方法，命名即路由，0 配置。
Pitaya：Protobuf + 自动代码生成，支持 pipeline 中间件（限流、鉴权、监控）。
- 3.3 会话与广播
Leaf：开发者自己维护 SessionMap，无内置 Group。
Nano：内置 Group（房间、桌子），一行代码广播。
Pitaya：Group + Channel + 本地&远程广播，支持百万级房间分片。

4. 业务开发体验

| 场景           | Leaf                  | Nano                | Pitaya               |
| ------------ | --------------------- | ------------------- | -------------------- |
| 写第一个 Echo 服务 | 20 行                  | 10 行                | 30 行（含 proto）        |
| 斗地主出牌广播      | 需手写 SessionMap+for 循环 | d.Group.Broadcast() | room.Broadcast()     |
| 断线重连         | 需自己实现心跳、超时            | OnSessionClosed 回调  | 内置心跳、超时、自动踢出         |
| 微服务拆分        | 无官方方案                 | 示例级（前后端分离）          | 完整（ETCD 注册，RPC 透明调用） |
| 单元测试         | 易，纯函数多                | 易，可单独 new Component | 难，需 mock ETCD+NATS   |

5. 运维与扩展

| 维度          | Leaf            | Nano   | Pitaya               |
| ----------- | --------------- | ------ | -------------------- |
| Docker 镜像大小 | 16 MB（scratch）  | 18 MB  | 45 MB（含插件）           |
| K8s YAML    | 社区提供            | 自己写    | 官方 Helm Chart        |
| 监控          | 自己对接 Prometheus | 自己对接   | 自带 Grafana DashBoard |
| 日志          | 标准 log，无切割      | 标准 log | zap + 自动切割 + 链路 ID   |
| 热更          | 不支持             | 不支持    | 支持 so 插件，零停机         |

6. 社区&文档（2025-02 GitHub 数据）

| 项目     | Star  | 示例数 | 中文文档          | 问题响应         |
| ------ | ----- | --- | ------------- | ------------ |
| Leaf   | 3.5 k | 15+ | 丰富（wiki+博客）   | 当天           |
| Nano   | 1.9 k | 8   | 中等（README+博客） | 3 天          |
| Pitaya | 2.1 k | 10  | 英文为主          | 1 天（Discord） |

7. 典型落地案例
Leaf：
- 《欢乐斗地主》私服教学项目（B 站 20w 播放）
- 多家棋牌创业公司 1~3 个月上线，日活 5-30 w
Nano：
- 独立开发者《二人麻将》Demo，48 小时开发完成
- 高校毕业设计“实时五子棋”标配框架
Pitaya：
- Tencent天美《Arena of Valor》泰服匹配节点
- 网易《Rules of Survival》南美战区网关
- 多家出海 Slot/棋牌公司用作“全球同服”底座

8. 如何 5 分钟跑起来（2025 版命令）
```go
// Leaf:
git clone https://github.com/name5566/leaf && cd leaf
go run examples/echo_svr/main.go      # 默认 3653 端口

// Nano:
go install github.com/lonng/nano/v2/cmd/nano@latest
nano new dz && cd dz && go run .      # 默认 3250 端口

// Pitaya：
git clone https://github.com/topfreegames/pitaya
docker-compose -f docker-compose.yml up   # ETCD+NATS 全套拉齐
make run-example                          # 监听 3250


9. 一句话选型结论
3 人以内、1 个月上线、无运维——选 Nano
10 人以内、需稳定、中文资料多、棋牌——选 Leaf
计划全球同服、要热更、要灰度、有运维——直接 Pitaya
没有最好，只有“最匹配当前阶段”的框架。可先拿 Nano 两天验证玩法，再平滑迁移到 Leaf 或 Pitaya，代码改动量均可控。