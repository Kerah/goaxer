package imdg

import "time"

type IGetter func() (interface {}, error)

type IInMemoryDataGrid interface {
	Get(key string, result interface {}, executionLimit, expire time.Duration, getter IGetter) (error)
}

type IStorage interface {
	Get(string) ([]byte, error)
	Set(key string, data []byte, expire time.Duration) (error)
	AddNode(id uint32, url string)
}

type ISerializer interface {
	Encode(value interface {}) ([]byte, error)
	Decode(data []byte) (map[string]interface {}, error)
}
