package round

func TwoDecimals(num float64) float64 {
	return float64(int(num*100)) / 100
}
