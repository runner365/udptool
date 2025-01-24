package main

import (
	"flag"
	"fmt"
	"net"
)

var (
	serverIp   = flag.String("s", "", "server ip address")
	serverPort = flag.String("p", "", "server port")
)

const UDP_BUF_SIZE = 1500

func init() {
	flag.Parse()
}

func onWork(conn *net.UDPConn) {
	buf := make([]byte, UDP_BUF_SIZE)
	for {

		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Printf("Received data from %s\n", addr.String())
		_, err = conn.WriteToUDP(buf[:n], addr)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
	}
}

func main() {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", *serverIp, *serverPort))
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer conn.Close()

	onWork(conn)
}
