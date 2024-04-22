package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gin-contrib/cors"
	middl "github.com/t1d333/smartlectures/internal/middleware"
	notesDel "github.com/t1d333/smartlectures/internal/notes/delivery/http"
	notesRep "github.com/t1d333/smartlectures/internal/notes/repository/pg"
	notesServ "github.com/t1d333/smartlectures/internal/notes/service"
	storage "github.com/t1d333/smartlectures/internal/storage"

	snippetsDel "github.com/t1d333/smartlectures/internal/snippets/delivery/http"
	snippetsRep "github.com/t1d333/smartlectures/internal/snippets/repository/pg"
	snippetsServ "github.com/t1d333/smartlectures/internal/snippets/service"

	dirsDel "github.com/t1d333/smartlectures/internal/dirs/delivery/http"
	dirsRep "github.com/t1d333/smartlectures/internal/dirs/repository/pg"
	dirsServ "github.com/t1d333/smartlectures/internal/dirs/service"

	recDel "github.com/t1d333/smartlectures/internal/recognizer/delivery/http"
	recServ "github.com/t1d333/smartlectures/internal/recognizer/service"

	"github.com/t1d333/smartlectures/pkg/logger"
)

func main() {
	logger := logger.NewLogger()

	router := gin.Default()

	router.Use(cors.Default())

	router.Use(middl.NewErrorHandler(logger))

	srv := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	logger.Infow("creating db connection...")

	pool, err := pgxpool.New(ctx, os.Getenv("DB_URL"))
	if err != nil {
		logger.Fatal("unable to create connection pool", err)
	}
	defer pool.Close()

	tmp := ""

	logger.Infow("checking db connection...")

	err = pool.QueryRow(context.Background(), "select '1'").Scan(&tmp)
	if err != nil {
		logger.Fatalf("queryRow failed: %v\n", err)
		os.Exit(1)
	}

	logger.Infow("db connection successfully")

	// notes

	logger.Info("creating storage client...")
	conn, err := grpc.Dial(
		"storage:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Fatalf("failed to create storage client %s", err)
	}
	defer conn.Close()
	logger.Info("storage client cratead successfully")

	storageClient := storage.NewStorageClient(conn)

	notesRep := notesRep.NewRepository(logger, pool)
	notesServ := notesServ.NewService(logger, notesRep, storageClient)
	notesDel := notesDel.NewDelivery(logger, router, notesServ)

	notesDel.RegisterHandler()

	// dirs
	dirsRep := dirsRep.NewRepository(logger, pool)
	dirsServ := dirsServ.NewService(logger, dirsRep, storageClient)
	dirsDel := dirsDel.NewDelivery(logger, router, dirsServ)

	dirsDel.RegisterHandler()

	// snippets

	snipRep := snippetsRep.NewRepository(logger, pool)
	snipServ := snippetsServ.NewService(logger, snipRep)
	snipDel := snippetsDel.NewDelivery(logger, router, snipServ)

	snipDel.RegisterHandler()

	// recognizer

	recognizerServ := recServ.NewService(logger, "recognizer:50051")
	recognizerDel := recDel.NewDelivery(logger, router, recognizerServ)

	recognizerDel.RegisterHandler()

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
