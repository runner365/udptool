package main

import (
	"flag"
	"fmt"
	"net"
	"time"
)

var (
	serverIp   = flag.String("s", "", "server ip address")
	serverPort = flag.Int("p", 0, "server port")
	clientIp   = flag.String("c", "", "client ip address")
	clientPort = flag.Int("d", 0, "client port")
	kbps       = flag.Int("k", 10, "kbits per second")
)

func init() {
	flag.Parse()
}

func updateBuffer(buffer []byte, seq uint32) {
	buffer[0] = byte(seq >> 24)
	buffer[1] = byte(seq >> 16)
	buffer[2] = byte(seq >> 8)
	buffer[3] = byte(seq)

	ts := uint32(time.Now().UnixMilli())

	buffer[4] = byte(ts >> 24)
	buffer[5] = byte(ts >> 16)
	buffer[6] = byte(ts >> 8)
	buffer[7] = byte(ts)
}

func getTs(buffer []byte) uint32 {
	return uint32(buffer[4])<<24 | uint32(buffer[5])<<16 | uint32(buffer[6])<<8 | uint32(buffer[7])
}

func getSeq(buffer []byte) uint32 {
	return uint32(buffer[0])<<24 | uint32(buffer[1])<<16 | uint32(buffer[2])<<8 | uint32(buffer[3])
}

func getBytes(kbps int, ms int64) int {
	return (kbps * 1024 / 8) * int(ms) / 1000
}

func onWork(conn *net.UDPConn) {
	sendBuffer := make([]byte, 1024)
	recvBuffer := make([]byte, 1500)
	seq := uint32(0)
	lostTotal := int64(0)
	lostPerSec := 0

	lastDbgMs := time.Now().UnixMilli()

	sendBytes := 2 * 1024
	sendTotal := 0
	avgDelay := 10.0

	for {
		updateBuffer(sendBuffer, seq)
		_, err := conn.Write(sendBuffer)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		n, remoteAddr, err := conn.ReadFromUDP(recvBuffer)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		if n <= 0 {
			fmt.Println("Error: empty data")
			return
		}
		sendBytes -= n
		sendTotal += n

		recvSeq := getSeq(recvBuffer)
		recvTs := getTs(recvBuffer)

		delay := float64(uint32(time.Now().UnixMilli()) - recvTs)
		avgDelay += (delay - avgDelay) / 8.0
		// fmt.Printf("sendBytes: %d, sendTotal: %d, delay:%f, avgDelay:%f\n", sendBytes, sendTotal, delay, avgDelay)
		if recvSeq != seq {
			lostTotal++
			lostPerSec++
		}
		seq++

		diffDbgMs := time.Now().UnixMilli() - lastDbgMs
		if diffDbgMs > 1000 {
			lastDbgMs = time.Now().UnixMilli()
			fmt.Println("Received data from ", remoteAddr.String(),
				", delay:", int(avgDelay),
				", lost:", lostTotal,
				", lost per sec:", lostPerSec,
				", send kbps:", sendTotal*8*1000/1024/int(diffDbgMs))
			lostPerSec = 0
			sendTotal = 0
		}

		for {
			if sendBytes > 0 {
				break
			}
			lastSendMs := time.Now().UnixMilli()
			time.Sleep(100 * time.Millisecond)
			sleepMs := time.Now().UnixMilli() - lastSendMs

			sendBytes += getBytes(*kbps, sleepMs)
		}

	}
}

func main() {
	var conn *net.UDPConn

	if len(*serverIp) == 0 {
		fmt.Println("Error: server ip address is required")
		return
	}
	if *serverPort == 0 {
		fmt.Println("Error: server port is required")
		return
	}

	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", *serverIp, *serverPort))
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	if len(*clientIp) != 0 && *clientPort != 0 {
		clientAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", *clientIp, *clientPort))
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		conn, err = net.DialUDP("udp", clientAddr, serverAddr)
	} else {
		conn, err = net.DialUDP("udp", nil, serverAddr)
	}
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer conn.Close()

	onWork(conn)
}
