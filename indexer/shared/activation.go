package shared

import "time"

var (
	Ap5ActivationTime uint64
)

func SetAp5ActivationTime(t time.Time) {
	Ap5ActivationTime = uint64(t.Unix())
}
