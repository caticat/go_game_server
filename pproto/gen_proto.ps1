# 生成protobuf

# 参数
$binProtoc=".\\protoc.exe"
$pathIn="E:\\pan\\go_game_server\\example\\proto_src"
$pathOut="E:\\pan\\go_game_server\\example\\proto"

# 导出
Get-ChildItem $pathIn\\*.proto | ForEach-Object -Process{
	Invoke-Expression "$binProtoc --go_out=$pathOut --proto_path=$pathIn $_"
}
