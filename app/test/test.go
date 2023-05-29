package main

import (
	"fmt"
	"net"
	"time"
)

// TODO: do not commit
func main() {
	host := "127.0.0.1"
	port := "8081"
	timeout := time.Duration(1 * time.Second)
	_, err := net.DialTimeout("tcp", host+":"+port, timeout)
	if err != nil {
		fmt.Printf("%s %s %s\n", host, "not responding", err.Error())
	} else {
		fmt.Printf("%s %s %s\n", host, "responding on port:", port)
	}
}
