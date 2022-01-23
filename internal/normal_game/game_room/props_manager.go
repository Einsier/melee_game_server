package game_room

import (
	pb "melee_game_server/api/proto"
	gt "melee_game_server/internal/normal_game/game_type"
	"sync"
)

/**
*@Author Sly
*@Date 2022/1/17
*@Version 1.0
*@Description:
 */

type PropsManager struct {
	props     map[int32]*gt.Prop
	idCounter int32
	lock      sync.Mutex
}

func NewPropsManager() *PropsManager {
	return &PropsManager{
		props:     make(map[int32]*gt.Prop),
		idCounter: 0,
	}
}

func (m *PropsManager) AddProp(p *gt.Prop) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.props[p.Id] = p
}

func (m *PropsManager) DeleteProp(id int32) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.props, id)
}

//EatProp 英雄吃道具,如果吃到返回true和道具的类型,如果没吃到返回false和-1
func (m *PropsManager) EatProp(id int32) (pb.PropType, bool) {
	m.lock.Lock()
	defer m.lock.Unlock()
	p, ok := m.props[id]
	if !ok {
		return -1, false
	}
	delete(m.props, id)
	return p.PropType, true
}

//AddRandomProp 可以作为time_event_callback
func (m *PropsManager) AddRandomProp() {

}
