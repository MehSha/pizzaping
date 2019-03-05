package main

import (
	"fmt"
	"time"

	"github.com/caarlos0/env"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type pgConfig struct {
	Host string `env:"POSTGRES_HOST" envDefault:"localhost"`
	DB   string `env:"POSTGRES_DB" envDefault:"pizza"`
	User string `env:"POSTGRES_USER" envDefault:"postgres"`
	Pass string `env:"POSTGRES_PASSWORD" envDefault:"postgres"`
}

func connectDB() *sqlx.DB {
	pgcfg := pgConfig{}
	err := env.Parse(&pgcfg)
	if err != nil {
		panic(err)
	}
	for {

		db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable", pgcfg.Host, pgcfg.DB, pgcfg.User, pgcfg.Pass))
		if err == nil {
			fmt.Println("DB connected successfully!")
			return db
		} else {
			fmt.Printf("connecting to DB was not successful: %s, \ntrying again\n", err)
			time.Sleep(time.Second * 3)
		}
	}
}
