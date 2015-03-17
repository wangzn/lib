package lib

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
  LogFile string
)

type DB struct {
  Name      string
  Link      *sql.DB
  Addrs     string
  LastSql   string
}

func (db *DB)init() {
  DBList = make(map[string]*DB)
}

func (db *DB)Set(addrs, name string) (r bool) {
  db, err := sql.Open("mysql", addrs)
  if err != nil {
    log.Fatal(err)
    return false
  }
  var d DB
  d = DB{name, db, addrs, ""}
  DBList[name] = &d
  return true
}

func (db *DB)Active(name string) (r bool) {
  if val, ok := DBList[name]; ok {
    Conn = val.Link
    return true
  } else {
    return false
  }
}

func (db *DB)Query(sql string) (r bool){
  rows, err := Conn.Query(sql)
  if err != nil {
    return false
  }
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

func (db *DB)Current() {
  return CurDB.name
}

func (db *DB)LogQuery(f string) {
  LogFile = f
}

func (db *DB)ListDB() {
  for i, j := range DBList {
    fmt.Println(i)
    fmt.Println(j)
  }
}

