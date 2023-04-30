package main

import (
	"bufio"
	"log"
	"net"
)

var (
	openConnections = make(map[net.Conn]bool)
	newConnection   = make(chan net.Conn)
	deadConnection  = make(chan net.Conn)
)

func main() {
	ln, err := net.Listen("tcp", ":9000")

	if err != nil {
		log.Fatal(err)
	}

	defer ln.Close()

	go func() {
		for {
			conn, err := ln.Accept()

			if err != nil {
				log.Fatal(err)
			}

			openConnections[conn] = true
			newConnection <- conn
		}
	}()

	for {
		select {
		case conn := <-newConnection:
			go broadcastMessage(conn)
		case conn := <-deadConnection:
			for item := range openConnections {
				if item == conn {
					break
				}
			}

			delete(openConnections, conn)
		}
	}
}

func broadcastMessage(conn net.Conn) {
	for {
		reader := bufio.NewReader(conn)
		message, err := reader.ReadString('\n')

		if err != nil {
			break
		}

		for item := range openConnections {
			if item != conn {
				item.Write([]byte(message))
			}
		}
	}

	deadConnection <- conn
}
