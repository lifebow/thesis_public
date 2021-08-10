package service

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"scriptbot/model"
	"strings"
	"time"
)

type ListScript struct {
	Challenge []model.Service
}
type Service struct {
	listScript ListScript
}
func InitScript() ListScript {
	err:= godotenv.Load()
	if err !=nil{
		log.Err(err)
	}
	listScript:=ListScript{}
	link:=os.Getenv("ScriptCheckFolder")
	//fmt.Printf("Link challenge:%v\n",link)
	listScript=CreateScript(link)
	//fmt.Printf("List challenge:\n")
	//for _,x := range listScript.Challenge {
	//	fmt.Printf(x.ServiceName+"\n")
	//}
	return listScript
}
func NewService() *Service {
	listScript:=InitScript()
	return &Service{listScript: listScript}
}
func (s *Service) Check(IP string,port string, servicename string)(bool){
	var result bool
	for _,service :=range s.listScript.Challenge {
		if(service.ServiceName==servicename){
			result=CheckService(IP,port,service)
			return result
		}
	}
	return result
}

func CheckService(IP string,port string, Service model.Service)(bool){
	rand.Seed(time.Now().UnixNano())
	numScript:=len(Service.Scripts)
	pathofScript:=Service.Scripts[rand.Intn(numScript)]
	fmt.Print("Use script path:"+ pathofScript+"\n")
	cmd:= exec.Command("python",pathofScript,IP,port)
	output,err:=cmd.Output()
	if err!=nil{
		fmt.Println(err)
	}
	result:=false
	if(strings.Contains(string(output),"Success!")){
		result=true
	}
	return result
}



func CreateScript(link string)  ListScript{
	files,_:=ioutil.ReadDir(link)
	list:=ListScript{}
	for _, f := range files {

		if f.IsDir() {
			tmp:=model.Service{ServiceName: f.Name()}
			tmp.Scripts=FindScript(link+"/"+f.Name())
			list.Challenge =append(list.Challenge,tmp)
		}
	}
	return list
}
func FindScript(link string)[]string{
	files,_:=ioutil.ReadDir(link)
	var script []string
	for _,f :=range files{
		if !f.IsDir(){
			script=append(script,link+"/"+f.Name())
		}
	}
	return script
}