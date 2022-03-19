package server

/**
*@Author Sly
*@Date 2022/3/17
*@Version 1.0
*@Description:
 */

import (
	"context"
	"encoding/json"
	"go.etcd.io/etcd/client/v3"
	"log"
	"melee_game_server/api/hall"
	"melee_game_server/configs"
	"melee_game_server/plugins/logger"
	"time"
)

var EtcdCli = NewEtcdCli()

func NewEtcdCli() *clientv3.Client {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("创建etcd client时出错!err:%s", err.Error())
	}

	kv := clientv3.NewKV(cli)
	ctx, _ := context.WithTimeout(context.TODO(), time.Second*2)
	_, err = kv.Get(ctx, "/test")
	if err != nil {
		log.Fatalf("连接etcd时出错!err:%s", err.Error())
	}

	logger.Infof("已创建etcd连接")
	return cli
}

func SetAccountToEtcd(gameId string, info *hall.GameAccountInfo) {
	lease := clientv3.NewLease(EtcdCli)

	ctx, _ := context.WithTimeout(context.TODO(), time.Second*10)
	grant, err := lease.Grant(ctx, 10)
	if err != nil {
		logger.Errorf("无法从etcd申请grant!")
		return
	}

	kv := clientv3.NewKV(EtcdCli)
	//向etcd放对局信息,过10s就会过期,所以需要尽快处理
	infoBytes, err := json.Marshal(info)
	if err != nil {
		logger.Errorf("向etcd上传结算信息时json转码出现问题:%s", err.Error())
		return
	}
	_, err = kv.Put(ctx, configs.AccountPath+"/"+gameId, string(infoBytes), clientv3.WithLease(grant.ID))
	if err != nil {
		logger.Errorf("向etcd上传结算信息时出现问题:%s", err.Error())
	}
}
