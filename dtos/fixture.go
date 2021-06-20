package dtos

type Fixtures struct {
	Stages []struct {
		Events []struct {
			MatchID string `json:"Eid"`
			T1      []struct {
				Name string `json:"Nm"`
				TBD  int    `json:"tbd"`
			} `json:"T1"`
			T2 []struct {
				Name string `json:"Nm"`
				TBD  int    `json:"tbd"`
			} `json:"T2"`
			IncsX     int   `json:"IncsX"`
			ComX      int   `json:"ComX"`
			LuX       int   `json:"LuX"`
			StatX     int   `json:"StatX"`
			SubsX     int   `json:"SubsX"`
			SDFowX    int   `json:"SDFowX"`
			SDInnX    int   `json:"SDInnX"`
			StartTime int64 `json:"Esd"`
			LuUT      int64 `json:"LuUT"`
			Eds       int   `json:"Eds"`
			Edf       int   `json:"Edf"`
			EO        int   `json:"EO"`
			Eact      int   `json:"Eact"`
			Stage     struct {
				StageName string `json:"Snm"`
			} `json:"Stg"`
		} `json:"Events"`
	} `json:"Stages"`
}
