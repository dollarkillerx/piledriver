build_agent:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o piledriver_core_darwin -ldflags "-s -w" client/client.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o piledriver_core_linux -ldflags "-s -w" client/client.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o piledriver_core.exe -ldflags "-s -w -H windowsgui" client/client.go
	upx piledriver_core_darwin
	upx piledriver_core_linux
	upx piledriver_core.exe

Generate:
	@echo 'Build GRPC'
	protoc -I rpc/proto/ rpc/proto/*.proto --go_out=plugins=grpc:rpc/proto/.

generate_test:
	@echo 'generate_test'
	protoc -I test_simple/proto/ test_simple/proto/*.proto --go_out=plugins=grpc:test_simple/proto/.
	protoc -I test_stream/proto/ test_stream/proto/*.proto --go_out=plugins=grpc:test_stream/proto/.


SSLKey:
	@echo 'SSLKey'
	openssl genrsa -out cert/server.key 2048
	openssl ecparam -genkey -name secp384r1 -out cert/serveryek
	openssl req -new -x509 -sha256 -key cert/serveryek -out cert/servermpe -days 3650  # plumber

