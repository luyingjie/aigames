# 🃏 斗地主游戏系统

一个基于 Go + Vue.js 开发的完整斗地主游戏系统，支持实时多人在线对战。

## 🎮 项目概述

本项目是一个全栈的斗地主游戏应用，包含用户认证、房间管理、实时游戏逻辑等完整功能。使用现代化的技术栈实现高性能的实时游戏体验。

### ✨ 主要特性

- 🎯 **完整的斗地主游戏逻辑**：叫地主、出牌、过牌等完整规则
- 🔒 **用户认证系统**：安全的注册登录机制，密码哈希存储
- 🏠 **房间管理系统**：创建/加入房间，支持公开和私人房间
- ⚡ **实时通信**：基于 WebSocket 的实时游戏状态同步
- 📱 **响应式界面**：现代化的 Vue.js 3 前端界面
- 💾 **数据持久化**：使用 BoltDB 进行本地数据存储

## 🚀 快速开始

### 环境要求

- Go 1.19+
- 现代浏览器（支持 WebSocket）

### 安装运行

1. **克隆项目**
   ```bash
   git clone <repository-url>
   cd aigames
   ```

2. **安装依赖**
   ```bash
   go mod tidy
   ```

3. **启动服务器**
   ```bash
   go run main.go
   ```

4. **访问游戏**
   - 打开浏览器访问：http://localhost:8080
   - 注册新用户或使用现有账号登录
   - 创建房间或加入现有房间开始游戏

### 服务端口

- **Web 服务**：http://localhost:8080
- **WebSocket 服务**：ws://localhost:3250/nano
- **数据库文件**：`data/game.db`

## 🏗️ 技术架构

### 后端技术栈

- **框架**：[nano](https://github.com/lonng/nano) - 高性能 Go WebSocket 游戏框架
- **数据库**：[BoltDB](https://github.com/etcd-io/bbolt) - 嵌入式键值数据库
- **架构模式**：依赖注入、服务化架构
- **通信协议**：WebSocket + JSON

### 前端技术栈

- **框架**：Vue.js 3 (Composition API)
- **样式**：原生 CSS + 响应式设计
- **通信**：WebSocket 客户端
- **构建**：无构建工具，直接浏览器运行

### 项目结构

```
aigames/
├── internal/               # 内部业务逻辑
│   ├── config/            # 配置管理
│   ├── database/          # 数据库操作
│   ├── handlers/          # WebSocket 处理器
│   │   ├── user.go        # 用户认证处理
│   │   ├── room.go        # 房间管理处理
│   │   └── game.go        # 游戏逻辑处理
│   ├── models/            # 数据模型
│   │   ├── user.go        # 用户模型
│   │   ├── room.go        # 房间模型
│   │   ├── card.go        # 扑克牌模型
│   │   ├── game.go        # 游戏状态模型
│   │   └── game_logic.go  # 游戏逻辑实现
│   └── services/          # 业务服务层
│       ├── user.go        # 用户服务
│       ├── room.go        # 房间服务
│       └── game.go        # 游戏服务
├── pkg/                   # 公共包
│   ├── logger/            # 日志工具
│   └── protocol/          # 通信协议
├── web/                   # 前端资源
│   ├── index.html         # 主页面
│   └── js/                # JavaScript 资源
├── data/                  # 数据目录
├── logs/                  # 日志目录
└── main.go               # 程序入口
```

## 🎲 游戏功能

### 用户系统

- **注册**：用户名、密码、年龄
- **登录**：安全认证，会话管理
- **密码加密**：SHA-256 哈希存储

### 房间系统

- **创建房间**：公开/私人房间选择
- **加入房间**：支持密码保护
- **准备状态**：玩家准备机制
- **房间列表**：实时更新的房间信息

### 游戏逻辑

#### 基本流程
1. **房间准备**：3名玩家全部准备后可开始
2. **发牌阶段**：每人17张牌，3张底牌
3. **叫地主**：轮流决定是否叫地主
4. **游戏阶段**：出牌和过牌操作

#### 支持的牌型
- **单牌**：任意单张牌
- **对子**：两张相同点数的牌
- **三张**：三张相同点数的牌
- **三带一**：三张+一张
- **三带二**：三张+一对
- **单顺**：连续的单牌（至少5张）
- **双顺**：连续的对子（至少3对）
- **飞机**：连续的三张
- **炸弹**：四张相同点数的牌
- **王炸**：大王+小王

#### 特殊规则
- **牌型比较**：按斗地主标准规则
- **炸弹规则**：炸弹可以压制其他牌型
- **回合制**：按顺序出牌，支持过牌

## 🔧 API 接口

### 用户接口

```javascript
// 用户注册
nano.request('user.Signup', {
    name: "用户名",
    password: "密码",
    age: 18
})

// 用户登录
nano.request('user.Login', {
    name: "用户名",
    password: "密码"
})
```

### 房间接口

```javascript
// 获取房间列表
nano.request('room.GetRoomList', {
    page: 1,
    size: 50,
    type: 0  // 0=公开房间, 1=私人房间
})

// 创建房间
nano.request('room.CreateRoom', {
    name: "房间名称",
    type: 0,  // 0=公开, 1=私人
    password: "密码"  // 私人房间密码
})

// 加入房间
nano.request('room.JoinRoom', {
    room_id: "房间ID",
    password: "密码"
})

// 设置准备状态
nano.request('room.SetReady', {
    room_id: "房间ID",
    ready: true
})
```

### 游戏接口

```javascript
// 获取游戏状态
nano.request('game.GetGameState', {
    room_id: "房间ID"
})

// 获取玩家手牌
nano.request('game.GetPlayerHand', {
    room_id: "房间ID"
})

// 叫地主
nano.request('game.CallLandlord', {
    room_id: "房间ID",
    call: true  // true=叫地主, false=不叫
})

// 出牌
nano.request('game.PlayCards', {
    room_id: "房间ID",
    cards: [...]  // 要出的牌
})

// 过牌
nano.request('game.PassTurn', {
    room_id: "房间ID"
})
```

## 🎨 界面预览

### 登录注册界面
- 简洁的用户认证表单
- 注册/登录模式切换
- 实时错误提示和成功反馈

### 房间列表界面
- 网格布局的房间卡片
- 房间状态显示（空闲/等待/游戏中）
- 创建房间和刷新功能

### 游戏界面
- 圆形游戏桌面设计
- 三个玩家位置显示
- 实时手牌展示和操作
- 游戏状态和回合提示

## 🛠️ 开发指南

### 添加新功能

1. **后端**：在 `internal/handlers/` 中添加新的处理器
2. **服务层**：在 `internal/services/` 中实现业务逻辑
3. **数据模型**：在 `internal/models/` 中定义数据结构
4. **协议定义**：在 `pkg/protocol/` 中定义请求/响应结构

### 数据库操作

```go
// 获取用户
user, err := userService.GetUser("username")

// 保存用户
err := userService.SaveUser(user)

// 创建房间
room, err := roomService.CreateRoom(id, name, owner, roomType, password)
```

### 配置管理

配置文件位置：`internal/config/config.go`

主要配置项：
- 服务器端口
- 数据库路径
- 日志级别
- 开发/生产模式

## 📝 日志系统

系统集成了完整的日志记录：

- **位置**：`logs/` 目录
- **级别**：DEBUG, INFO, WARN, ERROR, FATAL
- **格式**：时间戳 + 级别 + 消息
- **轮转**：按日期自动分割

## 🔐 安全考虑

- **密码加密**：使用 SHA-256 哈希存储
- **会话管理**：基于 WebSocket 会话状态
- **输入验证**：所有用户输入都经过验证
- **SQL 注入防护**：使用参数化查询

## 📊 性能特性

- **并发处理**：支持多用户同时在线
- **内存管理**：高效的游戏状态管理
- **网络优化**：WebSocket 长连接，减少延迟
- **数据库优化**：嵌入式 BoltDB，快速读写

## 🐛 故障排除

### 常见问题

1. **无法连接 WebSocket**
   - 检查端口 3250 是否被占用
   - 确认防火墙设置

2. **Web 页面无法访问**
   - 检查端口 8080 是否开放
   - 确认 `web/` 目录文件完整

3. **数据库错误**
   - 检查 `data/` 目录权限
   - 确认磁盘空间充足

### 调试模式

启用详细日志：
```bash
# 设置环境变量
export GAME_MODE=debug

# 或修改配置文件
// config.go
Mode: "debug"
```

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 📜 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

- [nano](https://github.com/lonng/nano) - 优秀的 Go 游戏框架
- [BoltDB](https://github.com/etcd-io/bbolt) - 高性能嵌入式数据库
- [Vue.js](https://vuejs.org/) - 渐进式 JavaScript 框架

## 📞 联系方式

如有问题或建议，请提交 Issue 或联系开发团队。

---

**开始你的斗地主之旅吧！** 🎉