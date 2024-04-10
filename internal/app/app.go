package app

import (
	bannerHandler "banner-service/internal/pkg/banner/http"
	"banner-service/internal/pkg/banner/repository"
	"banner-service/internal/pkg/banner/service"
	"banner-service/internal/pkg/cache"
	"banner-service/internal/pkg/config"
	"banner-service/internal/utils/jwter"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	logger *logrus.Logger
}

func NewApp(logger *logrus.Logger) *App {
	return &App{logger: logger}
}

func (a *App) Run() error {
	wd, err := os.Getwd()
	cfg, err := config.Load(wd + "/configs/config.yaml")
	if err != nil {
		a.logger.Fatalln(err)
	}

	a.logger.Println(cfg.DBPort)
	db, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName))

	if err != nil {
		err = fmt.Errorf("error happened in sql.Open: %w", err)
		a.logger.Fatalln(err)
		return err
	}
	defer db.Close()

	if err = db.Ping(context.Background()); err != nil {
		a.logger.Fatalln(err)
		return err
	}

	rc := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.DB,
	})
	defer rc.Close()

	cacheClient := cache.NewRedisClient(rc)

	bannerRepo := repository.NewBannerRepository(db)
	bannerService := service.NewBannerService(bannerRepo, cacheClient)
	bannerHandler := bannerHandler.NewBannerHandler(bannerService, a.logger)

	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	r.HandleFunc("/user_banner", bannerHandler.GetBanner).Methods("GET")
	r.HandleFunc("/banner", bannerHandler.GetBannerList).Methods("GET")
	r.HandleFunc("/banner", bannerHandler.AddBanner).Methods("POST")
	r.HandleFunc("/banner/{id:[0-9]+}", bannerHandler.UpdateBanner).Methods("PATCH")
	r.HandleFunc("/banner/{id:[0-9]+}", bannerHandler.DeleteBanner).Methods("DELETE")

	srv := http.Server{
		Handler:           r,
		Addr:              cfg.Address,
		ReadTimeout:       cfg.Timeout,
		WriteTimeout:      cfg.Timeout,
		IdleTimeout:       cfg.IdleTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}

	jwter.LoadSecret(cfg.JWTSecret)

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			a.logger.Error("listen and serve returned err: ", err)
		}
	}()

	a.logger.Info("server started")
	sig := <-quit
	a.logger.Debug("handle quit chanel: ", sig.String())
	a.logger.Info("server stopping...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		a.logger.Error("server shutdown returned an err: ", err)
		err = fmt.Errorf("error happened in srv.Shutdown: %w", err)

		return err
	}

	return nil
}
