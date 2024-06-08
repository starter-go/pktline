package pktline

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"
)

func TestEncodeDecode(t *testing.T) {

	list1 := make([]*Packet, 0)

	list1 = append(list1, &Packet{Type: TypeFlush})
	list1 = append(list1, &Packet{Type: TypeUndefine1})
	list1 = append(list1, &Packet{Type: TypeUndefine2})
	list1 = append(list1, &Packet{Type: TypeUndefine3})
	list1 = append(list1, &Packet{})
	list1 = append(list1, &Packet{Head: "this is a head-only packet"})
	list1 = append(list1, &Packet{Body: []byte("this is a body-only packet")})
	list1 = append(list1, &Packet{})
	list1 = append(list1, &Packet{Head: "h", Body: []byte("b")})

	buffer1 := bytes.NewBuffer(nil)
	enc := NewEncoder(buffer1)
	for _, p := range list1 {
		enc.Write(p)
	}
	enc.Close()

	raw := buffer1.Bytes()
	buffer2 := bytes.NewBuffer(raw)
	dec := NewDecoder(buffer2)
	for {
		p, err := dec.Read()
		if err != nil {
			if err != io.EOF {
				t.Error(err)
			}
			break
		}

		b, err := json.Marshal(p)
		str := string(b)
		t.Logf("packet: %s", str)
	}

}
