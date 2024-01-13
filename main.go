package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"product-api/handlers"
	"time"
)

func main() {

	l := log.New(os.Stdout, "product-api - ", log.LstdFlags)

	hh := handlers.NewHello(l)
	gb := handlers.NewGoodbye(l)

	ph := handlers.NewProducts(l)

	serveMux := http.NewServeMux()
	serveMux.Handle("/", hh)
	serveMux.Handle("/goodbye", gb)
	serveMux.Handle("/products", ph)
	serveMux.Handle("/products/", ph)

	s := &http.Server{
		Addr:         ":9090",
		Handler:      serveMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		l.Printf(" - 127.0.0.1:9090 listening...")
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}

	}()

	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, os.Kill)

	sig := <-signalChannel
	l.Println("Received terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)

	s.Shutdown(tc)
}
