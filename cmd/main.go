package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load env: %v",err)
	}

	if os.Getenv("MODE") == "dev" {
		log.Println("running in development mode")
		return
	}
	if os.Getenv("MODE") == "prod" {
		log.Println("running in production mode")
		return
	}
	log.Fatalf("MODE not set properly in env")
}

func main(){
	start()
}