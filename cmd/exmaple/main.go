package main

import (
	"fmt"
	"github.com/links-japan/kakaku"
)

func main() {
	price, _, err := kakaku.PriceWithTime(kakaku.BTC, kakaku.JPYC)
	if err != nil {
		panic(err)
	}

	fmt.Println(price)
}
