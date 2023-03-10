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

	if status == 1 {
		return nil, fmt.Errorf("unable to set %s", key)
	}

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
	binary.Read(c.conn, binary.LittleEndian, &resp.Status)
	binary.Read(c.conn, binary.LittleEndian, &valLen)
	resp.Value = make([]byte, valLen)
	binary.Read(c.conn, binary.LittleEndian, &resp.Value)

	if resp.Status == 1 {
		return nil, fmt.Errorf("no key found for %s", key)
	}
	return string(resp.Value), nil
}

func New(listenAddr string) (*Client, error) {
	conn, err := net.Dial("tcp", listenAddr)

	if err != nil {
		fmt.Printf("Unable to dial %v", err.Error())
	}

	return &Client{conn: conn}, nil
}
