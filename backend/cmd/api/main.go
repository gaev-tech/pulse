// @title           Pulse API
// @version         1.0
// @description     Трекер личных и командных задач
// @host            localhost
// @BasePath        /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/gaevivan/pulse/internal/infrastructure/config"
	"github.com/gaevivan/pulse/internal/infrastructure/email"
	infrajwt "github.com/gaevivan/pulse/internal/infrastructure/jwt"
	"github.com/gaevivan/pulse/internal/infrastructure/logger"
	"github.com/gaevivan/pulse/internal/infrastructure/postgres"

	authmiddleware "github.com/gaevivan/pulse/internal/handler/middleware"
	v1 "github.com/gaevivan/pulse/internal/handler/v1"

	repomagiclink "github.com/gaevivan/pulse/internal/repository/postgres/magic_link"
	repopat "github.com/gaevivan/pulse/internal/repository/postgres/pat"
	reporefreshtoken "github.com/gaevivan/pulse/internal/repository/postgres/refresh_token"
	repouser "github.com/gaevivan/pulse/internal/repository/postgres/user"

	"github.com/gaevivan/pulse/internal/repository/migrations"
	userusecase "github.com/gaevivan/pulse/internal/usecase/user"

	_ "github.com/gaevivan/pulse/api/docs"
)

func main() {
	cfg := config.Load()

	log, err := logger.New(cfg.Env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = log.Sync() }()

	log.Info("running migrations...")
	if err := postgres.Migrate(cfg.Database, migrations.FS); err != nil {
		log.Fatal("failed to run migrations", zap.Error(err))
	}
	log.Info("migrations applied")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := postgres.New(ctx, cfg.Database)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	log.Info("connected to database")

	// Repositories
	userRepo := repouser.New(pool)
	magicLinkRepo := repomagiclink.New(pool)
	refreshTokenRepo := reporefreshtoken.New(pool)
	patRepo := repopat.New(pool)

	// Infrastructure
	jwtManager := infrajwt.New(cfg.JWT.Secret)
	emailSender := email.NewResend(cfg.Resend.APIKey, cfg.Resend.FromEmail)

	// UseCases
	userUseCase := userusecase.New(
		userRepo,
		magicLinkRepo,
		refreshTokenRepo,
		jwtManager,
		emailSender,
		cfg.FrontendURL,
	)

	// Middleware
	authMW := authmiddleware.NewAuth(jwtManager, patRepo)

	// Handlers
	authHandler := v1.NewAuthHandler(userUseCase)

	server := &http.Server{
		Addr: ":" + cfg.Port,
		Handler: v1.NewRouter(v1.Deps{
			Auth:   authHandler,
			AuthMW: authMW,
		}),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info("server starting", zap.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server error", zap.Error(err))
		}
	}()

	<-quit
	log.Info("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("server forced to shutdown", zap.Error(err))
	}

	log.Info("server stopped")
}
