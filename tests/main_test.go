package tests

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

func TestMain(M *testing.M) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("No .env variables are loaded: %s", err.Error())
		return
	}
	code := M.Run()

	os.Exit(code)
}
