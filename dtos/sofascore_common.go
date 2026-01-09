package dtos

// SofaScore API common types

type SofaTournament struct {
	Name             string               `json:"name"`
	Slug             string               `json:"slug"`
	ID               int                  `json:"id"`
	Category         SofaCategory         `json:"category,omitempty"`
	UniqueTournament SofaUniqueTournament `json:"uniqueTournament,omitempty"`
}

type SofaUniqueTournament struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
	ID   int    `json:"id"`
}

type SofaCategory struct {
	Name   string `json:"name"`
	Slug   string `json:"slug"`
	ID     int    `json:"id"`
	Alpha2 string `json:"alpha2,omitempty"`
}

type SofaTeam struct {
	Name       string         `json:"name"`
	Slug       string         `json:"slug"`
	ShortName  string         `json:"shortName"`
	ID         int            `json:"id"`
	Country    SofaCategory   `json:"country,omitempty"`
	TeamColors SofaTeamColors `json:"teamColors,omitempty"`
}

type SofaTeamColors struct {
	Primary   string `json:"primary"`
	Secondary string `json:"secondary"`
	Text      string `json:"text"`
}

type SofaStatus struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

type SofaScore struct {
	Current    int `json:"current"`
	Display    int `json:"display"`
	Period1    int `json:"period1"`
	Period2    int `json:"period2"`
	Normaltime int `json:"normaltime"`
}

type SofaSeason struct {
	Name   string `json:"name"`
	Year   string `json:"year"`
	ID     int    `json:"id"`
	Editor bool   `json:"editor"`
}
