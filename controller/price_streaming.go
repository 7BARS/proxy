package controller

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"proxy/model"
	"strings"
)

type PriceStreaming struct {
	baseURL  string
	symbol   string
	stopCh   chan struct{}
	errCh    chan error
	tickCh   chan model.Tick
	client   *http.Client
	cancelFn context.CancelFunc
}

func NewPriceStreaming(baseURL, symbol string, client *http.Client, tickCh chan model.Tick, errCh chan error) *PriceStreaming {
	return &PriceStreaming{
		baseURL: baseURL,
		symbol:  symbol,
		stopCh:  make(chan struct{}),
		errCh:   errCh,
		tickCh:  tickCh,
		client:  client,
	}
}

func (s *PriceStreaming) Start(ctx context.Context) {
	log.Printf("start price streaming")
	childCtx, cancel := context.WithCancel(ctx)
	s.cancelFn = cancel
	go s.run(childCtx)
}

func (s *PriceStreaming) Stop() {
	log.Printf("stop price streaming")
	s.cancelFn()
}

func (s *PriceStreaming) Destroy() {
	log.Printf("destroy price streaming")
	s.cancelFn()
	close(s.errCh)
	close(s.tickCh)
}

func (s *PriceStreaming) run(ctx context.Context) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.baseURL, nil)
	if err != nil {
		s.pushError(err)
		return
	}
	const tickerStr = "ticker"
	query := req.URL.Query()
	if query.Get(tickerStr) == "" {
		query.Set(tickerStr, s.symbol)
		req.URL.RawQuery = query.Encode()
	}
	resp, err := s.client.Do(req)
	if err != nil {
		s.pushError(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.pushError(fmt.Errorf("status code:%v is not 200, url:%v, ", resp.StatusCode, req.URL.RawPath))
		return
	}

	log.Printf("start new streaming connection")
	reader := bufio.NewReader(resp.Body)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			data, err := reader.ReadBytes('\n')
			if err != nil {
				s.pushError(err)
				return
			}
			log.Printf("streaming message:%s", string(data))
			if strings.Contains(string(data), "msg_err") {
				var msgErr *model.MessageError
				err = json.Unmarshal(data, &msgErr)
				if err != nil {
					s.pushError(err)
					return
				}
				s.pushError(fmt.Errorf("%s", msgErr.Message))
				return
			}

			var tick *model.Tick
			err = json.Unmarshal(data, &tick)
			if err != nil {
				s.pushError(err)
				return
			}

			s.pushTick(*tick)
		}
	}
}

func (s *PriceStreaming) pushTick(tick model.Tick) {
	s.tickCh <- tick
}

func (s *PriceStreaming) pushError(err error) {
	s.errCh <- err
}
