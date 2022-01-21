package game_type

import (
	configs "melee_game_server/configs/normal_game_type_configs"
	"net"
	"sync"
)

/**
*@Author Sly
*@Date 2022/1/20
*@Version 1.0
*@Description:
 */

type Player struct {
	Id         int32
	HeroId     int32
	Nickname   string //玩家昵称
	Level      int    //玩家等级
	Score      int    //玩家得分
	Conn       *net.Conn
	Status     int32 //玩家当前状态
	statusLock sync.RWMutex
}

//NewPlayer 如果该Player还没有登录游戏,则heroId应该为-1
func NewPlayer(id int32) *Player {
	return &Player{
		Id:     id,
		HeroId: -1,
		Conn:   nil,
		Status: configs.PlayerNotRegisteredStatus,
	}
}

func (p *Player) SetStatus(s int32) {
	p.statusLock.Lock()
	defer p.statusLock.Unlock()
	p.Status = s
}

func (p *Player) GetStatus() int32 {
	p.statusLock.RLock()
	defer p.statusLock.RUnlock()
	return p.Status
}

func (p *Player) BindHeroId(hId int32) {
	p.HeroId = hId
}

func (p *Player) BindConn(conn *net.Conn) {
	p.Conn = conn
}
