package main

import (
	"net/http"
	"os"

	"github.com/sheelestun/WatchHRs-/internal/repository"
	"github.com/sheelestun/WatchHRs-/internal/service"
	"github.com/sheelestun/WatchHRs-/internal/web/handler"
	"github.com/sheelestun/WatchHRs-/internal/web/router"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	repo := repository.NewInMemoryStorage()
	serv := service.NewAPIService(repo)
	hand := handler.NewApiHandler(serv, []byte("a-string-secret-at-least-256-bits-long"))
	r := router.NewRouter(hand)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
		return
	}
}
