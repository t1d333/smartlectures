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

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/t1d333/smartlectures/internal/images"
	delivery "github.com/t1d333/smartlectures/internal/images/delivery/http"
	"github.com/t1d333/smartlectures/internal/images/repository/s3"
	"github.com/t1d333/smartlectures/internal/images/service"
	middl "github.com/t1d333/smartlectures/internal/middleware"
	"github.com/t1d333/smartlectures/pkg/logger"
)

func main() {
	appCfg, err := images.NewConfig(os.Getenv("CONFIG_PATH"))
	if err != nil {
		log.Fatal("failed to read config file", err)
	}

	logger := logger.NewLogger()

	// TODO: lower multipart form size
	router := gin.Default()

	router.Use(cors.Default())

	router.Use(middl.ErrorHandler)

	rep, err := s3.NewRepository(logger, appCfg)
	if err != nil {
		logger.Fatal("failed to create s3 repository", err)
	}

	serv := service.NewService(logger, rep)

	del := delivery.NewDelivery(logger, router, serv)

	del.RegisterHandler()

	fmt.Printf("%s:%d", appCfg.Address, appCfg.Port)
	srv := &http.Server{
		// Addr:    fmt.Sprintf("%s:%d", appCfg.Address, appCfg.Port),
		Addr: ":8000",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server Shutdown:", err)
	}

	<-ctx.Done()
	logger.Info("timeout of 5 seconds.")
	logger.Info("Server exiting")
}
