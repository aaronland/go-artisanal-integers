package client

import (
	"github.com/thisisaaronland/go-artisanal-integers"
	"io/ioutil"
	"net/http"
	"strconv"
)

type HTTPClient struct {
	artisanalinteger.Client
	address string
}

func NewHTTPClient(address string) (*HTTPClient, error) {

	cl := HTTPClient{
		address: address,
	}

	return &cl, nil
}

func (cl *HTTPClient) NextInt() (int64, error) {

	rsp, err := http.Get("http://" + cl.address)

	if err != nil {
		return -1, err
	}

	defer rsp.Body.Close()

	byte_i, err := ioutil.ReadAll(rsp.Body)

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
