package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/lptnkv/microservice_example/handlers"
)

func main() {
	logger := log.New(os.Stdout, "product-api", log.LstdFlags)
	productHandler := handlers.NewProductsHandler(logger)
	goodbyeHandler := handlers.NewGoodbye(logger)
	serveMux := http.NewServeMux()
	serveMux.Handle("/", productHandler)
	serveMux.Handle("/goodbye", goodbyeHandler)

	server := &http.Server{
		Addr:         ":9090",
		Handler:      serveMux,
		IdleTimeout:  120 * time.Second,
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt)
	signal.Notify(signalChan, os.Kill)

	sig := <-signalChan
	logger.Println("Received terminate, shutting down", sig)

	timeoutContext, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(timeoutContext)
}
