### 说明
游戏基础玩法的游戏逻辑+游戏房间布置


玩法包括:
* 根据地图地形(包括树木,池塘,房子,陷阱等)2D移动+近战攻击
* 拾取地图上的增益道具
* 拾取到子弹包后可以发射子弹
* 玩家血量管理

游戏运行逻辑:

玩家从匹配服务器处进行排队((mvp1+待完善))

匹配服务器凑齐足够的人数,通知game_server开启一个game_room,同时这个game_room为等待全部玩家进入状态(mvp1+待完善)

game_server开启game_room,并且把game_room的id返回给匹配服务器(mvp1+待完善)

匹配服务器将game_server和game_room的编号发给玩家(mvp1+待完善)

玩家使用从匹配服务器处得到game_room编号,加入对应的game_room(mvp1版本同一时刻仅有一个game_room)

game_room等待全部玩家进入(所有configs,
