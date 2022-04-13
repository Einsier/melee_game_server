package aoi

import (
	"melee_game_server/framework/entity"
	"melee_game_server/plugins/aoi/object"
	"melee_game_server/plugins/logger"
	"sync"
	"time"
)

/**
*@Author Sly
*@Date 2022/4/12
*@Version 1.0
*@Description:1.0版本 暂时将地图宽度和高度,格子宽度高度设置为整形,方便进行各种取模运算
大概思想来自:https://blog.csdn.net/qq_36748278/article/details/102787225
为了效率,很多固定的运算提前的进行了运算,例如每个格子的四周格子,每个格子的邻接格子等等
 */

type MapInfo struct {
	MX    int            //地图宽度
	MY    int            //地图高度
	GX    int            //每个格子的宽度
	GY    int            //每个格子的高度
	NGX   int            //X轴上有多少个格子
	NGY   int            //Y轴上有多少个格子
	NGrid int            //一共有多少个格子
	Grids []*object.Grid //被上面的博客带跑偏了= =感觉还是二维数组比较好写 而且性能也一样的= =
}

type HeroInfo struct {
	HeroNum			int
	Speed			float64				//每ms移动多少m
	UpdateDuration	time.Duration

	heroLock		sync.RWMutex
	heroes			map[int32]*Hero
}

type AOI struct {
	*MapInfo
	*HeroInfo
	ticker *time.Ticker
}

//GetGridByPos 获取当前位置应该属于哪个grid,并从中取出grid
func (mp *MapInfo)GetGridByPos(pos entity.Vector2)*object.Grid{
	return mp.Grids[(int(pos.Y) / mp.GY) * mp.NGX + (int(pos.X) / mp.GX)]
}

//GetGridById 通过id获取grid,如果越界返回nil,如果create = true说明如果不存在,创建一个
func (mp *MapInfo)GetGridById(id int)*object.Grid{
	if id < 0 || id >= mp.NGrid{
		return nil
	}
	return mp.Grids[id]
}

//GetSurround 获取周围的八个格子,注意对于例如最边上的一圈格子来说,可能存在邻居为nil的情况
func (mp *MapInfo)GetSurround(self *object.Grid)[8]*object.Grid{
	var ret [8]*object.Grid
	//分别为(-1,1) (0,1) (1,1) (-1,0) (1,0) (-1,-1) (0,-1) (1,-1),也就是从上到下,从左到右除了自己的集合,它们分别跟本格子id差了
	//		NGX-1  NGX	NGX+1  -1	  1		-NGX-1	-NGX   -NGX+1
	ret[0] = mp.GetGridById(self.Id+ mp.NGX-1)
	ret[1] = mp.GetGridById(self.Id+ mp.NGX)
	ret[2] = mp.GetGridById(self.Id+ mp.NGX+1)
	ret[3] = mp.GetGridById(self.Id-1)
	ret[4] = mp.GetGridById(self.Id+1)
	ret[5] = mp.GetGridById(self.Id- mp.NGX-1)
	ret[6] = mp.GetGridById(self.Id- mp.NGX)
	ret[7] = mp.GetGridById(self.Id- mp.NGX+1)
	return ret
}

//NewAOI 创建一个AOI模块
func NewAOI(heroesInitInfo *HeroesInitInfo,mx,my,gx,gy int,updateDuration time.Duration)*AOI{
	if mx<0 || my<0 || gx<0 || gy<0{
		return nil
	}

	aoi := new(AOI)
	//初始化map有关字段
	mapInfo := &MapInfo{
		MX:    mx,
		MY:    my,
		GX:    gx,
		GY:    gy,
		NGX:   mx/gx + 1,		//例如mx为10,gx为6,那么x方向应有两个格子
		NGY:   my/gy + 1,
		NGrid: (mx/gx + 1) * (my/gy + 1),
	}

	aoi.MapInfo = mapInfo
	//初始化grids,创建NGrid个数的格子
	id := 0
	aoi.Grids = make([]*object.Grid,mapInfo.NGrid)
	for y := 0; y < mapInfo.NGY; y++ {
		for x := 0; x < mapInfo.NGX; x++ {
			id = y*mapInfo.NGX + x
			aoi.Grids[id] = object.NewGrid(id,x,y)
		}
	}
	//初始化格子的周围的格子
	for y := 0; y < mapInfo.NGY; y++ {
		for x := 0; x < mapInfo.NGX; x++ {
			id = y*mapInfo.NGX + x
			aoi.Grids[id].Surround = aoi.GetSurround(aoi.Grids[id])
			aoi.Grids[id].Grids = aoi.Grids
		}
	}

	//初始化英雄
	heroInfo := &HeroInfo{
		HeroNum: len(heroesInitInfo.heroes),
		Speed:   heroesInitInfo.Speed,
		heroes:  make(map[int32]*Hero,len(heroesInitInfo.heroes)),
	}
	for i := 0; i < heroInfo.HeroNum; i++ {
		hero := heroesInitInfo.heroes[i]
		heroInfo.heroes[hero.Id] = NewHero(hero.Id,hero.Position,entity.Vector2Zero,heroInfo.Speed,aoi)
	}
	heroInfo.UpdateDuration = updateDuration
	aoi.HeroInfo = heroInfo

	return aoi
}

//DeleteHero 删除英雄
func (aoi *AOI)DeleteHero(id int32){
	aoi.heroLock.Lock()
	defer aoi.heroLock.Unlock()
	delete(aoi.heroes, id)
}

func (aoi *AOI)UpdateHeroPosition(info *HeroMovementInfo){
	aoi.heroLock.RLock()
	hero := aoi.heroes[info.Id]
	aoi.heroLock.RUnlock()

	if hero == nil{
		logger.Errorf("aoi收到了不存在的英雄位置更新信息:hero num:%d",info.Id)
		return
	}
	hero.UpdateMovement(info)
}

func (aoi *AOI)Work(){
	ticker := time.NewTicker(aoi.UpdateDuration)
	go func() {
		for range ticker.C{
			aoi.heroLock.RLock()
			for _, hero := range aoi.heroes {
				hero.UpdateMovement(nil)
			}
			aoi.heroLock.Unlock()
		}
	}()
	aoi.ticker = ticker
}

func (aoi *AOI)Stop(){
	aoi.ticker.Stop()
}

