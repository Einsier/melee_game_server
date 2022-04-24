package collision

import (
	"fmt"
	configs "melee_game_server/configs/normal_game_type_configs"
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
	Name   string         //方便debug
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

func NewRectangle(name string, position entity.Vector2, width, height float32) *Rectangle {
	return &Rectangle{
		Name:   name,
		LL:     position,
		UR:     entity.NewVector2(position.X+width, position.Y+height),
		Mid:    entity.NewVector2(position.X+width/2, position.Y+height/2),
		Width:  width,
		Height: height,
		Next:   nil,
	}
}

//NewRubyCollisionCheckRectangle 通过Ruby的运动之前的中心位置和运动之后的中心位置,计算出一个用于检测障碍的Rectangle
//运算逻辑见同步设计.md
func NewRubyCollisionCheckRectangle(name string, pre, after *entity.Vector2) *Rectangle {
	moveRight := pre.X <= after.X
	moveUp := pre.Y <= after.Y
	if moveRight && moveUp {
		//如果after在pre的右上角
		return NewRectangle(name, entity.NewVector2(pre.X-configs.HeroColliderXHalf, pre.Y-configs.HeroColliderYHalf),
			after.X-pre.X+configs.HeroColliderX, after.Y-pre.Y+configs.HeroColliderY)
	} else if moveRight && !moveUp {
		//如果after在pre的右下角
		return NewRectangle(name, entity.NewVector2(pre.X-configs.HeroColliderXHalf, after.Y-configs.HeroColliderYHalf),
			after.X-pre.X+configs.HeroColliderX, pre.Y-after.Y+configs.HeroColliderY)
	} else if !moveRight && moveUp {
		//如果after在pre的左上角
		return NewRectangle(name, entity.NewVector2(after.X-configs.HeroColliderXHalf, pre.Y-configs.HeroColliderYHalf),
			pre.X-after.X+configs.HeroColliderX, after.Y-pre.Y+configs.HeroColliderY)
	} else {
		//如果after在pre的左下角
		return NewRectangle(name, entity.NewVector2(after.X-configs.HeroColliderXHalf, after.Y-configs.HeroColliderYHalf),
			pre.X-after.X+configs.HeroColliderX, pre.Y-after.Y+configs.HeroColliderY)
	}
}

func NewRubyRectangleByMid(mid entity.Vector2, name string) *Rectangle {
	return &Rectangle{
		Name:   name,
		LL:     entity.NewVector2(mid.X-configs.HeroColliderXHalf, mid.Y-configs.HeroColliderYHalf),
		UR:     entity.NewVector2(mid.X+configs.HeroColliderXHalf, mid.Y+configs.HeroColliderYHalf),
		Mid:    mid,
		Width:  configs.HeroColliderX,
		Height: configs.HeroColliderY,
		Next:   nil,
	}
}

func (r *Rectangle) String() string {
	return fmt.Sprintf("[%s]LL:%v,Mid:%v,UR:%v,Width:%s,Height,%s", r.Name, &r.LL, &r.Mid, &r.UR,
		strconv.FormatFloat(float64(r.Width), 'f', 1, 32),
		strconv.FormatFloat(float64(r.Height), 'f', 1, 32))
}
