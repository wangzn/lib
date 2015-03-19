package cache

import (
  "wangzn/lib/util"
  "time"
  "fmt"
  "reflect"
)

const (
  URL   = iota
  FUNC  = iota
)

type CacheEntry struct {
  name        string
  t           int
  fn          interface{}
  param       []interface{}
  url         string
  value       []byte
  period      int64
  updateTime  int64
}

type CacheData map[string]*CacheEntry

var (
  data   CacheData
)

func Init() {
  data    = make(CacheData)
  go bgLoad()
}

func needUpdate(t int64, period int64) bool{
  now := time.Now().Unix()
  fmt.Println(now, t, period)
  if (now - t) >= period {
    return true
  }
  return false
}

func doUpdate(name string) bool {
  entry, ok := data[name]
  if ok {
    switch entry.t {
      case URL :
        v, err := util.HGet(entry.url)
        if err == nil {
          entry.value = v
          entry.updateTime = time.Now().Unix()
        }
      case FUNC :
        re := reflect.ValueOf(entry.fn)
        in := make([]reflect.Value, len(entry.param))
        for k, p := range entry.param {
          in[k] = reflect.ValueOf(p)
        }
        v := re.Call(in)
        if len(v) == 0 {
          return false
        }
        entry.value = v[0].Bytes()
        entry.updateTime = time.Now().Unix()
    }
    return true
  }
  return false
}

func bgLoad() {
  for {
    for k, v := range data {
      needUpdate := needUpdate(v.updateTime, v.period)
      if needUpdate {
        fmt.Println("do update ", k)
        go doUpdate(k)
      }
    }
    time.Sleep(time.Second * 1)
  }
}

func Get(name string) ([]byte, bool) {
  val, ok := data[name]
  if ok {
    if val.updateTime == 0 {
      doUpdate(name)
    }
    return val.value, ok
  } else {
    return []byte(""), ok
  }
}

func SetByFunc(name string, fn interface{}, period int64, param ... interface{}) (err error) {
  var entry CacheEntry
  entry.name = name
  entry.t = FUNC
  entry.fn   = fn
  entry.param = param
  entry.url  = ""
  entry.value = []byte("")
  entry.updateTime = 0
  entry.period = period
  data[name] = &entry
  return
}

func SetByUrl(name, url string, period int64) (err error) {
  var entry CacheEntry
  entry.name = name
  entry.t = URL
  entry.url  = url
  entry.value = []byte("")
  entry.updateTime = 0
  entry.period = period
  data[name] = &entry
  return
}

func Expire(name string) bool {
  v, ok := data[name]
  if ok {
    v.updateTime = 0
  }
  return ok
}

func UnSet(name string) bool {
  _, ok := data[name]
  if ok {
    delete(data, name)
    _, ok = data[name]
    return !ok
  }
  return ok
}

