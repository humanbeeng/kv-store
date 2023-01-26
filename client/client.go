package client

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/humanbeeng/kv-store/proto"
)

type Client struct {
	conn net.Conn
}

func (c *Client) Set(key []byte, value []byte) (any, error) {
	cmdset := &proto.CommandSet{Key: key, Value: value}
	_, err := c.conn.Write(cmdset.Bytes())

	if err != nil {
		return err, err
	}

	var status byte

	binary.Read(c.conn, binary.LittleEndian, &status)

	return nil, err
}

func (c *Client) Get(key []byte) (any, error) {
	cmdGet := proto.CommandGet{Key: key}
	_, err := c.conn.Write(cmdGet.Bytes())
	if err != nil {
		return nil, err
	}

	resp := &proto.GetResponse{}

	var valLen int32
	binary.Read(c.conn, binary.LittleEndian, &valLen)
	resp.Value = make([]byte, valLen)
	binary.Read(c.conn, binary.LittleEndian, &resp.Value)

	fmt.Printf("Client GET: %s", string(resp.Value))
	return string(resp.Value), nil
}

func New() (*Client, error) {
	conn, err := net.Dial("tcp", "localhost:8080")

	if err != nil {
		fmt.Printf("Unable to dial %v", err.Error())
	}

	return &Client{conn: conn}, nil
}
