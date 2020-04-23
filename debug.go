package main

/*
debug.go is used only for runtime exploring and should be ignored
*/

import (
	"fmt"

	"github.com/mikelsr/nahs/net"
)

func main() {
	n := net.NewNode(nil)
	fmt.Println(n)
}
