package normal_game_type_configs

/**
*@Author Sly
*@Date 2022/1/24
*@Version 1.0
*@Description:游戏地图有关的参数
 */

const (
	// MapWidth   = 240 //地图宽度,单位m
	// MapHeight  = 140 //地图高度,单位m
	MapWidth   = 480 //地图宽度,单位m
	MapHeight  = 240 //地图高度,单位m
	GridWidth  = 15  //网格宽度
	GridHeight = 10  //网格高度
)

const (
	QuadtreeMaxObjs  = 4 //四叉树达到5个应该分裂
	QuadtreeMaxLevel = 4 //四叉树层数达到4层应该不分裂
)
