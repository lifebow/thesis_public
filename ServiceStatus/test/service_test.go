package test

import (
	model "service_status/pkg"
	"service_status/pkg/service"
	"testing"
)

func Test1AddResult(t *testing.T) {//ServiceStatus update true
	var newService *service.Service
	newService=service.NewService()
	var a model.ResultData
	a = model.ResultData{
		TeamName:      "Team1",
		ServiceName:   "challengePwn",
		IP:            "127.0.0.1",
		Round:         1,
		Tick:          1,
		ScriptbotName: "ServiceStatus",
		IsSuccess:     true,
	}
	newService.AddResult(a)
	if(newService.ResultRounds[0].Teams[0].Challenges[0].Results[0].Results[0].IsSuccess!=true){
		t.Errorf("Output expect check Challenge is true instead of %v", newService.ResultRounds[0].Teams[0].Challenges[0].Results[0].Results[0].IsSuccess)
	}
}

func Test2AddResult(t *testing.T) {//ServiceStatus update false
	var newService *service.Service
	newService=service.NewService()
	var b model.ResultData
	b = model.ResultData{
		TeamName:      "Team1",
		ServiceName:   "challengePwn",
		IP:            "127.0.0.1",
		Round:         1,
		Tick:          1,
		ScriptbotName: "ServiceStatus",
		IsSuccess: false,
	}
	newService.AddResult(b)
	if(newService.ResultRounds[0].Teams[0].Challenges[0].Results[0].Results[0].IsSuccess!=false){
		t.Errorf("Output expect check Challenge is false instead of %v", newService.ResultRounds[0].Teams[0].Challenges[0].Results[0].Results[0].IsSuccess)
	}
}

func Test3AddResult(t *testing.T) { //ServiceStatus update result
	var newService *service.Service
	newService=service.NewService()
	var a model.ResultData
	a = model.ResultData{
		TeamName:      "Team1",
		ServiceName:   "challengePwn",
		IP:            "127.0.0.1",
		Round:         1,
		Tick:          1,
		ScriptbotName: "ServiceStatus",
		IsSuccess: false,
	}
	var b model.ResultData
	b = model.ResultData{
		TeamName:      "Team1",
		ServiceName:   "challengePwn",
		IP:            "127.0.0.1",
		Round:         1,
		Tick:          1,
		ScriptbotName: "ServiceStatus",
		IsSuccess: true,
	}
	newService.AddResult(a)
	newService.AddResult(b)
	if(newService.ResultRounds[0].Teams[0].Challenges[0].Results[0].Results[0].IsSuccess!=true){
		t.Errorf("Output expect check Challenge is true instead of %v", newService.ResultRounds[0].Teams[0].Challenges[0].Results[0].Results[0].IsSuccess)
	}
}

func Test1AddData(t *testing.T) { // 2 Scriptbot return data
	var newService *service.Service
	newService=service.NewService()
	var a model.ResultData
	a = model.ResultData{
		TeamName:      "Team1",
		ServiceName:   "challengePwn",
		IP:            "127.0.0.1",
		Round:         1,
		Tick:          1,
		ScriptbotName: "Scriptbot1",
		IsSuccess: false,
	}
	var b model.ResultData
	b = model.ResultData{
		TeamName:      "Team1",
		ServiceName:   "challengePwn",
		IP:            "127.0.0.1",
		Round:         1,
		Tick:          1,
		ScriptbotName: "Scriptbot2",
		IsSuccess: true,
	}
	newService.AddResult(a)
	newService.AddResult(b)
	if(len(newService.ResultRounds[0].Teams[0].Challenges[0].Results[0].Results)!=2){
		t.Errorf("Output expect number of Scriptbot check is 2 instead of %v", newService.ResultRounds[0].Teams[0].Challenges[0].Results[0].Results[0])
	}
}
func Test2AddData(t *testing.T) { //Scriptbot1 update result
	var newService *service.Service
	newService=service.NewService()
	var a model.ResultData
	a = model.ResultData{
		TeamName:      "Team1",
		ServiceName:   "challengePwn",
		IP:            "127.0.0.1",
		Round:         1,
		Tick:          1,
		ScriptbotName: "Scriptbot1",
		IsSuccess: false,
	}
	var b model.ResultData
	b = model.ResultData{
		TeamName:      "Team1",
		ServiceName:   "challengePwn",
		IP:            "127.0.0.1",
		Round:         1,
		Tick:          1,
		ScriptbotName: "Scriptbot1",
		IsSuccess: true,
	}
	newService.AddResult(a)
	newService.AddResult(b)
	if(newService.ResultRounds[0].Teams[0].Challenges[0].Results[0].Results[0].IsSuccess!=true){
		t.Errorf("Output expect result of Scriptbot1 check is true instead of %v", newService.ResultRounds[0].Teams[0].Challenges[0].Results[0].Results[0].IsSuccess)
	}
}
func Test3AddData(t *testing.T) { //Scriptbot1 update 2 challenge of a team
	var newService *service.Service
	newService=service.NewService()
	var a model.ResultData
	a = model.ResultData{
		TeamName:      "Team1",
		ServiceName:   "challengePwn",
		IP:            "127.0.0.1",
		Round:         1,
		Tick:          1,
		ScriptbotName: "Scriptbot1",
		IsSuccess: false,
	}
	var b model.ResultData
	b = model.ResultData{
		TeamName:      "Team1",
		ServiceName:   "challengePwn2",
		IP:            "127.0.0.1",
		Round:         1,
		Tick:          1,
		ScriptbotName: "Scriptbot1",
		IsSuccess: true,
	}
	newService.AddResult(a)
	newService.AddResult(b)
	if(len(newService.ResultRounds[0].Teams[0].Challenges)!=2){
		t.Errorf("Output expect number of challenge Team1 checked  is 2 instead of %v", len(newService.ResultRounds[0].Teams[0].Challenges))
	}
}

func Test1SetRound(t *testing.T) { //Scriptbot1 set newRound
	var newService *service.Service
	newService=service.NewService()
	var newTick model.NewTickData
	newTick=model.NewTickData{
		NumberTick: 5,
		Round: 1,
	}
	newService.SetRound(newTick)
	if(newService.CurrentRound!=1){
		t.Errorf("Output expect current round is 1 instead of %v", newService.CurrentRound)
	}
	if(newService.CurrentTick!=5){
		t.Errorf("Output expect current tick is 5 instead of %v", newService.CurrentTick)
	}
}
func Test2SetRound(t *testing.T) { //update newTick
	var newService *service.Service
	newService=service.NewService()
	var newTick model.NewTickData
	newTick=model.NewTickData{
		NumberTick: 5,
		Round: 1,
	}
	newService.SetRound(newTick)
	newTick.NumberTick=6
	newService.SetRound(newTick)
	if(newService.CurrentRound!=1){
		t.Errorf("Output expect current round is 1 instead of %v", newService.CurrentRound)
	}
	if(newService.CurrentTick!=6){
		t.Errorf("Output expect current tick is 5 instead of %v", newService.CurrentTick)
	}
}
func Test3SetRound(t *testing.T) { //update newRound, remove old roundData
	var newService *service.Service
	newService=service.NewService()
	var b model.ResultData
	b = model.ResultData{
		TeamName:      "Team1",
		ServiceName:   "challengePwn",
		IP:            "127.0.0.1",
		Round:         1,
		Tick:          1,
		ScriptbotName: "Scriptbot1",
		IsSuccess: true,
	}
	newService.AddData(b)
	b.Round=2
	b.Tick=3
	newService.AddData(b)
	b.Round=3
	b.Tick=5
	newService.AddData(b)
	b.Round=4
	b.Tick=7
	newService.AddData(b)
	var newTick model.NewTickData
	newTick=model.NewTickData{
		NumberTick: 9,
		Round: 5,
	}
	newService.SetRound(newTick)
	if(newService.CurrentRound!=5){
		t.Errorf("Output expect current round is 5 instead of %v", newService.CurrentRound)
	}
	if(newService.CurrentTick!=9){
		t.Errorf("Output expect current tick is 9 instead of %v", newService.CurrentTick)
	}
	if(len(newService.Rounds)!=3){
		t.Errorf("Output expect number round Data stored is 3 instead of %v", newService.CurrentTick)
	}
}