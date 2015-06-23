package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
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

type Outfit struct {
	Alias   string
	Name    string
	Members int
}

func (o *Outfit) String() string {
	if o == nil {
		return "None"
	}

	return fmt.Sprintf("[%v] %v (%v member%v)",
		o.Alias,
		o.Name,
		o.Members,
		plural(o.Members, "s"),
	)
}

type CensusChar struct {
	Name                        string
	Created                     time.Time
	LastLogin                   time.Time
	Played                      time.Duration
	Logins                      int
	Rank                        int
	Faction                     Faction
	Server                      Server
	Outfit                      *Outfit
	Score, Captures, Defenses   int
	Medals, Ribbons, Certs      int
	Kills /*, Assists*/, Deaths int
	KD                          float64
}

type rawChar struct {
	CharList []struct {
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
		Outfit *struct {
			Alias   string `json:"alias"`
			Name    string `json:"name"`
			Members string `json:"member_count"`
		} `json:"outfit_member,omitempty"`
		Stats struct {
			History []map[string]interface{} `json:"stat_history"`
		} `json:"stats"`
		WorldID Server `json:"world_id"`
	} `json:"character_list"`
	Returned int `json:"returned"`
}

func (rc *rawChar) findStat(stat string) int {
	for _, v := range rc.CharList[0].Stats.History {
		if v["stat_name"] == stat {
			val, _ := strconv.ParseInt(v["all_time"].(string), 10, 0)
			return int(val)
		}
	}

	panic("Stat not found: " + stat)
}

func GenReport(char string) (*CensusChar, error) {
	charURL := fmt.Sprintf(URLCensusChar, char)

	rsp, err := http.Get(charURL)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	d := json.NewDecoder(rsp.Body)

	var rc rawChar
	err = d.Decode(&rc)
	if err != nil {
		return nil, err
	}
	if rc.Returned == 0 {
		return nil, NoSuchCharacterErr
	}

	created, _ := strconv.ParseInt(rc.CharList[0].Times.Creation, 10, 64)
	lastLogin, _ := strconv.ParseInt(rc.CharList[0].Times.LastLogin, 10, 64)
	played, _ := strconv.ParseInt(rc.CharList[0].Times.MinPlayed, 10, 64)
	logins, _ := strconv.ParseInt(rc.CharList[0].Times.LoginCount, 10, 0)
	rank, _ := strconv.ParseInt(rc.CharList[0].BattleRank.Value, 10, 0)

	var outfit *Outfit
	if rc.CharList[0].Outfit != nil {
		outfitMembers, _ := strconv.ParseInt(rc.CharList[0].Outfit.Members, 10, 0)
		outfit = &Outfit{
			Alias:   rc.CharList[0].Outfit.Alias,
			Name:    rc.CharList[0].Outfit.Name,
			Members: int(outfitMembers),
		}
	}

	kills := rc.findStat("kills")
	deaths := rc.findStat("deaths")

	return &CensusChar{
		Name:      rc.CharList[0].Name.First,
		Created:   time.Unix(created, 0),
		LastLogin: time.Unix(lastLogin, 0),
		Played:    time.Duration(played) * time.Minute,
		Logins:    int(logins),
		Rank:      int(rank),
		Faction:   rc.CharList[0].Faction,
		Server:    rc.CharList[0].WorldID,
		Outfit:    outfit,
		Score:     rc.findStat("score"),
		Captures:  rc.findStat("facility_capture"),
		Defenses:  rc.findStat("facility_defend"),
		Medals:    rc.findStat("medals"),
		Ribbons:   rc.findStat("ribbons"),
		Certs:     rc.findStat("certs"),
		Kills:     kills,
		Deaths:    deaths,
		KD:        float64(kills) / float64(deaths),
	}, nil
}
