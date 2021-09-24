package models

import "github.com/vmihailenco/msgpack/v5"

type Tml struct {
	ID    string // 消息ID
	Close bool   // 关闭链接
	Start bool   // 启动链接
	Data  []byte // 消息体
}

func (t *Tml) ToBytes() []byte {
	marshal, err := msgpack.Marshal(t)
	if err != nil {
		panic(err)
	}

	return marshal
}

func (t *Tml) FromBytes(by []byte) error {
	return msgpack.Unmarshal(by, t)
}
