package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func signalHandler() <-chan bool {
	ch := make(chan bool)
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-s // consume channel
		log.Println("Received signal", sig)
		close(ch)
	}()
	return ch
}
