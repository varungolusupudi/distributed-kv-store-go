package main

import (
	"bufio"
	"distributed-kv-store-go/peer"
	"distributed-kv-store-go/store"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

type Config struct {
	Port  string   `json:"port"`
	Peers []string `json:"peers"`
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		msg, err := reader.ReadString('\n')

		if err != nil {
			log.Println("Client Disconnected....")
			return
		}

		fmt.Printf("Received %s \n", msg)

		args := strings.Split(strings.TrimSpace(msg), " ")

		if len(args) == 0 {
			conn.Write([]byte("ERR: Empty command\n"))
			return
		}

		if args[0] == "REPL" {
			if len(args) < 2 {
				conn.Write([]byte("ERR: Malformed REPL command\n"))
				return
			}
			peer.HandlePeerMessage(args[1:], conn)
			return
		}

		operation := strings.ToUpper(args[0])
		switch operation {
		case "GET":
			store.Get(args, conn)
		case "SET":
			store.Set(args, conn)
			peer.BroadcastToPeers("REPL " + msg)
		case "DEL":
			store.Del(args, conn)
			peer.BroadcastToPeers("REPL " + msg)
		case "EXPIRE":
			store.Expire(args, conn)
			peer.BroadcastToPeers("REPL " + msg)
		default:
			conn.Write([]byte("Unknown Operation....\n"))
		}

	}

}

func main() {
	fmt.Println("Node starting...")

	store.InitStore()

	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <config-file>")
	}
	configPath := os.Args[1]

	jsonFile, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Failed to open config file %s: %v", configPath, err)
	}

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully Opened config.json")

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var config Config
	json.Unmarshal(byteValue, &config)

	defer jsonFile.Close()

	ln, err := net.Listen("tcp", config.Port)
	if err != nil {
		log.Println("Error listening on port " + config.Port)
		return
	}

	peer.InitReplicator(config.Peers)

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
			return
		}

		go handleConnection(conn)
	}
}
