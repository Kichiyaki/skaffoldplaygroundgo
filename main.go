package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	log.Println("hello from Go")

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("skaffold"))
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func(srv *http.Server) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		<-ctx.Done()

		ctxShutdown, cancelCtxShutdown := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelCtxShutdown()
		if err := srv.Shutdown(ctxShutdown); err != nil {
			log.Println("srv.Shutdown:", err)
		}
	}(srv)

	log.Println("Server is listening on the port 8080")

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalln("srv.ListenAndServe:", err)
	}
}
