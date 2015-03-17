package DB

import (
  "database/sql"
  "log"
  "fmt"

  _ "github.com/go-sql-driver/mysql"

)

type Row []string

var (
  DBList map[string]*DB
  Conn  *sql.DB
  Result []Row
)

type DB struct {
  Link      *sql.DB
  Addrs     string
  LastSql   string
}

func init() {
  DBList = make(map[string]*DB)
}

func Set(addrs, name string) (r bool) {
  db, err := sql.Open("mysql", addrs)
  if err != nil {
    log.Fatal(err)
    return false
  }
  var d DB
  d = DB{db, addrs, ""}
  DBList[name] = &d
  return true
}

func Active(name string) (r bool) {
  if val, ok := DBList[name]; ok {
    Conn = val.Link
    return true
  } else {
    return false
  }
}

func Query(sql string)  {
  rows, _ := Conn.Query(sql)
  columns, _ := rows.Columns()
  count := len(columns)
  values := make([]interface{}, count)
  valuePtrs := make([]interface{}, count)

  for rows.Next() {
    for i, _ := range columns {
      valuePtrs[i] = &values[i]
    }
    rows.Scan(valuePtrs ...)
    for i, col := range columns {
      var v interface{}
      val := values[i]
      b, ok := val.([]byte)
      if (ok) {
        v = string(b)
      } else {
        v = val
      }
      fmt.Println(col, v)
    }
  }
}

func ListDB() {
  for i, j := range DBList {
    fmt.Println(i)
    fmt.Println(j)
  }
}

