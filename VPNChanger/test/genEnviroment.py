import sys
import os
if (sys.argv[1] is None):
    print("Give me num of client!")
else:
    num_Client=int(sys.argv[1])
    path=os.getcwd()
    newFolderPath=path+"/server"
    if os.path.exists(newFolderPath):
        os.system("rm -rf "+ newFolderPath)
    os.mkdir(newFolderPath)
    serverFile=open(newFolderPath+"/wg0.conf","w")
    serverString="""[Interface]
Address = 192.168.6.1/24
SaveConfig = true
ListenPort = 41194
PrivateKey = +CJ0jDJypgts/IgdtsrUoMs03QUAvPRJZcFFEFiHoU4=

"""
    client_string="""[Interface]
## This Desktop/client's private key ##
PrivateKey = {client_pkey}
## Client ip address ##
Address = 192.168.6.{i}/24
ListenPort=41194 
[Peer]
## Ubuntu 20.04 server public key ##
PublicKey = ensMjxfOEr1xPURhntlLXQSCsTYgSQfnzbNs5cr2q00=
## set ACL ##
AllowedIPs = 192.168.6.1/24
 
## Your Ubuntu 20.04 LTS server's public IPv4/IPv6 address and port ##
Endpoint = {serverIP}:41194
 
##  Key connection alive ##
PersistentKeepalive = 15
"""
    server_Stringadd="""[Peer]
PublicKey = {client_pubkey}
AllowedIPs = 192.168.6.{i}/32

"""
    dockercompose_file=open(path+"/docker-compose.yml","w") ##to change
    dockercompose_file.write("""version: "3.1"
services:
  server:
    build:
      context: .
      dockerfile: Dockerfile_server
    cap_add:
      - NET_ADMIN
    stdin_open: true
    volumes: 
      - ./server/:/etc/wireguard/
      - ./server_script/:/Documents/
    tty: true
    networks:
      testing_net:
          ipv4_address: 172.28.0.2
    ports: 
      - 41194:41194
""")
    dockercompose_file_Stringadd="""  client{i}:
    image: client_v1
    volumes: 
      - ./client_{i}/:/etc/wireguard/
      - ./client_script/:/Documents/
    networks:
      testing_net:
          ipv4_address: 172.28.0.{i2}
    cap_add:
      - NET_ADMIN
    stdin_open: true
    tty: true
"""
    client_script=""" "token": "client {i} token" """
    serverFile.write(serverString)
    serverIP="13.229.250.7"

    serverScript= """{"vpnClients":["""
    first="""{
      "name": "client1",
      "token": "client 1 token",
      "initIP": "192.168.6.2"
    }"""
    element= """
      "name": "client{i}",
      "token": "client {i} token",
      "initIP": "192.168.6.{i2}"
    """
    endServerScript="""]}"""
    server_scriptPath=path+"/server_script"
    server_script=open(server_scriptPath+"/server_config.json","w")
    serve_script_tmp=""
    for i in range(1,num_Client+1):
        newFolderPath=path+"/client_"+str(i)
        
        if i==1:
          serve_script_tmp=serverScript+first
        else:
          serve_script_tmp+=",{"+element.format(i=i,i2=i+1)+"}"
        if os.path.exists(newFolderPath):
            os.system("rm -rf "+ newFolderPath)
        os.mkdir(newFolderPath)
        #create key
        command='cd '+newFolderPath+'; wg genkey | tee privatekey | wg pubkey > publickey'
        os.system(command)
        publickeyFile=open(newFolderPath+"/publickey","r")
        privatekeyFile=open(newFolderPath+"/privatekey","r")
        #readkey
        publickey=publickeyFile.readline()
        privatekey=privatekeyFile.readline()
        #close
        publickeyFile.close()
        privatekeyFile.close()
        #write a conf for server
        serverFile.write(server_Stringadd.format(client_pubkey=publickey,i=i+1))
        client_wg0=open(newFolderPath+"/wg0.conf","w")
        client_wg0.write(client_string.format(client_pkey=privatekey,i=i+1,serverIP=serverIP))
        client_wg0.close()
        client_sc=open(newFolderPath+"/client_config.json","w")
        client_sc.write("{"+client_script.format(i=i)+"}")
        client_sc.close()



        
        dockercompose_file.write(dockercompose_file_Stringadd.format(i=i,i2=i+2))
    serve_script_tmp+=endServerScript
    server_script.write(serve_script_tmp)
    server_script.close()
    dockercompose_file.write("""networks:
    testing_net:
        ipam:
            driver: default
            config:
                - subnet: 172.28.0.0/16""")
    dockercompose_file.close()
    serverFile.close()
        