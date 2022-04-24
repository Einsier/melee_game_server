package entity

import (
	"fmt"
	"melee_game_server/api/client/proto"
	"melee_game_server/utils"
	"strconv"
)

/**
*@Author Sly
*@Date 2022/1/14
*@Version 1.0
*@Description:对Unity2D中的Vector2的服务器版本的封装
 */

type Vector2 struct {
	X float32
	Y float32
}

func NewVector2(x, y float32) Vector2 {
	return Vector2{x, y}
}

// 只能用于赋值,不要随意更改!
var (
	Vector2Up    = Vector2{X: 0, Y: 1}
	Vector2Down  = Vector2{X: 0, Y: -1}
	Vector2Left  = Vector2{X: -1, Y: 0}
	Vector2Right = Vector2{X: 1, Y: 0}
	Vector2Zero  = Vector2{X: 0, Y: 0}
	Vector2Unit  = Vector2{X: 1, Y: 1}
)

//V2toToHeroMovementType 用于将Vector2转换成HeroMovementType
var V2toToHeroMovementType = map[Vector2]proto.HeroMovementType{
	Vector2Left:  proto.HeroMovementType_HeroMoveLeftType,
	Vector2Right: proto.HeroMovementType_HeroMoveRightType,
	Vector2Up:    proto.HeroMovementType_HeroMoveUpType,
	Vector2Down:  proto.HeroMovementType_HeroMoveDownType,
	Vector2Zero:  proto.HeroMovementType_HeroStopType,
}
var HeroMovementTypeToV2 = map[proto.HeroMovementType]Vector2{
	proto.HeroMovementType_HeroMoveLeftType:  Vector2Left,
	proto.HeroMovementType_HeroMoveRightType: Vector2Right,
	proto.HeroMovementType_HeroMoveUpType:    Vector2Up,
	proto.HeroMovementType_HeroMoveDownType:  Vector2Down,
	proto.HeroMovementType_HeroStopType:      Vector2Zero,
}

func (v *Vector2) Add(v2 Vector2) Vector2 {
	return NewVector2(v.X+v2.X, v.Y+v2.Y)
}

func (v *Vector2) Subtract(v2 Vector2) Vector2 {
	return NewVector2(v.X-v2.X, v.Y-v2.Y)
}

func (v *Vector2) Multiply(v2 Vector2) Vector2 {
	return NewVector2(v.X*v2.X, v.Y*v2.Y)
}

//Divide 返回的是接收器/参数
func (v *Vector2) Divide(v2 Vector2) Vector2 {
	return NewVector2(v.X/v2.X, v.Y/v2.Y)
}

//MultiplyScalar 参数是float类型,测试对于Nanosecond,平均每次将int64转换成float32再转换成int64损失64...
func (v *Vector2) MultiplyScalar(s float32) Vector2 {
	return NewVector2(v.X*s, v.Y/v.Y*s)
}

func (v *Vector2) DivideScalar(s float32) Vector2 {
	return NewVector2(v.X/s, v.Y/s)
}

//String (X,Y)
func (v *Vector2) String() string {
	return fmt.Sprintf("(%s , %s)", strconv.FormatFloat(float64(v.X), 'f', 1, 32), strconv.FormatFloat(float64(v.Y), 'f', 1, 32))
}

//VectorEqual 两个Vector2是不是近似的相等...有时go浮点数运算会有小误差,所以只要高位足够相等即可
func VectorEqual(v1, v2 Vector2) bool {
	return utils.FloatEqual(v1.X, v2.X) && utils.FloatEqual(v1.Y, v2.Y)
}

//ToProto 做一个拷贝
func (v *Vector2) ToProto() *proto.Vector2 {
	return &proto.Vector2{
		X: v.X,
		Y: v.Y,
	}
}
