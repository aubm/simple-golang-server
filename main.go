package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: newServerMux(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err.Error())
		}
	}()
	log.Print("server started")

	// Block until a signal is received.
	s := <-c
	log.Println("got signal:", s)

	// emulate some graceful shutdown operations
	time.Sleep(60 * time.Second)

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("server shutdown failed:%+v", err)
	}
	log.Print("server exited properly")
}

func newServerMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/long", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s\n", r.Method, r.RequestURI)
		time.Sleep(10 * time.Second)
		_, _ = w.Write([]byte("That was a long process!"))
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s\n", r.Method, r.RequestURI)
		_, _ = w.Write([]byte("Hello from Golang server ❤️"))
	})

	return mux
}
