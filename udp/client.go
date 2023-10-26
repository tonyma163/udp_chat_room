package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	serverAddress = "127.0.0.1:6060"
)

func sendMessage(conn *net.UDPConn, room, message string) {
	fullMessage := room + " " + message
	_, err := conn.Write([]byte(fullMessage))
	if err != nil {
		fmt.Println("Error sending message: ", err)
		return
	}
}

func receiveMessage(conn *net.UDPConn, room string) {
	buffer := make([]byte, 1024)

	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from server: ", err)
			return
		}

		fullMessage := string(buffer[:n])
		parts := strings.SplitN(fullMessage, " ", 2)
		if len(parts) < 2 {
			continue
		}

		messageRoom := parts[0]
		messageContent := parts[1]

		if messageRoom == room {
			fmt.Printf("\n%s\n", messageContent)
			fmt.Print(">")
		}
	}
}

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", serverAddress)
	if err != nil {
		fmt.Println("Error resolving address: ", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("Error connecting to server: ", err)
		return
	}
	defer conn.Close()

	// Room Setting
	room := "room1"
	message := "join"
	sendMessage(conn, room, message)

	fmt.Println("Welcome to UDP Chat!")

	for {
		go receiveMessage(conn, room)
		fmt.Print("\nType message to send> ")
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')
		message = strings.TrimRight(message, "\n")

		sendMessage(conn, room, message)
	}
}
