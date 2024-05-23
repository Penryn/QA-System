build:
	go build main.go

build-linux:
	SET CGO_ENABLE=0
	SET GOOS=linux
	SET GOARCH=amd64
	go build main.go


.PHONY: build build-linux