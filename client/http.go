package client

import (
	"context"
	"fmt"
	"io"
	_ "log"
	"net/http"
	"net/url"
	"strconv"
)

type HTTPClient struct {
	Client
	url         *url.URL
	http_client *http.Client
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
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	http_cl := &http.Client{}

	cl := &HTTPClient{
		url:         u,
		http_client: http_cl,
	}

	return cl, nil
}

func (cl *HTTPClient) NextInt(ctx context.Context) (int64, error) {

	req, err := http.NewRequestWithContext(ctx, "GET", cl.url.String(), nil)

	if err != nil {
		return -1, fmt.Errorf("Failed to create integer request, %w", err)
	}

	rsp, err := cl.http_client.Do(req)

	if err != nil {
		return -1, fmt.Errorf("Failed to GET next integer, %w", err)
	}

	defer rsp.Body.Close()

	byte_i, err := io.ReadAll(rsp.Body)

	if err != nil {
		return -1, fmt.Errorf("Failed to read response, %w", err)
	}

	str_i := string(byte_i)

	i, err := strconv.ParseInt(str_i, 10, 64)

	if err != nil {
		return -1, fmt.Errorf("Failed to parse integer response, %w", err)
	}

	return i, nil
}
