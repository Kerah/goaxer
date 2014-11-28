package imdg

import (
	"time"
	"math/rand"
	"fmt"
	"testing")

func new_test_db() IStorage{
	storage := NewStorage(uint32(3))
	storage.AddNode(uint32(0), ":6379")
	storage.AddNode(uint32(1), ":6379")
	storage.AddNode(uint32(2), ":6379")
	return storage
}

func new_test_mw() IInMemoryDataGrid {
	db := new_test_db()
	var ser JSONSerializer
	dg := New(ser, db)
	return dg
}

type testResultData struct {
	Key string `json:"key"`
	Id  int `json:"id"`
	Count int `json:"count"`
}

func getResult() (interface {}, error){
	time.Sleep(time.Millisecond*30)
	var data testResultData
	data.Count = rand.Int()
	data.Key = fmt.Sprint("", data.Count)
	data.Id = data.Count
	return data, nil
}

func TestMiddleware(t *testing.T){
	mw := new_test_mw()
	var data testResultData
	delay := 1*time.Second
	expire := 1*time.Second
	err := mw.Get("rooney", &data, delay, expire, getResult)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMiddlewareDoubleResponse(t *testing.T){
	mw := new_test_mw()
	data_first := testResultData{
		Key: "1",
		Count: 1,
		Id: 1,
	}

	data_second := testResultData{
		Key: "2",
		Count: 2,
		Id: 2,
	}
	data_third := testResultData{
		Key: "3",
		Count: 3,
		Id: 3,
	}
	delay := 1*time.Second
	expire := 60*time.Second
	go mw.Get("fergusson", &data_first, delay, expire, getResult)
	go mw.Get("fergusson", &data_second, delay, expire, getResult)
	time.Sleep(time.Millisecond*100)
	if data_first.Key != data_second.Key {
		t.Errorf("unequal keys in responsed data: %s/%s", data_first.Key, data_second.Key)
	}
	if data_first.Count != data_second.Count {
		t.Errorf("unequal counts in responsed data: %d %d", data_first.Count, data_second.Count)
	}

	if data_first.Id != data_second.Id {
		t.Errorf("unequal ids in responsed data: %d %d", data_first.Id, data_second.Id)
	}
	mw.Get("fergusson", &data_third, delay, expire, getResult)
	if data_first.Key != data_third.Key {
		t.Errorf("unequal keys in responsed data: %s %s", data_first.Key, data_third.Key)
	}
}

func BenchmarkMilldewareMultiKeys(b *testing.B){
	err_cnts := 0
	mw := new_test_mw()
	delay := 1*time.Second
	expire := 60*time.Second
	keys := []string{
		"rojo",
		"fergusson",
		"cleverley",
		"ferdinand",
		"vidic",
		"robson",
		"di maria",
		"valencia",
		"rafael",
		"shaw",
		"evra",
		"de gea",
		"rooney",
		"van persie",
		"falcao",
	}
	key_len := len(keys)
	fnc := func (key string){
		var data testResultData
		err := mw.Get(key, &data, delay, expire, getResult)
		if err != nil {
			err_cnts += 1
		}

	}
	for  i:= 0; i < b.N; i++ {
		go fnc(keys[i%key_len])
	}

}
