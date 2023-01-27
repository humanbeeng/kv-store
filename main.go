package main

import (
	"flag"
	"log"

	"github.com/humanbeeng/kv-store/server"
)

func main() {

	var (
		listenAddr = flag.String("address", ":3000", "Address for the server to listen on.")
	)

	flag.Parse()

	s := server.NewServer(":" + *listenAddr)
	err := s.Start()
	if err != nil {
		log.Fatal(err.Error())
	}
}
