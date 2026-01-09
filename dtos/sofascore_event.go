package dtos

// SofaScore Events/Matches DTOs

type SofaEventsResponse struct {
	Events []SofaEvent `json:"events"`
}

type SofaEvent struct {
	Tournament          SofaTournament `json:"tournament"`
	Season              SofaSeason     `json:"season"`
	HomeTeam            SofaTeam       `json:"homeTeam"`
	AwayTeam            SofaTeam       `json:"awayTeam"`
	Status              SofaStatus     `json:"status"`
	HomeScore           SofaScore      `json:"homeScore"`
	AwayScore           SofaScore      `json:"awayScore"`
	ID                  int            `json:"id"`
	StartTimestamp      int64          `json:"startTimestamp"`
	Slug                string         `json:"slug"`
	RoundInfo           RoundInfo      `json:"roundInfo,omitempty"`
	CustomID            string         `json:"customId,omitempty"`
	HasGlobalHighlights bool           `json:"hasGlobalHighlights,omitempty"`
}

type RoundInfo struct {
	Round        int    `json:"round"`
	Name         string `json:"name,omitempty"`
	Slug         string `json:"slug,omitempty"`
	CupRoundType int    `json:"cupRoundType,omitempty"`
}

type SofaEventDetail struct {
	Event SofaEvent `json:"event"`
}
