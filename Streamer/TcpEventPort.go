package Streamer

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
)

type TcpEventPortConfig struct {
	Port     int
	HostName string
}

type TcpEventPort struct {
	sync.Mutex
	config      TcpEventPortConfig
	connections map[string]net.Conn
	outChannel  chan []byte
}

func CreateTcpEventPort(config TcpEventPortConfig) TcpEventPort {
	return TcpEventPort{
		config:      config,
		connections: make(map[string]net.Conn),
		outChannel:  make(chan []byte),
	}
}

func (e *TcpEventPort) PipeForward() <-chan []byte {
	return e.outChannel
}

func (e *TcpEventPort) Listen() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", e.config.HostName, e.config.Port))
	log.Printf("TCP event port registered at %s", listener.Addr())

	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		log.Printf("connection : [%s] <-- [%s]", conn.LocalAddr().String(), conn.RemoteAddr().String())

		if err != nil {
			log.Fatalf("error : %s", err.Error())
		} else {
			e.connections[conn.RemoteAddr().String()] = conn

			go func() {
				connectionKey := conn.RemoteAddr().String()

				e.Lock()
				scanner := bufio.NewReader(e.connections[connectionKey])
				e.Unlock()

				for {
					eventBytes, _, err := scanner.ReadLine()
					if err != nil {
						if err.Error() == "EOF" {
							e.Lock()
							e.connections[connectionKey].Close()
							delete(e.connections, connectionKey)
							e.Unlock()

							log.Printf("connection %s closed ", connectionKey)
							return
						} else {
							panic(err)
						}
					}
					e.outChannel <- eventBytes
				}
			}()
		}
	}
}
