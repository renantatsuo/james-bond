package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/renantatsuo/james-bond/internal/agent"
	"github.com/renantatsuo/james-bond/internal/env"
	"github.com/renantatsuo/james-bond/internal/ui"
)

var (
	apiKey = env.Get("API_KEY").Required().String().Parse()
)

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := agent.NewOpenAIClient(apiKey)
	a := agent.New(client)
	view := ui.New(a)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		cancel()
		view.Stop()
	}()

	view.Init(ctx)
}
