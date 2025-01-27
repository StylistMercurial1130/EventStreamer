package main

import (
	"EventStreamer/Streamer"
	"log"
)

type TcpEventPortConfig = Streamer.TcpEventPortConfig

func main() {
	eventPort := Streamer.CreateTcpEventPort(TcpEventPortConfig{
		Port:     8082,
		HostName: "localhost",
	})

	go eventPort.Listen()

	for {
		select {
		case eventBytes := <-eventPort.PipeForward():
			{
				msg := string(eventBytes)
				log.Printf("message received => %s", msg)
			}
		default:
			{
			}
		}
	}

}
