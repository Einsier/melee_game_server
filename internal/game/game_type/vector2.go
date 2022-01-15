package game_type

import (
	"fmt"
	"melee_game_server/utils"
)

/**
*@Author Sly
*@Date 2022/1/14
*@Version 1.0
*@Description:对Unity2D中的Vector2的服务器版本的封装
 */

type Vector2 struct {
	X float64
	Y float64
}

func NewVector2(x, y float64) Vector2 {
	return Vector2{x, y}
}

// 只能用于赋值,不要随意更改!
var (
	Vector2Up    = Vector2{X: 0, Y: 1}
	Vector2Down  = Vector2{X: 0, Y: -1}
	Vector2Left  = Vector2{X: 0, Y: -1}
	Vector2Right = Vector2{X: 1, Y: 0}
	Vector2Zero  = Vector2{X: 0, Y: 0}
	Vector2Unit  = Vector2{X: 1, Y: 1}
)

func (v *Vector2) Add(v2 Vector2) {
	v.X += v2.X
	v.Y += v2.Y
}

func (v *Vector2) Subtract(v2 Vector2) {
	v.X -= v2.X
	v.Y -= v2.Y
}

func (v *Vector2) Multiply(v2 Vector2) {
	v.X *= v2.X
	v.Y *= v2.Y
}

func (v *Vector2) Divide(v2 Vector2) {
	v.X /= v2.X
	v.Y /= v2.Y
}

func (v *Vector2) MultiplyScalar(s float64) {
	v.X *= s
	v.Y *= s
}

func (v *Vector2) DivideScalar(s float64) {
	v.X /= s
	v.Y /= s
}

func (v *Vector2) String() string {
	return fmt.Sprintf("%v:%v", v.X, v.Y)
}

func VectorEqual(v1, v2 Vector2) bool {
	return utils.FloatEqual(v1.X, v2.X) && utils.FloatEqual(v1.Y, v2.Y)
}
