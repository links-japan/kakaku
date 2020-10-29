package main

import "fmt"

func main() {
	liquid := newClient()
	p, err := liquid.getTicker(5)
	if err != nil {
		panic(err)
	}
	fmt.Println(p)
}
