package server

import (
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
	"game/controllers"
	"game/migration"
)

var server = controllers.Server{}

func Run() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}
	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
		migration.Load(server.DB)
		migration.Seed(server.DB)

	

	server.Run(":8082")

}