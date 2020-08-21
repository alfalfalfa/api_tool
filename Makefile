build-mac:
# 	packr build -o ./bin/mac/api_tool -ldflags="-w"
	go build -o ./bin/mac/api_tool -ldflags="-w"

build-linux:
	packr build -o ./bin/linux/api_tool -ldflags="-w"
#	go build -o ../../bin/mac/api_tool -ldflags="-w"
#	upx ../../bin/mac/api_tool
#	ls -altrh ../../bin/mac

