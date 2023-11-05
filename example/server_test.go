package example

import (
	"encoding/binary"
	"fmt"
	"testing"

	simpletcp "github.com/etfzy/simpleTcpServer"
)

type server struct {
	*simpletcp.RootEvent
}

func (s *server) onReact(data *[]byte, context *simpletcp.ConnContext) error {

	fmt.Println(data)
	return nil
}

func TestServer(t *testing.T) {
	event := server{}
	tcpsever := simpletcp.CreateSimpleServer("tcp", ":5000", &event)

	//设置receive
	tcpsever.NewReceive().
		SetFlag(1001).
		SetFlagLen(4).
		SetBorder(binary.LittleEndian).
		SetLengthLen(8).
		SetMaxLength(16 * 1024)
	tcpsever.NewResponse().
		SetFlag(2002).
		SetFlagLen(4).
		SetBorder(binary.LittleEndian).
		SetLengthLen(8).
		SetMaxLength(16 * 1024)

	tcpsever.Start()

}
