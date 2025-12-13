package main

import (
	"fmt"
	"log"

	"finalproject/internal/config"

	"github.com/joho/godotenv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load() // Загружает переменные из .env
	cfg, _ := config.Load()

	db, err := sqlx.Connect("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
	}
	defer db.Close()

	var now string
	err = db.Get(&now, "SELECT now()")
	if err != nil {
		log.Fatalf("Query error: %v", err)
	}
	fmt.Println("DB connected successfully! Time:", now)
}
