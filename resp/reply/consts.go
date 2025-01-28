package reply

type PongReply struct {
}

var pongbytes = []byte("+PONG\r\n")

func (r *PongReply) ToBytes() []byte {
	return pongbytes
}

func MakePongReply() *PongReply {
	return &PongReply{}
}

type OkReply struct {
}

var okbytes = []byte("+OK\r\n")

func (r *OkReply) ToBytes() []byte {
	return okbytes
}

var theOkReply = new(OkReply)

func MakeOkReply() *OkReply {
	return theOkReply
}

type NullBulkReply struct {
}

var nullbulkbytes = []byte("$-1\r\n")

func (n *NullBulkReply) ToBytes() []byte {
	return nullbulkbytes
}

func MakeNullBulkReply() *NullBulkReply {
	return &NullBulkReply{}
}

type EmptyMultiBulkReply struct {
}

var emptymultibulkbytes = []byte("*0\r\n")

func (e *EmptyMultiBulkReply) ToBytes() []byte {
	return emptymultibulkbytes
}

type NoReply struct {
}

var nobytes = []byte("")

func (n *NoReply) ToBytes() []byte {
	return nobytes
}
