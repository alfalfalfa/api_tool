init:
	go get -u "github.com/alfalfalfa/go-jsschema"
	go get -u "github.com/tealeg/xlsx"
    go get -u "github.com/docopt/docopt-go"
    go get -u "github.com/flosch/pongo2"
    go get -u "github.com/jinzhu/inflection"
    go get -u github.com/gobuffalo/packr/...

build:
	packr build -o ../../bin/mac/api_tool -ldflags="-w"
#	go build -o ../../bin/mac/api_tool -ldflags="-w"
#	upx ../../bin/mac/api_tool
#	ls -altrh ../../bin/mac

