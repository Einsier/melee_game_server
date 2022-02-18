package game_room

import (
	pb "melee_game_server/api/client/proto"
	configs "melee_game_server/configs/normal_game_type_configs"
	gt "melee_game_server/internal/normal_game/game_type"
	"sync"
)

/**
*@Author Sly
*@Date 2022/1/14
*@Version 1.1
*@Description:一开始创建好玩家人数的Hero并放入HeroesManager中,以后对heroes的读就不需要加锁了,而是放到英雄这一层加锁.提高效率,
游戏结束也不会在heroes这一层删除英雄,而是等游戏结束之后统一的删除
*/

//HeroesManager 用于管理游戏中的hero
type HeroesManager struct {
	heroes map[int32]*gt.Hero
	//lock   	sync.RWMutex //用于对heroes增删加锁
	idLock sync.Mutex //用于分配hero的锁
	//一开始是1,如果分配完了应该等于configs.MaxNormalGamePlayerNum
	idCount int32
}

//HeroInitPosition 英雄的初始位置,key为heroId,value为英雄的初始位置 todo:随机生成英雄位置
var HeroInitPosition = map[int32]gt.Vector2{
	int32(1):  {X: 0, Y: 0},
	int32(2):  {X: 0, Y: 0},
	int32(3):  {X: 0, Y: 0},
	int32(4):  {X: 0, Y: 0},
	int32(5):  {X: 0, Y: 0},
	int32(6):  {X: 0, Y: 0},
	int32(7):  {X: 0, Y: 0},
	int32(8):  {X: 0, Y: 0},
	int32(9):  {X: 0, Y: 0},
	int32(10): {X: 0, Y: 0},
}

//NewHeroesManager 创建好游戏人数个数的英雄,设置其id分别为1~游戏玩家最大个数,并且提前放入到HeroesManager中
func NewHeroesManager() *HeroesManager {
	heroes := make(map[int32]*gt.Hero)
	for i := int32(1); i <= configs.MaxNormalGamePlayerNum; i++ {
		h := gt.NewHero(i, HeroInitPosition[i])
		heroes[i] = h
	}
	hm := HeroesManager{
		heroes:  heroes,
		idLock:  sync.Mutex{},
		idCount: 0,
	}
	return &hm
}

//ArrangeHeroId 用于取出一个当前未被player注册的英雄的id,如果当前所有的hero都被注册了,则返回-1
func (hm *HeroesManager) ArrangeHeroId() int32 {
	hm.idLock.Lock()
	hm.idLock.Unlock()
	if hm.idCount >= configs.MaxNormalGamePlayerNum {
		return -1
	}
	hm.idCount++
	return hm.idCount
}

//AddHero 向heroes中增加一个英雄
/*func (hm *HeroesManager) AddHero(h *gt.Hero) {
	hm.lock.Lock()
	defer hm.lock.Unlock()
	hm.heroes[h.Id] = h
}*/

func (hm *HeroesManager) GetHero(heroId int32) *gt.Hero {
	return hm.heroes[heroId]
}

//GetHeroPosition 获取heroId的英雄的当前位置
func (hm *HeroesManager) GetHeroPosition(heroId int32) gt.Vector2 {
	/*hm.lock.RLock()
	defer hm.lock.RUnlock()*/
	return hm.heroes[heroId].GetPosition()
}

func (hm *HeroesManager) GetHeroDirection(heroId int32) gt.Vector2 {
	return hm.heroes[heroId].GetMoveStatus()
}

//DeleteHero 从heroes中删除一个英雄.由于使用的不是互斥锁,所以不能在游戏的过程中删除英雄!
func (hm *HeroesManager) DeleteHero(h *gt.Hero) {
	delete(hm.heroes, h.Id)
}

//UpdateHeroPositionInfo 更改hero的position,updateTime,moveStatus三个属性
func (hm *HeroesManager) UpdateHeroPositionInfo(heroId int32, position gt.Vector2, updateTime int64, movementType pb.HeroMovementType) {
	var t gt.Vector2
	switch movementType {
	case pb.HeroMovementType_HeroMoveDownType:
		t = gt.Vector2Down
	case pb.HeroMovementType_HeroMoveUpType:
		t = gt.Vector2Up
	case pb.HeroMovementType_HeroMoveRightType:
		t = gt.Vector2Right
	case pb.HeroMovementType_HeroMoveLeftType:
		t = gt.Vector2Left
	default:
		t = gt.Vector2Zero
	}
	hm.heroes[heroId].SetPositionInfo(position, t, updateTime)
}

//GetHeroesNearby 获取hid周围的需要广播的英雄,mvp1暂定广播给地图上的除了id为hid以外的其他所有英雄
func (hm *HeroesManager) GetHeroesNearby(hid int32) (ret []int32) {
	for _, hero := range hm.heroes {
		if hero.Id != hid {
			ret = append(ret, hero.Id)
		}
	}
	return
}

func (hm *HeroesManager) GetHeroStatus(hid int32) int32 {
	return hm.heroes[hid].GetStatus()
}
