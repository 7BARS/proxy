package view

import (
	"context"
	"log"
	"proxy/controller"
	"proxy/model"
)

type Stdout struct {
	cancelFn   context.CancelFunc
	controller *controller.Controller
	barsCh     chan model.Bar
}

const barsChLen = 50

func NewStdout(urls, symbols []string) *Stdout {
	barsCh := make(chan model.Bar, barsChLen)
	return &Stdout{
		barsCh:     barsCh,
		controller: controller.NewController(urls, symbols, barsCh),
	}
}

func (s *Stdout) Start(ctx context.Context) {
	log.Printf("start stdout")
	childCtx, cancelFn := context.WithCancel(ctx)
	s.cancelFn = cancelFn
	s.controller.Start(childCtx)
	s.run(childCtx)
}

func (s *Stdout) Stop() {
	log.Printf("stop stdout")
	s.cancelFn()
	s.controller.Stop()
}

func (s *Stdout) Destroy() {
	log.Printf("destroy stdout")
	s.cancelFn()
	s.controller.Destroy()
	close(s.barsCh)
}

func (s *Stdout) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case bar := <-s.barsCh:
			log.Print(bar)
		}
	}
}
