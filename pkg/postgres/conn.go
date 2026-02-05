package postgres

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	DbName   string
}

func MustGetPostgresqlClient(cfg Config) *sqlx.DB {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DbName)

	conn := sqlx.MustConnect("postgres", dsn)

	err := conn.Ping()
	if err != nil {
		panic("couldn't ping postgres client")
	}

	return conn
}

func Close(db *sql.DB) error {
	return db.Close()
}
