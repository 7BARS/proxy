package main

import (
	"context"
	"fmt"
	"proxy/view"
)

const (
	countOfExchange = 100
	defaultPort     = 20000
)

func main() {
	urls := make([]string, 0, countOfExchange)
	for i := defaultPort; i < defaultPort+countOfExchange; i++ {
		urls = append(urls, fmt.Sprintf("http://localhost:%v/streaming", i))
	}
	stdout := view.NewStdout(urls, []string{"BTC_USD"})
	stdout.Start(context.Background())
}
