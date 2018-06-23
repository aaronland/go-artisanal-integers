package client

// EXPERIMENTAL

import (
	"bufio"
	"github.com/aaronland/go-artisanal-integers"
	"net"
	"strconv"
	"strings"
)

type TCPClient struct {
	artisanalinteger.Client
	address string
}

func NewTCPClient(address string) (*TCPClient, error) {

	cl := TCPClient{
		address: address,
	}

	return &cl, nil
}

func (cl *TCPClient) NextInt() (int64, error) {

	conn, err := net.Dial("tcp", cl.address)

	if err != nil {
		return -1, err
	}

	str_i, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		return -1, err
	}

	str_i = strings.Trim(str_i, "\n")

	i, err := strconv.ParseInt(str_i, 10, 64)

	if err != nil {
		return -1, err
	}

	return i, err
}
