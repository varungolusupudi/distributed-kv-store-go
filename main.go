package main

import (
	"bufio"
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

		operation := strings.ToUpper(args[0])
		switch operation {
		case "GET":
			store.Get(args, conn)
		case "SET":
			store.Set(args, conn)
		case "DEL":
			store.Del(args, conn)
		case "EXPIRE":
			store.Expire(args, conn)
		default:
			conn.Write([]byte("Unknown Operation....\n"))
		}

	}

}

func main() {
	fmt.Println("Node starting...")

	store.InitStore()

	// Opening the config.json file
	jsonFile, err := os.Open("config.json")

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
