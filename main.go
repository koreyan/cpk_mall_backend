package main

import "cpk_mall/network"

func main() {
	r := network.MakeRouter()
	r.Run(":8080")
}
