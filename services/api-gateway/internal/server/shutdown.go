package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func SetupSignalContext(parent context.Context) context.Context {
	ctx, cancel := context.WithCancel(parent)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-ch
		cancel()
	}()

	return ctx
}
