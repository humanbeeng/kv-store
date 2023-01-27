package server

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/fatih/color"
	"github.com/humanbeeng/kv-store/cache"
	"github.com/humanbeeng/kv-store/proto"
)

type Server struct {
	storage    cache.Storer[string, any]
	ListenAddr string
}

func NewServer(listenAddr string) *Server {
	return &Server{
		storage:    cache.NewKVStore[string, any](),
		ListenAddr: listenAddr,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		log.Fatalf("error listening: %v", err.Error())
	}
	defer ln.Close()

	color.Blue("Server listening on %s", s.ListenAddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("error accepting connection %v", err.Error())
			return err
		}
		go s.handleConnection(conn)
	}

}

func (s *Server) handleConnection(conn net.Conn) {

	for {
		err := handleCommand(conn)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Invalid command : %v", err.Error())
			break
		}
	}

}

func handleCommand(conn net.Conn) error {
	var cmd proto.CommandAlias

	if err := binary.Read(conn, binary.LittleEndian, &cmd); err != nil {
		return err
	}

	switch cmd {
	case proto.CmdGet:
		return handleGetCommand(conn)
	case proto.CmdSet:
		return handleSetCommand(conn)
	default:
		{
			return fmt.Errorf("invalid command")
		}
	}
}

func handleGetCommand(conn net.Conn) error {
	cmdGet, err := parseGetCommand(conn)

	if err != nil {
		return fmt.Errorf(err.Error())
	}
	resp := &proto.GetResponse{}

	value, err := cache.StringStore.Get(fmt.Sprint(cmdGet.Key))
	if err != nil {
		return err
	}

	resp.Value = []byte(value)

	conn.Write(resp.Bytes())

	return nil
}

func handleSetCommand(conn net.Conn) error {
	cmdSet, err := parseSetCommand(conn)

	if err != nil {
		return err
	}
	resp := &proto.SetResponse{}
	err = cache.StringStore.Put(fmt.Sprint(cmdSet.Key), string(cmdSet.Value))
	if err != nil {
		resp.Status = 1
		conn.Write(resp.Bytes())
	}
	resp.Status = 1

	conn.Write(resp.Bytes())
	return nil
}

func parseSetCommand(r io.Reader) (*proto.CommandSet, error) {
	cmdSet := &proto.CommandSet{}

	var keyLen int32
	binary.Read(r, binary.LittleEndian, &keyLen)
	cmdSet.Key = make([]byte, keyLen)
	binary.Read(r, binary.LittleEndian, &cmdSet.Key)

	var valueLen int32
	binary.Read(r, binary.LittleEndian, &valueLen)
	cmdSet.Value = make([]byte, valueLen)
	binary.Read(r, binary.LittleEndian, &cmdSet.Value)

	return cmdSet, nil
}

func parseGetCommand(r io.Reader) (*proto.CommandGet, error) {
	cmdGet := &proto.CommandGet{}

	var keyLen int32
	binary.Read(r, binary.LittleEndian, &keyLen)
	cmdGet.Key = make([]byte, keyLen)
	binary.Read(r, binary.LittleEndian, &cmdGet.Key)

	return cmdGet, nil
}
