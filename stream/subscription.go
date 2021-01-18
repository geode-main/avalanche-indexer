package stream

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

type Subscription struct {
	socketPath string
	messages   chan []byte
}

func NewSubscription(network int, chain string) (*Subscription, error) {
	return &Subscription{
		messages:   make(chan []byte, 1024),
		socketPath: fmt.Sprintf("/tmp/%v-%v-decisions", network, chain),
	}, nil
}

func (s *Subscription) Start() error {
	conn, err := net.Dial("unix", s.socketPath)
	if err != nil {
		return err
	}
	defer conn.Close()

	var sz uint64

	for {
		if err := binary.Read(conn, binary.BigEndian, &sz); err != nil {
			log.Println("read error:", err)
			continue
		}

		msg := make([]byte, sz)
		if _, err := io.ReadFull(conn, msg); err != nil {
			log.Println("read error:", err)
			continue
		}

		s.messages <- msg
	}
}

func (s *Subscription) Messages() <-chan []byte {
	return s.messages
}
