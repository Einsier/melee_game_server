package aoi

import (
	configs "melee_game_server/configs/normal_game_type_configs"
	"melee_game_server/framework/entity"
	"melee_game_server/internal/normal_game/aoi/collision"
)

/**
*@Author Sly
*@Date 2022/4/21
*@Version 1.0
*@Description:normal game的游戏地图
 */

var gameMapObs = []*collision.Rectangle{
	collision.NewRectangle("Tree1", entity.NewVector2(19, 23.6), 1, 2),
	collision.NewRectangle("Tree2", entity.NewVector2(24, 48.6), 1, 2),
	collision.NewRectangle("Tree3", entity.NewVector2(26, 16.6), 1, 2),
	collision.NewRectangle("Tree4", entity.NewVector2(11, 12.6), 1, 2),
	collision.NewRectangle("Tree5", entity.NewVector2(45, 16.6), 1, 2),
	collision.NewRectangle("Tree6", entity.NewVector2(41, 42.6), 1, 2),
	collision.NewRectangle("Tree7", entity.NewVector2(7, 45.6), 1, 2),
	collision.NewRectangle("BigGreenTree1", entity.NewVector2(24, 120.5), 1, 1),
	collision.NewRectangle("BigGreenTree2", entity.NewVector2(31, 116.5), 1, 1),
	collision.NewRectangle("BigGreenTree3", entity.NewVector2(12, 103.5), 1, 1),
	collision.NewRectangle("BigGreenTree4", entity.NewVector2(29, 75.5), 1, 1),
	collision.NewRectangle("BigGreenTree5", entity.NewVector2(59, 110.5), 1, 1),
	collision.NewRectangle("BigGreenTree6", entity.NewVector2(29, 88.5), 1, 1),
	collision.NewRectangle("BigGreenTree7", entity.NewVector2(24, 113.5), 1, 1),
	collision.NewRectangle("BigGreenTree8", entity.NewVector2(71, 122.5), 1, 1),
	collision.NewRectangle("BigGreenTree9", entity.NewVector2(55, 128.5), 1, 1),
	collision.NewRectangle("BigYellowTree1", entity.NewVector2(214.5, 36.7), 1, 1),
	collision.NewRectangle("BigYellowTree2", entity.NewVector2(209.5, 46.7), 1, 1),
	collision.NewRectangle("BigYellowTree3", entity.NewVector2(169.5, 14.7), 1, 1),
	collision.NewRectangle("BigYellowTree4", entity.NewVector2(229.5, 36.7), 1, 1),
	collision.NewRectangle("BigYellowTree5", entity.NewVector2(207.5, 5.7), 1, 1),
	collision.NewRectangle("BigYellowTree6", entity.NewVector2(184.5, 45.7), 1, 1),
	collision.NewRectangle("House1", entity.NewVector2(45.5, 86.5), 6, 7),
	collision.NewRectangle("House2", entity.NewVector2(35.5, 26.5), 6, 7),
	collision.NewRectangle("House3", entity.NewVector2(195.5, 16.5), 6, 7),
	collision.NewRectangle("Pool1", entity.NewVector2(54, 26), 7, 22),
	collision.NewRectangle("Pool2", entity.NewVector2(33, 101), 11, 6),
	collision.NewRectangle("Pool3", entity.NewVector2(165, 24), 35, 15),
}

//NormalGameMapQt NormalGame全局四叉树,存储全局信息,不能增删改
var NormalGameMapQt *collision.Quadtree

//NormalGameMapRectangle NormalGame地图的Rectangle,存储地图的宽高等信息,避免计算.不能用于初始化除了 NormalGameMapQt 之外的四叉树
var NormalGameMapRectangle *collision.Rectangle

func init() {
	NormalGameMapRectangle = collision.NewRectangle("GlobalMap", entity.NewVector2(0, 0), configs.MapWidth, configs.MapHeight)
	NormalGameMapQt = collision.NewQuadtree(NormalGameMapRectangle, 1)
	for i := 0; i < len(gameMapObs); i++ {
		NormalGameMapQt.Insert(gameMapObs[i])
	}
}

//Test地图样子可以看测试地图.png
const (
	TestGameMapWidth  = 16
	TestGameMapHeight = 8
	TestGridWidth     = 4
	TestGridHeight    = 2
	TestHeroSpeed     = 0.001
)

var TestGameMap = []*collision.Rectangle{
	collision.NewRectangle("obj1", entity.NewVector2(0, 0), 1, 1),
	collision.NewRectangle("obj2", entity.NewVector2(2, 1), 1, 0.5),
	collision.NewRectangle("obj3", entity.NewVector2(4, 4), 4, 2),
	collision.NewRectangle("obj4", entity.NewVector2(8, 0), 8, 4),
	collision.NewRectangle("obj5", entity.NewVector2(8.01, 6.01), 0.5, 0.5),
	collision.NewRectangle("obj6", entity.NewVector2(8.01, 7.01), 0.5, 0.5),
	collision.NewRectangle("obj7", entity.NewVector2(11.01, 7.01), 0.5, 0.5),
	collision.NewRectangle("obj8", entity.NewVector2(11.01, 6.01), 0.5, 0.5),
	collision.NewRectangle("obj9", entity.NewVector2(12, 6), 2, 2),
	collision.NewRectangle("obj10", entity.NewVector2(3.5, 1.5), 1, 1),
}
var TestMapQT *collision.Quadtree
var TestMapRectangle *collision.Rectangle

func init() {
	TestMapRectangle = collision.NewRectangle("testMap", entity.NewVector2(0, 0), TestGameMapWidth, TestGameMapHeight)
	TestMapQT = collision.NewQuadtree(TestMapRectangle, 1)
	for i := 0; i < len(TestGameMap); i++ {
		TestMapQT.Insert(TestGameMap[i])
	}
}
