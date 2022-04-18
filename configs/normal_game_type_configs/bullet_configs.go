package normal_game_type_configs

/**
*@Author Sly
*@Date 2022/1/15
*@Version 1.0
*@Description:子弹信息的配置
 */

const BulletSpeed = float32(0.005) //子弹的飞行速度,单位为 m/ns,也就是5m/s
const BulletDuration = 2000000000  //子弹的持续时间(过了持续时间子弹自动狗带),单位为ns

const (
	BulletColliderX = 0.65 //子弹碰撞面积X
	BulletColliderY = 0.27 //子弹碰撞面积Y
)
