package main

import (
	"context"
	"log"
	"net"
	"syscall"
)

type TCPServer struct {
	ctx     context.Context
	address string

	stop chan bool
}

func NewTCPServer(ctx context.Context, address string) *TCPServer {
	return &TCPServer{
		ctx:     ctx,
		address: address,
	}
}

func (t *TCPServer) Start() error {
	t.stop = make(chan bool)

	listener, err := (&net.ListenConfig{}).Listen(t.ctx, "tcp4", t.address)
	if err != nil {
		return err
	}

	for {
		select {
		case <-t.stop:
			return listener.Close()
		default:
			connection, err := listener.Accept()
			if err != nil && err != syscall.EINVAL {
				log.Fatalf("Failed to accept the TCP Connection. Error: %s", err.Error())
				return err
			}
			defer connection.Close()

			go func() {
				connection.Write([]byte(string("ok")))
				connection.Close()
			}()
		}
	}
}

func (t *TCPServer) Stop() {
	close(t.stop)
}
