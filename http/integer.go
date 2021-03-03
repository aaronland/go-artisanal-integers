package http

import (
	"github.com/aaronland/go-artisanal-integers/service"
	gohttp "net/http"
	"strconv"
)

func IntegerHandler(svc service.Service) (gohttp.HandlerFunc, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		next, err := svc.NextInt()

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		str_next := strconv.FormatInt(next, 10)
		b := []byte(str_next)

		rsp.Header().Set("Content-Type", "text/plain")
		rsp.Header().Set("Content-Length", strconv.Itoa(len(b)))
		rsp.Header().Set("Access-Control-Allow-Origin", "*")

		rsp.Write(b)
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
