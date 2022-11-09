package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gopaljee023/gomicroservices/handlers"
)

func main() {

	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	ph := handlers.NewProduct(l)

	sm := http.NewServeMux()
	sm.Handle("/", ph)

	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     l,
	}

	//used goroutines so that s.ListenAndServe() won't bloked
	go func() {
		fmt.Println("Starting server on port 9090")
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Kill)
	signal.Notify(c, os.Interrupt)

	// Block until a signal is received.
	sig := <-c
	fmt.Println("Got signal:", sig)

	l.Println("Received terminate, graceful duration", sig)

	//Timeout to clean before closing all connection
	d := time.Now().Add(time.Second * 30)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	l.Println("cancel status is ", cancel)
	s.Shutdown(ctx)
}
