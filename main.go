package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func main() {
	defer func() {
		log.Println("bye")
	}()
	log.Println("hello from Go")

	db, err := newDB()
	if err != nil {
		log.Fatalln("newDB:", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			w.WriteHeader(500)
			return
		}

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

func newDB() (*bun.DB, error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(os.Getenv("DATABASE_DSN"))))

	db := bun.NewDB(sqldb, pgdialect.New())

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping: %w", err)
	}

	return db, nil
}
