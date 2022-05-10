package metrics

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// 定义游戏指标
	GaugeVecGameRoomPlayerCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "game_room_player_count",
		Help: "当前房间游戏人数",
	}, []string{"room"})

	// GaugeVecGameRoomPlayerCount.WithLabelValues("$room_id").Inc() 游戏指定房间玩家数量 +1
	// GaugeVecGameRoomPlayerCount.WithLabelValues("$room_id").Dec() 游戏指定房间玩家数量 -1
	// GaugeVecGameRoomPlayerCount.WithLabelValues("$room_id").Add(n) 游戏指定房间玩家数量 +n
	// GaugeVecGameRoomPlayerCount.WithLabelValues("$room_id").Sub(n) 游戏指定房间玩家数量 -n
	// GaugeVecGameRoomPlayerCount.WithLabelValues("$room_id").Set(n) 游戏指定房间玩家数量 n
)

func Start() {
	//将指标注册到 Prometheus 默认仓库
	prometheus.MustRegister(GaugeVecGameRoomPlayerCount)

	// Serve the default Prometheus metrics registry over HTTP on /metrics.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatalln(http.ListenAndServe(":8888", nil))
}
