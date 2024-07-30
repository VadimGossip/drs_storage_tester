package util

import (
	"crypto/rand"
	"math"
	"math/big"
)

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func RandInt(maxInt int) (int, error) {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(maxInt)))
	if err != nil {
		return 0, err
	}

	return int(nBig.Int64()), nil
}
