package main

import (
	"log"
	"time"

	"github.com/humanbeeng/kv-store/client"
	"github.com/humanbeeng/kv-store/server"
)

func main() {

	go func() {
		InitClient()
	}()

	s := server.NewServer("localhost:8080")
	s.Start()

}

func InitClient() {
	time.Sleep(time.Second * 2)

	client, err := client.New()
	if err != nil {
		log.Fatal(err.Error())
	}
	client.Set([]byte("hello"), []byte("there"))

	time.Sleep(time.Millisecond * 500)
	client.Get([]byte("hello"))
}
