package collision

import (
	configs "melee_game_server/configs/normal_game_type_configs"
	"melee_game_server/framework/entity"
	"testing"
)

/**
*@Author Sly
*@Date 2022/4/17
*@Version 1.0
*@Description:
 */

func TestQuadtree_PrintQuadtree(t *testing.T) {
	qt := NewQuadtree(NewRectangle("0", entity.Vector2{X: 0, Y: 0}, 100, 100), 1)
	qt.split()
	qt.Child[1].split()
	qt.Print()
}

func Init80Quadtree() *Quadtree {
	qt := NewQuadtree(NewRectangle("0", entity.Vector2{X: 0, Y: 0}, 80, 80), 1)
	r1 := NewRectangle("1", entity.NewVector2(0, 0), 10, 10)   //↙↙
	r2 := NewRectangle("2", entity.NewVector2(0, 10), 10, 10)  //↙
	r3 := NewRectangle("3", entity.NewVector2(10, 0), 10, 10)  //↙
	r4 := NewRectangle("4", entity.NewVector2(10, 10), 10, 10) //↙
	r5 := NewRectangle("5", entity.NewVector2(20, 10), 10, 10) //↙
	r6 := NewRectangle("6", entity.NewVector2(20, 20), 10, 10) //↙
	r7 := NewRectangle("7", entity.NewVector2(30, 30), 10, 10) //root
	qt.Insert(r1)
	qt.Insert(r2)
	qt.Insert(r3)
	qt.Insert(r4)
	qt.Insert(r5)
	qt.Insert(r6)
	qt.Insert(r7)
	return qt
}

func TestQuadtree_Insert(t *testing.T) {
	qt := Init80Quadtree()
	qt.Print()
	if qt.TotalObjs != 7 || qt.NObjs != 1 || qt.Child[0].TotalObjs != 6 ||
		qt.Child[0].NObjs != 5 || qt.Child[0].Child[0].NObjs != 1 || qt.Child[0].Child[0].TotalObjs != 1 {
		t.Fatalf("wrong insert")
	}
}

func TestQuadtree_CheckCollision(t *testing.T) {
	qt := Init80Quadtree()
	r1 := NewRectangle("-1", entity.NewVector2(0, 0), 5, 5)            //↙↙内,贴着边.应该撞上
	r2 := NewRectangle("-2", entity.NewVector2(10, 10), 5, 5)          //↙↙内,不贴边,应该撞上
	r3 := NewRectangle("-3", entity.NewVector2(20, 5), 5, 5)           //↙,贴着中心线,不应该撞上
	r4 := NewRectangle("-4", entity.NewVector2(0, 30), 10, 10)         //↙↖内,贴边,不应该撞上
	r5 := NewRectangle("-5", entity.NewVector2(30, 30), 20, 20)        //root内,应该撞上
	r6 := NewRectangle("-6", entity.NewVector2(40, 40), 40, 40)        //root内,不应该撞上
	r7 := NewRectangle("-7", entity.NewVector2(0, 0), 40, 40)          //整个↙,应该撞上
	r8 := NewRectangle("-8", entity.NewVector2(29.999, 0.001), 10, 10) //↙↘内,应该撞上
	m := map[*Rectangle]bool{
		r1: true,
		r2: true,
		r3: false,
		r4: false,
		r5: true,
		r6: false,
		r7: true,
		r8: true,
	}
	for r, ok := range m {
		if qt.CheckCollision(r) != ok {
			t.Fatalf("wrong:%s", r.Name)
		}
	}
}

func TestQuadtree_Illegal(t *testing.T) {
	qt := Init80Quadtree()
	r1 := NewRectangle("-1", entity.NewVector2(100, 0), 5, 5)     //x越界
	r2 := NewRectangle("-2", entity.NewVector2(-100, -100), 5, 5) //x,y均越界
	r3 := NewRectangle("-3", entity.NewVector2(0, 0), 500, 500)   //width过大
	m := map[*Rectangle]bool{
		r1: false,
		r2: false,
		r3: true,
	}
	for r, ok := range m {
		if qt.CheckCollision(r) != ok {
			t.Fatalf("wrong:%s", r.Name)
		}
	}
}

func NewTestRuby(mid entity.Vector2) *Rectangle {
	return NewRectangle("0",
		entity.NewVector2(mid.X-configs.HeroColliderXHalf, mid.Y-configs.HeroColliderYHalf),
		configs.HeroColliderX, configs.HeroColliderY)
}
func TestNewRubyCollisionCheckRectangle(t *testing.T) {
	var pre, after, r *Rectangle

	//测试直线向右走
	pre, after = NewTestRuby(entity.NewVector2(1, 0.5)), NewTestRuby(entity.NewVector2(2, 0.5))
	r = NewRubyCollisionCheckRectangle("0", &pre.Mid, &after.Mid)
	if r.Mid.Y != pre.Mid.Y {
		t.Fatalf("r.Mid:%v != pre.Mid:%v", r.Mid, pre.Mid)
	}
	if !entity.VectorEqual(r.LL, pre.LL) {
		t.Fatalf("r.LL:%v != pre.LL:%v", r.LL, pre.LL)
	}
	if !entity.VectorEqual(r.UR, after.UR) {
		t.Fatalf("r.UR:%v != pre.UR:%v", r.UR, after.UR)
	}
}
