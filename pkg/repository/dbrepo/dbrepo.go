package dbrepo

import (
	"database/sql"

	"github.com/jinyanomura/bookings/pkg/config"
	"github.com/jinyanomura/bookings/pkg/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB *sql.DB
}

func NewPostgresDBRepo(a *config.AppConfig, conn *sql.DB) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB: conn,
	}
}