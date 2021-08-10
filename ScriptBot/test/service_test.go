package test

import (
	"path/filepath"
	"scriptbot/internal/service"
	"scriptbot/model"
	"testing"
)

func Test1InitScript(t *testing.T) {
	var listScript service.ListScript
	listScript=service.InitScript()
	if len(listScript.Challenge)!=2{
		t.Errorf("Output expect len of listScript ? instead of %v", len(listScript.Challenge))
	}
	if listScript.Challenge[0].ServiceName!="challengePwn"{
		t.Errorf("Output expect name of firstChallenge ? instead of %v", len(listScript.Challenge[0].ServiceName))
	}
}
func Test2InitScript(t *testing.T) {
	var listScript service.ListScript
	listScript=service.InitScript()
	if len(listScript.Challenge)!=2{
		t.Errorf("Output expect len of listScript 2 instead of %v", len(listScript.Challenge))
	}
	if listScript.Challenge[1].ServiceName!="challengePwn2"{
		t.Errorf("Output expect name of firstChallenge challengePwn2 instead of %v", len(listScript.Challenge[1].ServiceName))
	}
}
func Test3InitScript(t *testing.T) {
	var listScript service.ListScript
	listScript=service.InitScript()
	if len(listScript.Challenge)!=2{
		t.Errorf("Output expect len of listScript 2 instead of %v", len(listScript.Challenge))
	}
	if listScript.Challenge[1].ServiceName!="challengePwn2"{
		t.Errorf("Output expect name of firstChallenge challengePwn2 instead of %v", len(listScript.Challenge[1].ServiceName))
	}
	if len(listScript.Challenge[0].Scripts)!=2{
		t.Errorf("Output expect name of scriptCheck challengePwn 2 instead of %v", len(listScript.Challenge[0].Scripts))
	}
}

func Test1FindScript(t *testing.T) {
	abs,err := filepath.Abs("./ServiceCheck1/challengePwn")
	if err != nil {
		t.Error("RelativePath wrong!")
	}
	result:=service.FindScript(abs)
	if len(result)!=2{
		t.Errorf("Output expect number of array Script 2 instead of %v", len(result))
	}
	link1,_:=filepath.Abs("./ServiceCheck1/challengePwn/check_1.py")
	if link1!=result[0]{
		t.Errorf("Output expect path of Script is \n  %v \n Instead of: \n of %v", link1,result[0])
	}
}
func Test2FindScript(t *testing.T) {
	abs,err := filepath.Abs("./ServiceCheck1/challengePwn2")
	if err != nil {
		t.Error("RelativePath wrong!")
	}
	result:=service.FindScript(abs)
	if len(result)!=2{
		t.Errorf("Output expect number of array Script 2 instead of %v", len(result))
	}
	link1,_:=filepath.Abs("./ServiceCheck1/challengePwn2/check_2.py")
	if link1!=result[1]{
		t.Errorf("Output expect path of Script is \n  %v \n Instead of: \n of %v", link1,result[1])
	}
}
func Test3FindScript(t *testing.T) {
	abs,err := filepath.Abs("./ServiceCheck2/challengeWeb")
	if err != nil {
		t.Error("RelativePath wrong!")
	}
	result:=service.FindScript(abs)
	if len(result)!=2{
		t.Errorf("Output expect number of array Script 2 instead of %v", len(result))
	}
	link1,_:=filepath.Abs("./ServiceCheck2/challengeWeb/check_2.py")
	if link1!=result[1]{
		t.Errorf("Output expect path of Script is \n  %v \n Instead of: \n of %v", link1,result[1])
	}
}

func Test1CreateScript(t *testing.T) {
	abs,err := filepath.Abs("./ServiceCheck1/")
	if err != nil {
		t.Error("RelativePath wrong!")
	}
	var result service.ListScript
	result=service.CreateScript(abs)
	if len(result.Challenge)!=2{
		t.Errorf("Output expect number of array Script 2 instead of %v", len(result.Challenge))
	}
	if result.Challenge[0].ServiceName!="challengePwn"{
		t.Errorf("Output expect name of firstChallenge challengePwn instead of %v", len(result.Challenge[0].ServiceName))
	}
}

func Test2CreateScript(t *testing.T) {
	abs,err := filepath.Abs("./ServiceCheck2/")
	if err != nil {
		t.Error("RelativePath wrong!")
	}
	var result service.ListScript
	result=service.CreateScript(abs)
	if len(result.Challenge)!=2{
		t.Errorf("Output expect number of array Script 2 instead of %v", len(result.Challenge))
	}
	if result.Challenge[1].ServiceName!="challengeWeb"{
		t.Errorf("Output expect name of firstChallenge challengeWeb instead of %v", len(result.Challenge[0].ServiceName))
	}
}

func Test3CreateScript(t *testing.T) {
	abs,err := filepath.Abs("./ServiceCheck2/")
	if err != nil {
		t.Error("RelativePath wrong!")
	}
	var result service.ListScript
	result=service.CreateScript(abs)
	if len(result.Challenge)!=2{
		t.Errorf("Output expect number of array Script 2 instead of %v", len(result.Challenge))
	}
	link1,_:=filepath.Abs("./ServiceCheck2/challengeWeb/check_2.py")
	if result.Challenge[1].Scripts[1]!=link1{
		t.Errorf("Output expect path of Script is \n  %v \n Instead of: \n of %v", link1,result.Challenge[1].Scripts[1])
	}
}

func TestADownCheckService(t *testing.T) {
	abs,_ := filepath.Abs("./ServiceCheck2/")
	var result service.ListScript
	result=service.CreateScript(abs)
	var service1 model.Service
	service1=result.Challenge[0]
	ip:="127.0.0.1"
	port:="2224"
	resultCheck:=service.CheckService(ip,port,service1)
	if resultCheck!=false{
		t.Errorf("Output expect CheckService is false  instead of:  %v", resultCheck)
	}
}
func TestAUpCheckService(t *testing.T) {
	abs,_ := filepath.Abs("./ServiceCheck1/")
	var result service.ListScript
	result=service.CreateScript(abs)
	var service1 model.Service
	service1=result.Challenge[0]
	ip:="127.0.0.1"
	port:="2223"
	resultCheck:=service.CheckService(ip,port,service1)
	if resultCheck!=true{
		t.Errorf("Output expect CheckService is false  instead of:  %v", resultCheck)
	}
}