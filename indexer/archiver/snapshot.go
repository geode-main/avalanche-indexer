package archiver

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"time"
)

type Snapshot struct {
	ID   string                 `json:"id"`
	Data map[string]interface{} `json:"data"`
	Meta Metadata               `json:"meta"`
}

type Metadata struct {
	AppName      string    `json:"app_name"`
	AppVersion   string    `json:"app_version"`
	ChainName    string    `json:"chain_name"`
	ChainNetwork string    `json:"chain_network"`
	ChainVersion string    `json:"chain_version"`
	Height       *int64    `json:"height"`
	Time         time.Time `json:"time"`
}

func NewSnapshot(id string) (*Snapshot, error) {
	if id == "" {
		return nil, errors.New("snapshot id is not provided")
	}

	return &Snapshot{
		ID:   id,
		Data: map[string]interface{}{},
		Meta: Metadata{},
	}, nil
}

func (s *Snapshot) Add(key string, data interface{}) {
	s.Data[key] = data
}

func (s *Snapshot) Encode(io io.Writer) error {
	writer, err := gzip.NewWriterLevel(io, gzip.BestCompression)
	if err != nil {
		return err
	}
	defer writer.Close()

	return json.NewEncoder(writer).Encode(s)
}
