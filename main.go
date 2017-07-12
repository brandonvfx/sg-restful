package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/meatballhat/negroni-logrus"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/urfave/negroni"
)

var Version string

func router(config clientConfig) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler(config))
	r.Handle("/favicon.ico", http.NotFoundHandler())

	authMiddleware := negroni.HandlerFunc(ShotgunAuthMiddleware(config))

	entityRoutes := mux.NewRouter()
	entityRoutes.Path("/{entity_type}/{id:[0-9]+}").HandlerFunc(entityGetHandler(config)).Methods("GET")
	entityRoutes.Path("/{entity_type}/{id:[0-9]+}").HandlerFunc(entityUpdateHandler(config)).Methods("PATCH")
	entityRoutes.Path("/{entity_type}/{id:[0-9]+}").
		HandlerFunc(entityDeleteHandler(config)).Methods("DELETE")
	entityRoutes.Path("/{entity_type}/{id:[0-9]+}/revive").
		HandlerFunc(entityReviveHandler(config)).Methods("POST")
	// entityRoutes.Path("/{entity_type}/{id:[0-9]+}/followers").
	// 	HandlerFunc(entityGetFollowersHandler(config)).Methods("GET")
	//entityRoutes.Path("/{entity_type}/{id:[0-9]+}/followers").
	//             HandlerFunc(entityAddFollowersHandler(config)).Methods("POST")
	//entityRoutes.Path("/{entity_type}/{id:[0-9]+}/followers/{user_type}/{user_id:[0-9]+}").
	//		       HandlerFunc(entityDeleteFollowersHandler(config)).Methods("DELETE")
	entityRoutes.Path("/{entity_type}").HandlerFunc(entityGetAllHandler(config)).Methods("GET")
	entityRoutes.Path("/{entity_type}").HandlerFunc(entityCreateHandler(config)).Methods("POST")

	// Adds auth on the sub router so that / can be accessed freely.
	r.PathPrefix("/{entity_type}").Handler(negroni.New(
		authMiddleware,
		negroni.Wrap(entityRoutes),
	))

	return r
}

func main() {
	f, err := os.OpenFile("sg-restful.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	// don't forget to close it
	defer f.Close()

	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{})

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
		if c.String("shotgun-host") == "" {
			log.Fatalln("Shotgun host not set.")
		}
		log.Infof("Shotgun Host: %v", c.String("shotgun-host"))
		config := newClientConfig(Version, c.String("shotgun-host"))

		r := router(config)
		corsMiddleware := cors.AllowAll()

		n := negroni.Classic()
		n.Use(negronilogrus.NewMiddleware())
		n.Use(corsMiddleware)
		n.UseHandler(r)
		n.Run(":" + c.String("port"))
	}
	app.Run(os.Args)
}
