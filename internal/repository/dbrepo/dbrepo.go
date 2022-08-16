package dbrepo

import (
	"database/sql"
	"github.com/margleb/booking/internal/config"
	"github.com/margleb/booking/internal/repository"
)

type postgresDbRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

// NewPostgresRep создает новое postgres подключение в репозитории
func NewPostgresRep(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDbRepo{
		App: a,
		DB:  conn,
	}
}
