package internal

import (
	"fmt"
	"log"
	"math/rand"
	"melee_game_server/internal/game/game_room"
	gt "melee_game_server/internal/game/game_type"
	"testing"
)

/**
*@Author Sly
*@Date 2022/1/14
*@Version 1.0
*@Description:
 */

func TestCountBulletId(t *testing.T) {
	i := gt.CountBulletId(1, 1)
	if i != int64(1<<32+1) {
		t.Fatalf("wrong")
	}
}

//getVector 例如把(0,0)拆成满足x = (-1,0) + (2,0) + (-1,0) = (0,0)这样的x集合,y同理,用于测试
func getVector(times int) (x []gt.Vector2, y []gt.Vector2) {
	for i := 0; i < times; i++ {
		xRand := rand.Float64()
		yRand := rand.Float64()
		xMove := xRand - float64(int64(xRand))
		yMove := yRand - float64(int64(yRand))
		x = append(x, gt.Vector2{X: xMove, Y: 0})
		y = append(y, gt.Vector2{X: 0, Y: yMove})
	}
	for i := 0; i < times; i++ {
		x = append(x, gt.Vector2{X: -x[i].X, Y: 0})
		y = append(y, gt.Vector2{X: 0, Y: -y[i].Y})
	}
	return
}

func TestHeroes(t *testing.T) {
	h1 := gt.NewHero(1, gt.Vector2Zero)
	h2 := gt.NewHero(2, gt.Vector2Zero)
	h3 := gt.NewHero(3, gt.Vector2Zero)
	h4 := gt.NewHero(4, gt.Vector2Zero)
	grm := game_room.NewGameRoomManager()
	grm.AddHero(h1)
	grm.AddHero(h2)
	grm.AddHero(h3)
	grm.AddHero(h4)

	if grm.GetHero(1).Id != 1 {
		t.Fatalf("get hero error")
	}

}

func TestGetVector(t *testing.T) {
	x, y := getVector(100000)
	v := gt.NewVector2(0, 0)
	for _, testV := range x {
		v.Add(testV)
	}
	for _, testV := range y {
		v.Add(testV)
	}
	if !gt.VectorEqual(v, gt.Vector2Zero) {
		log.Fatalf("fail! result:%v", v)
	}
	fmt.Printf("result:%v", v)
}
