package main

import (
	"fmt"
	"time"

	"github.com/nats-io/go-nats"
)

func main() {
	conn, e := nats.Connect("nats://localhost:4222")
	if e != nil {
		fmt.Errorf(e.Error())
	}
	msg, e := conn.Request("app.queue.test.create", []byte(`{"username":"admin", "password":"password", "name":"q-1"}`), 2*time.Second)
	if e != nil {
		fmt.Println(e.Error())
	}
	fmt.Println(string(msg.Data))

	msg, e := conn.Request("app.queue.test.delete", []byte(`{"username":"admin", "password":"password", "name":"q-1"}`), 2*time.Second)
	if e != nil {
		fmt.Println(e.Error())
	}
	fmt.Println(string(msg.Data))
}
