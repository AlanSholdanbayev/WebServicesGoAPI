package server

import (
	"context"
	"finalproject/internal/config"
	"finalproject/internal/handlers"
	"finalproject/internal/logger"
	"finalproject/internal/repository"
	"finalproject/internal/service"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Server struct {
	cfg     *config.Config
	log     *logger.LoggerWrapper
	httpSrv *http.Server
	db      *sqlx.DB
}

func New(cfg *config.Config, log *logger.LoggerWrapper) *Server {
	return &Server{
		cfg: cfg,
		log: log,
	}
}

func (s *Server) Run() error {
	s.log.Info().Msg("Connecting to database...")
	db, err := sqlx.Connect("postgres", s.cfg.DBUrl)
	if err != nil {
		s.log.Fatal().Err(err).Msg("DB connection failed")
		return err
	}
	s.db = db

	// Инициализация репозитория и сервисов
	userRepo := repository.NewPostgresUserRepo(db)
	userService := service.NewUserService(userRepo)

	// Передаём логгер в хэндлер
	userHandler := handlers.NewUserHandler(userService, s.cfg, s.log)

	// Регистрация маршрутов
	r := mux.NewRouter()
	userHandler.RegisterRoutes(r)

	s.httpSrv = &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%s", s.cfg.Port),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	s.log.Info().Msgf("Server started on port %s", s.cfg.Port)
	return s.httpSrv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpSrv != nil {
		_ = s.httpSrv.Shutdown(ctx)
	}
	if s.db != nil {
		_ = s.db.Close()
	}
	return nil
}
