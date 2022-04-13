package test

import (
	"fmt"
	configs "melee_game_server/configs/normal_game_type_configs"
	gt "melee_game_server/framework/entity"
	gr "melee_game_server/internal/normal_game/game_room"
	"sync"
	"testing"
	"time"
)

/**
*@Author Sly
*@Date 2022/1/20
*@Version 1.0
*@Description:
 */

func TestArrangeHeroId(t *testing.T) {
	pm := gr.NewHeroesManager()
	m := make(map[int32]int32)
	lock := sync.Mutex{}
	for i := 0; i < 1000; i++ {
		go func() {
			id := pm.ArrangeHeroId()
			lock.Lock()
			defer lock.Unlock()
			if _, ok := m[id]; ok {
				m[id]++
			} else {
				m[id] = 1
			}
		}()
	}
	time.Sleep(1 * time.Second)
	for i := int32(1); i <= configs.MaxNormalGamePlayerNum; i++ {
		if m[i] != 1 {
			t.Fatalf("id:%d not arrange\n", i)
		}
	}
	if m[int32(-1)] != 1000-configs.MaxNormalGamePlayerNum {
		t.Fatalf("not arrange enough -1\n")
	}
}

func TestVector2(t *testing.T) {
	v1 := gt.NewVector2(1, 2)
	fmt.Printf("%v\n", v1.Add(gt.NewVector2(1, 2)))
	fmt.Printf("%v\n", v1)
}

func MaxInt64(i, j int64) int64 {
	if i < j {
		return j
	}
	return i
}
func AbsInt64(i int64) int64 {
	if i < 0 {
		return -i
	}
	return i
}

func TestNanoSecond(t *testing.T) {
	max := int64(0)
	total := int64(0)
	const testTime = 10000
	for k := 0; k < testTime; k++ {
		i := time.Now().UnixNano()
		f := float64(i)
		miss := i - int64(f)
		max = MaxInt64(miss, max)
		total += AbsInt64(miss)
		time.Sleep(1 * time.Nanosecond)
	}
	fmt.Printf("max:%v,avg:%v\n", max, total/testTime)
}
