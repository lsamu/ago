git:
	./_git.sh

github:
	git config --global https.proxy http://127.0.0.1:1080
	git config --global https.proxy https://127.0.0.1:1080
	git config --global http.proxy 'socks5://127.0.0.1:1080'
	git config --global https.proxy 'socks5://127.0.0.1:1080'

	git push -f https://github.com/lsamu/ago.git main

	git config --global --unset http.proxy
	git config --global --unset https.proxy

github2:
	proxychains4 git push -f https://github.com/lsamu/ago.git main

lint:
	golint ./...

local-rest:
	go run -mod=vendor ./examples/rest/main.go

local-rpc-server:
	go run -mod=vendor ./examples/rpc.server/main.go

local-sock-server:
	go run -mod=vendor ./examples/sock.server/main.go

local-sock-client:
	go run -mod=vendor ./examples/sock.client/main.go

local-cron:
	go run -mod=vendor ./examples/cron/main.go

install:
	go get -u golang.org/x/lint/golint
	go get -u github.com/securego/gosec/v2/cmd/gosec
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get -u github.com/mgechev/revive


