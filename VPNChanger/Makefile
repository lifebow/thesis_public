all:
	rm -rf Output

	mkdir -p Output/Client
	go get github.com/gorilla/mux
	go build -o Output/Client/vpnChangerClient client.go
	cp -T vpnChangerClient.service Output/Client/vpnChangerClient.service
	cp -T deployClient.sh Output/Client/deploy.sh

	mkdir -p Output/Server
	go build -o Output/Server/vpnChangerServer server.go
	cp -T vpnChangerServer.service Output/Server/vpnChangerServer.service
	cp -T deployServer.sh Output/Server/deploy.sh
	rm test/server_script/vpnChangerServer
	cp -T Output/Server/vpnChangerServer test/server_script/vpnChangerServer

clean:
	rm -rf Output
