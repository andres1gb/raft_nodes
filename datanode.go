package raftnodes

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"sync"
)

type dataNode[DataType any] struct {
	mu   sync.Mutex
	data map[string]DataType
}

type dataItem[DataType any] struct {
	Key   string
	Value DataType
}

func newDataNode[DataType any]() *dataNode[DataType] {
	return &dataNode[DataType]{
		data: make(map[string]DataType),
	}
}

func (dn *dataNode[DataType]) Apply(data []byte) {
	var di dataItem[DataType]
	if err := json.Unmarshal(data, &di); err != nil {
		log.Println("wrong format for input", err)
		return
	}
	dn.mu.Lock()
	defer dn.mu.Unlock()
	dn.data[di.Key] = di.Value
}

func (dn *dataNode[DataType]) Snapshot() (io.ReadCloser, error) {
	dn.mu.Lock()
	defer dn.mu.Unlock()
	data, err := json.Marshal(&dn.data)
	if err != nil {
		return nil, err
	}
	return io.NopCloser(strings.NewReader(string(data))), nil
}

func (dn *dataNode[DataType]) Restore(reader io.ReadCloser) error {
	dn.mu.Lock()
	defer dn.mu.Unlock()
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &dn.data)
	if err != nil {
		return err
	}
	return reader.Close()
}
