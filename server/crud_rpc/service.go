package crud_rpc

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"go_crud/utils/etcd_center"
)

var ETCD *clientv3.Client

func Init() {
	ETCD = etcd_center.NewEtcdClient("etcd")
}
