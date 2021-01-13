package net

import (
	"fmt"
	"net"
	"testing"
)

func TestNetClient(t *testing.T) {
	conn, err := net.Dial("tcp", "192.168.11.60:445")
	if err != nil {
		fmt.Printf("error %s\n", err.Error())
	}
	defer conn.Close()
}
