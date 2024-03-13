package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	notesDel "github.com/t1d333/smartlectures/internal/notes/delivery/http"
	notesRep "github.com/t1d333/smartlectures/internal/notes/repository/pg"
	notesServ "github.com/t1d333/smartlectures/internal/notes/service"
	"github.com/t1d333/smartlectures/pkg/logger"
)

func main() {
	logger := logger.NewLogger()

	router := gin.Default()

	srv := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	// db connection

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	fmt.Println(os.Getenv("DB_URL"))
	pool, err := pgxpool.New(ctx, os.Getenv("DB_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	tmp := ""
	err = pool.QueryRow(context.Background(), "select '1'").Scan(&tmp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	// notes service
	notesRep := notesRep.NewRepository(logger, pool)
	notesServ := notesServ.NewService(logger, notesRep)
	notesDel := notesDel.NewDelivery(logger, router, notesServ)

	notesDel.RegisterHandler()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	<-ctx.Done()
	log.Println("timeout of 5 seconds.")
	log.Println("Server exiting")
}
