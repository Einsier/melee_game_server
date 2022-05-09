package normal_game_type_configs

/**
*@Author Sly
*@Date 2022/1/14
*@Version 1.0
*@Description:英雄信息的配置
 */

//当前英雄状态的枚举
const (
	HeroAlive = iota //英雄存活
	HeroDead         //英雄死亡
)

const HeroMaxHealth = 100  //最大血量
const HeroInitHealth = 100 //初始血量

const (
	HeroColliderX     = float32(1)        //Ruby的碰撞宽度
	HeroColliderY     = float32(0.5)      //Ruby的碰撞高度
	HeroColliderXHalf = HeroColliderX / 2 //碰撞宽度的一半
	HeroColliderYHalf = HeroColliderY / 2 //碰撞高度的一半
)

const HeroMoveSpeed = 0.008 //英雄移动速度,单位为m/ms,即8m/s

//Player的状态
const (
	PlayerNotRegisteredStatus = iota
	PlayerEnterGameStatus
	PlayerLeaveGameStatus
)
