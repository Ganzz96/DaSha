package main

import (
	"fmt"
	"net"
)

func main() {
	socket, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 6666,
	})
	if err != nil {
		panic(err)
	}

	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:7778")
	if err != nil {
		panic(err)
	}

	sent, err := socket.WriteToUDP([]byte("some_agent_id"), serverAddr)
	if err != nil {
		panic(err)
	}
	fmt.Println("LAL N bytes sent:", sent, serverAddr.String())

	buffer := make([]byte, 1024)

	for {
		read, _, err := socket.ReadFromUDP(buffer)
		if err != nil {
			panic(err)
		}
		fmt.Println("LAL N bytes read:", read, string(buffer))
	}
}
