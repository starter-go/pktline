package pktline

import (
	"fmt"
	"io"
)

type nopReaderWriter struct {
}

func (inst *nopReaderWriter) _impl() (io.ReadCloser, io.WriteCloser) {
	return inst, inst
}

func (inst *nopReaderWriter) Read(b []byte) (int, error) {
	return 0, fmt.Errorf("stream closed")
}

func (inst *nopReaderWriter) Write(b []byte) (int, error) {
	return 0, fmt.Errorf("stream closed")
}

func (inst *nopReaderWriter) Close() error {
	return fmt.Errorf("stream closed")
}
