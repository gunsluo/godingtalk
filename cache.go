package dingtalk

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"reflect"
	"sync"
	"time"
)

var (
	ErrAlreadyExpired = errors.New("Data already expired")
)

type Expirable interface {
	IsExpired() bool
}

type Cache interface {
	Set(data Expirable) error
	Get(data Expirable) error
}

type StringExpirable struct {
	Value     string `json:"value"`
	ExpiresIn int    `json:"expires_in"`
	Created   int64  `json:"created"`
}

func NewStringExpirable(value string, expiresIn int) *StringExpirable {
	return &StringExpirable{
		Value:     value,
		ExpiresIn: expiresIn | 7200,
		Created:   time.Now().Unix(),
	}
}

func (e *StringExpirable) IsExpired() bool {
	return time.Now().Unix() > e.Created+int64(e.ExpiresIn-60)
}

type FileCache struct {
	Path   string
	data   Expirable
	locker sync.RWMutex
}

func NewFileCache(path string) *FileCache {
	return &FileCache{
		Path: path,
	}
}

func (f *FileCache) Set(data Expirable) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	f.locker.Lock()
	defer f.locker.Unlock()
	if err := ioutil.WriteFile(f.Path, bytes, 0644); err != nil {
		return err
	}
	f.data = data

	return nil
}

func (f *FileCache) Get(data Expirable) error {
	f.locker.RLock()
	defer f.locker.RUnlock()
	if f.data != nil {
		if !f.data.IsExpired() {
			v := reflect.ValueOf(f.data).Elem()
			reflect.ValueOf(data).Elem().Set(v)
			return nil
		}
	}

	bytes, err := ioutil.ReadFile(f.Path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, data)
	if err != nil {
		return err
	}

	if data.IsExpired() {
		return ErrAlreadyExpired
	}

	return nil
}
