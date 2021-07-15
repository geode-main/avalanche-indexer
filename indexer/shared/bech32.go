package shared

import (
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/formatting"
)

var (
	HRP   string
	Asset string
)

func SetBech32HRP(network uint32) {
	HRP = constants.GetHRP(network)
}

func SetAvaxAssetID(id string) {
	Asset = id
}

func bech32addrs(addrs [][]byte) ([]string, error) {
	result := []string{}

	for _, addr := range addrs {
		addr, err := bech32fromBytes(addr)
		if err != nil {
			return nil, err
		}
		result = append(result, addr)
	}

	return result, nil
}

func bech32address(addr ids.ShortID) (string, error) {
	return formatting.FormatBech32(HRP, addr.Bytes())
}

func bech32fromBytes(addr []byte) (string, error) {
	return formatting.FormatBech32(HRP, addr)
}
