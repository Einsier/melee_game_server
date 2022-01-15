package game_room

import (
	"melee_game_server/configs/game_type_config"
	t "melee_game_server/internal/game/game_type"
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
	bullets     map[int64]*t.Bullet //存放本局游戏中的bullets
	oldBullets  []int64             //存放应该清理的Bullets
	newBullets  []int64             //存放新加入的Bullets
	refreshLock sync.Mutex          //新旧更替的lock
}

//AddBullets 将创建的bullet加入BulletsManager
func (bm *BulletsManager) AddBullets(b *t.Bullet) {
	bm.bullets[b.Id] = b
	bm.refreshLock.Lock()
	defer bm.refreshLock.Unlock()
	bm.newBullets = append(bm.newBullets, b.Id)
}

//RefreshBullets 定期删除需要清理的Bullets
func (bm *BulletsManager) RefreshBullets() {
	for {
		time.Sleep(game_type_config.BulletRefreshTime)
		bm.refreshLock.Lock()
		bm.oldBullets = bm.newBullets
		bm.newBullets = make([]int64, 0)
		bm.refreshLock.Unlock()

		//todo 由于是一个go程负责清除,所以不用加锁...?怎么更优化
		for _, id := range bm.oldBullets {
			delete(bm.bullets, id)
		}
	}
}

//DeleteBullets 删除子弹
func (bm *BulletsManager) DeleteBullets(id int64) {
	delete(bm.bullets, id)
}
