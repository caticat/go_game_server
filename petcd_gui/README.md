# ETCD客户端GUI

## 功能说明

- 使用`Fyne`实现的ETCD的GUI
- 支持简单的ETCD数据操作
- 注意:**由于是全字段`/`读取,操作数据量大时可能会出现问题**,没有过相应测试

## 打包

- `fyne package -os windows -icon assets/myapp.png`

### 打包`android`的apk

- 下载安装[Android Studio](https://developer.android.com/studio/index.html)
- 在`Android Studio`中下载安装NDK(Side by Side)
	- Tool
		- SDK Manager
			- SDK Tools
				- NDK(Side by Side)
- 添加环境变量
	- `PATH`,这个是否可以改成Android相关的环境变量名我还没试过
		- `C:\Users\pan\AppData\Local\Android\Sdk\platform-tools`
	- `ANDROID_NDK_HOME`
		- `C:\Users\pan\AppData\Local\Android\Sdk\ndk\25.2.9519653`
- 打包
	- `fyne package -os android -appID com.github.caticat.go_game_server.petcd_gui -icon assets/myapp.png`

#### 已知问题

- 安卓模拟器中运城程序后,切换其他应用程序再切换回来就黑屏了,不知道什么原因

### 打包`ios`,未测试

- 打包`fyne package -os ios -appID com.github.caticat.go_game_server.petcd_gui -icon assets/myapp.png`

## 测试

- `winpty docker run -it --rm --network pan_network bitnami/etcd:latest etcdctl --endpoints http://pan_etcd_c1:2379,http://pan_etcd_c2:2379,http://pan_etcd_c3:2379 get "" --prefix`
- `winpty docker run -it --rm --network pan_network bitnami/etcd:latest etcdctl --endpoints http://pan_etcd:2379 get "" --prefix`

## 版本

### v0.0.2

- 数据导入
	- 数据一般导入
	- 数据初始化导入
		- 清空以前的数据
		- 清空以前的版本信息
		- 数据库占用控件回收
- 数据导出
- 日志布点完善

### v0.0.1

- 查询
	- 全数据库读取
	- 快速定位key
- 添加
- 修改
- 删除
	- 单点删除
	- ~~子节点递归删除~~,不做了
- APP内部配置记录
- 多ETCD支持
	- 连接配置记录
	- 当前连接选择
- 中文支持
- 日志展示
- 认证
	- 账号名,密码认证
	- **不支持**证书认证,没用过,不做了
- 信息界面

## 其他

### 引用库`Fyne`

#### 地址

- [说明文档](https://developer.fyne.io/)
- [代码库](https://github.com/fyne-io/fyne/tree/master)

#### 库安装流程

- MingW-w64,`https://www.msys2.org/`
	- `MSYS2 MinGW 64-bit`
		- 安装库
			- `pacman -Syu`
			- `pacman -S git mingw-w64-x86_64-toolchain`
		- 添加path
			- `echo "export PATH=\$PATH:/c/Program\ Files/Go/bin:~/Go/bin" >> ~/.bashrc`
- 开启CGo
	- `go env -w CGO_ENABLED=1`
