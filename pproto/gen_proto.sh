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
	$binProtoc --proto_path=$pathIn --go_out=$pathOut --go_opt=paths=source_relative $file
done
