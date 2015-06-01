package main

import (
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/meatballhat/negroni-logrus"
)

var Version string
var SG_HOST string

func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.Handle("/favicon.ico", http.NotFoundHandler())

	auth_middleware := negroni.HandlerFunc(ShotgunAuthMiddleware)

	entityRoutes := mux.NewRouter()
	entityRoutes.Path("/{entity_type}/{id:[0-9]+}").HandlerFunc(entityGetHandler).Methods("GET")
	entityRoutes.Path("/{entity_type}/{id:[0-9]+}").HandlerFunc(entityUpdateHandler).Methods("PATCH")
	entityRoutes.Path("/{entity_type}/{id:[0-9]+}").HandlerFunc(entityDeleteHandler).Methods("DELETE")
	entityRoutes.Path("/{entity_type}").HandlerFunc(entityGetAllHandler).Methods("GET")
	entityRoutes.Path("/{entity_type}").HandlerFunc(entityCreateHandler).Methods("POST")

	// Adds auth on the sub router so that / can be accessed freely.
	r.PathPrefix("/{entity_type}").Handler(negroni.New(
		auth_middleware,
		negroni.Wrap(entityRoutes),
	))

	return r
}

func main() {

	app := cli.NewApp()
	app.Name = "sg-restful"
	app.Version = Version

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "port, p",
			Value:  "8000",
			Usage:  "Port to listen on",
			EnvVar: "PORT",
		},
		cli.StringFlag{
			Name:   "shotgun-host, s",
			Value:  "",
			Usage:  "Shotgun host",
			EnvVar: "SG_HOST",
		},
	}

	app.Action = func(c *cli.Context) {
		log.Infof("sg-restful Version: %v", Version)
		log.Infof("Shotgun Host: %v", c.String("shotgun-host"))
		SG_HOST = c.String("shotgun-host")
		r := Router()

		n := negroni.Classic()
		n.Use(negronilogrus.NewMiddleware())
		n.UseHandler(r)
		n.Run(":" + c.String("port"))
	}
	app.Run(os.Args)
}
