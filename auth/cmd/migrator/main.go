package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/viacheslavek/grpcauth/auth/internal/config"
	"log"
	"log/slog"
	"net/url"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	dbCfg := config.MustLoad().DB

	postgresURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(dbCfg.User, dbCfg.Password),
		Host:   fmt.Sprintf("%s:%d", dbCfg.Host, dbCfg.Port),
		Path:   dbCfg.DBName,
	}

	log.Println("current postgres url", slog.String("url", postgresURL.String()))

	db, err := sql.Open("postgres",
		postgresURL.String()+"?sslmode=disable")
	if err != nil {
		log.Fatalf("could not connect to the database: %v", err)
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			log.Println("failed to close db", err)
		}
	}(db)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("could not create migrate instance: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"test_db", driver)
	if err != nil {
		log.Fatalf("could not start migrate: %v", err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("could not run up migrations: %v", err)
	}
	log.Println("migrations ran successfully")
}
