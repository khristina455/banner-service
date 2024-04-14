package app

import (
	"banner-service/internal/utils/jwter"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"

	authHandler "banner-service/internal/pkg/auth/http"
	authRepository "banner-service/internal/pkg/auth/repository"
	authService "banner-service/internal/pkg/auth/sevice"
	bannerHandler "banner-service/internal/pkg/banner/http"
	bannerRepository "banner-service/internal/pkg/banner/repository"
	bannerService "banner-service/internal/pkg/banner/service"
	"banner-service/internal/pkg/cache"
	"banner-service/internal/pkg/config"
	"banner-service/internal/pkg/middleware"
)

type App struct {
	logger *logrus.Logger
}

func NewApp(logger *logrus.Logger) *App {
	return &App{logger: logger}
}

func (a *App) Run() error {
	wd, _ := os.Getwd()
	cfg, err := config.Load(wd + "/configs/config.yaml")
	if err != nil {
		a.logger.Fatalln(err)
	}

	db, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName))

	if err != nil {
		err = fmt.Errorf("error happened in sql.Open: %w", err)
		a.logger.Error(err)
		return err
	}
	defer db.Close()

	if err = db.Ping(context.Background()); err != nil {
		a.logger.Error(err)
		return err
	}

	rc := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.DB,
	})
	defer rc.Close()

	tokenManager := jwter.New(cfg.JWTSecret, cfg.JWTTTL)

	cacheClient := cache.NewRedisClient(rc, cfg.RedisTTL)

	bannerRepo := bannerRepository.NewBannerRepository(db)
	bannerService := bannerService.NewBannerService(bannerRepo, cacheClient)
	bannerHandler := bannerHandler.NewBannerHandler(bannerService, a.logger)

	authRepo := authRepository.NewAuthRepository(db)
	authService := authService.NewAuthService(authRepo)
	authHandler := authHandler.NewAuthHandler(authService, a.logger, tokenManager)

	mw := middleware.New(a.logger, tokenManager)

	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	r.Handle("/user_banner", mw.Auth(false, http.HandlerFunc(bannerHandler.GetBanner))).Methods("GET")
	r.Handle("/banner", mw.Auth(true, http.HandlerFunc(bannerHandler.GetBannerList))).Methods("GET")
	r.Handle("/banner", mw.Auth(true, http.HandlerFunc(bannerHandler.AddBanner))).Methods("POST")
	r.Handle("/banner/{id:[0-9]+}", mw.Auth(true,
		http.HandlerFunc(bannerHandler.UpdateBanner))).Methods("PATCH")
	r.Handle("/banner/{id:[0-9]+}", mw.Auth(true,
		http.HandlerFunc(bannerHandler.DeleteBanner))).Methods("DELETE")
	r.HandleFunc("/sign_in", authHandler.SignIn).Methods("POST")
	r.HandleFunc("/sign_up", authHandler.SignUp).Methods("POST")

	srv := http.Server{
		Handler:           r,
		Addr:              cfg.Address,
		ReadTimeout:       cfg.Timeout,
		WriteTimeout:      cfg.Timeout,
		IdleTimeout:       cfg.IDleTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}

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
