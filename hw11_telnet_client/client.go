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
	defer conn.Close()
	log.Println("Connected success")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)
	done := make(chan struct{})

	// Send
	go func() {
		if err := conn.Send(); err != nil {
			log.Println("Send error: ", err)
			closeDoneChannel(done)
			return
		}
	}()

	// Receive
	go func() {
		if err := conn.Receive(); err != nil {
			log.Println("Receive error: ", err)
			closeDoneChannel(done)
			return
		}
	}()

	for {
		select {
		case <-signalChan:
			return nil
		case <-done:
			return nil
		}
	}
}

func closeDoneChannel(done chan struct{}) {
	if _, ok := <-done; ok {
		close(done)
	}
}
