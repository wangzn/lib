package DB

import (
  "database/sql"
  "database/sql/driver"
  "errors"
  "log"
  "fmt"

  _ "github.com/go-sql-driver/mysql"

)

var (
  DBList map[string]*DB
)

type DB struct {
  Link      driver.Conn
  Addrs     string
  LastSql   string
}

func init() {
  DBList = make(map[string]*DB)
}

func Set(params ...string) (r bool, err error) {
  var addrs, name, charset string
  if len(params) < 2 || len(params) > 3 {
    return FALSE
  }
  addrs = params[0]
  name = params[1]
  if len(params) == 3 {
    charset = params[2]
  }
  db, err := sql.Open("mysql", addrs)
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()
  var d DB
  d = {db, addrs, ""}
  DBList[name] = &d

}

func ListDB() {
  for i, j in range DBlist {
    fmt.Println(i)
  }
}

