# ETCD客户端GUI

## Fyne

### 地址

- [说明文档](https://developer.fyne.io/)
- [代码库](https://github.com/fyne-io/fyne/tree/master)

### 库安装流程

- MingW-w64,`https://www.msys2.org/`
	- `MSYS2 MinGW 64-bit`
		- 安装库
			- `pacman -Syu`
			- `pacman -S git mingw-w64-x86_64-toolchain`
		- 添加path
			- `echo "export PATH=\$PATH:/c/Program\ Files/Go/bin:~/Go/bin" >> ~/.bashrc`
- 开启CGo
	- `go env -w CGO_ENABLED=1`

## go-app

### 说明

- go-app的方式实现失败,应该是etcd的引用库和gui的引用库版本冲突,没找到解决方法`https://github.com/maxence-charriere/go-app`
- 需要生成客户端合服务器才能执行

### 引用库

- `go get -u github.com/maxence-charriere/go-app/v9/pkg/app`

### 生成可执行程序

- 可以增加脚本执行生成

```bash
# build client
GOARCH=wasm GOOS=js go build -o ./build/web/app.wasm

# build server
go build -o ./build/
```

### 运行

- 运行服务器
	- `./build/${exe}`
- 调用
	- `http://localhost:${port}/`
