package helper

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/tests/v3/integration"
	"log"
	"time"
)

func GetKV(key string) bool {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.0.105:2379", "192.168.0.102:2379", "192.168.0.103:2379", "192.168.0.104:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer cli.Close() // 어떤 에러가 발생하더라도 마지막에 Close 된다.

	ctx, cancel := context.WithTimeout(context.Background(), integration.RequestWaitTimeout)
	resp, err := cli.Get(ctx, key)
	cancel()

	if err != nil {
		log.Fatal(err)
		return false
	}
	// Get 성공하면 Key, Value 쌍을 print 하고 True 반환
	// 그 외에는 False 반환
	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
		return true
	}
	return false
}

func PutKV(key string, value string) bool {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close() // 어떤 에러가 발생하더라도 마지막에 Close 된다.
	_, err = cli.Put(context.TODO(), key, value)
	// Put 성공하면 True 반환, 그 외에는 False 반환
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Put 성공\n")
		return true
	}
	return false
}
