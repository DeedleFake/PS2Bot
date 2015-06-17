package main

import (
	"github.com/jzelinskie/geddit"
	"log"
	"os"
	"text/template"
)

const (
	Username  = ""
	Password  = ""
	UserAgent = "Planetside 2 Stats Poster"

	URLDasanFall = "[[dasanfall]](http://stats.dasanfall.com/ps2/player/%s)"
	URLFisu      = "[[fisu]](http://ps2.fisu.pw/player/?name=%s)"
	URLPSU       = "[[psu]](http://www.planetside-universe.com/character-%s.php)"
	URLPlayers   = "[[players]](https://www.planetside2.com/players/#!/%s)"
	URLKillboard = "[[killboard]](https://www.planetside2.com/players/#!/%s/killboard)"

	PostReplyTemplate = `
**Some stats about {{.char_name_truecase}}.**

------

- Character created: {{.char_creation}}
- Last login: {{.char_login}}
- Time played: {{.char_playtime}} ({{.char_logins}} login{{.login_plural}})
- Battle rank: {{.char_rank}}
- Faction: {{.char_faction_en}}
- Server: {{.char_server}}
- Outfit: {{.char_outfit}}
- Score: {{.char_score}} | Captured: {{.char_captures}} | Defended: {{.char_defended}}
- Medals: {{.char_medals}} | Ribbons: {{.char_ribbons}} | Certs: {{.char_certs}}
- Kills: {{.char_kills}} | Assists: {{.char_assists}} | Deaths: {{.char_deaths}} | KDR: {{.char_kdr}}
- Links: {{.links_dasanfall}} {{.links_fisu}} {{.links_psu}} {{.links_players}} {{.links_killboard}}

------

^^This ^^post ^^was ^^made ^^by ^^a ^^bot.
^^Have ^^feedback ^^or ^^a ^^suggestion?
[^^\[pm ^^the ^^creator\]]
(https://np.reddit.com/message/compose/?to=microwavable_spoon&subject=PS2Bot%20Feedback)
^^| [^^\[see ^^my ^^code\]](https://github.com/plasticantifork/PS2Bot)
`
)

var (
	r    *geddit.LoginSession
	tmpl *template.Template
)

func init() {
	ls, err := geddit.NewLoginSession(Username, Password, UserAgent)
	if err != nil {
		log.Printf("Failed to login to Reddit: %v", err)
		os.Exit(1)
	}
	r = ls

	tmpl = template.Must(template.New("post-reply").Parse(PostReplyTemplate))
}
