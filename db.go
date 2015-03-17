package lib

import (
  "database/sql"
  "log"
  "fmt"
  _ "github.com/go-sql-driver/mysql"

)

type Row map[string][]uint8
type Result []Row

var (
  DBList map[string]*DB
  Conn  *sql.DB
<<<<<<< HEAD:db.go
  Result []Row
  LogFile string
=======
  CurDB *DB
>>>>>>> b1b8b555fbaa9dc4e94a347131dbd3f7c8b88649:db/db.go
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
    CurDB = val
    Conn = val.Link
    return true
  } else {
    return false
  }
}

<<<<<<< HEAD:db.go
func (db *DB)Query(sql string) (r bool){
=======
func Query(sql string) (r Result) {
  CurDB.LastSql = sql
>>>>>>> b1b8b555fbaa9dc4e94a347131dbd3f7c8b88649:db/db.go
  rows, err := Conn.Query(sql)
  if err != nil {
    return
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
    line := make(map[string][]uint8, count)
    for i, col := range columns {
      val := values[i]
      line[col] = val.([]byte)
    }
    r = append(r, line)
  }
  return r
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

