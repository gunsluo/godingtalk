package dingtalk

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"sync"
	"time"
)

var (
	ErrAlreadyExpired = errors.New("Data already expired")
)

const (
	KeySuiteAccessToken = "suite_access_token"
	KeySuiteTicket      = "suite_ticket"
)

type Cache interface {
	Set(ctx context.Context, data *KVExpirable) error
	Get(ctx context.Context, key string) (*KVExpirable, error)
}

type Persist = Cache

type KVExpirable struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	ExpiresIn int    `json:"expires_in"`
	Created   int64  `json:"created"`
}

func NewKVExpirable(key, value string, expiresIn int) *KVExpirable {
	return &KVExpirable{
		Key:       key,
		Value:     value,
		ExpiresIn: expiresIn | 7200,
		Created:   time.Now().Unix(),
	}
}

func (e *KVExpirable) IsExpired() bool {
	return time.Now().Unix() > e.Created+int64(e.ExpiresIn-60)
}

type FileCache struct {
	Path   string
	db     map[string]*KVExpirable
	locker sync.RWMutex
}

func NewFileCache(path string) *FileCache {
	return &FileCache{
		Path: path,
		db:   make(map[string]*KVExpirable),
	}
}

func (f *FileCache) Set(ctx context.Context, data *KVExpirable) error {
	f.locker.Lock()
	defer f.locker.Unlock()

	// read all data from file
	oBytes, err := ioutil.ReadFile(f.Path)
	var o []*KVExpirable
	if err != nil {
		o = append(o, data)
		f.db[data.Key] = data
	} else {
		if err := json.Unmarshal(oBytes, &o); err != nil {
			return err
		}

		idx := -1
		for i, item := range o {
			if item.Key == data.Key {
				idx = i
				o[i] = data
			}
			f.db[item.Key] = o[i]
		}

		if idx == -1 {
			o = append(o, data)
			f.db[data.Key] = data
		}
	}

	bytes, err := json.Marshal(o)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(f.Path, bytes, 0644); err != nil {
		return err
	}

	return nil
}

func (f *FileCache) Get(ctx context.Context, key string) (*KVExpirable, error) {
	f.locker.RLock()
	defer f.locker.RUnlock()

	if v, ok := f.db[key]; ok {
		if !v.IsExpired() {
			return v, nil
		}
	}

	bytes, err := ioutil.ReadFile(f.Path)
	if err != nil {
		return nil, err
	}

	var o []*KVExpirable
	err = json.Unmarshal(bytes, &o)
	if err != nil {
		return nil, err
	}

	idx := -1
	for i, item := range o {
		if item.Key == key {
			idx = i
			break
		}
	}

	if idx == -1 {
		return nil, ErrAlreadyExpired
	}
	data := o[idx]

	if data.IsExpired() {
		return data, ErrAlreadyExpired
	}

	return data, nil
}
