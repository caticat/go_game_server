# etcd封装

## 库

- `go get go.etcd.io/etcd/client/v3`
- `go mod edit -replace google.golang.org/grpc=google.golang.org/grpc@v1.26.0`
	- 这里需要限定grpc的版本,最新版本编译不过去
