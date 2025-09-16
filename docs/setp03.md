# 步骤 3: 实现 AI 玩家功能

本阶段的目标是为游戏添加简单的 AI 玩家。AI 逻辑非常简单：不叫地主，永远过牌。

## 设计思路

我们将采用基于 Goroutine 的异步模型来实现 AI 玩家，确保 AI 的行动不会阻塞游戏主线程。每个 AI 玩家将由一个独立的 Goroutine 控制。

1.  **扩展数据模型**: 区分普通玩家和 AI 玩家。
2.  **修改房间创建逻辑**: 允许在创建房间时指定加入的 AI 数量。
3.  **实现 AI 控制器**: 为每个 AI 启动一个独立的 Goroutine，负责监听游戏状态并自动执行操作。
4.  **整合游戏主流程**: 在游戏轮到 AI 玩家时，由游戏逻辑直接通知对应的 AI Goroutine，而不是等待网络消息。

## 具体实现步骤

### 1. 修改数据模型 (`internal/models/user.go`)

在 `User` 结构体中增加一个 `IsAI` 字段来标识该用户是否为 AI。

```go
// internal/models/user.go
type User struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    IsAI     bool   `json:"is_ai"` // 新增字段
    // ... 其他字段
}
```

### 2. 修改房间创建协议与逻辑

-   **协��层 (`pkg/protocol/user.go`)**: 在创建房间的请求消息 `C2SCreateRoom` 中，增加一个字段 `AICount`。
-   **接口层 (`internal/handlers/room.go`)**: 在处理创建房间请求的 Handler 中，读取 `AICount` 参数。
-   **服务层 (`internal/services/room.go`)**: 在创建房间的服务中，根据 `AICount` 的值，创建相应数量的 AI `User` 对象，并将其加入到房间中。AI 用户可以有简单的命名，如 "AI-1", "AI-2"。

### 3. 实现 AI 核心逻辑 (`internal/services/ai_controller.go`)

创建一个新的 `AIController` 服务来封装 AI 的行为。

-   **启动时机**: 当游戏在房间内开始时，`GameService` 会为每个 AI 玩家创建一个 `AIController` 实例并启动其主循环 Goroutine。
-   **通信方式**: 使用 Go 的 `channel` 实现游戏主逻辑与 AI Goroutine 之间的通信。当轮到 AI 操作时，`GameService` 向该 AI 对应的 `channel` 发送通知。
-   **AI 行为**:
    -   收到操作通知后，为了模拟“思考”，可以短暂延迟。
    -   **叫分阶段**: 自动发送“不叫”指令。
    -   **出牌阶段**: 自动发送“过牌”指令。

**AI Controller 伪代码:**
`go
package services

import (
    "time"
    "aigames/internal/models"
)

type AIController struct {
    player      *models.User
    gameService *GameService
    actionChan  chan bool
}

func NewAIController(player *models.User, gameService *GameService) *AIController {
    // ...
}

func (c *AIController) Start() {
    for range c.actionChan {
        time.Sleep(1 * time.Second) // 模拟思考
        // ... 根据游戏阶段执行“不叫”或“过牌”
    }
}

func (c *AIController) NotifyTurn() {
    c.actionChan <- true
}
`

### 4. 整合到游戏主流程 (`internal/services/game.go`)

`GameService` 需要维护一个 AI 控制器的映射表 (`map[string]*AIController`)。

-   **游戏开始时**: 遍历房间内所有玩家，如果是 AI，则创建并启动其 `AIController`。
-   **轮到玩家操作时**: 检查当前玩家是否为 AI。如果是，则通过 `channel` 通知其 `AIController` 行动；如果是真人玩家，则继续等待 WebSocket 消息。

`go
// In GameService, when it's a player's turn
currentPlayer := gs.GetCurrentPlayer()
if currentPlayer.IsAI {
    // 通知 AI Controller
    controller := gs.aiControllers[currentPlayer.ID]
    controller.NotifyTurn()
} else {
    // 等待真人玩家操作
}
`
