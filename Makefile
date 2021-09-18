build_agent:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o piledriver_core_darwin -ldflags "-s -w" client/client.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o piledriver_core_linux -ldflags "-s -w" client/client.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o piledriver_core.exe -ldflags "-s -w -H windowsgui" client/client.go
	upx piledriver_core_darwin
	upx piledriver_core_linux
	upx piledriver_core.exe

build_server:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o piledriver_server -ldflags "-s -w" server/server.go
	upx piledriver_server


SSLKey:
	@echo 'SSLKey'
	openssl genrsa -out cert/server.key 2048
	openssl ecparam -genkey -name secp384r1 -out cert/serveryek
	openssl req -new -x509 -sha256 -key cert/serveryek -out cert/servermpe -days 3650  # plumber

