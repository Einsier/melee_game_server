package internal

import (
	"fmt"
	"log"
	"math/rand"
	"melee_game_server/internal/game/game_room"
	gt "melee_game_server/internal/game/game_type"
	"sync/atomic"
	"testing"
	"time"
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

func TestSyncMap(t *testing.T) {
	//sm := sync.Map{}
	i1 := 1
	i2 := 2
	v1 := atomic.Value{}
	v1.Store(i1)
	v2 := atomic.Value{}
	v2.Store(i2)

	fmt.Printf("%v\n", v1.Load())
	fmt.Printf("%v\n", v2.Load())
}

func TestAtomic(t *testing.T) {
	var target int32 = 0
	for i := 0; i < 1000000; i++ {
		go func() {
			for {
				old := target
				after := old + 1
				swapped := atomic.CompareAndSwapInt32(&target, old, after)
				if swapped {
					return
				}
			}
		}()
		go func() {
			for {
				old := target
				after := old - 1
				swapped := atomic.CompareAndSwapInt32(&target, old, after)
				if swapped {
					return
				}
			}
		}()
	}
	time.Sleep(1 * time.Second)
}

func TestMap(t *testing.T) {
	m := make(map[string]string)
	for i := 0; i < 100000; i++ {
		go func() {
			old := m
			//这里进行从old中取数等操作
			_ = old["abc"]
			_ = old["abc"]
			_ = old["abc"]
			_ = old["abc"]
		}()
		go func() {
			m = make(map[string]string)
		}()
	}
	time.Sleep(1 * time.Second)
}
