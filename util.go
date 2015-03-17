package lib

import (
    "bufio"
    "io/ioutil"
    "os"
    )

func IsExist(path string) bool {
  _, err := os.Stat(path)
  if err == nil {
    return true
  }

  if os.IsNotExist(err) {
    return false
  }

  return false
}

func IsNotExist(path string) bool {
  return (!IsExist(path))
}

func FileGetContents(fn string) ([]byte, error) {
  fp, err := os.OpenFile(fn, os.O_RDONLY, os.ModePerm)
  if err != nil {
    return nil, err
  }
  defer fp.Close()
  reader := bufio.NewReader(fp)
  contents, _ := ioutil.ReadAll(reader)
  return contents, nil
}

func FilePutContents(fn string, con []byte) error {
  fp, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, os.ModePerm)
  if err != nil {
    return err
  }
  defer fp.Close()
  _, err = fp.Write(con)
  return err
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
