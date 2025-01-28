package resp

type Connection interface {
	//回复消息
	Write([]byte) error
	//得到Redis中DB的index
	GetDBIndex() int
	//切换Redis中指定DB
	SelectDB(int)
}
