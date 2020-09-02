package cache

import (
	"sync"
	"time"
)

type (
	LruCache struct {
		rwMutex   *sync.RWMutex
		cacheList map[string]*LruCacheRow

		evictList    *EvictList
		TotalLen     int
		currentLen   int
		endEvictList *EvictList
	}
	LruCacheRow struct {
		data interface{}

		expiredTime *time.Time
		evictList   *EvictList
	}
	EvictList struct {
		key  string
		prev *EvictList
		next *EvictList
	}
)

func NewLruCache(total int) Cache {
	if total < 10 {
		total = 10
	}
	newCache := &LruCache{
		cacheList: make(map[string]*LruCacheRow),
		rwMutex:   new(sync.RWMutex),
		TotalLen:  total,
	}
	caches = append(caches, newCache)
	return newCache
}

func (this *LruCache) Add(key string, val interface{}, timeout time.Duration) {
	expiredTime := time.Now().Add(timeout)
	this.rwMutex.Lock()
	if this.evictList == nil {
		this.evictList = &EvictList{
			key: key,
		}
		this.endEvictList = this.evictList
		this.currentLen++
	} else {
		if v, ok := this.cacheList[key]; ok {
			this.top(v.evictList)
			this.cacheList[key].data = val
			this.cacheList[key].expiredTime = &expiredTime
			this.rwMutex.Unlock()
			return
		}
		if this.TotalLen <= this.currentLen {
			this.delete(this.endEvictList.key)
		}
		this.evictList.prev = &EvictList{
			key:  key,
			next: this.evictList,
		}
		this.evictList = this.evictList.prev
		this.evictList.prev = nil
		this.currentLen++
	}
	this.cacheList[key] = &LruCacheRow{
		expiredTime: &expiredTime,
		data:        val,
		evictList:   this.evictList,
	}
	this.rwMutex.Unlock()
}

func (this *LruCache) top(evictList *EvictList) {
	if evictList == this.evictList {
		return
	}
	evictList.prev.next = evictList.next
	if evictList.next != nil {
		evictList.next.prev = evictList.prev
	} else {
		this.endEvictList = evictList.prev
	}
	this.evictList.prev = evictList
	evictList.prev = nil
	evictList.next = this.evictList
	this.evictList = evictList
}
func (this *LruCache) Set(key string, val interface{}) {
	this.rwMutex.Lock()
	if _, ok := this.cacheList[key]; ok {
		this.top(this.cacheList[key].evictList)
		this.cacheList[key].data = val
	}
	this.rwMutex.Unlock()
}

func (this *LruCache) delete(key string) {
	if v, flag := this.cacheList[key]; flag {

		if v.evictList.next == nil {
			this.endEvictList = v.evictList.prev
		} else {
			v.evictList.next.prev = v.evictList.prev
		}
		if v.evictList.prev == nil {
			this.evictList = v.evictList.next
		} else {
			v.evictList.prev.next = v.evictList.next
		}
		delete(this.cacheList, key)
		this.currentLen--
	}
}

func (this *LruCache) Delete(key string) {
	this.rwMutex.Lock()
	this.delete(key)
	this.rwMutex.Unlock()
}

func (this *LruCache) Get(key string) (interface{}, bool) {
	var flagMain bool
	var data interface{}
	this.rwMutex.RLock()
	if v, flag := this.cacheList[key]; flag {
		if v.expiredTime.After(time.Now()) {
			data = v.data
			flagMain = true
		}
	}
	this.rwMutex.RUnlock()
	if flagMain {
		this.rwMutex.Lock()
		this.top(this.cacheList[key].evictList)
		this.rwMutex.Unlock()
	}
	return data, flagMain
}

func (this *LruCache) IsExist(key string) bool {
	var flagMain bool
	this.rwMutex.RLock()
	if v, flag := this.cacheList[key]; flag {
		if v.expiredTime.After(time.Now()) {
			flagMain = true
		}
	}
	this.rwMutex.RUnlock()
	if flagMain {
		this.rwMutex.Lock()
		this.top(this.cacheList[key].evictList)
		this.rwMutex.Unlock()
	}
	return flagMain
}

func (this *LruCache) Flush() {
	this.rwMutex.Lock()
	this.cacheList = make(map[string]*LruCacheRow)
	this.endEvictList = nil
	this.evictList = nil
	this.currentLen = 0
	this.rwMutex.Unlock()
}

func (this *LruCache) GC() {
	this.rwMutex.Lock()
	for kk, vv := range this.cacheList {
		if vv.expiredTime.Sub(time.Now()) < -1 {
			this.delete(kk)
		}
	}
	this.rwMutex.Unlock()
}
