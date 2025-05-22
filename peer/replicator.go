package peer

import (
	"distributed-kv-store-go/store"
	"log"
	"net"
	"strings"
)

var peersList []string

func InitReplicator(peers []string) {
	peersList = peers
}

func BroadcastToPeers(command string) {
	for _, peer := range peersList {
		conn, err := net.Dial("tcp", peer)
		if err != nil {
			log.Printf("Error connecting to peer %s\n", peer)
			continue
		}

		log.Printf("Connected to peer %s\n", peer)
		_, err = conn.Write([]byte(command + "\n"))
		if err != nil {
			log.Printf("Write failed to %s: %v", peer, err)
		}
		conn.Close()
	}
}

func HandlePeerMessage(args []string, conn net.Conn) {
	operation := strings.ToUpper(args[0])
	args = args[1:]
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
