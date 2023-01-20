package main

import (
	"log"
	"net/http"
	"todolist/app"
)

func main() {
	mux := app.MakeHandler("./test.db") // flag.args 로 변환 가능함
	defer mux.Close()

	log.Println("App started")
	err := http.ListenAndServe(":3000", mux)
	if err != nil {
		panic(err)
	}
}
