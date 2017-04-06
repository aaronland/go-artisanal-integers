package server

import (
	"bufio"
	"github.com/thisisaaronland/go-artisanal-integers"
	"log"
	"net"
	"strconv"
)

type TCPServer struct {
	artisanalinteger.Server
	address string
}

func NewTCPServer(address string) (*TCPServer, error) {

	server := TCPServer{
		address: address,
	}

	return &server, nil
}

func (s *TCPServer) ListenAndServe(service artisanalinteger.Service) error {

	listener, err := net.Listen("tcp", s.address)

	if err != nil {
		return err
	}

	for {

		conn, err := listener.Accept()

		if err != nil {
			log.Println(err)
			continue
		}

		go func() {

			defer conn.Close()

			i, err := service.NextInt()

			if err != nil {
				return
			}

			str_i := strconv.FormatInt(i, 10)

			bufout := bufio.NewWriter(conn)
			bufout.WriteString(str_i)
			bufout.Flush()
		}()
	}
}
