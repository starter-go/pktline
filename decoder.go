package pktline

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

// NewDecoder 创建一个解码器，它从 src 读取一个数据包，并解码
func NewDecoder(src io.Reader) ReadCloser {
	if src == nil {
		src = &nopReaderWriter{}
	}
	inst := &decoder{
		src:             src,
		enableAutoClose: false,
	}
	return inst
}

////////////////////////////////////////////////////////////////////////////////

const (
	packetSizeLength = 4
)

type decoder struct {
	src             io.Reader
	enableAutoClose bool
}

func (inst *decoder) _impl() ReadCloser {
	return inst
}

func (inst *decoder) Close() error {
	nop := &nopReaderWriter{}
	stream := inst.src
	inst.src = nop
	if inst.enableAutoClose {
		cl, ok := stream.(io.Closer)
		if ok {
			return cl.Close()
		}
	}
	return nil
}

func (inst *decoder) Read() (*Packet, error) {

	// read size
	size, err := inst.readPacketSize()
	if err != nil {
		return nil, err
	}

	// prepare packet
	pack := &Packet{}
	if size < packetSizeLength {
		pack.Type = inst.sizeToType(size)
		return pack, nil
	}
	pack.Type = TypeData

	// read data
	buffer := make([]byte, size-packetSizeLength)
	err = inst.readInSize(buffer)
	if err != nil {
		return nil, err
	}

	// parse head & body
	idx := bytes.IndexByte(buffer, 0)
	if idx < 0 {
		// no end of string
		pack.Head = string(buffer)
	} else {
		pack.Head = string(buffer[0:idx])
		pack.Body = buffer[idx+1:]
	}

	head := pack.Head
	if head != "" {
		pack.Head = strings.TrimSpace(head)
	}

	return pack, nil
}

func (inst *decoder) sizeToType(size int) PacketType {
	switch size {
	case 0:
		return TypeFlush
	case 1:
		return TypeUndefine1
	case 2:
		return TypeUndefine2
	case 3:
		return TypeUndefine3
	}
	return TypeData
}

func (inst *decoder) readPacketSize() (int, error) {

	bufferRaw := [packetSizeLength]byte{}
	buffer := bufferRaw[:]
	value := 0

	err := inst.readInSize(buffer)
	if err != nil {
		return 0, err
	}

	for i := 0; i < packetSizeLength; i++ {
		b := buffer[i]
		n := 0
		if '0' <= b && b <= '9' {
			n = int(b - '0')
		} else if 'a' <= b && b <= 'f' {
			n = int(b-'a') + 0x0a
		} else if 'A' <= b && b <= 'F' {
			n = int(b-'A') + 0x0a
		} else {
			str := string(buffer)
			return 0, fmt.Errorf("bad format of pkt-len text: '%s'", str)
		}
		value = (value << 4) | n
	}
	return value, nil
}

func (inst *decoder) readInSize(buf []byte) error {
	size1 := len(buf)
	size2, err := io.ReadAtLeast(inst.src, buf, size1)
	if err != nil {
		return err
	}
	if size1 != size2 {
		return fmt.Errorf("bad size of data [want:%d have:%d]", size1, size2)
	}
	return nil
}
