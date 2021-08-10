package model

type Target struct {
	ServiceName string
	IP 			string
	Port		int
	TeamName	string
	Round 		int
	Tick		int
}

type Service struct {
	ServiceName string
	Scripts []string
}
