package Streamer

type EventPort interface {
	// listen for events
	listen()
	PipeForward() <-chan []byte
}
