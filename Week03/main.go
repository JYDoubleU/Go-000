package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello World"))
	})

	server := &http.Server{Addr: ":8080", Handler: nil}
	g.Go(func() error {
		return server.ListenAndServe()
	})

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	g.Go(func() error {
		select {
		case <-sigCh:
		case <-ctx.Done():
		}
		return server.Close()
	})

	log.Println("main: http server started on :8080")

	if err := g.Wait(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}

	log.Println("main: http server closed")
}
