package main

type RankCurrent		struct {
	Combat				int `json:"Combat"`
	Trade				int `json:"Trade"`
	Explore				int `json:"Explore"`
	CQC					int `json:"CQC"`
	Federation			int `json:"Federation"`
	Empire				int `json:"Empire"`
}

type RankProgress		struct {
	Combat				int `json:"Combat"`
	Trade				int `json:"Trade"`
	Explore				int `json:"Explore"`
	CQC					int `json:"CQC"`
	Federation			int `json:"Federation"`
	Empire				int `json:"Empire"`
}

type RankVerbose		struct {
	Combat				string `json:"Combat"`
	Trade				string `json:"Trade"`
	Explore				string `json:"Explore"`
	CQC					string `json:"CQC"`
	Federation			string `json:"Federation"`
	Empire				string `json:"Empire"`
}

type CmdrRankLog 		struct {
	Msgnum				int `json:"msgnum"`
	Msg					string `json:"msg"`
	Current				RankCurrent `json:"ranks"`
	Progress			RankProgress `json:"progress"`
	Verbose				RankVerbose `json:"ranksVerbose"`
}

type Credits 			struct {
	Balance				int `json:"balance"`
	Loan				int `json:"loan"`
	Date				string `json:"data"`
}

type CmdrCreditLog		struct {
	Msgnum				int `json:"msgnum"`
	Msg					string `json:"msg"`
	Credits				[]Credits `json:"credits"`
}

type Item				struct {
	Ntype				int `json:"type"`
	Type				string `json:"type"`
	Name				string `json:"name"`
	Qty					int `json:"qty"`
}

type CmdrInventoryLog	struct {
	Msgnum				int `json:"msgnum"`
	Msg					string `json:"msg"`
	Items				[]Item `json:"materials"`
}

type FlightLogEntry		struct {
	ShipId				int `json:"shipId"`
	System				string `json:"system"`
	SystemId			int `json:"systemId"`
	FirstDiscover		bool `json:"firstDiscover"`
	Date				string `json:"date"`
}

type CmdrFlightLog		struct {
	Msgnum				int `json:"msgnum"`
	Msg					string `json:"msg"`
	StartDate			string `json:"startDateTime"`
	EndDate				string `json:"endDateTime"`
	Logs				[]FlightLogEntry `json:"logs"`
	LastPos				CmdrLastPosition
}

type CmdrLastPosition	struct {
	Msgnum				int `json:"msgnum"`
	Msg					string `json:"msg"`
	System				string `json:"system"`
	SystemId			int `json:"systemId"`
	FirstDiscover		bool `json:"firstDiscover"`
	Date				string `json:"date"`
	Profile				string `json:"url"`
}

type CmdrLog			struct {
	Name				string
	Rank				CmdrRankLog
	Credits				CmdrCreditLog
	Data				CmdrInventoryLog
	Materials			CmdrInventoryLog
	FlightLog			CmdrFlightLog
}
