package client

import (
	"fmt"
)

type IndexClient struct {
	rpc
}

var containerTypes = map[string]string{
	"X": "tx",
	"C": "block",
	"P": "block",
}

func (c IndexClient) GetLastAccepted(chain string) (*Container, error) {
	url := fmt.Sprintf("%s/%s/%s", c.endpoint, chain, containerTypes[chain])

	data, err := c.callRaw(url, "index.getLastAccepted", map[string]string{"encoding": "cb58"})
	if err != nil {
		return nil, err
	}

	container := &Container{}
	return container, c.decode(data, container)
}

func (c IndexClient) GetContainerRange(chain string, startIndex int, numToFetch int) (*ContainersResponse, error) {
	url := fmt.Sprintf("%s/%s/%s", c.endpoint, chain, containerTypes[chain])
	args := map[string]interface{}{
		"encoding":   "hex",
		"startIndex": startIndex,
		"numToFetch": numToFetch,
	}

	data, err := c.callRaw(url, "index.getContainerRange", args)
	if err != nil {
		return nil, err
	}

	resp := &ContainersResponse{}
	return resp, c.decode(data, resp)
}
