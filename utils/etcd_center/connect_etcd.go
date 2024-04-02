package etcd_center

import (
	"context"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"math/rand"
	"time"
)

func NewEtcdClient(etcdName string) *clientv3.Client {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   viper.GetStringSlice(etcdName + ".Endpoints"),
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Printf("Error creating etcd client: %v \n", err)
	}
	return cli
}

func GetAddrFromEtcd(clientName string, etcdClient *clientv3.Client) string {

	// 从etcd获取服务地址
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := etcdClient.Get(ctx, clientName)
	//fmt.Println(resp)
	cancel()
	if err != nil {
		log.Printf("failed to get service address from etcd: %v \n", err)
	}

	if len(resp.Kvs) == 0 {
		log.Println("service address not found in etcd")
	}

	//// 解析服务地址
	//serviceAddr := string(resp.Kvs[0].Value)

	// 随机选择一个服务地址，负载均衡
	rand.Seed(time.Now().UnixNano())
	serviceAddr := string(resp.Kvs[rand.Intn(len(resp.Kvs))].Value)

	return serviceAddr
}

func RegisterAddrToEtcd(serviceName string, serviceAddr string, cli *clientv3.Client) {

	// 注册服务到etcd
	leaseResp, err := cli.Grant(context.Background(), 10)
	if err != nil {
		log.Printf("Error granting lease: %v \n", err)
	}

	// 设置键值对，其中键通常是服务名称，值是服务地址
	putResp, err := cli.Put(context.Background(), serviceName, serviceAddr, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		log.Printf("Error putting service to etcd: %v \n", err)
	}
	log.Println(putResp, serviceName, serviceAddr)

	// 保持心跳，以续约租约
	keepAliveChan, err := cli.KeepAlive(context.Background(), leaseResp.ID)
	if err != nil {
		log.Printf("Error keeping lease alive: %v \n", err)
	}

	for {
		select {
		case ka := <-keepAliveChan:
			if ka == nil {
				log.Println("Lease expired or KeepAlive channel closed")
				return
			}
			//log.Println(ka.String())
		}
	}
}
