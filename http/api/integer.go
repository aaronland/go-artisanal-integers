package api

import (
	"github.com/aaronland/go-artisanal-integers/service"
	"net/http"
	"strconv"
)

func IntegerHandler(s service.Service) (http.HandlerFunc, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		next, err := s.NextInt(ctx)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		str_next := strconv.FormatInt(next, 10)
		b := []byte(str_next)

		rsp.Header().Set("Content-Type", "text/plain")
		rsp.Header().Set("Content-Length", strconv.Itoa(len(b)))
		rsp.Header().Set("Access-Control-Allow-Origin", "*")

		rsp.Write(b)
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
