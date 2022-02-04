package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	userhandler "github.com/frycm/user-service/cmd/user-service/server/handler/user"
	userdb "github.com/frycm/user-service/internal/db/user"
	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

type Conf struct {
	Addr string `conf:"default::8080"`

	ReadTimeout     time.Duration `conf:"default:5s"`
	WriteTimeout    time.Duration `conf:"default:5s"`
	IdleTimeout     time.Duration `conf:"default:120s"`
	ShutdownTimeout time.Duration `conf:"default:10s"`
}

func Serve(ctx context.Context, cfg Conf) error {
	userDB := userdb.Client{}

	userHandler := userhandler.NewHandler(&userDB)

	r := mux.NewRouter()
	r.HandleFunc("/api/users", userHandler.Create) //.Methods(http.MethodPost)

	srv := http.Server{
		Addr:         cfg.Addr,
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		ErrorLog:     zap.NewStdLog(zap.L()),
	}

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- srv.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if err != http.ErrServerClosed {
			return fmt.Errorf("server error: %w", err)
		}
	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("gracefull server shutdown failed: %w", err)
		}
	}

	return nil
}
