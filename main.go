package main

import (
	"context"
	"dapr-im/chat"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func init() {
	log.SetFlags(0)
}
func main() {
	start()
}

// run initializes the chatServer and then
// starts a http.Server for the passed in address.
func start() error {
	l, err := net.Listen("tcp", "localhost:80")
	if err != nil {
		return err
	}
	log.Printf("listening on http://%v", l.Addr())

	cs := chat.NewChatServer()
	s := &http.Server{
		Handler:      cs,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	errc := make(chan error, 1)
	go func() {
		errc <- s.Serve(l)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		log.Printf("failed to serve: %v", err)
	case sig := <-sigs:
		log.Printf("exit & terminating: %v", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return s.Shutdown(ctx)
}
