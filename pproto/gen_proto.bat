:: 生成protobuf

@echo off

:: 参数
set "binProtoc=protoc.exe"
set "pathIn=E:\\pan\\go_game_server\\example\\proto_src"
set "pathOut=E:\\pan\\go_game_server\\example\\proto"

:: 导出
for %%i in (%pathIn%\\*.proto) do (
	echo export %%i
	%binProtoc% --proto_path=%pathIn% --go_out=%pathOut% --go_opt=paths=source_relative %%i
)

@REM pause
