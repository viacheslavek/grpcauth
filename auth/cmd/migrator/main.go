package main

import (
	"database/sql"
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {

	// TODO: очередной техдолг: надо брать путь до базы, пароль и тд как в парсинге конфига, но из
	// флага или переменных окружения + добавить для этого валидацию, а пока так работает, пойду делать
	// ручки бд

	db, err := sql.Open("postgres",
		"postgres://slava:password_from_env@localhost:5432/test_db?sslmode=disable")
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
