package collision

import (
	"fmt"
	configs "melee_game_server/configs/normal_game_type_configs"
	"melee_game_server/framework/entity"
	"strings"
)

/**
*@Author Sly
*@Date 2022/4/17
*@Version 1.0
*@Description:四叉树结构,借鉴https://github.com/JamesLMilner/quadtree-go/blob/master/quadtree.go
相比原博客,结合了学长的代码以及自己的思路,将每个四叉树节点的地图资源的组织形式改为指针链表,避免拷贝带来的浪费.
同时给定英雄,检索英雄可能碰撞的物体时,采用返回链表头指针的切片(而不是返回所有的可能碰撞的物体的指针)的方式,最大限度的减少拷贝
由于是无锁的单线程,所以同一时刻只有一个线程对四叉树进行增删改查.
*/

type Quadtree struct {
	Self      *Rectangle  //表示自己的位置信息,同时也作为头结点,指向本四叉树节点存放的地图资源
	Level     int         //自己的层数
	NObjs     int         //当前内部的资源数量
	TotalObjs int         //自己+子孙结点内部的资源数量
	Child     []*Quadtree //分别是0:左下角	1:左上角	2:右上角	3:右下角
}

//NewQuadtree 新建一颗四叉树,注意 self 字段最好使用NewRectangle新建一个Rectangle,因为涉及到链表
func NewQuadtree(self *Rectangle, level int) *Quadtree {
	return &Quadtree{
		Self:  self,
		Level: level,
		Child: make([]*Quadtree, 0, 4),
		NObjs: 0,
	}
}

//split 分裂一手
func (qt *Quadtree) split() {
	if len(qt.Child) == 4 {
		return
	}

	halfWidth, halfHeight := qt.Self.Width/2, qt.Self.Height/2
	level := qt.Level + 1
	qt.Child = append(qt.Child, NewQuadtree(NewRectangle(qt.Self.Name+"↙", qt.Self.LL, halfWidth, halfHeight), level))                                                         //左下角
	qt.Child = append(qt.Child, NewQuadtree(NewRectangle(qt.Self.Name+"↖", entity.NewVector2(qt.Self.LL.X, qt.Self.LL.Y+halfHeight), halfWidth, halfHeight), level))           //左上角
	qt.Child = append(qt.Child, NewQuadtree(NewRectangle(qt.Self.Name+"↗", entity.NewVector2(qt.Self.LL.X+halfWidth, qt.Self.LL.Y+halfHeight), halfWidth, halfHeight), level)) //右上角
	qt.Child = append(qt.Child, NewQuadtree(NewRectangle(qt.Self.Name+"↘", entity.NewVector2(qt.Self.LL.X+halfWidth, qt.Self.LL.Y), halfWidth, halfHeight), level))            //右下角
}

//getIndex 判断r应该在qt的哪个孩子中.如果没有一个孩子可以完美的放下,那么返回-1,表示应该在qt本身中.
func (qt *Quadtree) getIndex(r *Rectangle) int {
	inBot := r.UR.Y < qt.Self.Mid.Y
	inUp := r.LL.Y > qt.Self.Mid.Y
	inLeft := r.UR.X < qt.Self.Mid.X
	inRight := r.LL.X > qt.Self.Mid.X

	if inBot && inLeft {
		return 0
	} else if inUp && inLeft {
		return 1
	} else if inUp && inRight {
		return 2
	} else if inBot && inRight {
		return 3
	}

	return -1
}

func (qt *Quadtree) Insert(r *Rectangle) {
	qt.TotalObjs++

	if len(qt.Child) != 0 {
		//如果本节点有孩子的话,插入到本节点的孩子中
		idx := qt.getIndex(r)
		if idx != -1 {
			qt.Child[idx].Insert(r)
			return
		}
	}

	//没有子节点 or 该节点应该插入qt自身
	qt.NObjs++
	r.Next = qt.Self.Next //先把r插入到自身中
	qt.Self.Next = r
	if (qt.NObjs > configs.QuadtreeMaxObjs) && (qt.Level < configs.QuadtreeMaxLevel) {
		//看学长和github上的项目都有个小bug,就是例如MaxObjs是4,但是本节点中刚好有5个穿过中心线的资源,那么每次新加入资源都要重新让所有的资源重新分配到四个孩子中,就很狗.
		//所以这里做了一点点逻辑优化:如果当前已经满了,并且没有达到最大Level的话,并且当前是没有孩子的,那么分裂一手,并且把资源分配给孩子们
		//但是如果之前分配过,那么说明当前父节点中保存的地图资源肯定是没法完美的放到子节点中的资源,那么不需要重新检查.
		//并且如果有孩子且当前r可以放进任何一个孩子中,那么本函数中的第一个if就已经这么干了,根本不会执行到这里.所以这里可以不做任何处理
		if len(qt.Child) == 0 {
			//如果没有分裂,那么分裂一下,然后把地图资源分配给孩子
			qt.split()
			pre, obj := qt.Self, qt.Self.Next
			for obj != nil {
				if objIdx := qt.getIndex(obj); objIdx == -1 {
					//应该仍然放到父节点中
					obj, pre = obj.Next, pre.Next
				} else {
					pre.Next = obj.Next //从父节点的链表中去除
					qt.Child[objIdx].Insert(obj)
					obj = pre.Next
					qt.NObjs--
				}
			}
		}
	}
}

//CheckCollision 判断给定的r(通常是英雄)是否和四叉树中的地图资源相撞
func (qt *Quadtree) CheckCollision(r *Rectangle) bool {
	for obj := qt.Self.Next; obj != nil; obj = obj.Next {
		//跟本节点中的地图资源比较
		if obj.CollisionWith(r) {
			//logger.Testf("%s collision with %s", obj.Name, r.Name)
			return true
		}
	}
	//如果有孩子,那么需要跟本节点中的资源比完,再比孩子
	if len(qt.Child) == 4 {
		if idx := qt.getIndex(r); idx != -1 {
			return qt.Child[idx].CheckCollision(r)
		}
	}
	return false
}

//Print 打印四叉树
func (qt *Quadtree) Print() {
	qt.doPrint("")
}
func (qt *Quadtree) doPrint(prefix string) {
	objId := make([]string, 0)
	for obj := qt.Self.Next; obj != nil; obj = obj.Next {
		objId = append(objId, obj.Name)
	}
	fmt.Printf("%s%+v,objs:[%s]\n", prefix, qt.Self, strings.Join(objId, ","))
	if len(qt.Child) == 4 {
		qt.Child[0].doPrint("\t" + prefix)
		qt.Child[1].doPrint("\t" + prefix)
		qt.Child[2].doPrint("\t" + prefix)
		qt.Child[3].doPrint("\t" + prefix)
	}
}
