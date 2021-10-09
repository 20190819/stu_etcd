package main

import (
	"context"
	"etcd/exception"
	"fmt"
	"time"

	c3 "github.com/coreos/etcd/clientv3"
)

var config c3.Config
var client *c3.Client

const k1 = "/cron/jobs/job1"

func main() {

	config = c3.Config{
		Endpoints:   []string{"http://127.0.0.1:4379", "http://127.0.0.1:2379", "http://127.0.0.1:3379"},
		DialTimeout: time.Second * 5,
	}

	client, err := c3.New(config)
	exception.Handler(err)

	kv := c3.NewKV(client)
	resp, err := kv.Put(context.TODO(), k1, "bye", c3.WithPrevKV())
	exception.Handler(err)

	//获取版本信息
	fmt.Println("Revision", resp.Header.Revision)
	if resp.PrevKv != nil {
		fmt.Println("key:", string(resp.PrevKv.Key))
		fmt.Println("Value:", string(resp.PrevKv.Value))
		fmt.Println("Version:", rune(resp.PrevKv.Version))
	}

	// 读取键值对

	readResp, err := kv.Get(context.Background(), k1)
	exception.Handler(err)
	fmt.Println("读取信息：", readResp.Kvs)

}
