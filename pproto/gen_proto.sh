#!/bin/bash

# 生成protobuf

# 参数
binProtoc="./protoc.exe"
pathIn="/E/pan/go_game_server/example/proto_src"
pathOut="/E/pan/go_game_server/example/proto"

# proto源文件
protoFiles=$(ls $pathIn/*.proto)

# 导出
for file in "${protoFiles[@]}"; do
	$binProtoc --go_out=$pathOut --proto_path=$pathIn $file
done
