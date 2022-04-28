package normal_game_type_configs

import "time"

/**
*@Author Sly
*@Date 2022/1/21
*@Version 1.0
*@Description:
 */

//game_room当前的状态
const (
	NormalGameIdleStatus           = iota //被创建,但没有被初始化
	NormalGameInitStatus                  //已经被初始化
	NormalGameWaitPlayerStatus            //等待Player的到来
	NormalGameStartStatus                 //全部Player已经到来,游戏开始
	NormalGameGameDestroyingStatus        //全部Player狗带,游戏结束
)

//FrameSyncSlice 状态同步的间隔,当前暂定每秒同步20次
var FrameSyncSlice = 80 * time.Millisecond
