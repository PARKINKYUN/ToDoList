package main

import (
	"justin/todos/app"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	mux := app.MakeHandler(os.Getenv("DATABASE_URL"))
	defer mux.Close()

	log.Println("App started")
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		panic(err)
	}
}
