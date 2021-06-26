package dtos

type StageData struct {
	Stages []Stage `json:"Stages"`
}

type Stage struct {
	StageName string  `json:"Snm"`
	Events    []Event `json:"Events"`
}

type Event struct {
	Eid    string `json:"Eid"`
	Tr1    string `json:"Tr1"`
	Tr2    string `json:"Tr2"`
	Trh1   string `json:"Trh1"`
	Trh2   string `json:"Trh2"`
	Tr1OR  string `json:"Tr1OR"`
	Tr2OR  string `json:"Tr2OR"`
	T1     []Team `json:"T1"`
	T2     []Team `json:"T2"`
	IncsX  int    `json:"IncsX"`
	ComX   int    `json:"ComX"`
	LuX    int    `json:"LuX"`
	StatX  int    `json:"StatX"`
	SubsX  int    `json:"SubsX"`
	SDFowX int    `json:"SDFowX"`
	SDInnX int    `json:"SDInnX"`
	Esd    int64  `json:"Esd"`
	LuUT   int64  `json:"LuUT"`
	Eds    int64  `json:"Eds"`
	Edf    int64  `json:"Edf"`
	EO     int    `json:"EO"`
}

type Team struct {
	Nm   string `json:"Nm"`
	ID   string `json:"ID"`
	Tbd  int    `json:"tbd"`
	Img  string `json:"Img"`
	Gd   int    `json:"Gd"`
	CoNm string `json:"CoNm"`
	CoId string `json:"CoId"`
}
