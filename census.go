package main

import (
	"net/http"
)

const (
	URLCensusChar     = "http://census.daybreakgames.com/s:vAPP/get/ps2:v2/character/?name.first=%s&c:case=false&c:resolve=stat_history,faction,world,outfit_member_extended"
	URLCensusCharStat = "http://census.daybreakgames.com/s:vAPP/get/ps2:v2/characters_stat?character_id=%s&c:limit=5000"
)

type Server string

func (s Server) String() string {
	switch s {
	case "1":
		return "Connery (US West)"
	case "17":
		return "Emerald (US East)"
	case "10":
		return "Miller (EU)"
	case "13":
		return "Cobalt (EU)"
	case "25":
		return "Briggs (AU)"
	case "19":
		return "Jaeger"
	}

	panic("Unknown server: " + s)
}
