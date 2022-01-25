package internal

import (
	"fmt"
	"log"
	"math/rand"
	gt "melee_game_server/internal/normal_game/game_type"
	"strconv"
	"sync"
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

func TestAtomic2(t *testing.T) {
	target := int32(0)
	m := sync.Map{}
	for i := 0; i < 100000; i++ {
		go func() {
			newTarget := atomic.AddInt32(&target, 1)
			if _, ok := m.Load(newTarget); ok {
				t.Fatalf("get dup:%d", target)
			}
			m.Store(newTarget, 1)
		}()
	}
	time.Sleep(time.Second)
	for i := 1; i < 100001; i++ {
		if _, ok := m.Load(int32(i)); !ok {
			t.Fatalf("not get:%d", i)
		}
	}
}

func TestMap(t *testing.T) {
	m := make(map[string]string)
	mLock := sync.Mutex{}
	m["target"] = "target value"
	for i := 0; i < 10000; i++ {
		go func() {
			s, ok := m["target"]
			if !ok || s != "target value" {
				t.Errorf("change target value into: %s", s)
				return
			}
		}()
		go func() {
			key := strconv.Itoa(rand.Int())
			value := strconv.Itoa(rand.Int())
			mLock.Lock()
			defer mLock.Unlock()
			m[key] = value
		}()
	}
	time.Sleep(100 * time.Millisecond)
}
