package main

import (
	"fmt"
	"math/big"
	"time"

	"github.com/yudppp/isutools/profile"
	"github.com/yudppp/isutools/utils/throttle"
)

func main() {
	profile.StartCPU(time.Second*5, true)
	loopGcd()
}

var printThrottle = throttle.New(time.Second * 1)

func loopGcd() {
	size := int64(4000)
	for a := int64(0); a < size; a++ {
		for b := a; b < size; b++ {
			// logging per time
			printThrottle.Do(func() {
				fmt.Println(a, b)
			})
			gcd(big.NewInt(a), big.NewInt(b))
		}
	}
}

func gcd(a, b *big.Int) *big.Int {
	if a.Cmp(b) < 0 {
		return gcd(b, a)
	}
	for b.Cmp(big.NewInt(0)) != 0 {
		r := new(big.Int).Mod(a, b)
		a = b
		b = r
	}
	return a
}
