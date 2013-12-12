package main

import (
	"github.com/elcuervo/minimalweather/minimalweather"
	"github.com/yvasiyarov/gorelic"

	"log"
	"net/http"
	"os"
)

func main() {
	agent := gorelic.NewAgent()
	agent.Verbose = true
	agent.NewrelicLicense = os.Getenv("NEW_RELIC_LICENSE_KEY")
	agent.Run()

	handler := minimalweather.Handler()
	port := ":" + os.Getenv("PORT")
	http.Handle("/", handler)

	log.Println("Listening in", port)
	err := http.ListenAndServe(port, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
