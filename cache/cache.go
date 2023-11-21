package cache

import (
	"errors"
	"sync"
)

type CacheService interface {
	Set(string, any)
	Get(string) (any, error)
	Invalidate(string)
}

type CacheServiceImp struct{
	data sync.Map
}

var (
	ErrDataNotInCache error = errors.New("data not in cache")
)

func NewCacheService() CacheService {
	return &CacheServiceImp{}
}

func (l *CacheServiceImp)Set(key string, data any) {
	l.data.Store(key, data)
}

func (l *CacheServiceImp)Get(key string) (any, error) {
	v, ok := l.data.Load(key)
	if !ok {
		return nil, ErrDataNotInCache
	}
	return v, nil
}

func (l *CacheServiceImp)Invalidate(key string) {
	l.data.Delete(key)
}