package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/t1d333/smartlectures/cmd/auth/config"
	delivery "github.com/t1d333/smartlectures/internal/auth/delivery/http"
	"github.com/t1d333/smartlectures/internal/auth/repository"
	"github.com/t1d333/smartlectures/internal/auth/service"
	middl "github.com/t1d333/smartlectures/internal/middleware"
	"github.com/t1d333/smartlectures/pkg/logger"
)

func ExampleClient() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "dragonfly:6379",
		Password: "password",
		DB:       0,
	})

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
}

func main() {
	logger := logger.NewLogger()

	cfg, err := config.NewConfig(os.Getenv("CONFIG_PATH"))
	if err != nil {
		logger.Fatal(err)
	}

	router := gin.Default()

	router.Use(middl.NewErrorHandler(logger))
	router.Use(cors.Default())

	srv := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.DragonflyURL,
		Password: cfg.DragonflyPassword,
		DB:       0,
	})

	if status := redisClient.Ping(context.Background()); status.Err() != nil {
		logger.Fatal(status.Err())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.PostgresURL)
	if err != nil {
		logger.Fatal(err)
	}

	rep := repository.New(logger, redisClient, pool)
	svc := service.New(logger, rep)
	del := delivery.NewDelivery(logger, router, svc)

	del.RegisterHandler()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutdown Server ...")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server Shutdown:", err)
	}

	<-ctx.Done()
	logger.Info("timeout of 5 seconds.")
	logger.Info("Server exiting")
}
