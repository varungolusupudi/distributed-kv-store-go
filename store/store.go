package store

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

var (
	globalCache map[string]interface{}
	mutex       sync.RWMutex
	ttlMap      map[string]time.Time
)

func InitStore() {
	globalCache = make(map[string]interface{})
	ttlMap = make(map[string]time.Time)

	go func() {
		for {
			time.Sleep(1 * time.Second)

			mutex.Lock()
			for k, v := range ttlMap {
				if time.Now().After(v) {
					delete(ttlMap, k)
				}
			}
			mutex.Unlock()
		}
	}()
}

func Set(args []string, conn net.Conn) {
	if len(args) != 3 {
		conn.Write([]byte("ERR usage: SET <key> <value>\n"))
		return
	}

	key := args[1]
	value := args[2]

	mutex.Lock()
	globalCache[key] = value
	mutex.Unlock()

	conn.Write([]byte("OK\n"))
}

func Get(args []string, conn net.Conn) {
	if len(args) != 2 {
		conn.Write([]byte("ERR usage: GET <key>\n"))
		return
	}

	key := args[1]

	mutex.RLock()
	expiry, exists := ttlMap[key]
	value, ok := globalCache[key]
	mutex.RUnlock()

	if exists && time.Now().After(expiry) {
		conn.Write([]byte("Key Expired\n"))
		mutex.Lock()
		delete(ttlMap, key)
		delete(globalCache, key)
		mutex.Unlock()
		return
	}

	if !ok {
		conn.Write([]byte("nil\n"))
		return
	}

	conn.Write([]byte(fmt.Sprintf("%v", value)))
}

func Del(args []string, conn net.Conn) {
	if len(args) != 2 {
		conn.Write([]byte("ERR usage: DEL <key>\n"))
		return
	}

	key := args[1]

	mutex.Lock()
	delete(globalCache, key)
	mutex.Unlock()

	conn.Write([]byte("OK\n"))
}

func Expire(args []string, conn net.Conn) {
	if len(args) != 3 {
		conn.Write([]byte("ERR usage: EXPIRE <key> <seconds>\n"))
		return
	}

	key := args[1]

	seconds, err := strconv.Atoi(args[2])
	if err != nil {
		conn.Write([]byte("ERR usage: EXPIRE <key> <seconds>\n"))
		return
	}

	mutex.Lock()
	ttlMap[key] = time.Now().Add(time.Duration(seconds) * time.Second)
	mutex.Unlock()

	conn.Write([]byte("OK\n"))
}
