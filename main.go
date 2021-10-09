package main

import (
	"context"
	"etcd/exception"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

var config clientv3.Config

const pre = "/cron/jobs/"

func main() {

	config = clientv3.Config{
		Endpoints:   []string{"http://127.0.0.1:4379", "http://127.0.0.1:2379", "http://127.0.0.1:3379"},
		DialTimeout: time.Second * 5,
	}

	client, err := clientv3.New(config)
	exception.Handler(err)

	kv := clientv3.NewKV(client)
	// kv.Delete(context.TODO(), pre, clientv3.WithPrefix())
	for i := 0; i < 10; i++ {
		_, err := kv.Put(context.TODO(), fmt.Sprintf("%s%d", pre, i), "bye"+time.Now().String(), clientv3.WithPrevKV())
		exception.Handler(err)
	}
	fmt.Println("add success")

	go func() {
		for {
			kv.Put(context.TODO(), "/cron/jobs/7", "localhost:8080")
			kv.Delete(context.TODO(), "/cron/jobs/7")
			time.Sleep(3 * time.Second)
		}
	}()

	getResp, err := kv.Get(context.TODO(), "/cron/jobs/7")
	exception.Handler(err)

	// 当前etcd集群事务ID, 单调递增的
	watchStartRevision := getResp.Header.Revision + 1

	// 创建一个watcher
	watcher := clientv3.NewWatcher(client)

	// 启动监听
	fmt.Println("从该版本向后监听:", watchStartRevision)

	ctx, cancelFunc := context.WithCancel(context.TODO())
	time.AfterFunc(20*time.Second, func() {
		cancelFunc()
	})
	ops := []clientv3.OpOption{clientv3.WithRev(watchStartRevision)}
	watchRespChan := watcher.Watch(ctx, "/cron/jobs/7", ops...)

	// 处理kv变化事件
	for resp := range watchRespChan {
		for _, ev := range resp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				fmt.Println("handler put")
			case mvccpb.DELETE:
				fmt.Println("handler delete")
			}
		}
	}
}
