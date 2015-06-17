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

type Faction string

func (f Faction) String() string {
	switch f {
	case "1":
		return "Vanu Sovriegnty"
	case "2":
		return "New Conglomerate"
	case "3":
		return "Terran Republic"
	}

	panic("Unknown faction: " + f)
}

type CensusChar struct {
	Name      string
	Created   time.Time
	LastLogin time.Time
	Played    time.Duration
	Logins    int
	Rank      int
	Faction   Faction
	Server    Server
	Outfit    string
	Score     int
	Captures  int
	Defenses  int
	Medals    int
	Ribbons   int
	Certs     int
	Kills     int
	Assists   int
	Deaths    int
	KD        float64
}

func GenReport(char string) (*CensusChar, error) {
	charURL := fmt.Sprintf(URLCensusChar, char)

	rsp, err := http.Get(charURL)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	d := json.NewDecoder(rsp.Body)

	var rawChar struct {
		CharacterList []struct {
			Name struct {
				First string `json:"first"`
			} `json:"name"`
			Faction Faction `json:"faction_id"`
			Times   struct {
				Creation   string `json:"creation"`
				LastLogin  string `json:"last_login"`
				LoginCount string `json:"login_count"`
				MinPlayed  string `json:"minutes_played"`
			} `json:"times"`
			BattleRank struct {
				Value string `json:"value"`
			} `json:"battle_rank"`
		} `json:"character_list"`
		Returned int `json:"returned"`
		Stats    struct {
			History []map[string]interface{} `json:"stat_history"`
		} `json:"stats"`
	}

	findStat := func(stat string) int {
		for _, v := range rawChar.Stats.History {
			if v["stat_name"] == stat {
				val, _ := strconv.ParseInt(v["all_time"].(string), 10, 0)
				return int(val)
			}
		}

		panic("Stat not found: " + stat)
	}

	err = d.Decode(&rawChar)
	if err != nil {
		return nil, err
	}
	if rawChar.Returned == 0 {
		return nil, NoSuchCharacterErr
	}

	return nil, NoSuchCharacterErr
}
