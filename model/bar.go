package model

import "fmt"

type Bar struct {
	// ohlcv full bar
	Name      string
	Close     float64
	TimeStamp uint64
}

func (b Bar) String() string {
	return fmt.Sprintf("minute bar ticker:%v, time:%d, price:%f", b.Name, b.TimeStamp, b.Close)
}
