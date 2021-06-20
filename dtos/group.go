package dtos

type TableData struct {
	Stages []struct {
		StageName      string `json:"Snm"`
		TournamentName string `json:"Cnm"`
		LeagueTable    struct {
			L []struct {
				Tables []struct {
					LTT  int           `json:"LTT"`
					Team []TeamData    `json:"team"`
					PhrX []interface{} `json:"phrX"`
				} `json:"Tables"`
			} `json:"L"`
		} `json:"LeagueTable"`
	} `json:"Stages"`
}

type TeamData struct {
	Rank   int           `json:"rnk"`
	Win    int           `json:"win"`
	Wreg   int           `json:"wreg"`
	Wot    int           `json:"wot"`
	Wap    int           `json:"wap"`
	Draw   int           `json:"drw"`
	Lost   int           `json:"lst"`
	Lreg   int           `json:"lreg"`
	Lot    int           `json:"lot"`
	Lap    int           `json:"lap"`
	Gf     int           `json:"gf"`
	Ga     int           `json:"ga"`
	Gd     int           `json:"gd"`
	Points int           `json:"pts"`
	Ptsn   string        `json:"ptsn"`
	Pld    int           `json:"pld"`
	Tid    string        `json:"Tid"`
	Name   string        `json:"Tnm"`
	Phr    []int         `json:"phr,omitempty"`
	Ipr    int           `json:"Ipr"`
	Com    []interface{} `json:"com"`
	Pa     int           `json:"pa"`
	Pf     int           `json:"pf"`
	Img    string        `json:"Img"`
}
