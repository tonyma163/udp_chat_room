package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

const (
	serverAddress = "127.0.0.1:6060"
)

type Client struct {
	Address *net.UDPAddr
	Room    string
}

var (
	clients     = make(map[string]*Client)
	clientMutex sync.Mutex
)

func setClientRoom(addr *net.UDPAddr, room string) {
	clientMutex.Lock()
	defer clientMutex.Unlock()

	client, ok := clients[addr.String()]
	if !ok {
		clients[addr.String()] = &Client{Address: addr, Room: room}
	} else {
		client.Room = room
	}
}

func broadcastMessage(conn *net.UDPConn, sender *net.UDPAddr, message string) {
	clientMutex.Lock()
	defer clientMutex.Unlock()

	senderClient := clients[sender.String()]

	if senderClient == nil {
		fmt.Println("Client not found: ", sender)
		return
	}

	fullMessage := ""
	if message == "join" {
		fullMessage = senderClient.Room + " " + senderClient.Address.String() + " has joined!"
	} else {
		fullMessage = senderClient.Room + " " + senderClient.Address.String() + ": " + message
	}

	for _, client := range clients {
		if client.Room == senderClient.Room && client.Address != senderClient.Address {
			_, err := conn.WriteToUDP([]byte(fullMessage), client.Address)
			if err != nil {
				fmt.Println("Error sending message to ", client.Address, ": ", err)
			}
		}
	}
}

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", serverAddress)
	if err != nil {
		fmt.Println("Error resolving address: ", err)
		return
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Error while listening: ", err)
		return
	}
	defer conn.Close()

	fmt.Printf("Server is running on %s\n", serverAddress)

	buffer := make([]byte, 1024)

	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP: ", err)
			continue
		}

		message := string(buffer[:n])
		fmt.Printf("Received message from %s: %s\n", addr, message)

		parts := strings.SplitN(message, " ", 2)
		if len(parts) >= 2 {
			room := parts[0]
			message = parts[1]
			setClientRoom(addr, room)
		}

		go broadcastMessage(conn, addr, message)
	}
}
