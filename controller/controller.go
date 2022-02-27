package controller

import (
	"context"
	"log"
	"net/http"
	"proxy/internal/barbuilder"
	"proxy/model"
)

type Controller struct {
	barsCh     chan model.Bar
	tickCh     chan model.Tick
	errCh      chan error
	urls       []string
	symbols    []string
	streaming  []PriceStreamSubscriber
	cancelFn   context.CancelFunc
	barBuilder *barbuilder.BarBuilder
	client     *http.Client
}

const (
	tickChLen = 100
	errChLen  = 10
)

func NewController(urls, symbols []string, barsCh chan model.Bar) *Controller {
	barbuilder := barbuilder.NewBarBuilder(barsCh)
	return &Controller{
		tickCh:     make(chan model.Tick, tickChLen),
		errCh:      make(chan error, errChLen),
		client:     http.DefaultClient,
		barBuilder: barbuilder,
		urls:       urls,
		symbols:    symbols,
	}
}

func (c *Controller) Start(ctx context.Context) {
	log.Printf("start controller")
	childCtx, cancel := context.WithCancel(ctx)
	c.cancelFn = cancel
	go c.run(childCtx)
	for _, url := range c.urls {
		streaming := NewPriceStreaming(url, c.symbols[0], http.DefaultClient, c.tickCh, c.errCh)
		streaming.Start(childCtx)
		c.streaming = append(c.streaming, streaming)
	}
}

func (c *Controller) ReadBarCh() model.Bar {
	return <-c.barsCh
}

func (c *Controller) Stop() {
	log.Printf("stop controller")
	for _, streaming := range c.streaming {
		streaming.Stop()
	}
	c.cancelFn()
}

func (c *Controller) Destroy() {
	log.Printf("destroy controller")
	for _, streaming := range c.streaming {
		streaming.Destroy()
	}
	c.cancelFn()
	close(c.tickCh)
	close(c.errCh)
	close(c.barsCh)
}

func (c *Controller) run(ctx context.Context) {
	for {
		select {
		case tick := <-c.tickCh:
			c.barBuilder.AppendTick(tick)
		case err := <-c.errCh:
			log.Printf("exchange return error: %v", err)
		case <-ctx.Done():
			return
		}
	}
}
