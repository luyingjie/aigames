package constants

// 响应状态码
const (
	StatusSuccess      = 200 // 成功
	StatusBadRequest   = 400 // 请求错误
	StatusUnauthorized = 401 // 未授权
	StatusForbidden    = 403 // 禁止访问
	StatusNotFound     = 404 // 未找到
	StatusServerError  = 500 // 服务器错误
)

// 用户状态
const (
	UserStatusOffline = 0 // 离线
	UserStatusOnline  = 1 // 在线
	UserStatusInGame  = 2 // 游戏中
)

// 游戏阶段
const (
	GamePhaseWaiting  = 0 // 等待玩家
	GamePhaseBidding  = 1 // 叫地主阶段
	GamePhasePlaying  = 2 // 出牌阶段
	GamePhaseFinished = 3 // 游戏结束
)

// 玩家角色
const (
	RoleNone     = 0 // 未分配角色
	RoleLandlord = 1 // 地主
	RoleFarmer   = 2 // 农民
)

// 扑克牌花色
const (
	SuitSpades   = 1 // 黑桃 ♠
	SuitHearts   = 2 // 红桃 ♥
	SuitDiamonds = 3 // 方片 ♦
	SuitClubs    = 4 // 梅花 ♣
)

// 扑克牌点数
const (
	Rank3           = 3  // 3
	Rank4           = 4  // 4
	Rank5           = 5  // 5
	Rank6           = 6  // 6
	Rank7           = 7  // 7
	Rank8           = 8  // 8
	Rank9           = 9  // 9
	Rank10          = 10 // 10
	RankJ           = 11 // J
	RankQ           = 12 // Q
	RankK           = 13 // K
	RankA           = 14 // A
	Rank2           = 15 // 2
	RankLittleJoker = 16 // 小王
	RankBigJoker    = 17 // 大王
)

// 出牌类型
const (
	PlayTypeSingle            = 1  // 单张
	PlayTypePair              = 2  // 对子
	PlayTypeTriple            = 3  // 三张
	PlayTypeTriplePair        = 4  // 三带一对
	PlayTypeTripleSingle      = 5  // 三带一
	PlayTypeStraight          = 6  // 顺子(5张以上连续)
	PlayTypePairStraight      = 7  // 连对(3对以上连续)
	PlayTypeTripleStraight    = 8  // 飞机不带翅膀
	PlayTypeTriplePairPlane   = 9  // 飞机带对子
	PlayTypeTripleSinglePlane = 10 // 飞机带单张
	PlayTypeFourDualPair      = 11 // 四带两对
	PlayTypeFourDualSingle    = 12 // 四带两单
	PlayTypeBomb              = 13 // 炸弹
	PlayTypeJokerBomb         = 14 // 王炸
)

// WebSocket消息类型
const (
	// 系统消息
	MsgTypeSystem    = "system"
	ActionConnect    = "connect"    // 连接建立
	ActionDisconnect = "disconnect" // 连接断开
	ActionHeartbeat  = "heartbeat"  // 心跳检测

	// 房间相关
	MsgTypeRoom      = "room"
	ActionJoinRoom   = "join_room"   // 加入房间
	ActionLeaveRoom  = "leave_room"  // 离开房间
	ActionRoomUpdate = "room_update" // 房间状态更新

	// 游戏相关
	MsgTypeGame       = "game"
	ActionGameStart   = "game_start"   // 游戏开始
	ActionGameEnd     = "game_end"     // 游戏结束
	ActionBidLandlord = "bid_landlord" // 叫地主
	ActionPlayCards   = "play_cards"   // 出牌
	ActionPass        = "pass"         // 过牌
	ActionGameState   = "game_state"   // 游戏状态同步

	// 聊天相关
	MsgTypeChat          = "chat"
	ActionSendMessage    = "send_message"    // 发送消息
	ActionReceiveMessage = "receive_message" // 接收消息
)

// 数据库存储桶名称
const (
	BucketUsers     = "users"      // 用户数据
	BucketGames     = "games"      // 游戏记录
	BucketRooms     = "rooms"      // 房间数据
	BucketAIPlayers = "ai_players" // AI玩家
	BucketChats     = "chats"      // 聊天记录
	BucketConfigs   = "configs"    // 系统配置
)

// AI配置相关
const (
	AIPersonalityAggressive   = "aggressive"   // 激进型
	AIPersonalityConservative = "conservative" // 保守型
	AIPersonalityBalanced     = "balanced"     // 平衡型
)

// AI提示词类型
const (
	PromptGameStrategy = "game_strategy" // 游戏策略提示词
	PromptBidding      = "bidding"       // 叫地主提示词
	PromptPlaying      = "playing"       // 出牌提示词
	PromptChat         = "chat"          // 聊天提示词
	PromptPersonality  = "personality"   // 性格描述提示词
)

// 默认配置值
const (
	DefaultRoomCapacity   = 3    // 默认房间容量(斗地主3人)
	DefaultGameTimeout    = 1800 // 默认游戏超时时间(秒)
	DefaultBiddingTimeout = 30   // 默认叫地主超时时间(秒)
	DefaultPlayTimeout    = 60   // 默认出牌超时时间(秒)
	DefaultAIThinkTime    = 3    // 默认AI思考时间(秒)
	DefaultChatFrequency  = 0.3  // 默认AI聊天频率
)
