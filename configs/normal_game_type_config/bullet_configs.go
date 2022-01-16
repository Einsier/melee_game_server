package normal_game_type_config

/**
*@Author Sly
*@Date 2022/1/15
*@Version 1.0
*@Description:子弹信息的配置
 */

const BulletSpeed = 300.0       //子弹的飞行速度
const BulletDuration = 20000000 //子弹的持续时间(过了持续时间子弹自动狗带)
const BulletRefreshTime = 4     //子弹需要定时刷新的时间,单位为秒,过了这段时间子弹服务器认为已经失效

const (
	BulletColliderX       = 0.6485714 //子弹碰撞面积X
	BulletColliderY       = 0.2771429 //子弹碰撞面积Y
	BulletColliderOffsetX = 0.0       //
	BulletColliderOffsetY = 0.0
)
