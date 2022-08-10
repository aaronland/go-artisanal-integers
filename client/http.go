package client

import (
	"context"
	"io"
	_ "log"
	"net/http"
	"net/url"
	"strconv"
)

type HTTPClient struct {
	Client
	url *url.URL
}

func init() {

	ctx := context.Background()

	schemes := []string{
		"http",
		"https",
	}

	for _, prefix := range schemes {

		err := RegisterClient(ctx, prefix, NewHTTPClient)

		if err != nil {
			panic(err)
		}
	}
}

func NewHTTPClient(ctx context.Context, uri string) (Client, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	cl := &HTTPClient{
		url: u,
	}

	return cl, nil
}

func (cl *HTTPClient) NextInt(ctx context.Context) (int64, error) {

	rsp, err := http.Get(cl.url.String())

	if err != nil {
		return -1, err
	}

	defer rsp.Body.Close()

	byte_i, err := io.ReadAll(rsp.Body)

	if err != nil {
		return -1, err
	}

	str_i := string(byte_i)

	i, err := strconv.ParseInt(str_i, 10, 64)

	if err != nil {
		return -1, err
	}

	return i, err
}
