package ipc

import (
	"encoding/binary"
	"io"
	"net"
	"time"
)

const (
	defaultTimeout = time.Second * 10

	TypeDesicions = "decisions"
	TypeConsensus = "consensus"
)

type Socket struct {
	conn net.Conn
}

type Message struct {
	Data []byte
	Time time.Time
}

type WriterFn func(Message) error

func DialTimeout(socketPath string, timeout time.Duration) (*Socket, error) {
	conn, err := net.DialTimeout("unix", socketPath, timeout)
	if err != nil {
		return nil, err
	}

	return &Socket{conn: conn}, nil
}

func Dial(socketPath string) (*Socket, error) {
	return DialTimeout(socketPath, defaultTimeout)
}

func (sock *Socket) Close() error {
	if sock.conn != nil {
		return sock.conn.Close()
	}
	return nil
}

func (sock *Socket) Start(writer WriterFn) error {
	var (
		sz uint64
	)

	for {
		if err := binary.Read(sock.conn, binary.BigEndian, &sz); err != nil {
			return err
		}

		data := make([]byte, sz)

		_, err := io.ReadFull(sock.conn, data)
		if err != nil {
			return err
		}

		writer(Message{
			Data: data,
			Time: time.Now().UTC(),
		})
	}
}
