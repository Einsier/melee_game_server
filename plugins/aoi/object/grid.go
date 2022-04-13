package object

import "sync"

/**
*@Author Sly
*@Date 2022/4/12
*@Version 1.0
*@Description:使用int32表示本grid存放的内容
 */


type Grid struct {
	Id			int
	XIdx		int
	YIdx		int
	mu			sync.RWMutex
	objs		map[int32]struct{}		//key为英雄id
	Surround	[8]*Grid				//周围八个格子
	Grids		[]*Grid					//所有的格子
}

//NewGrid 初始化除了Neighbor和Surround两个字段的格子
func NewGrid(id,xIdx,yIdx int)*Grid{
	return &Grid{
		Id:       id,
		XIdx:     xIdx,
		YIdx:     yIdx,
		mu:       sync.RWMutex{},
		objs:     make(map[int32]struct{}),
	}
}

//GetObjs 获取当前格子内有哪几个物体(即英雄)
func (g *Grid)GetObjs()[]int32{
	g.mu.RLock()
	defer g.mu.RUnlock()
	ret := make([]int32,len(g.objs))
	i := 0
	for id, _ := range g.objs{
		ret[i] = id
		i++
	}
	return ret
}

//ObjLeave 英雄从本格子离开到另一个格子,d表示去的格子是哪个
func(g *Grid)ObjLeave(objId int32,to *Grid){
	//本格子删除物体
	g.mu.Lock()
	delete(g.objs, objId)
	g.mu.Unlock()

	//根据方向,将物体添加到对应的格子中
	to.mu.Lock()
	to.objs[objId] = struct{}{}
	to.mu.Unlock()
}
