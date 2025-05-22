package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	fmt.Println("Client Starting....")

	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	fmt.Println("Client Connection Established.")

	fmt.Printf("Enter your message: ")

	go func() {
		reader := bufio.NewReader(conn)
		for {
			msg, err := reader.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf(msg)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		if scanner.Scan() {
			msg := scanner.Text() + "\n"
			_, err := conn.Write([]byte(msg))
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
