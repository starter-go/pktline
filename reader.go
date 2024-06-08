package pktline

import "io"

// Reader 数据包读取接口
type Reader interface {
	Read() (*Packet, error)
}

// ReadCloser 数据包读取（及关闭）接口
type ReadCloser interface {
	Reader
	io.Closer
}
