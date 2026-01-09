package dtos

// SofaScore Standings/Tables DTOs

type SofaStandingsResponse struct {
	Standings []SofaStanding `json:"standings"`
}

type SofaStanding struct {
	Tournament         SofaTournament    `json:"tournament"`
	Type               string            `json:"type"`
	Name               string            `json:"name"`
	Descriptions       []string          `json:"descriptions,omitempty"`
	TieBreakingRule    TieBreakingRule   `json:"tieBreakingRule,omitempty"`
	Rows               []SofaStandingRow `json:"rows"`
	ID                 int               `json:"id"`
	UpdatedAtTimestamp int64             `json:"updatedAtTimestamp"`
}

type TieBreakingRule struct {
	Text string `json:"text"`
	ID   int    `json:"id"`
}

type SofaStandingRow struct {
	Team          SofaTeam `json:"team"`
	Position      int      `json:"position"`
	Matches       int      `json:"matches"`
	Wins          int      `json:"wins"`
	ScoresFor     int      `json:"scoresFor"`
	ScoresAgainst int      `json:"scoresAgainst"`
	ID            int      `json:"id"`
	Losses        int      `json:"losses"`
	Draws         int      `json:"draws"`
	Points        int      `json:"points"`
	// Promotion can be either string or object, we'll ignore it for now
	// Promotion     string   `json:"promotion,omitempty"`
}
