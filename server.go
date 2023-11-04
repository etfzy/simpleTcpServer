// Code generated by "sloth -out=SimpleTcpSever -mod=./model -fun=get,set"; DO NOT EDIT.
package simpletcp

import (
	"errors"
	"fmt"
	"net"
	mempools "simpletcp/mempools"
	"simpletcp/proto"

	"go.uber.org/zap"
)

type SimpleConf struct {
	protocol      string
	addr          string
	log           *zap.Logger
	contentScopes []uint64
	proto         proto.Proto
	memps         mempools.MemPools
	event         Event
}

type TcpServer interface {
	Start() error
	SetLog(l *zap.Logger) TcpServer
	SetReceive(r *proto.Tlv) TcpServer
	SetResponse(r *proto.Tlv) TcpServer
	EnableContentMemPools(scopes []uint64) TcpServer
}

func CreateSimpleServer(protocol, addr string, event Event) TcpServer {
	return &SimpleConf{
		protocol: protocol,
		addr:     addr,
		log:      createDefaultLog(),
		proto:    proto.NewProto(),
		event:    event,
	}
}

func (s *SimpleConf) SetLog(l *zap.Logger) TcpServer {
	s.log = l
	return s
}

func (s *SimpleConf) SetReceive(r *proto.Tlv) TcpServer {
	s.proto.SetRecv(r)
	return s
}

func (s *SimpleConf) SetResponse(r *proto.Tlv) TcpServer {
	s.proto.SetResp(r)
	return s
}

// 设置返回内容的预期长度，底层buffer会根据该范围进行sync pool
func (s *SimpleConf) EnableContentMemPools(scopes []uint64) TcpServer {
	s.contentScopes = scopes
	return s
}

func (s *SimpleConf) checkConf() error {

	recv := s.proto.GetRecv()
	if recv == nil {
		return errors.New("receive config can not be nil!")
	}

	err := recv.CheckTlv()
	if err != nil {
		return err
	}

	resp := s.proto.GetResp()
	if resp == nil {
		return nil
	}

	err = resp.CheckTlv()
	if err != nil {
		return err
	}

	return nil
}

func (s *SimpleConf) getProtoLen() []uint64 {
	result := make([]uint64, 0, 4)
	recv := s.proto.GetRecv()
	result = append(result, recv.GetFlagLen())
	result = append(result, recv.GetLengthLen())

	resp := s.proto.GetResp()
	if resp == nil {
		return result
	}

	result = append(result, resp.GetFlagLen())
	result = append(result, resp.GetLengthLen())

	return nil
}

func (s *SimpleConf) createConnection(conn net.Conn) *connection {
	return &connection{
		context: createConnContext(s.proto, s.memps, conn),
		proto:   s.proto,
		log:     s.log,
		event:   s.event,
		memps:   s.memps,
	}
}

func (s *SimpleConf) Start() error {

	err := s.checkConf()
	if err != nil {
		return err
	}

	lens := s.getProtoLen()
	s.memps = mempools.CreateMems(lens, s.contentScopes)

	listen, err := net.Listen(s.protocol, s.addr)
	if err != nil {
		return err
	}

	defer listen.Close()

	fmt.Println("tcp server listen ", s.addr)

	for {
		conn, err := listen.Accept()
		if err != nil {
			s.log.Error("connection accept error", zap.String("context", "Accept"), zap.String("error", err.Error()))
			continue
		}

		connection := s.createConnection(conn)
		go connection.process()
	}
}
