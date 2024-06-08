package pktline

import "io"

// Writer 数据包写入接口
type Writer interface {
	Write(p *Packet) error
}

// WriteCloser 数据包写入（及关闭）接口
type WriteCloser interface {
	Writer
	io.Closer
}
