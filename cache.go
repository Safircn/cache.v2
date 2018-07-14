package cache

import (
  "time"
  "sync"
)

func init() {
  caches = make([]GC, 0)
}

var (
  caches []GC
  _      Cache = (*CacheBox)(nil)
)

type GC interface {
  GC()
}

type Cache interface {
  Add(key string, val interface{}, timeout time.Duration)
  Set(key string, val interface{})
  Get(key string) (interface{}, bool)
  Delete(key string)
  IsExist(key string) bool
  Flush()
}

type (
  CacheBox struct {
    rwMutex   *sync.RWMutex
    cacheList map[string]*cacheRow
  }
  cacheRow struct {
    expiredTime *time.Time
    data        interface{}
  }
)

func (this *CacheBox) Add(key string, val interface{}, timeout time.Duration) {
  if (timeout == 0) {
    return
  }
  this.rwMutex.Lock()
  expiredTime := time.Now().Add(timeout)
  this.cacheList[key] = &cacheRow{
    expiredTime: &expiredTime,
    data:        val,
  }
  this.rwMutex.Unlock()
}

func (this *CacheBox) Set(key string, val interface{}) {
  this.rwMutex.Lock()
  if cacheData, ok := this.cacheList[key]; ok {
    this.cacheList[key] = &cacheRow{
      expiredTime: cacheData.expiredTime,
      data:        val,
    }
  }
  this.rwMutex.Unlock()
}

func (this *CacheBox) Delete(key string) {
  this.rwMutex.Lock()
  if _, flag := this.cacheList[key]; flag {
    delete(this.cacheList, key)
  }
  this.rwMutex.Unlock()
}

func (this *CacheBox) Get(key string) (interface{}, bool) {
  this.rwMutex.RLock()
  defer this.rwMutex.RUnlock()
  if v, flag := this.cacheList[key]; flag {
    if v.expiredTime.After(time.Now()) {
      return v.data, true
    }
  }
  return nil, false
}

func (this *CacheBox) IsExist(key string) bool {
  this.rwMutex.RLock()
  defer this.rwMutex.RUnlock()
  if v, flag := this.cacheList[key]; flag {
    if v.expiredTime.After(time.Now()) {
      return true
    }
  }
  return false
}

func (this *CacheBox) Flush() {
  this.rwMutex.Lock()
  this.cacheList = make(map[string]*cacheRow)
  this.rwMutex.Unlock()
}

func NewCache() Cache {
  newCache := &CacheBox{
    cacheList: make(map[string]*cacheRow),
    rwMutex:   new(sync.RWMutex),
  }
  caches = append(caches, newCache)
  return newCache
}

func Run(Interval int) {
  if len(caches) == 0 {
    return
  }
  if Interval == 0 {
    Interval = 1
  }

  for {
    time.Sleep(time.Minute * time.Duration(Interval))
    Gc()
  }
}

func (this *CacheBox) GC() {
  this.rwMutex.Lock()
  for kk, vv := range this.cacheList {
    if vv.expiredTime.Sub(time.Now()) < -1 {
      delete(this.cacheList, kk)
    }
  }
  this.rwMutex.Unlock()
}

func Gc() {
  for k := range caches {
    caches[k].GC()
  }
}

func RunTntervalTenM() {
  Run(10)
}
