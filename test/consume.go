package main

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/go-nats"
)

func main() {
	conn, e := nats.Connect("nats://localhost:4222")
	if e != nil {
		fmt.Errorf(e.Error())
	}
	topic := fmt.Sprintf("%x", sha256.Sum256([]byte("q-1")))
	fmt.Println(topic)
	conn.Subscribe(topic, func(m *nats.Msg) {
		fmt.Printf(string(m.Data))
	})
	msg, e := conn.Request("app.queue.test.req.consume", []byte(`{"username":"admin", "password":"password", "name":"q-1"}`), 2*time.Second)
	for e != nil {
		_, e = conn.Request("app.queue.test.req.consume", []byte(`{"username":"admin", "password":"password", "name":"q-1"}`), 2*time.Second)
		fmt.Println(e.Error())
	}
	fmt.Println(string(msg.Data))
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
	fmt.Println(string(msg.Data))
}
