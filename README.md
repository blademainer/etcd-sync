# etcd异步同步
- 支持自定义fixed（固定key）、prefix（前缀匹配），range（范围）的同步
- 支持配置多个目标源

## 功能
- 第一次初始化时，会同步master配置的fixed、prefix、range等配置到slaves
- 启动完成后会启动goroutine来watch etcd的变更，有接收到etcd事件会将事件推送到每个slave的channel内
- 每个slave会启动一个goroutine来接收channel内的事件，每个`PUT`、`DELETE`事件都会`PUT`、`DELETE`对相应的key进行更新、删除操作
- 如果跟某个slave的网络出现超时，则会把slave设置为`disconnected`，健康检查会每隔10秒重连，重连成功会从master同步到这个刚刚恢复的slave

## 注意事项
1. etcd-sync目前不允许存在多个master，如果互为主备，会出现死循环复制
2. 如果网络出现问题，可能会存在脏数据（在slave断开的这段时间，master删除了某条数据，恢复之后即使全量同步也无法删除这条数据）

# 测试
启动docker
```shell
export HostIP="127.0.0.1"
docker run --restart always -d -v /usr/share/ca-certificates/:/etc/ssl/certs -p 14001:4001 -p 12380:2380 -p 12379:2379 \
 -e ETCDCTL_API=3 \
 --name etcd quay.io/coreos/etcd:v3.2.18 \
 etcd \
 -name etcd0 \
 -advertise-client-urls http://${HostIP}:2379,http://${HostIP}:4001 \
 -listen-client-urls http://0.0.0.0:2379,http://0.0.0.0:4001 \
 -initial-advertise-peer-urls http://${HostIP}:12380 \
 -listen-peer-urls http://0.0.0.0:2380 \
 -initial-cluster-token etcd-cluster-1 \
 -initial-cluster etcd0=http://${HostIP}:12380 \
 -initial-cluster-state new

docker run --restart always -d -v /usr/share/ca-certificates/:/etc/ssl/certs -p 24001:4001 -p 22380:2380 -p 22379:2379 \
 -e ETCDCTL_API=3 \
 --name etcd2 quay.io/coreos/etcd:v3.2.18 \
 etcd \
 -name etcd0 \
 -advertise-client-urls http://${HostIP}:2379,http://${HostIP}:4001 \
 -listen-client-urls http://0.0.0.0:2379,http://0.0.0.0:4001 \
 -initial-advertise-peer-urls http://${HostIP}:22380 \
 -listen-peer-urls http://0.0.0.0:2380 \
 -initial-cluster-token etcd-cluster-1 \
 -initial-cluster etcd0=http://${HostIP}:22380 \
 -initial-cluster-state new

docker run --restart always -d -v /usr/share/ca-certificates/:/etc/ssl/certs -p 34001:4001 -p 32380:2380 -p 32379:2379 \
 -e ETCDCTL_API=3 \
 --name etcd3 quay.io/coreos/etcd:v3.2.18 \
 etcd \
 -name etcd0 \
 -advertise-client-urls http://${HostIP}:2379,http://${HostIP}:4001 \
 -listen-client-urls http://0.0.0.0:2379,http://0.0.0.0:4001 \
 -initial-advertise-peer-urls http://${HostIP}:32380 \
 -listen-peer-urls http://0.0.0.0:2380 \
 -initial-cluster-token etcd-cluster-1 \
 -initial-cluster etcd0=http://${HostIP}:32380 \
 -initial-cluster-state new

启动程序...

docker exec -it etcd etcdctl put /com/a/go_test '{"name":"zhangsan"}'
docker exec -it etcd etcdctl put /com/xiongyingqi 'xxxxxxx'
docker exec -it etcd etcdctl put /com/xc 'xcccc'
docker exec -it etcd etcdctl put /com/bdddzz/xxx 'dddd'
```