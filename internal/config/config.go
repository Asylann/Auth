package config

import (
	"fmt"
	"os"
)

type Config struct {
	DataConnection string
	Port           string
	JWTSecret      string
}

func LoadConfig() Config {
	databaseStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("PG_USER"), os.Getenv("PG_PASSWORD"), os.Getenv("PG_HOST"), os.Getenv("PG_PORT"), os.Getenv("DATABASE"))
	return Config{DataConnection: databaseStr, Port: os.Getenv("PORT"), JWTSecret: os.Getenv("JWT_SECRET")}
}
