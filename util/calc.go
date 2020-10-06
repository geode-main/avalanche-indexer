package util

func PercentOf(value, total int64) float64 {
	return (float64(value) * 100.0) / float64(total)
}
