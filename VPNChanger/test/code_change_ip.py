import requests
import time
import json
import random
import sys
def sendRequest(ip,data):
    url = "http://{}:8085/updateIPs".format(ip)
    headers = {"Connection": "close", "Upgrade-Insecure-Requests": "1", "User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36", "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9", "Sec-Fetch-Site": "same-origin", "Sec-Fetch-Mode": "navigate", "Sec-Fetch-User": "?1", "Sec-Fetch-Dest": "document", "Referer": "https://websec.fr/", "Accept-Encoding": "gzip, deflate", "Accept-Language": "en-US,en;q=0.9"}
    return requests.post(url, headers=headers, data=data,verify=False)
def init_Input(numPlayer):
    result=[]
    for i in range(3,numPlayer+2):
        tmp={"oldIP":"192.168.6.{}".format(i),"newIP":"192.168.6.{}".format(i),"keyIP":""}
        result.append(tmp)
    return result
def change_IP(data):
    tmp=[]
    for i in data:
        currentIP=i['newIP']
        tmp.append(currentIP)
    tmp_index=list(range(len(data)))
    result=[]
    for i in range(len(data)):
        while (1):
            index=random.choice(tmp_index)
            if index!=i:
                break
        player=data[i]
        player['oldIP']=player['newIP']
        player['newIP']=tmp[index]
        tmp_index.remove(index)
        result.append(player)
    return result
IPs={
    "numUpdated":0,
    "VPNServer":"192.168.6.1",
    "Controller+ServiceStatus":"192.168.6.2",
    "Scriptbot1":{
        "numUpdated":0,
        "currentIP":"192.168.6.3"
    },
    "Scriptbot2":{
        "numUpdated":0,
        "currentIP":"192.168.6.4"
    },
    "Team1":{
        "numUpdated":0,
        "currentIP":"192.168.6.5"
    },
    "Team2":{
        "numUpdated":0,
        "currentIP":"192.168.6.6"
    }

}
def main():
    ip="54.251.197.225"
    data={"players":[]}
    data['players']=init_Input(5)
    data['players']=change_IP(data['players'])
    data["preSharedkey"]="secretKey"
    print(data)
    count=0
    for i in data['players']:
        count+=1
        print("Team{} change IP from {} to {}".format(count,i['oldIP'],i['newIP']))
        for key in IPs.keys():
            if key=="numUpdated" or key=="VPNServer" or "Controller+ServiceStatus":
                continue
            if IPs[key]["currentIP"] == i['oldIP']:
                if IPs[key]["numUpdated"]==IPs["numUpdated"]:
                    IPs[key]["numUpdated"]+=1
                    IPs[key]["currentIP"]=i['newIP']
    IPs["numUpdated"]+=1
    start=time.time()
    a=sendRequest(ip, json.dumps(data))
    print(a)
    end=time.time()
    count=0
    print("Scriptbot1 new IP: {}\n".format(IPs["Scriptbot1"]["currentIP"]))
    print("Scriptbot2 new IP: {}\n".format(IPs["Scriptbot2"]["currentIP"]))
    print("Team1 new IP: {}\n".format(IPs["Team1"]["currentIP"]))
    print("Team2 new IP: {}\n".format(IPs["Team2"]["currentIP"]))

    print("Time: {}".format(end-start))
if __name__ == "__main__":
    main()