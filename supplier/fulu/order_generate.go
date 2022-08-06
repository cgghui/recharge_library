package fulu

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"time"
)

const EmptyOrderNo = ""

func CreateOrderID() string {
	return fmt.Sprintf("%s%d", time.Now().Format(TimeFormat2)[2:], RangeRand(10000000, 99999999))
}

func RangeRand(min, max int64) int64 {
	if min < 0 {
		f64Min := math.Abs(float64(min))
		i64Min := int64(f64Min)
		result, _ := rand.Int(rand.Reader, big.NewInt(max+1+i64Min))
		return result.Int64() - i64Min
	}
	result, _ := rand.Int(rand.Reader, big.NewInt(max-min+1))
	return min + result.Int64()
}
