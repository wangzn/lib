package util

import (
    "bufio"
    "io/ioutil"
    "os"
    "net/http"
    "strings"

    )

func FExist(path string) bool {
  _, err := os.Stat(path)
  if err == nil {
    return true
  }

  if os.IsNotExist(err) {
    return false
  }
  return false
}

func FNExist(path string) bool {
  return (!FExist(path))
}

func HRun(method, url, reqBody string) ([]byte, error) {
  var req *http.Request
  var err error
  client := &http.Client{}
  if len(reqBody) > 0 {
    req, err = http.NewRequest(method, url, strings.NewReader(reqBody))
    req.ContentLength = int64(len(reqBody))
  } else {
    req, err = http.NewRequest(method, url, nil)
  }
  resp, err := client.Do(req)
  if err != nil {
    return []byte(""), err
  }
  defer resp.Body.Close()
  contents, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return []byte(""), err
  }
  return contents, err
}

func HGet(url string) ([]byte, error) {
  resp, err := http.Get(url)
  if err != nil {
    return []byte(""), err
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return []byte(""), err
  }
  return body, err
}

func HPost(url, reqBody, bodyType string) ([]byte, error) {
  if bodyType == "" {
    bodyType = "application/x-www-form-urlencoded"
  }
  resp, err := http.Post(url, bodyType, strings.NewReader(reqBody))
  if err != nil {
    return []byte(""), err
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return []byte(""), err
  }
  return body, nil
}

func HPut(url, reqBody string) ([]byte, error) {
  return HRun("PUT", url, reqBody)
}

func HDelete(url, reqBody string) ([]byte, error) {
  return HRun("DELETE", url, reqBody)
}

func FileGetContents(fn string) (contents []byte, err error) {
  if strings.HasPrefix(fn, "http://") {
    return HGet(fn)
  } else {
    fp, err := os.OpenFile(fn, os.O_RDONLY, os.ModePerm)
    if err != nil {
      return []byte(""), err
    }
    defer fp.Close()
    reader := bufio.NewReader(fp)
    con, _ := ioutil.ReadAll(reader)
    contents = con
  }
  return
}

func FilePutContents(fn string, con []byte) ([]byte, int, error) {
  var size int
  if strings.HasPrefix(fn, "http://") {
    body, err := HPost(fn, string(con), "")
    return body, len(body), err
  }
  fp, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, os.ModePerm)
  if err != nil {
    return []byte(""), 0, err
  }
  defer fp.Close()
  size, err = fp.Write(con)
  return []byte(""), size, err
}

func FileAppendContents(fn string, con []byte) error {
  fp, err := os.OpenFile(fn, os.O_APPEND|os.O_CREATE, os.ModePerm)
  if err != nil {
    return err
  }
  defer fp.Close()
  _, err = fp.Write(con)
  return err
}
