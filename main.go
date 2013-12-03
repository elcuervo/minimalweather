package main

import (
	"github.com/elcuervo/minimalweather/minimalweather"
	"log"
	"net/http"
	"os"
)

func main() {
	handler := minimalweather.Handler()
	port := ":" + os.Getenv("PORT")
	http.Handle("/", handler)

	log.Println("Listening in", port)
	err := http.ListenAndServe(port, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
