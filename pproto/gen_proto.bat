:: 生成protobuf

@echo off

:: 参数
set "binProtoc=protoc.exe"
set "pathIn=E:\\pan\\go_game_server\\example\\proto_src"
set "pathOut=E:\\pan\\go_game_server\\example\\proto"

:: 导出
for %%i in (%pathIn%\\*.proto) do (
	echo export %%i
	%binProtoc% --go_out=%pathOut% --proto_path=%pathIn% %%i
)

@REM pause
