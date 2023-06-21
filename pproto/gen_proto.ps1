# 生成protobuf

# 参数
$binProtoc=".\\protoc.exe"
$pathIn="E:\\pan\\go_game_server\\example\\proto_src"
$pathOut="E:\\pan\\go_game_server\\example\\proto"

# 导出
Get-ChildItem $pathIn\\*.proto | ForEach-Object -Process{
	Invoke-Expression "$binProtoc --proto_path=$pathIn --go_out=$pathOut --go_opt=paths=source_relative $_"
}
