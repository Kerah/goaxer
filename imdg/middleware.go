package imdg

import (
	"time"
	"errors"
	"fmt"
	"encoding/json"
)

type fetching struct {
	Result    interface{}
	Getter    IGetter
	Expire    time.Duration
	ISFetcher bool
	Response  chan bool
	IsTimeout bool
	Fail      error

	Next   *fetching
	Prev   *fetching

}


type DataGrid struct {
	serializer ISerializer
	storage    IStorage
	fetching   map[string]*fetching
}

func New(serializer ISerializer, storage IStorage) *DataGrid {
	dg := DataGrid{
		serializer: serializer,
		storage: storage,
		fetching: make(map[string]*fetching),
	}

	return &dg
}

func (dg *DataGrid) Get(key string, result interface{}, executionLimit, expire time.Duration, getter IGetter) (error) {
	fetcher := dg.getFetcher(key, result, expire, getter)
	if fetcher.ISFetcher {
		go dg.fetch(key, fetcher, executionLimit)
	}
	select {
	case _ = <-fetcher.Response:
		if fetcher.Fail != nil {
			return fetcher.Fail
		}
	case <-time.After(executionLimit):
		fetcher.IsTimeout = true
		fetcher.Fail = errors.New("execution timeout")
		dg.clearFetcher(fetcher)
		return fetcher.Fail
	}
	return nil
}

func (dg *DataGrid) failed(err error, root *fetching) {
	row := root
	for row != nil {
		row.Fail = err
		row.Prev = nil
		row = root.Next
	}
}

func (dg *DataGrid) response(data []byte, root *fetching) {
	row := root
	for row != nil {
		err_dec := json.Unmarshal(data, row.Result)
		if err_dec != nil {
			row.Fail = err_dec
		}
		row.Prev = nil
		row = row.Next
	}
}

func (dg *DataGrid) fetch(key string, fetcher *fetching, executionLimit time.Duration) {
	data, err := dg.storage.Get(key)

	if err != nil {
		res, err := fetcher.Getter()
		if err != nil {
			fetcher.Fail = err
			dg.failed(err, fetcher)

		}else {
			data, err := json.Marshal(res)
			if err != nil {
				dg.failed(err, fetcher)
			}else {
				dg.response(data, fetcher)
				if err == nil {
					go dg.storage.Set(key, data, fetcher.Expire)
				}else {
					fmt.Printf("Error encoding to json\n")
				}
			}
		}


	} else {
		dg.response(data, fetcher)
	}
	delete(dg.fetching, key)
	fetcher.Response <- true
}

func (dg *DataGrid) clearFetcher(fetcher *fetching) {
	if fetcher.ISFetcher {

	} else {
		nxt := fetcher.Next
		prev := fetcher.Prev
		if prev != nil {
			prev.Next = nxt
		}
		if nxt != nil {
			nxt.Prev = prev
		}
		fetcher.Next = nil
		fetcher.Prev = nil
	}
}


func (dg *DataGrid) getFetcher(key string, result interface{}, expire time.Duration, getter IGetter) *fetching {
	fetcher := &fetching{
		Result: result,
		Expire: expire,
		Getter: getter,
		Fail:   nil,
	}

	if prev, ok := dg.fetching[key]; ok {
		prev.Next = fetcher
		fetcher.Prev = prev
		fetcher.Response = prev.Response
	} else {
		fetcher.Response = make(chan bool)
		fetcher.ISFetcher = true
	}
	dg.fetching[key] = fetcher
	return fetcher
}

