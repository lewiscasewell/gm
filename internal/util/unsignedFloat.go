package util

func UnsignFloat64(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}
