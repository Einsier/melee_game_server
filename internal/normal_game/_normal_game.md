### 说明
游戏基础玩法的游戏逻辑+游戏房间布置


玩法包括:
* 根据地图地形(包括树木,池塘,房子,陷阱等)2D移动+近战攻击
* 拾取地图上的增益道具
* 拾取到子弹包后可以发射子弹,攻击其他玩家
* 玩家血量管理

游戏运行逻辑:

玩家从匹配服务器处进行排队((mvp1+待完善))

匹配服务器凑齐足够的人数,通知game_server开启一个game_room,将所有的playerId发给game_server(mvp1+待完善)

game_server使用一个端口开启一个game_room,game_room开始监听该端口并获取kcp包,初始等待全部玩家进入状态,并且把game_room的id返回给匹配服务器(mvp1+待完善)

匹配服务器将game_server和game_room的编号发给玩家(mvp1+待完善)

玩家使用从匹配服务器处得到game_room编号,加入对应的game_room(mvp1版本同一时刻仅有一个game_room)

game_room在玩家加入(玩家发送PlayerEnterGameRequest)时给玩家分配heroId,记录玩家的联系方式到net_manager中

game_room等待全部(configs.MaxNormalGamePlayerNum个)玩家进入,并且校验身份后开始游戏,具体操作是发送一个GameStartBroadcast给所有玩家

前端和后端进行四种格式的kcp通信,在game_room这一层里面实现对局和网络的解耦,即 game_room同时记录所有玩家的联系方式到net_manager中,用于通信,
并且内部存有一个game_manager对象,以对当前的游戏地图,玩家,道具等进行管理.在game_room这一层 使用事件驱动(从被分配的game_room端口监听kcp)模型,
当有kcp报文来时译码,并根据报文的内容选择相应的消息处理函数进行处理.

![image-20220116234555998](C:\Users\Administrator\Desktop\melee_game_server\internal\normal_game\image-20220116234555998.png)
