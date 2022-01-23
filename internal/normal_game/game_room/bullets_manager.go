package game_room

import (
	"melee_game_server/configs/normal_game_type_configs"
	gt "melee_game_server/internal/normal_game/game_type"
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

//AddBullets 将创建的bullet加入BulletsManager
func (bm *BulletsManager) AddBullets(b *gt.Bullet) {
	bm.bullets.Store(b.Id, b)
	bm.refreshLock.Lock()
	defer bm.refreshLock.Unlock()
	bm.newBullets = append(bm.newBullets, b.Id)
}

//RefreshBullets 定期删除需要清理的Bullets
//todo 改成注册到TimeEvent中的形式
func (bm *BulletsManager) RefreshBullets() {
	for {
		time.Sleep(normal_game_type_configs.BulletRefreshTime)
		bm.refreshLock.Lock()
		bm.oldBullets = bm.newBullets
		bm.newBullets = make([]int64, 0)
		bm.refreshLock.Unlock()

		for _, id := range bm.oldBullets {
			bm.bullets.Delete(id)
		}
	}
}

//DeleteBullets 删除子弹
func (bm *BulletsManager) DeleteBullets(bid int64) {
	bm.bullets.Delete(bid)
}

func (bm *BulletsManager) CheckBulletHitHero(bid int64, heroPosition gt.Vector2) bool {
	b, ok := bm.bullets.Load(bid)
	if !ok {
		return false
	}
	return b.(*gt.Bullet).IsBulletHitHero(heroPosition)
}
