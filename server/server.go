package main

import (
	"fmt"
	"net"
	"strings"
)

type Msg struct {
	id  int
	txt string
}

var ClientList map[int]net.Conn = make(map[int]net.Conn)
var msgCh chan Msg = make(chan Msg, 100)

func main() {
	ln := StartServer(":1000")
	defer ln.Close()

	go ChatBoard()

	for id := 1; len(ClientList) < 100; id++ {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("[Server] Unknown Client failed to connect")
			continue
		}
		go OnClientConnect(conn, id)
	}
	OnServerEnd()
}
func StartServer(endpoint string) net.Listener {
	ln, err := net.Listen("tcp", endpoint)
	if err != nil {
		fmt.Println("[Server] Server Can't Start. Exit Program")
		panic(err)
	}
	fmt.Println("[Server] Start tcp listening")
	return ln
}
func OnServerEnd() {
	SendAll("[Global] Server is going to be closed")
}

func SendAll(msg string) {
	for id, client := range ClientList {
		_, err := client.Write([]byte(msg))
		if err != nil {
			fmt.Printf("[Server] Client %v connection lost\n", id)
			defer OnClientDisconneted(id)
		}
	}
}

func OnClientConnect(conn net.Conn, id int) {
	defer conn.Close()
	ClientList[id] = conn
	fmt.Printf("[Server] Client %v is connected\n", id)
	SendAll(fmt.Sprintf("[Global] Client %v is connected\n", id))
	readBuf := make([]byte, 4092)
	var msg string
	for {
		n, err := conn.Read(readBuf)
		if err != nil {
			fmt.Printf("[Server] Client %v connection lost\n", id)
			OnClientDisconneted(id)
			break
		}
		msg = string(readBuf[:n])
		if len(msg) > 0 {
			msg = strings.TrimRight(msg, "\r\n")
			if msg == "!exit" {
				break
			}
			fmt.Printf("[Client %v] %s\n", id, msg)
			msgCh <- Msg{id, msg}
		}
	}
	fmt.Printf("[Server] Client %v exit\n", id)
}

func OnClientDisconneted(id int) {
	delete(ClientList, id)
	SendAll(fmt.Sprintf("[Global] Client %v disconnected", id))
}

func ChatBoard() {
	for msg := range msgCh {
		txt := fmt.Sprintf("[Client %v] %s\n", msg.id, msg.txt)
		SendAll(txt)
		/*
			for id, client := range ClientList {
				_, err := client.Write([]byte(txt))
				if err != nil {
					fmt.Printf("[Server] Client %v connection lost\n", id)
					//delete(ClientList, id)
					OnClientDisconneted(id)
				}
			}
		*/
	}
}
