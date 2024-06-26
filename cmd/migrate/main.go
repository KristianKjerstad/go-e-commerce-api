package main

import (
	"fmt"
	"log"
	"os"

	"github.com/KristianKjerstad/go-e-commerce-api/config"
	"github.com/KristianKjerstad/go-e-commerce-api/db"
	mysqlCfg "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	db, err := db.NewMySQLStorage(mysqlCfg.Config{
		User:                 config.EnvConfig.DBUser,
		Passwd:               config.EnvConfig.DBPassword,
		Addr:                 config.EnvConfig.DBAddress,
		DBName:               config.EnvConfig.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})

	if err != nil {
		log.Fatal(err)
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"mysql",
		driver,
	)
	fmt.Println(m)
	if err != nil {
		log.Fatal(err)
	}

	cmd := os.Args[(len(os.Args) - 1)]
	fmt.Println(cmd)
	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
	if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}

}
