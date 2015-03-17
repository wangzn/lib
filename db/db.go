package DB

import (
  "database/sql"
  "log"
  "fmt"

  _ "github.com/go-sql-driver/mysql"

)

var (
  DBList map[string]*DB
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
  defer db.Close()
  var d DB
  d = DB{db, addrs, ""}
  DBList[name] = &d
  return true
}

func ListDB() {
  for i, j := range DBList {
    fmt.Println(i)
    fmt.Println(j)
  }
}

