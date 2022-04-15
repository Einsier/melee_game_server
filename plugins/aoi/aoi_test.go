package aoi

import (
	"melee_game_server/framework/entity"
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
				Position:  entity.NewVector2(0.1, 0), //在(0,0)处
				Direction: entity.Vector2Right,       //往右走
				Time:      time.Now(),
			},
			{
				Id:        1,
				Position:  entity.NewVector2(4.9, 0), //在(0,40)处
				Direction: entity.Vector2Left,        //往左走
				Time:      time.Now(),
			},
		},
	}
	aoi := NewAOI(testHeroInitInfo, 40, 80, 1, 1, 1*time.Second)
	aoi.Work()
	time.Sleep(6 * time.Second)
	aoi.Stop()
}
