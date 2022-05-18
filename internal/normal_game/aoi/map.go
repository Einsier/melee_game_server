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
	collision.NewRectangle("Tree1-3", entity.NewVector2(15, 92.1), 1, 4),
	collision.NewRectangle("Tree4-6", entity.NewVector2(28, 84.77), 1, 4),
	collision.NewRectangle("Tree7-10", entity.NewVector2(36, 45.6), 1, 5.5),
	collision.NewRectangle("Tree11-14", entity.NewVector2(43, 13.5), 1, 5.5),
	collision.NewRectangle("Tree15-16", entity.NewVector2(37, 5.6), 1, 2.5),
	collision.NewRectangle("Tree17-18", entity.NewVector2(80, 54.6), 1, 2.5),
	collision.NewRectangle("Tree19-20", entity.NewVector2(87, 34.6), 1, 2.5),
	collision.NewRectangle("Tree21-22", entity.NewVector2(119, 60.1), 1, 2.5),
	collision.NewRectangle("Tree23-28", entity.NewVector2(217, 45.6), 1, 8.5),
	collision.NewRectangle("Tree29-44", entity.NewVector2(77, 100.6), 7, 6.5),
	collision.NewRectangle("Tree45-60", entity.NewVector2(46.7, 114.1), 7, 6.5),

	collision.NewRectangle("GreeTree1~3", entity.NewVector2(14, 40.5), 5, 1),
	collision.NewRectangle("GreeTree4~9", entity.NewVector2(28, 61.5), 11, 1),
	collision.NewRectangle("GreeTree10~13", entity.NewVector2(53, 94.5), 7, 1),
	collision.NewRectangle("GreeTree14~17", entity.NewVector2(69, 71.5), 7, 1),
	collision.NewRectangle("GreeTree18~23", entity.NewVector2(73, 46.5), 11, 1),
	collision.NewRectangle("GreeTree24~27", entity.NewVector2(102, 23.5), 1, 4),
	collision.NewRectangle("GreeTree28~30", entity.NewVector2(104, 56.5), 5, 1),
	collision.NewRectangle("GreeTree31~36", entity.NewVector2(108, 68.5), 3, 2),
	collision.NewRectangle("GreeTree37~40", entity.NewVector2(117, 79.5), 1, 4),
	collision.NewRectangle("GreeTree41~43", entity.NewVector2(128, 82.5), 5, 1),
	collision.NewRectangle("GreeTree44~45", entity.NewVector2(130, 73.5), 1, 2),
	collision.NewRectangle("GreeTree46~47", entity.NewVector2(128, 66.5), 1, 2),
	collision.NewRectangle("GreeTree48~52", entity.NewVector2(136, 30.5), 9, 1),
	collision.NewRectangle("GreeTree53~57", entity.NewVector2(154, 52.5), 9, 1),
	collision.NewRectangle("GreeTree58~61", entity.NewVector2(161, 81.5), 7, 1),
	collision.NewRectangle("GreeTree62~65", entity.NewVector2(172, 63.5), 1, 4),
	collision.NewRectangle("GreeTree66~70", entity.NewVector2(182, 46.5), 9, 1),
	collision.NewRectangle("GreeTree71~73", entity.NewVector2(189, 67.5), 5, 1),
	collision.NewRectangle("GreeTree74~105", entity.NewVector2(49, 76.5), 15, 4),
	collision.NewRectangle("GreeTree106~111", entity.NewVector2(185, 110.5), 5, 2),
	collision.NewRectangle("GreeTree112~117", entity.NewVector2(205, 102.5), 5, 2),
	collision.NewRectangle("GreeTree118~123", entity.NewVector2(216, 114.5), 5, 2),
	collision.NewRectangle("GreeTree124~135", entity.NewVector2(249, 19.5), 5, 4),
	collision.NewRectangle("GreeTree136~147", entity.NewVector2(129, 8.5), 5, 4),
	collision.NewRectangle("GreeTree148~159", entity.NewVector2(144, 114.5), 11, 2),
	collision.NewRectangle("GreeTree160~161", entity.NewVector2(100, 73.5), 1, 2),

	collision.NewRectangle("House1", entity.NewVector2(11.5, 23.5), 6, 7),
	collision.NewRectangle("House2", entity.NewVector2(74.5, 12.5), 6, 7),
	collision.NewRectangle("House3", entity.NewVector2(96.5, 117.5), 6, 7),
	collision.NewRectangle("House4", entity.NewVector2(202.5, 28.5), 6, 7),
	collision.NewRectangle("House5", entity.NewVector2(100.5, 81.5), 6, 7),

	collision.NewRectangle("Pool1", entity.NewVector2(54, 26), 7, 22),
	collision.NewRectangle("Pool2", entity.NewVector2(33, 101), 11, 6),
	collision.NewRectangle("Pool3", entity.NewVector2(165, 24), 25, 15),
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
