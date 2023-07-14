# 框架

[TOC]

## 待制作

- [x] 日志完善
	- [x] 日志等级
	- [x] 支持的日志类型
		- [x] 控制台日志
		- [x] 文件日志(按小时拆分)
- [x] 通信完善
	- [x] 客户端服务器通信
	- [x] 服务器间通信
	- [x] socket参数补充完善
		- [x] 没有选项要写...
	- [x] 配置整理
	- [x] 与ETCD整合
- [x] etcd支持
	- [x] 添加数据
		- [x] 永久数据
		- [x] 进程存在有效数据
	- [x] 删除数据
	- [x] 监听数据
	- [x] etcd配置文件初始化
	- [ ] GUI修改数据支持
		- [x] go-app的方式实现失败,应该是etcd的引用库和gui的引用库版本冲突,没找到解决方法`https://github.com/maxence-charriere/go-app`
	- [ ] GUI
		- [x] 查询
			- [x] 当前只有选择查询
		- [x] 添加
		- [x] 删除
			- [x] 单点删除
			- [x] ~~子节点递归删除~~,不做了
		- [x] 修改
		- [x] 脱离配置文件
			- [x] 程序内部保存数据
		- [x] 多ETCD连接
		- [x] 中文支持
		- [x] 日志展示
		- [x] 认证
		- [ ] 代码整理
- [x] 配置文件支持
	- 没有封装,使用的时候自己写
- [x] 命令行参数
	- 没有封装,使用的时候自己写
- [x] Protobuf支持
	- [x] 使用`pproto`目录下的任意脚本
		- `gen_proto.bat`
		- `gen_proto.ps1`
		- `gen_proto.sh`
	- [x] 修改对应脚本的参数生成`proto`代码文件

## 扩展安装

- petcd
	- `go get go.etcd.io/etcd/client/v3`
- petcd_gui
	- `go get -u fyne.io/fyne/v2`
	- 已废弃
		- `go get -u github.com/maxence-charriere/go-app/v9/pkg/app`

## 常用命令

- 整理`go.work`格式
	- `go work edit -fmt`
- 替换库引用
	- `go mod edit -replace github.com/coreos/bbolt@v1.3.4=go.etcd.io/bbolt@v1.3.4`
	- `go mod edit -replace google.golang.org/grpc=google.golang.org/grpc@v1.26.0`
