package simpletcp

import (
	"context"
	"errors"
	"net"
	mempools "simpletcp/mempools"
	"simpletcp/proto"
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

	//先写入flag
	protoResp.WriteFlag(buffer)

	//偏移后写入长度
	protoResp.WriteLength(uint64(len(*resp)), buffer)

}
