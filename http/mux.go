package http

import (
	"fmt"
	"github.com/aaronland/go-artisanal-integers"
	gohttp "net/http"
	gourl "net/url"
	"strings"
)

func NewServeMux(s artisanalinteger.Service, u *gourl.URL) (*gohttp.ServeMux, error) {

	handler, err := IntegerHandler(s)

	if err != nil {
		return nil, err
	}

	path := u.Path

	if !strings.HasPrefix(path, "/") {
		path = fmt.Sprintf("/%s", path)
	}

	mux := gohttp.NewServeMux()
	mux.Handle(path, handler)

	return mux, nil
}
