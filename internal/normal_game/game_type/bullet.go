package game_type

import (
	config "melee_game_server/configs/normal_game_type_configs"
	"melee_game_server/framework/entity"
)

/**
*@Author Sly
*@Date 2022/1/14
*@Version 1.0
*@Description:子弹有关逻辑
 */

type Bullet struct {
	Id             int64          //本局游戏的独一无二的id
	HeroId         int32          //发射的Hero的id
	BulletIdByHero int32          //某个Hero发射的第几个子弹(由前端提供)
	CreateTime     int64          //客户端记录的发射的时间
	DestroyTime    int64          //应该销毁的时间(防止占用过多内存)
	CreatePosition entity.Vector2 //客户端记录的发射的位置
	Direction      entity.Vector2 //客户端记录的发射的方向
}

//NewBullet 创建一个bullet
func NewBullet(heroId, bulletIdByHero int32, createTime int64, createPosition, direction entity.Vector2) *Bullet {
	return &Bullet{
		HeroId:         heroId,
		BulletIdByHero: bulletIdByHero,
		Id:             CountBulletId(heroId, bulletIdByHero),
		CreateTime:     createTime,
		DestroyTime:    createTime + config.BulletDuration,
		CreatePosition: createPosition,
		Direction:      direction,
	}
}

//GetPosition 通过创建的时间,创建时的方向,飞行的速度算出子弹大概的位置
func (b *Bullet) GetPosition(t int64) *entity.Vector2 {
	duration := t - b.CreateTime
	//Direction中的x,y至多只能有一个不为1
	//todo 待优化...这里有更高效的算法吗
	if b.Direction.X != 0 {
		return &entity.Vector2{
			//运算逻辑:子弹初始位置x + 子弹发射方向的x(为1/-1) * 子弹速度 * 子弹飞行时间
			X: b.CreatePosition.X + b.Direction.X*config.BulletSpeed*float64(duration),
			Y: b.CreatePosition.Y,
		}
	}
	return &entity.Vector2{
		X: b.CreatePosition.X,
		Y: b.CreatePosition.Y + b.Direction.Y*config.BulletSpeed*float64(duration),
	}
}

//CountBulletId 通过hero的id以及本子弹是用户发射的第几个子弹,计算出独一无二的本局游戏的子弹
func CountBulletId(heroId, BulletIdByHero int32) int64 {
	return int64(heroId)<<32 + int64(BulletIdByHero)
}

//IsBulletHitHero todo 子弹是否击中英雄,heroMid英雄的中心坐标
func (b *Bullet) IsBulletHitHero(heroMid entity.Vector2) bool {
	return true
}
