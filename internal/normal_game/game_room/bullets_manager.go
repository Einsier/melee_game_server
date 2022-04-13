package game_room

import (
	"melee_game_server/framework/entity"
	gt "melee_game_server/internal/normal_game/game_type"
	"melee_game_server/plugins/logger"
	"sync"
	"time"
)

/**
*@Author Sly
*@Date 2022/1/14
*@Version 1.0
*@Description:用于管理游戏中产生的bullets
 */

type BulletsManager struct {
	bullets     sync.Map   //存放本局游戏中的bullets,key为int64,bulletId,value为*Bullet
	oldBullets  []int64    //存放应该清理的Bullets
	newBullets  []int64    //存放新加入的Bullets
	refreshLock sync.Mutex //新旧更替的lock
}

func NewBulletsManager() *BulletsManager {
	return &BulletsManager{
		bullets:     sync.Map{},
		oldBullets:  make([]int64, 0),
		newBullets:  make([]int64, 0),
		refreshLock: sync.Mutex{},
	}
}

//AddBullets 将创建的bullet加入BulletsManager
func (bm *BulletsManager) AddBullets(b *gt.Bullet) {
	bm.bullets.Store(b.Id, b)
	bm.refreshLock.Lock()
	defer bm.refreshLock.Unlock()
	bm.newBullets = append(bm.newBullets, b.Id)

	//todo 测试用,待删除
	hid, bid := countHidBid(b.Id)
	t := time.Now().Format(time.StampNano)
	logger.Testf("[%s]存放了hero[%d]的第[%d]颗子弹", t, hid, bid)
}

//DeleteBullets 删除子弹
func (bm *BulletsManager) DeleteBullets(bid int64) {
	bm.bullets.Delete(bid)
}

func (bm *BulletsManager) CheckBulletHitHero(bid int64, heroPosition entity.Vector2) bool {
	b, ok := bm.bullets.Load(bid)
	if !ok {
		return false
	}
	return b.(*gt.Bullet).IsBulletHitHero(heroPosition)
}

//countHidBid 根据子弹的64位的id,反向计算出发射子弹的英雄id和英雄发射的第几颗子弹
func countHidBid(bid int64) (hid, bIdByHero int32) {
	hid = int32(bid >> 32)
	bIdByHero = int32((bid << 32) >> 32)
	return
}
