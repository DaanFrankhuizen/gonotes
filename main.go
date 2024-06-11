package main

import (
	"log"

	"github.com/gonotes/api"
	"github.com/gonotes/storage"
)

func main() {
	store, err := storage.NewPostgresStore()

	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(":8080", store)
	server.Run()

}
