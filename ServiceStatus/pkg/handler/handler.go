package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	model "service_status/pkg"
	"service_status/pkg/service"

	"github.com/joho/godotenv"
)

type Handler struct {
	normalService *service.Service
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
func NewHandler() *Handler {
	newService := service.NewService()
	newService.StartDaemon()
	return &Handler{
		normalService: newService,
	}
}

func responseWithJson(writer http.ResponseWriter, status int, object interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	json.NewEncoder(writer).Encode(object)
}
func (h *Handler) GetResult(writer http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	failOnError(err, "Failed to read body request")
	var result model.ResultData
	err = json.Unmarshal(body, &result)
	result.ScriptbotName = request.Header.Get("Name")
	if err != nil {
		failOnError(err, "Failed to pass params")
		responseWithJson(writer, http.StatusBadRequest, map[string]string{"message": "Wrong format"})
	}
	fmt.Printf("Receive message from %v about Team: %v , Service: %v, IsSucces: %v\n", result.ScriptbotName, result.TeamName, result.ServiceName, result.IsSuccess)
	h.normalService.AddData(result)
	responseWithJson(writer, http.StatusOK, map[string]string{"message": "Success"})

}

func (h *Handler) NewTick(writer http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	failOnError(err, "Failed to read body request")
	var data model.NewTickData
	err = json.Unmarshal(body, &data)
	if err != nil {
		failOnError(err, "Failed to pass params")
		responseWithJson(writer, http.StatusBadRequest, map[string]string{"message": "Wrong format"})
	}
	h.normalService.SetRound(data)

	err = godotenv.Load()
	failOnError(err, "Falied to load enviroment")
	controllerHost := os.Getenv("ControllerHost")

	//get listTeam,
	var listTeam model.ListTeam
	getListTeamLink := fmt.Sprintf("http://%v/%v", controllerHost, "getListTeam")
	resp, err := http.Get(getListTeamLink)
	body, err = ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &listTeam)
	if err != nil {
		failOnError(err, "Failed to pass params")
		responseWithJson(writer, http.StatusBadRequest, map[string]string{"message": "Wrong format"})
	}
	//getListService
	var listService model.ListService
	getListServiceLink := fmt.Sprintf("http://%v/%v", controllerHost, "getListService")
	resp, err = http.Get(getListServiceLink)
	body, err = ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &listService)

	h.normalService.PublicService(listTeam, listService)
	responseWithJson(writer, http.StatusOK, map[string]string{"message": "Success"})

}
