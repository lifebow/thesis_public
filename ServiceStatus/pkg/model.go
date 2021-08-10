package queue

type ResultData struct {
	TeamName      string `json:"TeamName"`
	ServiceName   string `json:"ServiceName"`
	IP            string `json:"IP"`
	Round         int    `json :"Round"`
	Tick          int    `json:"Tick"`
	IsSuccess     bool   `json::"IsSuccess"`
	ScriptbotName string
	Port          int `json:"Port"`
}
type NewTickData struct {
	NumberTick int `json:"currentTick"`
	Round      int `json:"currentRound"`
}

type Round struct {
	RoundNumber int 	`json:"roundNumber"`
	Teams       []Team 	`json:"teams"`
}
type Team struct {
	TeamName   string		`json:"teamName"`
	IP         string		`json:"IP"`
	Challenges []Challenge	`json:"challenges"`
}
type Challenge struct {
	ServiceName string		`json:"serviceName"`
	Port        int			`json:"port"`
	Results     []TickResult `json:"results"`
}
type TickResult struct {
	Tick    int				`json:"tick"`
	Results []ResultbyBot	`json:"results"`
}
type ResultbyBot struct {
	ScripbotName string		`json:"scripbotName"`
	IsSuccess    bool		`json:"isSuccess"`
}
type ServiceInfo struct {
	ServiceName string
	Port        int
}
type TeamInfo struct {
	TeamName string
	IP       string
}
type TeamData struct {
	TeamName string `json:"TeamName"`
	IP       string `json:"IP"`
}
type ListTeam struct {
	Teams []TeamData `json:"Teams"`
}
type ServiceData struct {
	ServcieName string `json:"ServiceName"`
	Port        int    `json:"Port"`
}
type ListService struct {
	Services []ServiceData `json:"Services"`
}
type TargetToQueue struct {
	ServiceName string `json:"ServiceName"`
	IP          string `json:"IP"`
	Port        int    `json:"Port"`
	TeamName    string `json:"TeamName"`
	Round       int    `json:"Round"`
	Tick        int    `json:"Tick"`
}
