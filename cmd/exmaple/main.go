package main

import (
	"fmt"
	"github.com/links-japan/kakaku"
)

func main() {
	price, err := kakaku.BTCToJPY()
	if err != nil {
		panic(err)
	}

	fmt.Println(price)
}
