# etcd封装

## 库

- `go get go.etcd.io/etcd/client/v3`
- `go mod edit -replace google.golang.org/grpc=google.golang.org/grpc@v1.26.0`
	- 这里需要限定grpc的版本,最新版本编译不过去

## 已知问题

- [ ] 配置文件中添加账号密码后,关闭连接时会出现报错,暂时不知道怎么解决,看起来没什么影响,先不管了,错误信息:

```json
{"level":"warn","ts":"2023-07-14T10:10:03.866+0800","caller":"clientv3/retry_interceptor.go:62","msg":"retrying of unary invoker failed","target":"endpoint://client-39f6d824-64c2-476a-9b32-e57350d0c8b3/localhost:60001","attempt":0,"error":"rpc error: code = Canceled desc = context canceled"}
{"level":"error","ts":"2023-07-14T10:10:03.866+0800","caller":"clientv3/retry_interceptor.go:114","msg":"clientv3/retry_interceptor: getToken failed","error":"context canceled","stacktrace":"go.etcd.io/etcd/clientv3.(*Client).streamClientInterceptor.func1\n\tC:/Users/pan/go/pkg/mod/go.etcd.io/etcd@v3.3.27+incompatible/clientv3/retry_interceptor.go:114\ngoogle.golang.org/grpc.(*ClientConn).NewStream\n\tC:/Users/pan/go/pkg/mod/google.golang.org/grpc@v1.26.0/stream.go:148\ngithub.com/coreos/etcd/etcdserver/etcdserverpb.(*leaseClient).LeaseKeepAlive\n\tC:/Users/pan/go/pkg/mod/github.com/coreos/etcd@v3.3.27+incompatible/etcdserver/etcdserverpb/rpc.pb.go:6590\ngo.etcd.io/etcd/clientv3.(*retryLeaseClient).LeaseKeepAlive\n\tC:/Users/pan/go/pkg/mod/go.etcd.io/etcd@v3.3.27+incompatible/clientv3/retry.go:152\ngo.etcd.io/etcd/clientv3.(*lessor).resetRecv\n\tC:/Users/pan/go/pkg/mod/go.etcd.io/etcd@v3.3.27+incompatible/clientv3/lease.go:472\ngo.etcd.io/etcd/clientv3.(*lessor).recvKeepAliveLoop\n\tC:/Users/pan/go/pkg/mod/go.etcd.io/etcd@v3.3.27+incompatible/clientv3/lease.go:437"}
```
