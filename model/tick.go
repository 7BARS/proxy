package model

type Tick struct {
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	TimeStamp uint64  `json:"time_stamp"`
}
