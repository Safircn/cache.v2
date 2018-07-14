package cache

import (
  "testing"
  "time"
  "fmt"
  "strconv"
)

func TestRun(t *testing.T) {
  cache := NewLruCache(1000)

  cache.Add("1", 1, time.Minute)
  cache.Add("2", 1, time.Minute)
  cache.Add("3", 1, time.Minute)
  cache.Add("4", 1, time.Minute)
  cache.Add("5", 1, time.Minute)
  cache.Add("6", 1, time.Minute)
  cache.Add("7", 1, time.Minute)
  cache.Add("8", 1, time.Minute)
  cache.Add("9", 1, time.Minute)
  cache.Add("10", 1, time.Minute)

  cache.Add("11", 1, time.Minute)
  cache.Add("12", 1, time.Minute)

  cache.Delete("6")
  cache.Delete("8")
  cache.Delete("4")
  fmt.Printf("%+v\n", cache.(*LruCache))
  cache.Set("5", 11)
  //cache.Set("3",11)
  fmt.Printf("%+v\n", cache.(*LruCache))
  cache.Set("3", 11)

  cache.Add("15", 1, time.Minute)
  cache.Add("16", 1, time.Minute)
  cache.Add("17", 1, time.Minute)
  cache.Add("18", 1, time.Minute)
  cache.Add("19", 1, time.Minute)
  cache.Add("20", 1, time.Minute)
  cache.Add("20", 1, time.Minute)
  cache.Add("20", 1, time.Minute)
  cache.Add("20", 1, time.Minute)
  cache.Add("20", 1, time.Minute)
  cache.Add("20", 1, time.Minute)


  for i:=0;i<1000000;i++{
    cache.Add(strconv.Itoa(i+99), 1, time.Minute)
  }
  cache.Get("1000000")
  cache.Get("1000001")
  cache.Get("1000002")
  cache.Get("1000003")
  cache.Delete("1000003")
  cache.Flush()
  fmt.Printf("%+v\n", cache.(*LruCache))
  printEvictList(cache.(*LruCache).evictList,cache.(*LruCache).currentLen)
}

func printEvictList(evictList *EvictList, currentNum int) {
  if evictList != nil {
    num := 0
    for e := evictList; e != nil; e = e.next {
      //fmt.Printf("key:%s  pren:%p current:%p next:%p \n", e.key, e.prev, e, e.next)
      if e.next != nil {
        if e.next.prev != e {
          fmt.Println("!!!!!!!!!!!!!! next")
        }
      }
      if e.prev != nil {
        if e.prev.next != e {
          fmt.Println("!!!!!!!!!!!!!! prev")
        }
      }
      num++
    }
    if num != currentNum {
      fmt.Println("!!!!!!!!!!!!!! currentNum")
    }
  }
}
