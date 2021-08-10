package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	model "service_status/pkg"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"gopkg.in/robfig/cron.v2"
)

type Service struct {
	CurrentRound int
	CurrentTick  int
	Rounds       []model.Round
	queueIt      *model.QueueConnection
	ResultRounds []model.Round
}

func NewService() *Service {
	rounds := make([]model.Round, 0)
	queueconnect, err := model.NewQueue()
	if err != nil {
		log.Warn().Err(err)
	}
	return &Service{CurrentRound: 0, Rounds: rounds, queueIt: queueconnect}
}
func (s *Service) StartDaemon() {
	err := godotenv.Load()
	if err != nil {
		log.Err(err)
	}
	timeOfTick := os.Getenv("TimeOfTick")
	rankingHost := os.Getenv("RankingHost")
	go func() {
		go func() {
			c := cron.New()
			timestr := fmt.Sprintf("@every 0h0m%vs", timeOfTick)
			c.AddFunc(timestr, func() {
				timer, err := strconv.Atoi(timeOfTick)
				fmt.Printf("Timer %v\n", timer*400)
				if err != nil {
					log.Err(err).Msgf("Failed to convert string of time to number")
				}
				time.Sleep(time.Duration(timer*400) * time.Millisecond)
				log.Info().Msgf("Summarize the first!\n")
				for _, round := range s.Rounds {
					if round.RoundNumber != s.CurrentRound {
						continue
					}
					for _, team := range round.Teams {
						for _, service := range team.Challenges {
							for _, tickResult := range service.Results {
								if tickResult.Tick != s.CurrentTick {
									continue
								}
								a := model.ResultData{
									TeamName:      team.TeamName,
									ServiceName:   service.ServiceName,
									IP:            team.IP,
									Round:         s.CurrentRound,
									Tick:          s.CurrentTick,
									ScriptbotName: "ServiceStatus",
								}
								success := true
								for _, resultbyBot := range tickResult.Results {
									IsSuccess := resultbyBot.IsSuccess
									if IsSuccess == false {
										success = false
									}
								}
								a.IsSuccess = success
								if success == false {
									log.Info().Msgf("Send emergency Team: %v , Service: %v \n", a.TeamName, a.ServiceName)
									r := model.TargetToQueue{
										ServiceName: a.ServiceName,
										IP:          team.IP,
										TeamName:    team.TeamName,
										Port:        a.Port,
										Round:       s.CurrentRound,
										Tick:        s.CurrentTick,
									}
									msg := Message{
										Result: r,
									}
									s.queueIt.PublishEmergencyJob(msg)
								}
								s.AddResult(a)
							}
						}
					}
				}
				time.Sleep(time.Duration(timer*300) * time.Millisecond)
				//check 1
				fmt.Printf("Timer %v\n", timer*(400+300))
				log.Info().Msgf("The Second summary\n")
				for _, round := range s.Rounds {
					if round.RoundNumber != s.CurrentRound {
						continue
					}
					for _, team := range round.Teams {
						for _, service := range team.Challenges {
							for _, tickResult := range service.Results {
								if tickResult.Tick != s.CurrentTick {
									continue
								}
								a := model.ResultData{
									TeamName:      team.TeamName,
									ServiceName:   service.ServiceName,
									IP:            team.IP,
									Round:         s.CurrentRound,
									Tick:          s.CurrentTick,
									ScriptbotName: "ServiceStatus",
								}
								success := true
								for _, resultbyBot := range tickResult.Results {
									IsSuccess := resultbyBot.IsSuccess
									if IsSuccess == false {
										success = false
									}
								}
								a.IsSuccess = success
								s.AddResult(a)
							}
						}
					}
				}

				time.Sleep(time.Duration(timer*100) * time.Millisecond)
				fmt.Printf("Timer %v\n", timer*(400+300+100))
				fmt.Println("Send data to Ranking!")
				var data model.Round
				for _,resultRound :=range s.ResultRounds{
					if resultRound.RoundNumber==s.CurrentRound{
						data=resultRound
					}
				}
				jsonValue,_:=json.Marshal(data)
				getListServiceLink := fmt.Sprintf("http://%v/%v", rankingHost, "resultRound")
				times := 0
				for times < 3 {
					req,err :=http.NewRequest("POST",getListServiceLink,bytes.NewBuffer(jsonValue))
					if err !=nil{
						log.Err(err)
					}
					client := &http.Client{}
					resp, err := client.Do(req)
					log.Info().Msgf("Send data to Ranking! : %s",getListServiceLink)
					if (err != nil) || (resp.StatusCode != 200) {
						log.Err(err).Msg("Send response to Ranking failed, retry in 0.5s")
						times++
						time.Sleep(time.Millisecond * 500)
					} else {
						break
					}
				}
				//increase current tick
				s.CurrentTick++
				//fmt.Println(s.ResultRounds)
			})
			c.Start()
			stop := make(chan int)
			<-stop
		}()
	}()
}
func (s *Service) AddResult(data model.ResultData) bool {

	hasRound := false
	var roundIndex int
	if len(s.ResultRounds) != 0 {
		for i, round := range s.ResultRounds {
			if round.RoundNumber == data.Round {
				hasRound = true
				roundIndex = i
			}
		}
	}
	if hasRound == false {
		resultbyBot := model.ResultbyBot{
			IsSuccess:    data.IsSuccess,
			ScripbotName: data.ScriptbotName,
		}
		resultbyBots := make([]model.ResultbyBot, 0)
		resultbyBots = append(resultbyBots, resultbyBot)
		result := model.TickResult{
			Tick:    data.Tick,
			Results: resultbyBots,
		}
		results := make([]model.TickResult, 0)
		results = append(results, result)

		challenge := model.Challenge{
			ServiceName: data.ServiceName,
			Port:        data.Port,
			Results:     results,
		}
		challenges := make([]model.Challenge, 0)
		challenges = append(challenges, challenge)

		team := model.Team{
			TeamName:   data.TeamName,
			IP:         data.IP,
			Challenges: challenges,
		}
		teams := make([]model.Team, 0)
		teams = append(teams, team)

		round := model.Round{
			RoundNumber: data.Round,
			Teams:       teams,
		}

		s.ResultRounds = append(s.ResultRounds, round)
		return true
	}

	hasTeam := false
	var teamIndex int
	if len(s.ResultRounds[roundIndex].Teams) != 0 {
		for i, team := range s.ResultRounds[roundIndex].Teams {
			if team.TeamName == data.TeamName {
				hasTeam = true
				teamIndex = i
			}
		}
	}
	if hasTeam == false {
		resultbyBot := model.ResultbyBot{
			IsSuccess:    data.IsSuccess,
			ScripbotName: data.ScriptbotName,
		}
		resultbyBots := make([]model.ResultbyBot, 0)
		resultbyBots = append(resultbyBots, resultbyBot)
		result := model.TickResult{
			Tick:    data.Tick,
			Results: resultbyBots,
		}
		results := make([]model.TickResult, 0)
		results = append(results, result)

		challenge := model.Challenge{
			ServiceName: data.ServiceName,
			Port:        data.Port,
			Results:     results,
		}
		challenges := make([]model.Challenge, 0)
		challenges = append(challenges, challenge)

		team := model.Team{
			TeamName:   data.TeamName,
			IP:         data.IP,
			Challenges: challenges,
		}
		s.ResultRounds[roundIndex].Teams = append(s.ResultRounds[roundIndex].Teams, team)
		return true
	}

	hasChallenge := false
	var challengeIndex int
	if len(s.ResultRounds[roundIndex].Teams[teamIndex].Challenges) != 0 {
		for i, challenge := range s.ResultRounds[roundIndex].Teams[teamIndex].Challenges {
			if challenge.ServiceName == data.ServiceName {
				hasChallenge = true
				challengeIndex = i
			}
		}
	}
	if hasChallenge == false {
		resultbyBot := model.ResultbyBot{
			IsSuccess:    data.IsSuccess,
			ScripbotName: data.ScriptbotName,
		}
		resultbyBots := make([]model.ResultbyBot, 0)
		resultbyBots = append(resultbyBots, resultbyBot)
		result := model.TickResult{
			Tick:    data.Tick,
			Results: resultbyBots,
		}
		results := make([]model.TickResult, 0)
		results = append(results, result)

		challenge := model.Challenge{
			ServiceName: data.ServiceName,
			Port:        data.Port,
			Results:     results,
		}
		s.ResultRounds[roundIndex].Teams[teamIndex].Challenges = append(s.ResultRounds[roundIndex].Teams[teamIndex].Challenges, challenge)
		return true
	}

	hasTickResult := false
	var tickResultIndex int
	if len(s.ResultRounds[roundIndex].Teams[teamIndex].Challenges[challengeIndex].Results) != 0 {
		for i, result := range s.ResultRounds[roundIndex].Teams[teamIndex].Challenges[challengeIndex].Results {
			if result.Tick == data.Tick {
				hasTickResult = true
				tickResultIndex = i
			}
		}
	}
	if hasTickResult == false {
		resultbyBot := model.ResultbyBot{
			IsSuccess:    data.IsSuccess,
			ScripbotName: data.ScriptbotName,
		}
		resultbyBots := make([]model.ResultbyBot, 0)
		resultbyBots = append(resultbyBots, resultbyBot)
		result := model.TickResult{
			Tick:    data.Tick,
			Results: resultbyBots,
		}
		s.ResultRounds[roundIndex].Teams[teamIndex].Challenges[challengeIndex].Results = append(
			s.ResultRounds[roundIndex].Teams[teamIndex].Challenges[challengeIndex].Results,
			result)
		return true
	}

	hasResultbyBot := false
	var resultbyBotIndex int
	if len(s.ResultRounds[roundIndex].Teams[teamIndex].Challenges[challengeIndex].Results[tickResultIndex].Results) != 0 {
		for i, result := range s.ResultRounds[roundIndex].Teams[teamIndex].Challenges[challengeIndex].Results[tickResultIndex].Results {
			if result.ScripbotName == data.ScriptbotName {
				hasResultbyBot = true
				resultbyBotIndex = i
			}
		}
	}
	if hasResultbyBot == false {
		resultbyBot := model.ResultbyBot{
			IsSuccess:    data.IsSuccess,
			ScripbotName: data.ScriptbotName,
		}
		s.ResultRounds[roundIndex].Teams[teamIndex].Challenges[challengeIndex].Results[tickResultIndex].Results = append(
			s.ResultRounds[roundIndex].Teams[teamIndex].Challenges[challengeIndex].Results[tickResultIndex].Results,
			resultbyBot)
		return true
	}
	s.ResultRounds[roundIndex].Teams[teamIndex].Challenges[challengeIndex].Results[tickResultIndex].Results[resultbyBotIndex].ScripbotName = data.ScriptbotName
	s.ResultRounds[roundIndex].Teams[teamIndex].Challenges[challengeIndex].Results[tickResultIndex].Results[resultbyBotIndex].IsSuccess = data.IsSuccess
	return true
}
func (s *Service) AddData(data model.ResultData) bool {

	hasRound := false
	var roundIndex int
	if len(s.Rounds) != 0 {
		for i, round := range s.Rounds {
			if round.RoundNumber == data.Round {
				hasRound = true
				roundIndex = i
			}
		}
	}
	if hasRound == false {
		resultbyBot := model.ResultbyBot{
			IsSuccess:    data.IsSuccess,
			ScripbotName: data.ScriptbotName,
		}
		resultbyBots := make([]model.ResultbyBot, 0)
		resultbyBots = append(resultbyBots, resultbyBot)
		result := model.TickResult{
			Tick:    data.Tick,
			Results: resultbyBots,
		}
		results := make([]model.TickResult, 0)
		results = append(results, result)

		challenge := model.Challenge{
			ServiceName: data.ServiceName,
			Port:        data.Port,
			Results:     results,
		}
		challenges := make([]model.Challenge, 0)
		challenges = append(challenges, challenge)

		team := model.Team{
			TeamName:   data.TeamName,
			IP:         data.IP,
			Challenges: challenges,
		}
		teams := make([]model.Team, 0)
		teams = append(teams, team)

		round := model.Round{
			RoundNumber: data.Round,
			Teams:       teams,
		}

		s.Rounds = append(s.Rounds, round)
		return true
	}

	hasTeam := false
	var teamIndex int
	if len(s.Rounds[roundIndex].Teams) != 0 {
		for i, team := range s.Rounds[roundIndex].Teams {
			if team.TeamName == data.TeamName {
				hasTeam = true
				teamIndex = i
			}
		}
	}
	if hasTeam == false {
		resultbyBot := model.ResultbyBot{
			IsSuccess:    data.IsSuccess,
			ScripbotName: data.ScriptbotName,
		}
		resultbyBots := make([]model.ResultbyBot, 0)
		resultbyBots = append(resultbyBots, resultbyBot)
		result := model.TickResult{
			Tick:    data.Tick,
			Results: resultbyBots,
		}
		results := make([]model.TickResult, 0)
		results = append(results, result)

		challenge := model.Challenge{
			ServiceName: data.ServiceName,
			Port:        data.Port,
			Results:     results,
		}
		challenges := make([]model.Challenge, 0)
		challenges = append(challenges, challenge)

		team := model.Team{
			TeamName:   data.TeamName,
			IP:         data.IP,
			Challenges: challenges,
		}
		s.Rounds[roundIndex].Teams = append(s.Rounds[roundIndex].Teams, team)
		return true
	}

	hasChallenge := false
	var challengeIndex int
	if len(s.Rounds[roundIndex].Teams[teamIndex].Challenges) != 0 {
		for i, challenge := range s.Rounds[roundIndex].Teams[teamIndex].Challenges {
			if challenge.ServiceName == data.ServiceName {
				hasChallenge = true
				challengeIndex = i
			}
		}
	}
	if hasChallenge == false {
		resultbyBot := model.ResultbyBot{
			IsSuccess:    data.IsSuccess,
			ScripbotName: data.ScriptbotName,
		}
		resultbyBots := make([]model.ResultbyBot, 0)
		resultbyBots = append(resultbyBots, resultbyBot)
		result := model.TickResult{
			Tick:    data.Tick,
			Results: resultbyBots,
		}
		results := make([]model.TickResult, 0)
		results = append(results, result)

		challenge := model.Challenge{
			ServiceName: data.ServiceName,
			Port:        data.Port,
			Results:     results,
		}
		s.Rounds[roundIndex].Teams[teamIndex].Challenges = append(s.Rounds[roundIndex].Teams[teamIndex].Challenges, challenge)
		return true
	}

	hasTickResult := false
	var tickResultIndex int
	if len(s.Rounds[roundIndex].Teams[teamIndex].Challenges[challengeIndex].Results) != 0 {
		for i, result := range s.Rounds[roundIndex].Teams[teamIndex].Challenges[challengeIndex].Results {
			if result.Tick == data.Tick {
				hasTickResult = true
				tickResultIndex = i
			}
		}
	}
	if hasTickResult == false {
		resultbyBot := model.ResultbyBot{
			IsSuccess:    data.IsSuccess,
			ScripbotName: data.ScriptbotName,
		}
		resultbyBots := make([]model.ResultbyBot, 0)
		resultbyBots = append(resultbyBots, resultbyBot)
		result := model.TickResult{
			Tick:    data.Tick,
			Results: resultbyBots,
		}
		s.Rounds[roundIndex].Teams[teamIndex].Challenges[challengeIndex].Results = append(
			s.Rounds[roundIndex].Teams[teamIndex].Challenges[challengeIndex].Results,
			result)
		return true
	}

	hasResultbyBot := false
	var resultbyBotIndex int
	if len(s.Rounds[roundIndex].Teams[teamIndex].Challenges[challengeIndex].Results[tickResultIndex].Results) != 0 {
		for i, result := range s.Rounds[roundIndex].Teams[teamIndex].Challenges[challengeIndex].Results[tickResultIndex].Results {
			if result.ScripbotName == data.ScriptbotName {
				hasResultbyBot = true
				resultbyBotIndex = i
			}
		}
	}
	if hasResultbyBot == false {
		resultbyBot := model.ResultbyBot{
			IsSuccess:    data.IsSuccess,
			ScripbotName: data.ScriptbotName,
		}
		s.Rounds[roundIndex].Teams[teamIndex].Challenges[challengeIndex].Results[tickResultIndex].Results = append(
			s.Rounds[roundIndex].Teams[teamIndex].Challenges[challengeIndex].Results[tickResultIndex].Results,
			resultbyBot)
		return true
	}
	s.Rounds[roundIndex].Teams[teamIndex].Challenges[challengeIndex].Results[tickResultIndex].Results[resultbyBotIndex].ScripbotName = data.ScriptbotName
	s.Rounds[roundIndex].Teams[teamIndex].Challenges[challengeIndex].Results[tickResultIndex].Results[resultbyBotIndex].IsSuccess = data.IsSuccess
	return true
}

func (s *Service) SetRound(data model.NewTickData) bool {
	s.CurrentRound = data.Round
	s.CurrentTick = data.NumberTick
	// delete old round if round < data.Round-3
	if len(s.Rounds) < 3 {
		return true
	}
	for data.Round-s.Rounds[0].RoundNumber > 3 {
		s.Rounds = s.Rounds[1:]
	}
	return true
}

type Message struct {
	Result interface{} `json:"Result"`
}

func (s *Service) PublicService(listTeam model.ListTeam, listService model.ListService) {
	for _, team := range listTeam.Teams {
		for _, service := range listService.Services {
			r := model.TargetToQueue{
				ServiceName: service.ServcieName,
				IP:          team.IP,
				TeamName:    team.TeamName,
				Port:        service.Port,
				Round:       s.CurrentRound,
				Tick:        s.CurrentTick,
			}
			msg := Message{
				Result: r,
			}
			err := s.queueIt.PublishNormalJob(msg)
			if err != nil {
				log.Warn().Msgf("Failed to publish normal_check: TeamName - %v, ServiceName - %v \n", r.TeamName, r.ServiceName)
			}
			log.Info().Msgf("Successes to publish normal_check: TeamName - %v, ServiceName - %v \n", r.TeamName, r.ServiceName)
		}
	}
}
