package imdg

import (
	redis "gopkg.in/redis.v1"
	"hash/fnv"
	"errors"
	"time"
)

type RedisStorage struct {
	nodes map[uint32]*redis.Client
	nodes_cnt uint32
}

func NewStorage(num_nodes uint32) *RedisStorage {
	return &RedisStorage{
		nodes: make(map[uint32]*redis.Client),
		nodes_cnt: num_nodes,
	}
}

func (rds *RedisStorage) hex(key string) uint32 {
	hex := fnv.New32()
	hex.Write([]byte(key))
	digest := hex.Sum32()
	return digest
}

func (rds *RedisStorage) getNodeId(key string) uint32 {
	digest := rds.hex(key)
	ident := digest%rds.nodes_cnt
	return ident
}

func (rds *RedisStorage) GetNode(key string) (*redis.Client, error) {
	ident := rds.getNodeId(key)
	if node, ok := rds.nodes[ident]; ok {
		return node, nil
	}
	return nil, errors.New("Node for this key not found")
}

func (rds *RedisStorage) AddNode(id uint32, url string){
	opt := redis.Options{
		Addr: url,
	}
	cli := redis.NewTCPClient(&opt)
	rds.nodes[id] = cli
}

func (rds *RedisStorage) Get(key string) ([]byte, error) {
	node, err := rds.GetNode(key)
	if err != nil {
		return nil, err
	}
	data, err := node.Get(key).Result()
	if err != nil {
		return nil, err
	}
	return []byte(data), nil
}

func (rds *RedisStorage) Set(key string, data []byte, expire time.Duration) (error){
	node, err := rds.GetNode(key)
	if err != nil {
		return err
	}
	return node.SetEx(key, expire, string(data)).Err()
}
