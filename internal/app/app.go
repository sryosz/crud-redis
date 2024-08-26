package app

import (
	"log/slog"
	httpapp "microservice-redis/internal/app/http"
)

type App struct {
	HttpApp *httpapp.App
}

func New(log *slog.Logger) *App {
	httpApp := httpapp.NewApp(log)
	return &App{
		HttpApp: httpApp,
	}
}
