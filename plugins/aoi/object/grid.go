package object

/**
*@Author Sly
*@Date 2022/4/12
*@Version 1.0
*@Description:使用int32表示本grid存放的内容
 */

type Grid struct {
	Id    int
	XIdx  int                //横向的第几个格子（从0开始）
	YIdx  int                //纵向的第几个格子（从0开始）
	Objs  map[int32]struct{} //key为英雄id
	Grids [][]*Grid          //所有的格子
}

//NewGrid 初始化除了Neighbor和Surround两个字段的格子
func NewGrid(id, xIdx, yIdx int, grids [][]*Grid) *Grid {
	return &Grid{
		Id:    id,
		XIdx:  xIdx,
		YIdx:  yIdx,
		Objs:  make(map[int32]struct{}),
		Grids: grids,
	}
}

//GetObjs 获取当前格子内有哪几个物体(即英雄)
func (g *Grid) GetObjs() []int32 {
	ret := make([]int32, len(g.Objs))
	i := 0
	for id, _ := range g.Objs {
		ret[i] = id
		i++
	}
	return ret
}
