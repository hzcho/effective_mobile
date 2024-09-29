package main

import (
	"flag"
	"fmt"
	"song_lib/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
)

func main() {
	var migrationPath, configPath string
	var up, down bool

	flag.StringVar(&migrationPath, "m", "", "path to the migration dir")
	flag.StringVar(&configPath, "c", "", "path to the config dir")
	flag.BoolVar(&up, "up", false, "apply migrations")
	flag.BoolVar(&down, "down", false, "rollback migrations")
	flag.Parse()

	log := logrus.New()

	if !up && !down {
		log.Fatal("You must specify either --up or --down")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.DBName, cfg.DB.SSLMode)

	m, err := migrate.New("file://"+migrationPath, dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	if up {
		if err := m.Up(); err != nil {
			if err != migrate.ErrNoChange {
				log.Fatal(err)
			}
		}
		fmt.Println("Migrations applied")
	} else if down {
		if err := m.Down(); err != nil {
			if err != migrate.ErrNoChange {
				log.Fatal(err)
			}
		}
		fmt.Println("Migrations rolled back")
	}
}
