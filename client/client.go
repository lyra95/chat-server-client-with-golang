package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	//"time"
)

func main() {
	//t, _ := time.ParseDuration("10m")
	//conn, err := net.DialTimeout("tcp", "3.36.128.78:65507", t)
	// conn, err := net.Dial("tcp", "ec2-3-36-128-78.ap-northeast-2.compute.amazonaws.com:65507")
	conn, err := net.Dial("tcp", ":1000")
	if err != nil {
		fmt.Println("[Client] Can't connect to server")
	}
	defer conn.Close()
	fmt.Println("[Client] Connected to server")
	go ChatBoard(conn)
	reader := bufio.NewReader(os.Stdin)
	for {
		msg, _ := reader.ReadString('\n')
		_, err := conn.Write([]byte(msg))
		if err != nil {
			fmt.Println("[Client] Connection lost.")
			break
		}
	}
	fmt.Println("[Client] exit")
}

func ChatBoard(conn net.Conn) {
	readBuf := make([]byte, 4092)
	for {
		n, err := conn.Read(readBuf)
		if err != nil {
			log.Println("[Client] Connection lost")
			break
		}
		msg := string(readBuf[:n])
		if len(msg) > 0 {
			msg = strings.TrimRight(msg, "\r\n")
			log.Println(msg)
		}
	}
}
