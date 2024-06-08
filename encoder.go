package pktline

import (
	"bytes"
	"fmt"
	"io"
)

// NewEncoder 创建一个编码器，它把传入的包编码后，写入 dst
func NewEncoder(dst io.Writer) WriteCloser {
	if dst == nil {
		dst = &nopReaderWriter{}
	}
	inst := &encoder{
		dst:             dst,
		enableAutoClose: false,
	}
	return inst
}

type encoder struct {
	dst             io.Writer
	enableAutoClose bool
}

func (inst *encoder) _impl() WriteCloser {
	return inst
}

func (inst *encoder) Close() error {
	nop := &nopReaderWriter{}
	stream := inst.dst
	inst.dst = nop
	if inst.enableAutoClose {
		cl, ok := stream.(io.Closer)
		if ok {
			return cl.Close()
		}
	}
	return nil
}

func (inst *encoder) Write(p *Packet) error {

	head := p.Head
	body := p.Body

	if head == "" && body == nil {
		return inst.writeWithPacketType(p.Type)
	}

	buffer := bytes.NewBuffer([]byte{'x', 'x', 'x', 'x'})

	if head != "" {
		buffer.WriteString(head)
		buffer.WriteRune('\n')
	}

	if body != nil {
		buffer.WriteByte(0)
		buffer.Write(body)
	}

	raw := buffer.Bytes()
	size := len(raw)
	err := inst.updatePacketSize(size, raw[0:4])
	if err != nil {
		return err
	}
	return inst.writeInSize(raw)
}

func (inst *encoder) writeInSize(b []byte) error {
	size1 := len(b)
	size2, err := inst.dst.Write(b)
	if err != nil {
		return err
	}
	if size1 != size2 {
		return fmt.Errorf("bad size of data [want:%d have:%d]", size1, size2)
	}
	return nil
}

func (inst *encoder) writeWithPacketType(t PacketType) error {
	switch t {
	case TypeFlush:
		return inst.writeInSize([]byte{'0', '0', '0', '0'})
	case TypeUndefine1:
		return inst.writeInSize([]byte{'0', '0', '0', '1'})
	case TypeUndefine2:
		return inst.writeInSize([]byte{'0', '0', '0', '2'})
	case TypeUndefine3:
		return inst.writeInSize([]byte{'0', '0', '0', '3'})
	}
	return nil
}

func (inst *encoder) updatePacketSize(size int, buffer []byte) error {

	// check size
	const (
		min = 4
		max = 0xffff
	)
	if size < min {
		return fmt.Errorf("bad pkt-len: %d", size)
	}
	if size > max {
		return fmt.Errorf("bad pkt-len: %d", size)
	}

	// check buffer
	bufferLen := len(buffer)
	if bufferLen != 4 {
		return fmt.Errorf("bad pkt-len buffer size: %d", bufferLen)
	}

	// update
	value := size
	for i := 3; i >= 0; i-- {
		n := value & 0x0f
		value = (value >> 4)
		b := byte(0)
		if 0 <= n && n <= 9 {
			b = byte('0' + n)
		} else {
			b = byte(n - 0x0a + 'a')
		}
		buffer[i] = b
	}
	return nil
}
