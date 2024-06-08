package pktline

// PacketType 表示包类型
type PacketType int

// 定义包类型
const (
	TypeUndefine  PacketType = 500
	TypeFlush     PacketType = 0
	TypeUndefine1 PacketType = 1
	TypeUndefine2 PacketType = 2
	TypeUndefine3 PacketType = 3
	TypeData      PacketType = 200
)

// Packet 表示一个 pktline 数据包
type Packet struct {
	Type PacketType
	Head string
	Body []byte
}
