package tools

import (
  "github.com/Safircn/cache.v2"
  "time"
)

type (
  CacheData struct {
    New      func(key string) interface{}
    cacheBox cache.Cache
    timeout  time.Duration
  }

  CacheDataI interface {
    Get(key string) (interface{})
    Delete(key string)
    IsExist(key string) bool
    Flush()
  }
)

var (
  _ CacheDataI = (*CacheData)(nil)
)

func NewCacheData(newFunc func(string) interface{}, timeout time.Duration) CacheDataI {
  return &CacheData{
    New:      newFunc,
    cacheBox: cache.NewCache(),
    timeout:  timeout,
  }
}

func (cd *CacheData) Get(key string) (interface{}) {
  if dObj, flag := cd.cacheBox.Get(key); flag {
    return dObj
  }
  dObj := cd.New(key)
  cd.cacheBox.Add(key, dObj, cd.timeout)
  return dObj
}

func (cd *CacheData) Delete(key string) {
  cd.cacheBox.Delete(key)
}

func (cd *CacheData) IsExist(key string) bool {
  return cd.cacheBox.IsExist(key)
}

func (cd *CacheData) Flush() {
  cd.cacheBox.Flush()
}
