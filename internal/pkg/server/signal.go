package server

import (
	"os"
	"os/signal"
	"syscall"
)

func SetupSignalHandler() <-chan int {
	handlerChan := make(chan int)

	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-c
		close(handlerChan)
		<-c
		os.Exit(1)
	}()

	return handlerChan
}
