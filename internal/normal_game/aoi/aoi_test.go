package aoi

import (
	configs "melee_game_server/configs/normal_game_type_configs"
	"melee_game_server/framework/entity"
	"melee_game_server/internal/normal_game/aoi/collision"
	"strconv"
	"sync"
	"testing"
	"time"
)

/**
*@Author Sly
*@Date 2022/4/12
*@Version 1.0
*@Description:
 */

//TestMeet 英雄0从(0,0.1)[0][0]向右走,英雄1从(0,4.9)[0][4]向左走,速度为1m/s(aka 0.001m/ms),
//应该是	1s:		 (0,1.1)[0][1]			 (0,3.9)[0][3]
//		2s:		 (0,2.1)[0][2]->进入视野	 (0,2.9)[0][2]
//		3s:		 (0,3.1)[0][3]			 (0,1.9)[0][1]->离开视野
//		4s:		 (0,4.1)[0][4]	 		 (0,0.9)[0][0]
func TestSingleMeet(t *testing.T) {
	var testHeroInitInfo = &HeroesInitInfo{
		Speed: 0.001,
		heroes: []*HeroMoveMsg{
			{
				Id:        0,
				Position:  entity.NewVector2(0.1, 0),
				Direction: entity.Vector2Right, //往右走
				Time:      time.Now(),
			},
			{
				Id:        1,
				Position:  entity.NewVector2(4.9, 0),
				Direction: entity.Vector2Left, //往左走
				Time:      time.Now(),
			},
		},
	}
	aoi := NewAOI(testHeroInitInfo, 40, 80, 1, 1, 1*time.Second, nil, nil)
	aoi.Work()
	time.Sleep(6 * time.Second)
	aoi.Stop()
}

func TestSingleMeet2(t *testing.T) {
	//2 3
	//0 1
	var testHeroInitInfo = &HeroesInitInfo{
		Speed: 0.001,
		heroes: []*HeroMoveMsg{
			{
				//左下角出发
				Id:        0,
				Position:  entity.NewVector2(0.1, 0.1),
				Direction: entity.Vector2Zero,
				Time:      time.Now(),
			},
			{
				//右下角出发
				Id:        1,
				Position:  entity.NewVector2(4.9, 0.1),
				Direction: entity.Vector2Zero,
				Time:      time.Now(),
			},
			{
				//左上角出发
				Id:        2,
				Position:  entity.NewVector2(0.1, 4.9),
				Direction: entity.Vector2Zero,
				Time:      time.Now(),
			},
			{
				//右上角出发
				Id:        3,
				Position:  entity.NewVector2(4.9, 4.9),
				Direction: entity.Vector2Zero,
				Time:      time.Now(),
			},
		},
	}
	aoi := NewAOI(testHeroInitInfo, 5, 5, 1, 1, 1*time.Second, nil, nil)
	aoi.Work()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		i := 50
		for i != 0 {
			time.Sleep(100 * time.Millisecond)
			aoi.PutMove(&HeroMoveMsg{
				Id:        0,
				Position:  entity.Vector2Zero,
				Direction: entity.Vector2Up,
				Time:      time.Now(),
			})
			aoi.PutMove(&HeroMoveMsg{
				Id:        1,
				Position:  entity.Vector2Zero,
				Direction: entity.Vector2Left,
				Time:      time.Now(),
			})
			aoi.PutMove(&HeroMoveMsg{
				Id:        2,
				Position:  entity.Vector2Zero,
				Direction: entity.Vector2Down,
				Time:      time.Now(),
			})
			aoi.PutMove(&HeroMoveMsg{
				Id:        3,
				Position:  entity.Vector2Zero,
				Direction: entity.Vector2Left,
				Time:      time.Now(),
			})
			time.Sleep(100 * time.Millisecond)
			aoi.PutMove(&HeroMoveMsg{
				Id:        0,
				Position:  entity.Vector2Zero,
				Direction: entity.Vector2Right,
				Time:      time.Now(),
			})
			aoi.PutMove(&HeroMoveMsg{
				Id:        1,
				Position:  entity.Vector2Zero,
				Direction: entity.Vector2Up,
				Time:      time.Now(),
			})
			aoi.PutMove(&HeroMoveMsg{
				Id:        2,
				Position:  entity.Vector2Zero,
				Direction: entity.Vector2Right,
				Time:      time.Now(),
			})
			aoi.PutMove(&HeroMoveMsg{
				Id:        3,
				Position:  entity.Vector2Zero,
				Direction: entity.Vector2Down,
				Time:      time.Now(),
			})
			i--
		}
		wg.Done()
	}()
	wg.Wait()
	aoi.RemoveHero(0)
	aoi.RemoveHero(1)
	aoi.RemoveHero(2)
	aoi.RemoveHero(3)
	aoi.Stop()
}

//TestAOI_Grid 梦幻联动
func TestAOI_Grid(t *testing.T) {
	var testHeroInitInfo = &HeroesInitInfo{
		Speed: 0.001, //1s走1m
		heroes: []*HeroMoveMsg{
			{
				Id:        1,
				Position:  entity.NewVector2(0.1, 7),
				Direction: entity.Vector2Right, //往右走
				Time:      time.Now(),
			},
			//{
			//	Id:        2,
			//	Position:  entity.NewVector2(4.9, 0.1),
			//	Direction: entity.Vector2Left, //往左走
			//	Time:      time.Now(),
			//},
		},
	}
	aoi := NewAOI(testHeroInitInfo, TestGameMapWidth, TestGameMapHeight, TestGridWidth, TestGridHeight, 500*time.Millisecond, nil, TestMapQT)
	//aoi.qt.Print()
	go aoi.Work()
	time.Sleep(10 * time.Second)
}

func TestNewRandomHeroesInitInfo(t *testing.T) {
	const heroNum = 100
	info := NewRandomHeroesInitInfo(heroNum, 0.001, TestMapQT)
	if info.heroes == nil || len(info.heroes) != heroNum {
		t.Fatalf("wrong num")
	}
	heroQt := collision.NewQuadtree(collision.NewRectangle("heroQt", entity.NewVector2(0, 0), configs.MapWidth, configs.MapHeight), 1)
	for i := 0; i < heroNum; i++ {
		heroMid := collision.NewRubyRectangleByMid(info.heroes[i].Position, "info-"+strconv.Itoa(i+1))
		if TestMapQT.CheckCollision(heroMid) || heroQt.CheckCollision(heroMid) {
			t.Fatalf("hero%d collision", i+1)
		}
		heroQt.Insert(heroMid)
	}
	heroQt.Print()
}
