package util

import (
	"crypto/sha1"
	"fmt"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/hashing"

	"github.com/tuvistavie/securerandom"
)

func UUID() string {
	str, err := securerandom.Uuid()
	if err != nil {
		panic(err)
	}
	return str
}

func Hex(n int) string {
	str, err := securerandom.Hex(n)
	if err != nil {
		panic(err)
	}
	return str
}

func StringSHA1(input string) string {
	h := sha1.New()
	h.Write([]byte(input))

	return fmt.Sprintf("%x", h.Sum(nil))
}

func AvalancheHash256(data []byte) []byte {
	return hashing.ComputeHash256(data)
}

func AvalancheID(data []byte) (string, error) {
	id, err := ids.ToID(AvalancheHash256(data))
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

func AvalancheIDFromString(data string) (string, error) {
	return AvalancheID([]byte(data))
}
