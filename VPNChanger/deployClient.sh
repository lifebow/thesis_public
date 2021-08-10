#! /bin/sh


# This script will run inside vm
# This script must be run with root permission
# This script will setup client vpn changer service


cp /VPNChanger/Client/vpnChangerClient.service /etc/systemd/system

systemctl daemon-reload
systemctl enable vpnChangerClient.service