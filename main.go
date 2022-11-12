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
	"github.com/gorilla/mux"
)

func main() {

	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	ph := handlers.NewProduct(l)

	sm := mux.NewRouter()

	//curl get:
	/*curl localhost:9090 */
	getRouter := sm.Methods("GET").Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)

	//curl command:
	/*
		curl localhost:9090/1 -XPUT -d
		'{"name":"juice with capcunio ", "description":" coffee one", "price":1.45, "sku":"akkdc"}'
	*/
	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProducts)
	putRouter.Use(ph.MiddlewareValidateProduct) ////this will be executed before the ph.UpdateProducts

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.AddProduct)
	postRouter.Use(ph.MiddlewareValidateProduct) //this will be executed before the ph.AddProduct

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
