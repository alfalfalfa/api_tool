build: build-mac build-linux build-win
build-mac:
	GOOS=darwin GOARCH=amd64 go build "-ldflags=-s -w -buildid=" -trimpath -o ./bin/mac/api_tool

build-linux:
	GOOS=linux GOARCH=amd64 go build "-ldflags=-s -w -buildid=" -trimpath -o ./bin/linux/api_tool

build-win:
	GOOS=windows GOARCH=amd64 go build "-ldflags=-s -w -buildid=" -trimpath -o ./bin/win/api_tool

upx:
	upx --lzma ./bin/mac/api_tool
	upx --lzma ./bin/linux/api_tool
	upx --lzma ./bin/win/api_tool
