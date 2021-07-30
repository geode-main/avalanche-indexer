package util

import (
	"bytes"
	"strings"
)

func TxMemo(data []byte) string {
	data = bytes.ReplaceAll(data, []byte("\x00"), []byte(""))
	return strings.ToValidUTF8(string(data), "")
}

func StringPtr(val string) *string {
	return &val
}

func Uint64Prt(val uint64) *uint64 {
	return &val
}

func BoolPtr(val bool) *bool {
	return &val
}
