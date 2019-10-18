package env

import (
	"os"

	"github.com/payfazz/go-apt/pkg/fazzkv"
)

type env struct {
	base map[string]string
}

func (e *env) Set(key string, value interface{}) error {
	e.base[key] = value.(string)
	return nil
}

func (e *env) Get(key string) (string, error) {
	r := os.Getenv(key)
	if r == "" {
		if _, ok := e.base[key]; !ok {
			return "", nil
		}
		r = e.base[key]
	}
	return r, nil
}

func (e *env) Delete(key string) error {
	temp := make(map[string]string)
	for k, v := range e.base {
		if v == key {
			continue
		}
		temp[k] = v
	}
	e.base = temp
	return nil
}

func (e *env) Truncate() error {
	e.base = make(map[string]string)
	return nil
}

func NewFazzEnv() fazzkv.Store {
	return &env{
		base: make(map[string]string),
	}
}
