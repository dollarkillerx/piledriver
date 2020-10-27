Generate:
	@echo 'Build GRPC'
	protoc -I rpc/proto/ rpc/proto/*.proto --go_out=plugins=grpc:rpc/proto/.

SSLKey:
	@echo 'SSLKey'
	openssl genrsa -out cert/server.key 2048
	openssl ecparam -genkey -name secp384r1 -out cert/serveryek
	openssl req -new -x509 -sha256 -key cert/serveryek -out cert/servermpe -days 3650  # plumber