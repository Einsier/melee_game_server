package codec

import (
	"melee_game_server/api/client/proto"
	"melee_game_server/plugins/logger"
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
	case *proto.HeroBulletLaunchBroadcast:
		broadcast := proto.Broadcast{
			BroadcastCode:             proto.BroadcastCode_HeroBulletLaunchBroadcastCode,
			HeroBulletLaunchBroadcast: msg.(*proto.HeroBulletLaunchBroadcast),
		}
		topMsg.TopMessageType = proto.TopMessageType_BroadcastType
		topMsg.Broadcast = &broadcast
		return topMsg
	case *proto.HeroSwordAttackBroadcast:
		broadcast := proto.Broadcast{
			BroadcastCode:            proto.BroadcastCode_HeroSwordAttackBroadcastCode,
			HeroSwordAttackBroadcast: msg.(*proto.HeroSwordAttackBroadcast),
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
	case *proto.HeroDeadBroadcast:
		broadcast := proto.Broadcast{
			BroadcastCode:     proto.BroadcastCode_HeroDeadBroadcastCode,
			HeroDeadBroadcast: msg.(*proto.HeroDeadBroadcast),
		}
		topMsg.TopMessageType = proto.TopMessageType_BroadcastType
		topMsg.Broadcast = &broadcast
		return topMsg
	case *proto.HeroBulletDestroyBroadcast:
		broadcast := proto.Broadcast{
			BroadcastCode:              proto.BroadcastCode_HeroBulletDestroyBroadcastCode,
			HeroBulletDestroyBroadcast: msg.(*proto.HeroBulletDestroyBroadcast),
		}
		topMsg.TopMessageType = proto.TopMessageType_BroadcastType
		topMsg.Broadcast = &broadcast
		return topMsg
	default:
		logger.TestErrf("收到了错误的译码请求,%T不是可正确译码的类型", msg)
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
	case *proto.PlayerQuitGameRequest:
		req := proto.Request{
			RequestCode:           proto.RequestCode_PlayerQuitGameRequestCode,
			PlayerQuitGameRequest: msg.(*proto.PlayerQuitGameRequest),
		}
		topMsg.TopMessageType = proto.TopMessageType_RequestType
		topMsg.Request = &req
		return topMsg
	case *proto.HeroMovementChangeRequest:
		req := proto.Request{
			RequestCode:               proto.RequestCode_HeroMovementChangeRequestCode,
			HeroMovementChangeRequest: msg.(*proto.HeroMovementChangeRequest),
		}
		topMsg.TopMessageType = proto.TopMessageType_RequestType
		topMsg.Request = &req
		return topMsg
	case *proto.HeroPositionReportRequest:
		req := proto.Request{
			RequestCode:               proto.RequestCode_HeroPositionReportRequestCode,
			HeroPositionReportRequest: msg.(*proto.HeroPositionReportRequest),
		}
		topMsg.TopMessageType = proto.TopMessageType_RequestType
		topMsg.Request = &req
		return topMsg
	case *proto.HeroBulletLaunchRequest:
		req := proto.Request{
			RequestCode:             proto.RequestCode_HeroBulletLaunchRequestCode,
			HeroBulletLaunchRequest: msg.(*proto.HeroBulletLaunchRequest),
		}
		topMsg.TopMessageType = proto.TopMessageType_RequestType
		topMsg.Request = &req
		return topMsg
	default:
		logger.TestErrf("收到了错误的译码请求,%T不是可正确译码的类型", msg)
	}
	return nil
}

func EncodeUnicast(msg interface{}) *proto.TopMessage {
	topMsg := new(proto.TopMessage)
	switch msg.(type) {
	case *proto.HeroFrameSyncUnicast:
		uni := proto.Unicast{
			UnicastCode:          proto.UnicastCode_HeroFrameSyncUnicastCode,
			HeroFrameSyncUnicast: msg.(*proto.HeroFrameSyncUnicast),
		}
		topMsg.TopMessageType = proto.TopMessageType_UnicastType
		topMsg.Unicast = &uni
		return topMsg
	case *proto.HeroLeaveSightUnicast:
		uni := proto.Unicast{
			UnicastCode:           proto.UnicastCode_HeroLeaveSightUnicastCode,
			HeroLeaveSightUnicast: msg.(*proto.HeroLeaveSightUnicast),
		}
		topMsg.TopMessageType = proto.TopMessageType_UnicastType
		topMsg.Unicast = &uni
		return topMsg
	default:
		logger.TestErrf("EncodeUnicast:收到了错误的译码请求,%T不是可正确译码的类型", msg)
		return nil
	}
}
