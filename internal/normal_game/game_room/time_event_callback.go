package game_room

/**
*@Author chenjiajia
*@Date 2022/1/21
*@Version 1.0
*@Description:
 */

func UpdateHeroPositionTimeEventCallback(gm *NormalGameRoom) {
	/*		for _, hero := range gm.heroManager.heroes {
			if hero.GetStatus() == configs.HeroDead {
				continue
			}
			nowTime := time.Now().UnixNano()
			timeElapse := nowTime - hero.updateTime //距上次更新的间隔
			hero.updateTime = nowTime
			distance := float64(timeElapse) * configs.HeroMoveSpeed //移动的距离
			newX := hero.GetPosition().X + distance*hero.GetMoveStatus().X
			newY := hero.GetMoveStatus().Y + distance*hero.GetMoveStatus().Y
			gm.heroManager.MoveHeroPosition(hero.Id, gt.NewVector2(newX, newY))
		}*/
}
