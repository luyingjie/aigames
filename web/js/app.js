const { createApp, ref, computed, onMounted, onUnmounted } = Vue;

const app = createApp({
    setup() {
        // 响应式数据
        const currentView = ref('login'); // login, rooms, game
        const loading = ref(false);
        const error = ref('');
        const success = ref('');

        // 用户认证
        const isLogin = ref(true);
        const currentUser = ref(null);
        const authForm = ref({
            name: '',
            password: '',
            age: 18
        });

        // 房间相关
        const rooms = ref([]);
        const currentRoom = ref(null);
        const showCreateRoomModal = ref(false);
        const showJoinRoomModal = ref(false);
        const joinRoomPassword = ref('');
        const selectedRoom = ref(null);
        const createRoomForm = ref({
            name: '',
            type: '0',
            password: ''
        });

        // 游戏相关
        const gameState = ref(null);
        const playerHand = ref([]);
        const selectedCards = ref([]);
        const playerReady = ref(false);
        const nano = ref(null);
        const nanoInitialized = ref(false);

        // 计算属性
        const isMyTurn = computed(() => {
            if (!gameState.value || !currentUser.value) return false;
            const myPlayer = gameState.value.players?.find(p => p && p.username === currentUser.value.name);
            return myPlayer && myPlayer.position === gameState.value.current_turn;
        });

        const myPlayerPosition = computed(() => {
            if (!gameState.value || !currentUser.value) return -1;
            const myPlayer = gameState.value.players?.find(p => p && p.username === currentUser.value.name);
            return myPlayer ? myPlayer.position : -1;
        });

        const allPlayersReady = computed(() => {
            if (!gameState.value || !gameState.value.players) return false;
            const players = gameState.value.players.filter(p => p !== null);
            return players.length === 3 && players.every(p => p.is_ready);
        });

        // nano相关方法
        const initNano = () => {
            return new Promise((resolve, reject) => {
                // 如果已经初始化，直接返回
                if (nanoInitialized.value) {
                    resolve();
                    return;
                }
                
                if (!window.nano) {
                    reject(new Error('nano client not loaded'));
                    return;
                }
                
                nano.value = window.nano;

                nano.value.init({
                    host: "127.0.0.1",
                    port: 3250,
                    path: "/nano",
                    log: false,
                    user: {},
                    handshakeCallback: function() {
                        console.log('nano握手完成');
                    }
                }, function() {
                    console.log('nano连接成功');
                    nanoInitialized.value = true;
                    resolve();
                });
            });
        };

        const request = (route, data) => {
            return new Promise((resolve, reject) => {
                if (!nano.value || !nanoInitialized.value) {
                    reject(new Error('nano client not initialized'));
                    return;
                }
                nano.value.request(route, data, function(response) {
                    resolve(response);
                });
            });
        };

        // 认证相关方法
        const submitAuth = async () => {
            if (!authForm.value.name || !authForm.value.password) {
                error.value = '请填写完整信息';
                return;
            }

            loading.value = true;
            error.value = '';
            success.value = '';

            try {
                // 确保nano已经初始化
                await initNano();
                
                const route = isLogin.value ? 'user.Login' : 'user.Signup';
                const response = await request(route, authForm.value);

                if (response.code === 200) {
                    success.value = response.message;
                    if (isLogin.value) {
                        currentUser.value = response.data;
                        currentView.value = 'rooms';
                        await refreshRooms();
                    } else {
                        success.value = '注册成功，请登录';
                        isLogin.value = true;
                        authForm.value.password = '';
                    }
                } else {
                    error.value = response.message || '操作失败';
                }
            } catch (err) {
                error.value = '网络错误：' + err.message;
            } finally {
                loading.value = false;
            }
        };

        const toggleAuthMode = () => {
            isLogin.value = !isLogin.value;
            error.value = '';
            success.value = '';
            authForm.value.password = '';
        };

        const logout = () => {
            currentUser.value = null;
            currentView.value = 'login';
            currentRoom.value = null;
            gameState.value = null;
            playerHand.value = [];
        };

        // 房间相关方法
        const refreshRooms = async () => {
            loading.value = true;
            try {
                // 确保nano已经初始化
                await initNano();
                
                const response = await request('room.GetRoomList', {
                    page: 1,
                    size: 50,
                    type: 0
                });

                if (response.code === 200) {
                    rooms.value = response.data.rooms || [];
                } else {
                    error.value = response.message || '获取房间列表失败';
                }
            } catch (err) {
                error.value = '网络错误：' + err.message;
            } finally {
                loading.value = false;
            }
        };

        const createRoom = async () => {
            if (!createRoomForm.value.name) {
                error.value = '请输入房间名称';
                return;
            }

            loading.value = true;
            error.value = '';

            try {
                // 确保nano已经初始化
                await initNano();
                
                const response = await request('room.CreateRoom', {
                    name: createRoomForm.value.name,
                    type: parseInt(createRoomForm.value.type),
                    password: createRoomForm.value.password
                });

                if (response.code === 200) {
                    currentRoom.value = response.data;
                    currentView.value = 'game';
                    showCreateRoomModal.value = false;
                    createRoomForm.value = { name: '', type: '0', password: '' };
                    await getGameState();
                    startGameStatePolling();
                } else {
                    error.value = response.message || '创建房间失败';
                }
            } catch (err) {
                error.value = '网络错误：' + err.message;
            } finally {
                loading.value = false;
            }
        };

        const joinRoom = async (room) => {
            selectedRoom.value = room;
            if (room.has_password) {
                showJoinRoomModal.value = true;
                joinRoomPassword.value = '';
            } else {
                await doJoinRoom();
            }
        };

        const doJoinRoom = async () => {
            loading.value = true;
            error.value = '';

            try {
                // 确保nano已经初始化
                await initNano();
                
                const response = await request('room.JoinRoom', {
                    room_id: selectedRoom.value.id,
                    password: joinRoomPassword.value
                });

                if (response.code === 200) {
                    currentRoom.value = response.data;
                    currentView.value = 'game';
                    showJoinRoomModal.value = false;
                    await getGameState();
                    startGameStatePolling();
                } else {
                    error.value = response.message || '加入房间失败';
                }
            } catch (err) {
                error.value = '网络错误：' + err.message;
            } finally {
                loading.value = false;
            }
        };

        const leaveRoom = async () => {
            if (!currentRoom.value) return;

            try {
                // 确保nano已经初始化
                await initNano();
                
                await request('room.LeaveRoom', {
                    room_id: currentRoom.value.id
                });
                currentRoom.value = null;
                currentView.value = 'rooms';
                gameState.value = null;
                playerHand.value = [];
                stopGameStatePolling();
                await refreshRooms();
            } catch (err) {
                error.value = '离开房间失败：' + err.message;
            }
        };

        // 游戏相关方法
        const getGameState = async () => {
            if (!currentRoom.value) return;

            try {
                // 确保nano已经初始化
                await initNano();
                
                const response = await request('game.GetGameState', {
                    room_id: currentRoom.value.id
                });

                if (response.code === 200) {
                    gameState.value = response.data;

                    // 获取玩家手牌
                    const handResponse = await request('game.GetPlayerHand', {
                        room_id: currentRoom.value.id
                    });

                    if (handResponse.code === 200) {
                        playerHand.value = handResponse.data.cards || [];
                    }

                    // 检查玩家准备状态
                    const myPlayer = gameState.value.players?.find(p => p && p.username === currentUser.value.name);
                    if (myPlayer) {
                        playerReady.value = myPlayer.is_ready;
                    }
                }
            } catch (err) {
                console.error('获取游戏状态失败:', err);
            }
        };

        const toggleReady = async () => {
            try {
                // 确保nano已经初始化
                await initNano();
                
                const response = await request('room.SetReady', {
                    room_id: currentRoom.value.id,
                    ready: !playerReady.value
                });

                if (response.code === 200) {
                    playerReady.value = !playerReady.value;
                    await getGameState();
                } else {
                    error.value = response.message || '设置准备状态失败';
                }
            } catch (err) {
                error.value = '网络错误：' + err.message;
            }
        };

        const startGame = async () => {
            try {
                // 确保nano已经初始化
                await initNano();
                
                const response = await request('room.StartGame', {
                    room_id: currentRoom.value.id
                });

                if (response.code === 200) {
                    await getGameState();
                } else {
                    error.value = response.message || '开始游戏失败';
                }
            } catch (err) {
                error.value = '网络错误：' + err.message;
            }
        };

        const callLandlord = async (call) => {
            try {
                // 确保nano已经初始化
                await initNano();
                
                const response = await request('game.CallLandlord', {
                    room_id: currentRoom.value.id,
                    call: call
                });

                if (response.code === 200) {
                    await getGameState();
                } else {
                    error.value = response.message || '叫地主失败';
                }
            } catch (err) {
                error.value = '网络错误：' + err.message;
            }
        };

        const toggleCardSelection = (index) => {
            const pos = selectedCards.value.indexOf(index);
            if (pos > -1) {
                selectedCards.value.splice(pos, 1);
            } else {
                selectedCards.value.push(index);
            }
        };

        const playCards = async () => {
            if (selectedCards.value.length === 0) return;

            const cards = selectedCards.value.map(index => playerHand.value[index]);

            try {
                // 确保nano已经初始化
                await initNano();
                
                const response = await request('game.PlayCards', {
                    room_id: currentRoom.value.id,
                    cards: cards
                });

                if (response.code === 200) {
                    selectedCards.value = [];
                    await getGameState();
                } else {
                    error.value = response.message || '出牌失败';
                }
            } catch (err) {
                error.value = '网络错误：' + err.message;
            }
        };

        const passTurn = async () => {
            try {
                // 确保nano已经初始化
                await initNano();
                
                const response = await request('game.PassTurn', {
                    room_id: currentRoom.value.id
                });

                if (response.code === 200) {
                    await getGameState();
                } else {
                    error.value = response.message || '过牌失败';
                }
            } catch (err) {
                error.value = '网络错误：' + err.message;
            }
        };

        // 工具方法
        const formatCard = (card) => {
            const suits = { 1: '♠', 2: '♥', 3: '♦', 4: '♣', 5: '王' };
            const values = {
                3: '3', 4: '4', 5: '5', 6: '6', 7: '7', 8: '8', 9: '9', 10: '10',
                11: 'J', 12: 'Q', 13: 'K', 14: 'A', 15: '2',
                16: '小王', 17: '大王'
            };

            if (card.suit === 5) {
                return values[card.value];
            }
            return suits[card.suit] + values[card.value];
        };

        const getCardColor = (card) => {
            if (card.suit === 2 || card.suit === 3) { // 红桃，方块
                return 'red';
            }
            return 'black';
        };

        // 定时刷新游戏状态
        let gameStateInterval = null;
        const startGameStatePolling = () => {
            if (gameStateInterval) clearInterval(gameStateInterval);
            gameStateInterval = setInterval(() => {
                if (currentView.value === 'game' && currentRoom.value) {
                    getGameState();
                }
            }, 2000);
        };

        const stopGameStatePolling = () => {
            if (gameStateInterval) {
                clearInterval(gameStateInterval);
                gameStateInterval = null;
            }
        };

        // 生命周期钩子
        onMounted(async () => {
            console.log('应用已挂载');
            // 在应用挂载后初始化nano客户端
            try {
                await initNano();
                console.log('Nano客户端初始化完成');
            } catch (err) {
                console.error('Nano客户端初始化失败:', err);
            }
        });

        onUnmounted(() => {
            stopGameStatePolling();
        });

        // 暴露给模板的属性和方法
        return {
            currentView,
            loading,
            error,
            success,
            isLogin,
            currentUser,
            authForm,
            rooms,
            currentRoom,
            showCreateRoomModal,
            showJoinRoomModal,
            joinRoomPassword,
            createRoomForm,
            gameState,
            playerHand,
            selectedCards,
            playerReady,
            isMyTurn,
            myPlayerPosition,
            allPlayersReady,
            submitAuth,
            toggleAuthMode,
            logout,
            refreshRooms,
            createRoom,
            joinRoom,
            doJoinRoom,
            leaveRoom,
            toggleReady,
            startGame,
            callLandlord,
            toggleCardSelection,
            playCards,
            passTurn,
            formatCard,
            getCardColor
        };
    }
});

app.mount('#app');