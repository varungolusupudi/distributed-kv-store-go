package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
)

type Config struct {
	Port  string   `json:"port"`
	Peers []string `json:"peers"`
}

func handleConnection(conn net.Conn) {
	// TODO
}

func main() {
	fmt.Println("Nodes starting...")

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
		// TODO: Gracefully exit?
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleConnection(conn)
	}
}
