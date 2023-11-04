package simpletcp

import (
	"bytes"
	"context"
	"errors"
	"net"

	mempools "github.com/etfzy/simpleTcpServer/mempools"
	"github.com/etfzy/simpleTcpServer/proto"
)

type ConnContext struct {
	Context context.Context
	Conn    net.Conn
	proto   proto.Proto
	memps   mempools.MemPools
}

func createConnContext(proto proto.Proto, memps mempools.MemPools, c net.Conn) *ConnContext {
	return &ConnContext{
		Context: context.Background(),
		Conn:    c,
		proto:   proto,
		memps:   memps,
	}
}

func (c *ConnContext) send(resp *[]byte) error {
	protoResp := c.proto.GetResp()
	if protoResp == nil {
		return errors.New("proto response config is nil, can not send data!")
	}

	if resp == nil {
		return errors.New("input is nil")
	}

	if len(*resp) == 0 {
		return errors.New("input length is zero")
	}

	contentLen := protoResp.GetFlagLen() + protoResp.GetLengthLen() + uint64(len(*resp))

	//这里先清零，再使用
	buffer := c.memps.GetContentMems(contentLen)
	newBuffer := (*buffer)
	newBuffer = newBuffer[:0]

	//先写入flag
	err := protoResp.WriteFlag(buffer)
	if err != nil {
		return err
	}

	bytes.Buffer
	//偏移后写入长度
	err := protoResp.WriteLength(uint64(len(*resp)), buffer)
	if err != nil {
		return err
	}
	buffer
}
