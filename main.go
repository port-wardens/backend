package main

import (
	"fmt"
	"os"

	"github.com/jayendramadaram/port-wardens/auth"
	"github.com/jayendramadaram/port-wardens/model"
	"github.com/jayendramadaram/port-wardens/server"
	"github.com/jayendramadaram/port-wardens/store"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	Secret := os.Getenv("SECRET")
	if Secret == "" {
		panic("SECRET environment variable empty")
	}

	dns := os.Getenv("DATABASE_URL")
	if dns == "" {
		panic("DATABASE_URL environment variable empty")
	}
	db, err := model.NewDB(dns)
	if err != nil {
		panic(err)
	}

	str := store.NewStore(db)
	authModule := auth.NewAuth(db)
	LetsFuckingGoo := server.NewServer(str, authModule, Secret)

	LetsFuckingGoo.Run()
}
