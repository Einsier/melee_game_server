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

const HeroMaxHealth = 20  //最大血量
const HeroInitHealth = 20 //初始血量

const (
	HeroColliderX       = 0.6722848 //Ruby的碰撞宽度
	HeroColliderY       = 0.3633368 //Ruby的碰撞高度
	HeroColliderOffsetX = 0.0       //Ruby的碰撞偏移X
	HeroColliderOffsetY = 0.0       //Ruby的碰撞偏移Y
)

const HeroMoveSpeed = 0.008 //英雄移动速度,单位为m/ms,即8m/s

//Player的状态
const (
	PlayerNotRegisteredStatus = iota
	PlayerEnterGameStatus
	PlayerLeaveGameStatus
)
