package game_room

/**
*@Author chenjiajia
*@Date 2022/1/21
*@Version 1.0
*@Description:
 */

func UpdateHeroPositionTimeEventCallback(gm *NormalGameRoom) {
	/*	for _, hero := range gm.heroManager.heroes {
		if hero.status == configs.HeroDead {
			continue
		}
		nowTime := time.Now().UnixNano()
		timeElapse := nowTime - hero.UpdateTime //距上次更新的间隔
		hero.UpdateTime = nowTime
		distance := float64(timeElapse) * configs.HeroMoveSpeed //移动的距离
		newX := hero.GetPosition().X + distance*hero.GetDirection().X
		newY := hero.GetDirection().Y + distance*hero.GetDirection().Y
		gm.heroManager.MoveHeroPosition(hero.Id, gt.NewVector2(newX, newY))
	}*/
}
