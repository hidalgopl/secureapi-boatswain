package chefscheduler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/hidalgopl/secureapi-boatswain/internal/chef"
	"github.com/hidalgopl/secureapi-boatswain/internal/statusendpoints"
)

func Main() {
	r := mux.NewRouter()

	chefServer, err := chef.NewServer()
	if err != nil {
		panic(err)
	}
	chefServer.AttachRoutes(r)

	statusServer := statusendpoints.NewServer()
	statusServer.AttachRoutes(r)
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", 8080),
		Handler:      r,
		ReadTimeout:  time.Duration(60) * time.Second,
		WriteTimeout: time.Duration(60) * time.Second,
	}

	err = httpServer.ListenAndServe()
	log.Println(err)
}
