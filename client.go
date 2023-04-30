package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

func main() {
	connection, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		log.Fatal(err)
	}

	defer connection.Close()

	fmt.Printf("Name: ")

	reader := bufio.NewReader(os.Stdin)
	username, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	username = strings.Trim(username, " \r\n")
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
	fmt.Println("************TCP CHAT************")
	go read(connection)
	write(connection, username)
}

func read(connection net.Conn) {
	for {
		reader := bufio.NewReader(connection)
		message, err := reader.ReadString('\r')

		if err == io.EOF {
			connection.Close()
			fmt.Println("Connection closed.")
			os.Exit(0)
		}

		fmt.Println(message)
	}
}

func write(connection net.Conn, username string) {
	for {
		reader := bufio.NewReader(os.Stdin)
		message, err := reader.ReadString('\n')

		if err != nil {
			break
		}

		message = fmt.Sprintf("%s: %s", username, strings.Trim(message, " \r"))
		connection.Write([]byte(message))

	}
}
