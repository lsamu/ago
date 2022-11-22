git:
	./_git.sh

github:
	git push -f https://github.com/lsamu/ago.git main

lint:
	golint ./...

local-rest:
	go run -mod=vendor ./examples/rest/main.go

local-rpc-server:
	go run -mod=vendor ./examples/rpc.server/main.go

local-sockio-server:
	go run -mod=vendor ./examples/sockio.server/main.go

local-sockio-client:
	go run -mod=vendor ./examples/sockio.client/main.go

local-cron:
	go run -mod=vendor ./examples/cron/main.go

install:
	go get -u golang.org/x/lint/golint
	go get -u github.com/securego/gosec/v2/cmd/gosec
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get -u github.com/mgechev/revive

local-ago:
	go run -mod=vendor ./tool/main.go datetime

build-ago:
	go build -mod=vendor -o ./tool/ago ./tool/main.go

run-ago:
	./tool/ago datetime help


