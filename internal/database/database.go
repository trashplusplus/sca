package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {

	godotenv.Load("../config/.env")

	dbuser := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	dbpass := os.Getenv("DB_PASS")

	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", dbuser, dbname, dbpass)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	sqlBytes, err := os.ReadFile("../init.sql")
	if err != nil {
		return nil, fmt.Errorf("не вдалося прочитати init.sql: %w", err)
	}

	sqlStatements := strings.Split(string(sqlBytes), ";")
	for _, stmt := range sqlStatements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		_, err := db.Exec(stmt)
		if err != nil {
			return nil, fmt.Errorf("помилка виконання запиту: %s, %w", stmt, err)
		}
	}

	fmt.Println("init.sql успішно виконано")

	log.Println("Database connected")
	return db, err

}
