package handler

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type handler struct {
	DB *sql.DB
}

func NewApp(initTable bool) *handler {
	h := connectDb()

	if initTable {
		h.createTableExpenses()
	}

	return h
}

func connectDb() *handler {
	url := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal().Msgf("Connect Database Error: %s", err.Error())
	}

	log.Info().Msg("Connect Database Success.")
	return &handler{db}
}

func (h *handler) createTableExpenses() {
	sqlCreateTbExpenses := `
	
	CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);`

	_, err := h.DB.Exec(sqlCreateTbExpenses)
	if err != nil {
		log.Fatal().Msgf("Create Table expenses Error: %s", err.Error())
	}

	log.Info().Msg("Create Table expenses Success.")
}
