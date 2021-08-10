package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"scriptbot/internal/config"
	"scriptbot/internal/queue"
	"scriptbot/internal/service"
	"scriptbot/model"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"
)

func main() {
	cfg := config.MustReadWorkerConfig()
	queueConnection, err := queue.NewQueue(cfg)
	if err != nil {
		log.Warn().Msg("failed to create new queue")
	}
	defer queueConnection.Close()

	worker := service.NewService()
	msgs, err := queueConnection.ConsumeJob()
	if err != nil {
		log.Warn().Msg("Failed to consume job")
	}

	stopChan := make(chan bool)
	go func() {
		for msg := range msgs {
			var target model.Target
			data := gjson.GetBytes(msg.Body, "Result").Value()
			dataJson, _ := json.Marshal(data)
			err := json.Unmarshal(dataJson, &target)
			if err != nil {
				log.Warn().Msg("Error decoding JSON")
			}

			log.Info().Msg(fmt.Sprintf("Received new job %v - IP: %v \n", target.ServiceName, target.IP))
			result := worker.Check(target.IP, strconv.Itoa(target.Port), target.ServiceName)
			//send request to service status

			message := map[string]interface{}{
				"TeamName":    target.TeamName,
				"ServiceName": target.ServiceName,
				"IP":          target.IP,
				"Round":       target.Round,
				"Tick":        target.Tick,
				"Port":        target.Port,
				"IsSuccess":   result,
			}
			byteRepresentation, _ := json.Marshal(message)

			times := 0
			for times < 3 {
				req, err := http.NewRequest("POST", cfg.UrlLink, bytes.NewBuffer(byteRepresentation))
				if err != nil {
					log.Err(err)
				}
				req.Header.Set("jwt_token", cfg.Jwt)
				client := &http.Client{}
				resp, err := client.Do(req)
				log.Info().Msgf("Send data to service status with Scriptbotname : %v", cfg.Jwt)
				if (err != nil) || (resp.StatusCode != 200) {
					log.Err(err).Msg("Send response to ServiceStatus failed, retry in 0.5s")
					times++
					time.Sleep(time.Millisecond * 500)
				} else {
					msg.Ack(true)
					break
				}
			}

		}
	}()
	msgs2, err:= queueConnection.ConsumeEmergencyJob()
	go func() {
		for msg1 := range msgs2 {
			var target model.Target
			data := gjson.GetBytes(msg1.Body, "Result").Value()
			dataJson, _ := json.Marshal(data)
			err := json.Unmarshal(dataJson, &target)
			if err != nil {
				log.Warn().Msg("Error decoding JSON")
			}

			log.Info().Msg(fmt.Sprintf("Received new emergency job %v - IP: %v \n", target.ServiceName, target.IP))
			result := worker.Check(target.IP, strconv.Itoa(target.Port), target.ServiceName)

			//send request to service status
			msg1.Ack(true)
			message := map[string]interface{}{
				"TeamName":    target.TeamName,
				"ServiceName": target.ServiceName,
				"IP":          target.IP,
				"Round":       target.Round,
				"Tick":		   target.Tick,
				"IsSuccess":   result,
			}
			byteRepresentation, err := json.Marshal(message)
			times := 0
			for times < 3 {
				req,err :=http.NewRequest("POST",cfg.UrlLink,bytes.NewBuffer(byteRepresentation))
				if err !=nil{
					log.Err(err)
				}
				req.Header.Set("jwt_token",cfg.Jwt)
				client := &http.Client{}
				resp, err := client.Do(req)
				log.Info().Msgf("Send emergency data to service status with Scriptbot token : %v",cfg.Jwt)
				if (err != nil) || (resp.StatusCode != 200) {
					log.Err(err).Msg("Send response to ServiceStatus failed, retry in 0.5s")
					times++
					time.Sleep(time.Millisecond * 500)
				} else {
					break
				}
			}
		}
	}()
	<-stopChan
}
