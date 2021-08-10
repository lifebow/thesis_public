#! /bin/sh


# This script will run inside vm
# This script must be run with root permission
# This script will setup server vpn changer service


cp /VPNChanger/Server/vpnChangerServer.service /etc/systemd/system

systemctl daemon-reload
systemctl enable vpnChangerServer.service