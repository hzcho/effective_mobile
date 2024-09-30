package app

import (
	"context"
	"fmt"
	"song_lib/internal/app/server"
	"song_lib/internal/config"
	"song_lib/internal/group"
	"song_lib/internal/handler"
	"song_lib/internal/repository"
	"song_lib/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
)

type App struct {
	server *server.Server
	pool   *pgxpool.Pool
}

func NewApp(ctx context.Context, cfg *config.Config, log *logrus.Logger) *App {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.DBName)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		panic(err)
	}

	repos := repository.NewRepositories(pool, log)
	usecases := usecase.NewUsecases(repos, log)
	groups := group.NewGroups(usecases, log)
	router := gin.New()
	router.Use(ginlogrus.Logger(log))
	handler.InitRoutes(router, *groups)
	server := server.NewServer(&cfg.Server, router)

	return &App{
		server: server,
		pool:   pool,
	}
}

func (a *App) Start() {
	a.server.Run()
}

func (a *App) Stop(ctx context.Context) {
	a.server.Stop(ctx)
	a.pool.Close()
}
