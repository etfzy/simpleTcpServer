package simpletcp

import (
	"net"
)

type Event interface {
	onOpen(conn net.Conn) error
	onClose(conn *ConnContext)
	onReact(data *[]byte, conn *ConnContext) error
}

type event struct{}

// 连接生成时的回调
func (e *event) onOpen(conn *ConnContext) error {
	return nil
}

// 连接关闭时的回调
func (e *event) onClose(conn *ConnContext) {}

// 请求进入时的回调
func (e *event) onReact(data *[]byte, conn *ConnContext) error {
	return nil
}
