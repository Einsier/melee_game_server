package game_room

/**
*@Author Sly
*@Date 2022/1/14
*@Version 1.0
*@Description:用于管理一次游戏的对局,包括玩家管理,对局中的prop(道具)管理以及子弹管理
 */

type GameRoom struct {
	hm HeroesManager
	pm PropManager
	bm BulletsManager
	nm NetManager
}
