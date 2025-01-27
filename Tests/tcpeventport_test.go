package tests

import (
	"EventStreamer/Streamer"
	"log"
	"net"
	"testing"
)

type TcpEventPortConfig = Streamer.TcpEventPortConfig

func TestTcpEventPort(t *testing.T) {
	log.Print("== instantiating event port ==")

	var tcpEventPort = Streamer.CreateTcpEventPort(TcpEventPortConfig{
		Port:     8082,
		HostName: "localhost",
	})
	go tcpEventPort.Listen()

	var conn, err = net.Dial("tcp", "localhost:8082")

	if err != nil {
		log.Fatalf("conneciton dial error : %s", err.Error())
		t.Fail()
	} else {
		_, err = conn.Write([]byte("hello\n"))

		if err != nil {
			log.Fatalf("connection send error : %s", err.Error())
			t.Fail()
		} else {
			for {
				select {
				case data := <-tcpEventPort.PipeForward():
					{
						if string(data) != "hello" {
							log.Printf("TestTcpEventPort FAIL : expected hello got %s", data)
							t.Fail()
						}
						log.Printf("got event : %s", string(data))
						conn.Close()
						return
					}
				default:
					{
					}
				}
			}
		}
	}
}
