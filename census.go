package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	URLCensusChar     = "http://census.daybreakgames.com/s:vAPP/get/ps2:v2/character/?name.first=%s&c:case=false&c:resolve=stat_history,faction,world,outfit_member_extended"
	URLCensusCharStat = "http://census.daybreakgames.com/s:vAPP/get/ps2:v2/characters_stat?character_id=%s&c:limit=5000"
)

var (
	NoSuchCharacterErr = errors.New("No such character")
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

type CensusChar struct {
}

func GenReport(char string) (*CensusChar, error) {
	charURL := fmt.Sprintf(URLCensusChar, char)

	rsp, err := http.Get(charURL)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	d := json.NewDecoder(rsp.Body)

	var cc struct {
		CharacterList []struct {
		} `json:"character_list"`
		Returned int `json:"returned"`
	}

	err = d.Decode(&cc)
	if err != nil {
		return nil, err
	}
	if cc.Returned == 0 {
		return nil, NoSuchCharacterErr
	}

	return nil, NoSuchCharacterErr
}
