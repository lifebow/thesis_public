package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type VPNClient struct {
	Name       string
	Token      string
	OldIP      string
	CurrentIP  string
	LastUpdate int64
}
type DataClient struct {
	Clients []VPNClient
}

var globalData DataClient
var preShareKey string

func LoadInitData(configPath string) {
	// open config file
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	// get json data
	var data map[string]interface{}
	json.Unmarshal(configData, &data)

	for _, jClient := range data["vpnClients"].([]interface{}) {
		client := VPNClient{
			Name:       jClient.(map[string]interface{})["name"].(string),
			Token:      jClient.(map[string]interface{})["token"].(string),
			OldIP:      jClient.(map[string]interface{})["initIP"].(string),
			CurrentIP:  jClient.(map[string]interface{})["initIP"].(string),
			LastUpdate: 0,
		}

		globalData.Clients = append(globalData.Clients, client)
	}
	preShareKey = data["preSharedkey"].(string)
}

func after(value string, a string) string {
	// Get substring after a string.
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:len(value)]
}

func writeTofile(currentIP string, newIP string, num int) {
	file, err := ioutil.ReadFile("/etc/wireguard/wg0.conf")
	if num != 0 {
		file, err = ioutil.ReadFile("/etc/wireguard/wg0_new.conf")
	}
	currentIP = currentIP + "/32"
	newIP = newIP + "/32"
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	var content = string(file)
	output := strings.Split(content, "\n")
	for index, _ := range output {
		if strings.Contains(output[index], currentIP) {
			if strings.Contains(output[index], "#LastUpdate=") {
				last := after(output[index], "#LastUpdate=")
				value, _ := strconv.ParseInt(last, 10, 64)
				currentTime := time.Now().UnixNano()
				if currentTime-value > 2000000000 {
					now := strconv.FormatInt(currentTime, 10)
					output[index] = strings.Replace(output[index], last, now, -1)
					output[index] = strings.Replace(output[index], currentIP, newIP, -1)
				}
			} else {
				output[index] += " #LastUpdate=" + strconv.FormatInt(time.Now().UnixNano(), 10)
				output[index] = strings.Replace(output[index], currentIP, newIP, -1)
			}

		}
	}
	new_output := strings.Join(output[:], "\n")
	err3 := ioutil.WriteFile("/etc/wireguard/wg0_new.conf", []byte(new_output), 0644)
	if err3 != nil {
		panic("write file problem")
	}
}

func restart() {
	//down
	cmd := exec.Command("wg-quick", "down", "wg0")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	os.Rename("/etc/wireguard/wg0.conf", "/etc/wireguard/wg0_backup.conf")
	os.Rename("/etc/wireguard/wg0_new.conf", "/etc/wireguard/wg0.conf")
	cmd2 := exec.Command("wg-quick", "up", "wg0")
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr
	err2 := cmd2.Run()
	if err2 != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	fmt.Print("restartWG")
}

func updateIP(currentIP string, newIP string, num int, wg *sync.WaitGroup) {
	defer wg.Done()
	//write to global
	var token string
	for i, team := range globalData.Clients {
		if team.CurrentIP == currentIP {
			currentTime := time.Now().UnixNano()
			// fmt.Printf("Update client" + team.Name + " current ip: "+ team.CurrentIP+ "to: "+ c )
			if currentTime-team.LastUpdate > 2000000000 {
				globalData.Clients[i].OldIP = team.CurrentIP
				token = team.Token
				globalData.Clients[i].CurrentIP = newIP
				globalData.Clients[i].LastUpdate = currentTime
			}
		}
	}
	data, err := json.Marshal(map[string]string{
		"method": "update",
		"oldIP":  currentIP,
		"newIP":  newIP,
		"token":  token,
	})
	if err != nil {
		log.Fatalf("Failed to create json data\n")
	}
	url := "http://" + currentIP + ":8085" + "/updateIP"
	for i := 0; i < 3; i++ {
		_, err = http.Post(url, "application/json", bytes.NewBuffer(data))
		if err != nil {
			log.Printf("Failed to send request to target IP: %v\n", currentIP)
			log.Printf("Waiting 0.5 second to resend!")
			time.Sleep(500 * time.Millisecond)
		} else {
			break
		}
	}

	//send request restart to client
	restartUrl := "http://" + currentIP + ":8085" + "/restartClient"
	//need to replace
	secretKey, err := json.Marshal(map[string]string{
		"token": token,
	})

	for i := 0; i < 3; i++ {
		_, err := http.Post(restartUrl, "application/json", bytes.NewBuffer(secretKey))
		if err != nil {
			log.Printf("Failed to send request to restart IP: %v\n", currentIP)
			log.Fatal(err)
			log.Printf("Waiting 0.5 second to resend!\n")
			time.Sleep(500 * time.Millisecond)
		} else {
			break
		}
	}
	return
}

type Player struct {
	OldIP string `json:"oldIP"`
	NewIP string `json:"newIP"`
}
type Data struct {
	Players     []Player `json:"players"`
	preShareKey string   `json:"preShareKey"`
}

func updateIPs(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	var data Data
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}
	if data.preShareKey != preShareKey {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	var wg sync.WaitGroup
	// log.Println(data.Players)
	for index, player := range data.Players {
		wg.Add(1)
		oldIP := player.OldIP
		newIP := player.NewIP
		fmt.Println(oldIP + "->" + newIP)
		go updateIP(oldIP, newIP, index, &wg)
		writeTofile(oldIP, newIP, index)
	}
	wg.Wait()
	fmt.Fprintf(rw, "OK!")
	restart()

}
func main() {
	LoadInitData("/Documents/server_config.json")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/updateIPs", updateIPs).Methods("POST")
	log.Fatal(http.ListenAndServe(":8085", router))
}
