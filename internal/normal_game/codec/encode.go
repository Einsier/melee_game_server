package codec

import (
	"melee_game_server/api/proto"
)

/**
*@Author Sly
*@Date 2022/1/20
*@Version 1.0
*@Description:将Notify,Response,Unicast封装成TopMessage类型
 */

//Encode 将*Response等转换成*proto.TopMessage,用于发送
func Encode(msg interface{}) *proto.TopMessage {
	topMsg := new(proto.TopMessage)
	switch msg.(type) {
	//response
	case *proto.PlayerEnterGameResponse:
		resp := proto.Response{
			ResponseCode:            proto.ResponseCode_PlayerEnterGameResponseCode,
			PlayerEnterGameResponse: msg.(*proto.PlayerEnterGameResponse),
		}
		topMsg.TopMessageType = proto.TopMessageType_ResponseType
		topMsg.Response = &resp
		return topMsg
	case *proto.PlayerQuitGameResponse:
		resp := proto.Response{
			ResponseCode:           proto.ResponseCode_PlayerQuitGameResponseCode,
			PlayerQuitGameResponse: msg.(*proto.PlayerQuitGameResponse),
		}
		topMsg.TopMessageType = proto.TopMessageType_ResponseType
		topMsg.Response = &resp
		return topMsg
	case *proto.HeroGetPropResponse:
		resp := proto.Response{
			ResponseCode:        proto.ResponseCode_HeroGetPropResponseCode,
			HeroGetPropResponse: msg.(*proto.HeroGetPropResponse),
		}
		topMsg.TopMessageType = proto.TopMessageType_ResponseType
		topMsg.Response = &resp
		return topMsg
	case *proto.PlayerHeartBeatResponse:
		resp := proto.Response{
			ResponseCode:            proto.ResponseCode_PlayerHeartBeatResponseCode,
			PlayerHeartBeatResponse: msg.(*proto.PlayerHeartBeatResponse),
		}
		topMsg.TopMessageType = proto.TopMessageType_ResponseType
		topMsg.Response = &resp
		return topMsg

	//broadcast
	case *proto.HeroChangeHealthBroadcast:
		broadcast := proto.Broadcast{
			BroadcastCode:             proto.BroadcastCode_HeroChangeHealthBroadcastCode,
			HeroChangeHealthBroadcast: msg.(*proto.HeroChangeHealthBroadcast),
		}
		topMsg.TopMessageType = proto.TopMessageType_BroadcastType
		topMsg.Broadcast = &broadcast
		return topMsg
	case *proto.HeroMovementChangeBroadcast:
		broadcast := proto.Broadcast{
			BroadcastCode:               proto.BroadcastCode_HeroMovementChangeBroadcastCode,
			HeroMovementChangeBroadcast: msg.(*proto.HeroMovementChangeBroadcast),
		}
		topMsg.TopMessageType = proto.TopMessageType_BroadcastType
		topMsg.Broadcast = &broadcast
		return topMsg
	case *proto.HeroPositionReportBroadcast:
		broadcast := proto.Broadcast{
			BroadcastCode:               proto.BroadcastCode_HeroPositionReportBroadcastCode,
			HeroPositionReportBroadcast: msg.(*proto.HeroPositionReportBroadcast),
		}
		topMsg.TopMessageType = proto.TopMessageType_BroadcastType
		topMsg.Broadcast = &broadcast
		return topMsg
	case *proto.HeroPropDeleteBroadcast:
		broadcast := proto.Broadcast{
			BroadcastCode:           proto.BroadcastCode_HeroPropDeleteBroadcastCode,
			HeroPropDeleteBroadcast: msg.(*proto.HeroPropDeleteBroadcast),
		}
		topMsg.TopMessageType = proto.TopMessageType_BroadcastType
		topMsg.Broadcast = &broadcast
		return topMsg
	case *proto.HeroPropAddBroadcast:
		broadcast := proto.Broadcast{
			BroadcastCode:        proto.BroadcastCode_HeroPropAddBroadcastCode,
			HeroPropAddBroadcast: msg.(*proto.HeroPropAddBroadcast),
		}
		topMsg.TopMessageType = proto.TopMessageType_BroadcastType
		topMsg.Broadcast = &broadcast
		return topMsg
	case *proto.HeroAttackBroadcast:
		broadcast := proto.Broadcast{
			BroadcastCode:       proto.BroadcastCode_HeroAttackBroadcastCode,
			HeroAttackBroadcast: msg.(*proto.HeroAttackBroadcast),
		}
		topMsg.TopMessageType = proto.TopMessageType_BroadcastType
		topMsg.Broadcast = &broadcast
		return topMsg
	case *proto.GameStartBroadcast:
		broadcast := proto.Broadcast{
			BroadcastCode:      proto.BroadcastCode_GameStartBroadcastCode,
			GameStartBroadcast: msg.(*proto.GameStartBroadcast),
		}
		topMsg.TopMessageType = proto.TopMessageType_BroadcastType
		topMsg.Broadcast = &broadcast
		return topMsg
	}
	return nil
}

//EncodeRequest 将Request类型转换成TopMessage类型,用于排除网络的测试
func EncodeRequest(msg interface{}) *proto.TopMessage {
	topMsg := new(proto.TopMessage)
	switch msg.(type) {
	case *proto.PlayerEnterGameRequest:
		req := proto.Request{
			RequestCode:            proto.RequestCode_PlayerEnterGameRequestCode,
			PlayerEnterGameRequest: msg.(*proto.PlayerEnterGameRequest),
		}
		topMsg.TopMessageType = proto.TopMessageType_RequestType
		topMsg.Request = &req
		return topMsg
	}
	return nil
}
