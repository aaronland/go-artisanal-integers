package application

import (
	"flag"
	"fmt"
	"github.com/aaronland/go-artisanal-integers"
	"github.com/aaronland/go-artisanal-integers/server"
	"github.com/aaronland/go-artisanal-integers/service"
	"log"
	"net/url"
	"os"
	"strings"
)

func ServerApplication(eng artisanalinteger.Engine) error {

	var proto = flag.String("protocol", "http", "The protocol for the server to implement. Valid options are: http,tcp.")
	var host = flag.String("host", "localhost", "The hostname to listen for requests on")
	var port = flag.Int("port", 8080, "The port number to listen for requests on")
	var path = flag.String("path", "/", "The path to listen for requests on")

	var last = flag.Int("set-last-int", 0, "Set the last known integer.")
	var offset = flag.Int("set-offset", 0, "Set the offset used to mint integers.")
	var increment = flag.Int("set-increment", 0, "Set the increment used to mint integers.")

	flag.Parse()

	flag.VisitAll(func(fl *flag.Flag) {

		name := fl.Name
		env := name

		env = strings.ToUpper(env)
		env = strings.Replace(env, "-", "_", -1)
		env = fmt.Sprintf("INTEGER_%s", env)

		val, ok := os.LookupEnv(env)

		if ok {
			log.Printf("set -%s flag (%s) from %s environment variable\n", name, val, env)
			flag.Set(name, val)
		}
	})

	if *last != 0 {

		err := eng.SetLastInt(int64(*last))

		if err != nil {
			return err
		}
	}

	if *increment != 0 {

		err := eng.SetIncrement(int64(*increment))

		if err != nil {
			return err
		}
	}

	if *offset != 0 {

		err := eng.SetOffset(int64(*offset))

		if err != nil {
			return err
		}
	}

	svc, err := service.NewArtisanalService("simple", eng)

	if err != nil {
		return err
	}

	u := new(url.URL)

	u.Scheme = *proto
	u.Host = fmt.Sprintf("%s:%d", *host, *port)
	u.Path = *path

	_, err = url.Parse(u.String())

	if err != nil {
		return err
	}

	s, err := server.NewArtisanalServer(*proto, u)

	if err != nil {
		return err
	}

	log.Println("Listen on", s.Address())

	return s.ListenAndServe(svc)
}
