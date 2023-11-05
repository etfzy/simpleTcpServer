package simpletcp

import (
	"io"

	mem "github.com/etfzy/mempool/base"
	mempools "github.com/etfzy/simpleTcpServer/mempools"
	"github.com/etfzy/simpleTcpServer/proto"
	"go.uber.org/zap"
)

type connection struct {
	context *ConnContext
	proto   proto.Proto
	memps   mempools.MemPools
	log     *zap.Logger
	event   Event
}

func (cc *connection) process() {
	defer cc.context.Conn.Close()
	//回调处理
	err := cc.event.onOpen(cc.context)
	if err != nil {
		cc.log.Error("connection close", zap.String("context", "onOpen"), zap.String("error", err.Error()))
		cc.close()
		return
	}

	for {
		//解析flag
		err := cc.flag()
		if err != nil {
			return
		}

		//解析长度
		length, err := cc.length()
		if err != nil {
			return
		}

		//获取内容
		content, err := cc.content(length)
		defer cc.memps.PutContentMems(content)
		if err != nil {
			return
		}

		//回调处理
		err = cc.event.onReact(content.Buf(), cc.context)
		if err != nil {
			cc.log.Info("connection close", zap.String("context", "onReact"), zap.String("error", err.Error()))
			cc.close()
			return
		}
	}
}

func (cc *connection) flag() error {
	//flag 解析处理
	flaglen := cc.proto.GetRecv().GetFlagLen()
	buf := cc.memps.GetProtoMems(flaglen)
	defer cc.memps.PutProtoMems(buf)

	err := cc.read(buf.Buf())

	if err != nil {
		if err != io.EOF {
			cc.log.Error("connection close", zap.String("context", "read"), zap.String("error", err.Error()))
			cc.close()
			return err
		} else {
			cc.log.Info("connection close", zap.String("context", "EOF"), zap.String("error", "EOF"))
			cc.close()
			return err
		}
	}

	err = cc.proto.GetRecv().ReadFlag(buf)
	if err != nil {
		cc.log.Error("connection close", zap.String("context", "ParseFlag"), zap.String("error", err.Error()))
		cc.close()
		return err
	}

	return nil
}

func (cc *connection) length() (uint64, error) {
	//length 处理
	lenLen := cc.proto.GetRecv().GetLengthLen()
	buf := cc.memps.GetProtoMems(lenLen)
	defer cc.memps.PutProtoMems(buf)

	err := cc.read(buf.Buf())
	if err != nil {
		if err != io.EOF {
			cc.log.Error("connection close", zap.String("context", "read"), zap.String("error", err.Error()))
			cc.close()
			return 0, err
		} else {
			cc.log.Info("connection close", zap.String("context", "EOF"), zap.String("error", "EOF"))
			cc.close()
			return 0, err
		}
	}

	length, err := cc.proto.GetRecv().ReadLength(buf)
	if err != nil {
		cc.log.Error("connection close", zap.String("context", "ParseLength"), zap.String("error", err.Error()))
		cc.close()
		return length, err
	}

	return length, nil
}

func (cc *connection) content(length uint64) (*mem.Buffer[byte], error) {
	content := cc.memps.GetContentMems(length)
	err := cc.read(content.Buf())

	if err != nil {
		if err != io.EOF {
			cc.log.Error("connection close", zap.String("context", "read"), zap.String("error", err.Error()))
			cc.close()
			return nil, err
		} else {
			cc.log.Info("connection close", zap.String("context", "EOF"), zap.String("error", "EOF"))
			cc.close()
			return nil, err
		}
	}

	return content, nil
}

func (cc *connection) read(target *[]byte) error {

	_, err := io.ReadFull(cc.context.Conn, *target)

	return err
}

func (cc *connection) close() {
	cc.context.Conn.Close()
	cc.event.onClose(cc.context)
	return
}
