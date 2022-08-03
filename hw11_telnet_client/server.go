package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Run(conn TelnetClient) error {
	log.Println("Start connect")
	if err := conn.Connect(); err != nil {
		return fmt.Errorf("connect error: %w", err)
	}

	log.Println("Connected success")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)
	done := make(chan struct{})

	// Send
	go func() {
		defer close(done)
		if err := conn.Send(); err != nil {
			log.Println("Send error: ", err)
			return
		}
	}()

	// Close
	go func() {
		defer conn.Close()
		select {
		case <-signalChan:
			close(done)
			return
		case <-done:
			return
		}
	}()

	// Receive
	go func() {
		defer func() {
			if _, ok := <-done; ok {
				close(done)
			}
		}()
		if err := conn.Receive(); err != nil {
			log.Println("Receive error: ", err)
		}
	}()

	<-done

	return nil
}
