package server

// EXPERIMENTAL

import (
        "fmt"
	"github.com/facebookgo/grace/gracehttp"	
	"github.com/thisisaaronland/go-artisanal-integers"
	"log"
	"net/http"
	"strconv"		
)

type HTTPServer struct {
	artisanalinteger.Server
	service artisanalinteger.Service
	host string
	port int
}

func NewHTTPServer (service artisanalinteger.Service) (*HTTPServer, error) {

     server := HTTPServer{
     	    service: service,
	    host: "localhost",
	    port: 8080
     }

     return &server, nil
}

func (s *HTTPServer) Listen() error {

	handler := func(rsp http.ResponseWriter, req *http.Request) {

		next, err := s.service.NextInt()

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		str_next := strconv.FormatInt(next, 10)

		b := []byte(str_next)

		rsp.Header().Set("Content-Type", "text/plain")
		rsp.Header().Set("Content-Length", strconv.Itoa(len(b)))

		rsp.Write(b)
	}

	endpoint := fmt.Sprintf("%s:%d", s.host, s.port)
	log.Println("listening on ", endpoint)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	err := gracehttp.Serve(&http.Server{Addr: endpoint, Handler: mux})
	return err
}
