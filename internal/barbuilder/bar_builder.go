package barbuilder

import (
	"proxy/model"
	"time"
)

type BarBuilder struct {
	barsCh      chan model.Bar
	minuterPrev int
}

func NewBarBuilder(barsCh chan model.Bar) *BarBuilder {
	return &BarBuilder{
		barsCh: barsCh,
	}
}

func (bb *BarBuilder) AppendTick(tick model.Tick) {
	parsedTime := time.Unix(int64(tick.TimeStamp), 0)
	_, minuteParsed, _ := parsedTime.Clock()
	if bb.minuterPrev != 0 && minuteParsed != 0 && minuteParsed != bb.minuterPrev {
		bb.pushBar(model.Bar{
			Name:      tick.Name,
			Close:     tick.Price,
			TimeStamp: tick.TimeStamp,
		})
	}
	bb.minuterPrev = minuteParsed
}

func (bb *BarBuilder) pushBar(bar model.Bar) {
	bb.barsCh <- bar
}
