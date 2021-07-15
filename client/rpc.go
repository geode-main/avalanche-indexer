package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/rpc/v2/json2"
	"github.com/sirupsen/logrus"
)

const (
	prefixInfo     = "/ext/info"
	prefixPlatform = "/ext/P"
	prefixAvm      = "/ext/bc/X"
	prefixEvm      = "/ext/bc/C/rpc"
	prefixIpc      = "/ext/ipcs"
	prefixIndex    = "/ext/index"
)

type rpc struct {
	endpoint string
	client   *http.Client
	logger   *logrus.Logger
}

func initRPC(endpoint string, prefix string) rpc {
	return rpc{
		endpoint: fmt.Sprintf("%s%s", endpoint, prefix),
		client: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

func (c rpc) callRaw(url string, method string, args interface{}) ([]byte, error) {
	data, err := json2.EncodeClientRequest(method, args)
	if err != nil {
		return nil, err
	}
	reqBody := bytes.NewReader(data)

	req, err := http.NewRequest(http.MethodPost, url, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	ts := time.Now()

	resp, err := c.client.Do(req)
	defer c.logRequest(method, args, resp, err, time.Since(ts))

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func (c rpc) call(method string, args interface{}, out interface{}) error {
	data, err := c.callRaw(c.endpoint, method, args)
	if err != nil {
		return err
	}
	return c.decode(data, out)
}

func (c rpc) decode(data []byte, out interface{}) error {
	return json2.DecodeClientResponse(bytes.NewReader(data), out)
}

func (c rpc) logRequest(method string, args interface{}, resp *http.Response, err error, duration time.Duration) {
	entry := logrus.
		WithField("method", method).
		WithField("duration", duration.Milliseconds())

	if args != nil {
		entry = entry.WithField("args", args)
	}

	if resp != nil {
		entry = entry.WithField("status", resp.StatusCode)
	}

	if err == nil {
		entry.Debug("rpc call")
	} else {
		entry.WithError(err).Error("rpc call failed")
	}
}
