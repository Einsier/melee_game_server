package collision

import (
	"fmt"
	"melee_game_server/framework/entity"
	"strconv"
)

/**
*@Author Sly
*@Date 2022/4/17
*@Version 1.0
*@Description:地图中的实体,具有实际的位置和大小
 */

//Rectangle 长方体
type Rectangle struct {
	Id     int            //如果是地图资源的话,加一个id字段,方便debug,如果不是地图资源此字段可以为0
	GridX  int            //在Grids中的X
	GridY  int            //在Grids中的Y
	LL     entity.Vector2 //左下角的位置(lower left)
	UR     entity.Vector2 //右上角的位置(upper right)
	Mid    entity.Vector2 //中心的位置
	Width  float32        //宽度
	Height float32        //高度
	Next   *Rectangle     //用于索引下一个 Rectangle
}

//CollisionWith 判断是否和另一个长方体发生碰撞
func (me *Rectangle) CollisionWith(other *Rectangle) bool {
	//如果me的最大X够不到other的X左边/other的最大X够不到me的最左边/me的最大Y够不到other的最下面/other的最大Y够不到me的最下面,那么没有碰撞
	return !(me.UR.X <= other.LL.X || other.UR.X <= me.LL.X || me.UR.Y <= other.LL.Y || other.UR.Y <= me.LL.Y)
}

func NewRectangle(id int, position entity.Vector2, width, height float32) *Rectangle {
	return &Rectangle{
		Id:     id,
		LL:     position,
		UR:     entity.NewVector2(position.X+width, position.Y+height),
		Mid:    entity.NewVector2(position.X+width/2, position.Y+height/2),
		Width:  width,
		Height: height,
		Next:   nil,
	}
}

func (r *Rectangle) String() string {
	return fmt.Sprintf("position:%v,Mid:%v,Width:%s,Height,%s", &r.LL, &r.Mid,
		strconv.FormatFloat(float64(r.Width), 'f', 1, 32),
		strconv.FormatFloat(float64(r.Height), 'f', 1, 32))
}
