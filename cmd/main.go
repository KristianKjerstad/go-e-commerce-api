package main

import (
	"database/sql"
	"log"

	"github.com/KristianKjerstad/go-e-commerce-api/cmd/api"
	"github.com/KristianKjerstad/go-e-commerce-api/cmd/config"
	"github.com/KristianKjerstad/go-e-commerce-api/cmd/db"
	"github.com/go-sql-driver/mysql"
)

func main() {

	db, err := db.NewMySQLStorage(mysql.Config{
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
	initStorage(db)

	server := api.NewAPIServer(":8080", nil)
	if err := server.Run(); err != nil {
		log.Fatal(err)

	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(("DB: connected"))
}
