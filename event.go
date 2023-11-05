package simpletcp

type Event interface {
	onOpen(context *ConnContext) error
	onClose(context *ConnContext)
	onReact(data *[]byte, context *ConnContext) error
}

type RootEvent struct{}

// 连接生成时的回调
func (e *RootEvent) onOpen(context *ConnContext) error {
	return nil
}

// 连接关闭时的回调
func (e *RootEvent) onClose(context *ConnContext) {}

// 请求进入时的回调
func (e *RootEvent) onReact(data *[]byte, context *ConnContext) error {
	return nil
}
