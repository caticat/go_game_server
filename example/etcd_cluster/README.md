# ETCD的Docker集群配置

- 启动ETCD服务器

```bash
docker run -d \
	--name pan_etcd_c1 \
	--network pan_network \
	--network-alias pan_etcd_c1 \
	--env ALLOW_NONE_AUTHENTICATION=yes \
	--env ETCD_NAME=pan_etcd_c1 \
	--env ETCD_INITIAL_ADVERTISE_PEER_URLS=http://pan_etcd_c1:2380 \
	--env ETCD_LISTEN_PEER_URLS=http://0.0.0.0:2380 \
	--env ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:2379 \
	--env ETCD_ADVERTISE_CLIENT_URLS=http://pan_etcd_c1:2379 \
	--env ETCD_INITIAL_CLUSTER_TOKEN=pan-etcd-cluster \
	--env ETCD_INITIAL_CLUSTER=pan_etcd_c1=http://pan_etcd_c1:2380,pan_etcd_c2=http://pan_etcd_c2:2380,pan_etcd_c3=http://pan_etcd_c3:2380 \
	--env ETCD_INITIAL_CLUSTER_STATE=new \
	-v F:/pan/var/docker/etcd_c1/data:/bitnami/etcd/data \
	-p 60001:2379 \
	bitnami/etcd:latest

docker run -d \
	--name pan_etcd_c2 \
	--network pan_network \
	--network-alias pan_etcd_c2 \
	--env ALLOW_NONE_AUTHENTICATION=yes \
	--env ETCD_NAME=pan_etcd_c2 \
	--env ETCD_INITIAL_ADVERTISE_PEER_URLS=http://pan_etcd_c2:2380 \
	--env ETCD_LISTEN_PEER_URLS=http://0.0.0.0:2380 \
	--env ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:2379 \
	--env ETCD_ADVERTISE_CLIENT_URLS=http://pan_etcd_c2:2379 \
	--env ETCD_INITIAL_CLUSTER_TOKEN=pan-etcd-cluster \
	--env ETCD_INITIAL_CLUSTER=pan_etcd_c1=http://pan_etcd_c1:2380,pan_etcd_c2=http://pan_etcd_c2:2380,pan_etcd_c3=http://pan_etcd_c3:2380 \
	--env ETCD_INITIAL_CLUSTER_STATE=new \
	-v F:/pan/var/docker/etcd_c2/data:/bitnami/etcd/data \
	-p 60002:2379 \
	bitnami/etcd:latest

docker run -d \
	--name pan_etcd_c3 \
	--network pan_network \
	--network-alias pan_etcd_c3 \
	--env ALLOW_NONE_AUTHENTICATION=yes \
	--env ETCD_NAME=pan_etcd_c3 \
	--env ETCD_INITIAL_ADVERTISE_PEER_URLS=http://pan_etcd_c3:2380 \
	--env ETCD_LISTEN_PEER_URLS=http://0.0.0.0:2380 \
	--env ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:2379 \
	--env ETCD_ADVERTISE_CLIENT_URLS=http://pan_etcd_c3:2379 \
	--env ETCD_INITIAL_CLUSTER_TOKEN=pan-etcd-cluster \
	--env ETCD_INITIAL_CLUSTER=pan_etcd_c1=http://pan_etcd_c1:2380,pan_etcd_c2=http://pan_etcd_c2:2380,pan_etcd_c3=http://pan_etcd_c3:2380 \
	--env ETCD_INITIAL_CLUSTER_STATE=new \
	-v F:/pan/var/docker/etcd_c3/data:/bitnami/etcd/data \
	-p 60003:2379 \
	bitnami/etcd:latest
```

- 使用`docker-compose`,`docker-compose.yml`没用过这个
	- 使用方法

		```bash
		curl -LO https://raw.githubusercontent.com/bitnami/containers/main/bitnami/etcd/docker-compose.yml
		docker-compose up
		```

```yml
# Copyright VMware, Inc.
# SPDX-License-Identifier: APACHE-2.0

version: '2'

services:
  etcd:
    image: docker.io/bitnami/etcd:3.5
    environment:
      - ALLOW_NONE_AUTHENTICATION=yes
    volumes:
      - etcd_data:/bitnami/etcd
volumes:
  etcd_data:
    driver: local
```

```yml
version: '2'

services:
  etcd1:
    image: docker.io/bitnami/etcd:3
    environment:
      - ALLOW_NONE_AUTHENTICATION=yes
      - ETCD_NAME=etcd1
      - ETCD_INITIAL_ADVERTISE_PEER_URLS=http://etcd1:2380
      - ETCD_LISTEN_PEER_URLS=http://0.0.0.0:2380
      - ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:2379
      - ETCD_ADVERTISE_CLIENT_URLS=http://etcd1:2379
      - ETCD_INITIAL_CLUSTER_TOKEN=etcd-cluster
      - ETCD_INITIAL_CLUSTER=etcd1=http://etcd1:2380,etcd2=http://etcd2:2380,etcd3=http://etcd3:2380
      - ETCD_INITIAL_CLUSTER_STATE=new
  etcd2:
    image: docker.io/bitnami/etcd:3
    environment:
      - ALLOW_NONE_AUTHENTICATION=yes
      - ETCD_NAME=etcd2
      - ETCD_INITIAL_ADVERTISE_PEER_URLS=http://etcd2:2380
      - ETCD_LISTEN_PEER_URLS=http://0.0.0.0:2380
      - ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:2379
      - ETCD_ADVERTISE_CLIENT_URLS=http://etcd2:2379
      - ETCD_INITIAL_CLUSTER_TOKEN=etcd-cluster
      - ETCD_INITIAL_CLUSTER=etcd1=http://etcd1:2380,etcd2=http://etcd2:2380,etcd3=http://etcd3:2380
      - ETCD_INITIAL_CLUSTER_STATE=new
  etcd3:
    image: docker.io/bitnami/etcd:3
    environment:
      - ALLOW_NONE_AUTHENTICATION=yes
      - ETCD_NAME=etcd3
      - ETCD_INITIAL_ADVERTISE_PEER_URLS=http://etcd3:2380
      - ETCD_LISTEN_PEER_URLS=http://0.0.0.0:2380
      - ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:2379
      - ETCD_ADVERTISE_CLIENT_URLS=http://etcd3:2379
      - ETCD_INITIAL_CLUSTER_TOKEN=etcd-cluster
      - ETCD_INITIAL_CLUSTER=etcd1=http://etcd1:2380,etcd2=http://etcd2:2380,etcd3=http://etcd3:2380
      - ETCD_INITIAL_CLUSTER_STATE=new
```

- 测试构建是否成功

```bash
# 查看集群节点
winpty docker run -it --rm \
	--network pan_network \
	bitnami/etcd:latest \
	etcdctl \
	--endpoints http://pan_etcd_c1:2379 \
	--write-out table \
	member list

# 测试写数据
winpty docker run -it --rm \
	--network pan_network \
	bitnami/etcd:latest \
	etcdctl \
	--endpoints http://pan_etcd_c1:2379 \
	put //hello "pan's world"

# 测试本地读数据
winpty docker run -it --rm \
	--network pan_network \
	bitnami/etcd:latest \
	etcdctl \
	--endpoints http://pan_etcd_c1:2379 \
	get //hello

# 测试远端读数据
winpty docker run -it --rm \
	--network pan_network \
	bitnami/etcd:latest \
	etcdctl \
	--endpoints http://pan_etcd_c2:2379 \
	get //hello
winpty docker run -it --rm \
	--network pan_network \
	bitnami/etcd:latest \
	etcdctl \
	--endpoints http://pan_etcd_c3:2379 \
	get //hello
```
