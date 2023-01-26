package proto

import (
	"bytes"
	"encoding/binary"
)

type Command interface {
	CommandAlias() CommandAlias
}

type CommandAlias byte

const (
	CmdNone CommandAlias = iota
	CmdGet
	CmdSet
	CmdDel
	CmdUp
)

type CommandGet struct {
	Key []byte
}

type CommandSet struct {
	Key   []byte
	Value []byte
}

type GetResponse struct {
	Value []byte
}

type SetResponse struct {
	Status byte
}

type CommandNone struct{}

func (cmdGet *CommandGet) CommandAlias() CommandAlias {
	return CmdGet
}

func (cmdGet *CommandSet) CommandAlias() CommandAlias {
	return CmdSet
}

func (cmdNone *CommandNone) CommandAlias() CommandAlias {
	return CmdNone
}

func (cmdSet *CommandSet) Bytes() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, CmdSet)
	keyLen := int32(len(cmdSet.Key))
	binary.Write(buf, binary.LittleEndian, keyLen)
	binary.Write(buf, binary.LittleEndian, cmdSet.Key)

	valueLen := int32(len(cmdSet.Value))
	binary.Write(buf, binary.LittleEndian, valueLen)
	binary.Write(buf, binary.LittleEndian, cmdSet.Value)

	return buf.Bytes()
}

func (cmdGet *CommandGet) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, CmdGet)

	keyLen := int32(len(cmdGet.Key))
	binary.Write(buf, binary.LittleEndian, keyLen)
	binary.Write(buf, binary.LittleEndian, cmdGet.Key)

	return buf.Bytes()
}

func (resp *GetResponse) Bytes() []byte {
	buf := new(bytes.Buffer)

	valLen := len(resp.Value)
	binary.Write(buf, binary.LittleEndian, int32(valLen))
	binary.Write(buf, binary.LittleEndian, resp.Value)
	return buf.Bytes()
}

func (resp *SetResponse) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, resp.Status)
	return buf.Bytes()
}
