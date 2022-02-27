package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"proxy/model"
	"strconv"
	"time"
)

const (
	defaultPort      = 20000
	routeStriming    = "/streaming"
	countOfExchanges = 100
	startPrice       = 100.0
	tickTimeRangeMs  = 2000
	timeErrorRangeS  = 300
)

func main() {
	rand.Seed(time.Now().UnixNano())
	http.HandleFunc(routeStriming, tick)

	for i := defaultPort; i < defaultPort+countOfExchanges; i++ {
		go runServer(i)
	}

	select {}
}

func runServer(port int) {
	portStr := strconv.Itoa(port)
	log.Printf("exchange is starting, url: http://localhost:%s%s", portStr, routeStriming)

	if err := http.ListenAndServe(":"+portStr, nil); err != nil {
		log.Printf("exchange is stop with error:%v", err)
		log.Printf("service is stop")
		os.Exit(1)
	}
}

func tick(w http.ResponseWriter, req *http.Request) {
	ticker := req.URL.Query().Get("ticker")
	if ticker != "BTC_USD" {
		data, err := json.Marshal(model.MessageError{
			Message: fmt.Sprintf("ticker:%v does not exist\"}", ticker),
		})
		if err != nil {
			log.Printf("cannot marshal error:%v", err)
			return
		}

		w.Write(append(data, '\n'))
		return
	}

	log.Printf("new connection to streaming")
	price := startPrice
	stopTime := timeToError()
	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Printf("HTTP buffer flushing is not implemented")
		return
	}
	flusher.Flush()
	for !time.Now().After(stopTime) {
		time.Sleep(tickTime())

		price += priceChange()
		data, err := json.Marshal(model.Tick{
			Name:      ticker,
			Price:     price,
			TimeStamp: uint64(time.Now().Unix()),
		})
		if err != nil {
			data, err = json.Marshal(model.MessageError{
				Message: err.Error(),
			})
			if err != nil {
				log.Printf("cannot marshal error:%v", err)
				return
			}

			w.Write(append(data, '\n'))
		}
		w.Write(append(data, '\n'))
		flusher.Flush()
	}
	fmt.Fprintf(w, "{\"msg_err\":\"exchange is stop\"}\n")
	data, err := json.Marshal(model.MessageError{
		Message: "exchange is stop",
	})
	if err != nil {
		log.Printf("cannot marshal error:%v", err)
		return
	}

	w.Write(append(data, '\n'))
	log.Printf("connection is closed")
}

func tickTime() time.Duration {
	r := rand.Intn(tickTimeRangeMs)

	return time.Duration(r) * time.Millisecond
}

func timeToError() time.Time {
	r := rand.Intn(timeErrorRangeS)

	return time.Now().Add(time.Duration(r * int(time.Second)))
}

func priceChange() float64 {
	price := 0.0
	price += rand.Float64()

	if rand.Int()%2 == 0 {
		return -price
	}

	return price
}
