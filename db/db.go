package DB

import (
  "database/sql"
  "log"
  "fmt"
  "reflect"
  "net/url"
  "strings"

  _ "github.com/go-sql-driver/mysql"

)

type Row map[string][]uint8
type Result []Row
type AssocData map[string]interface{}

var (
  DBList map[string]*DB
  Conn  *sql.DB
  LogFile string
  CurDB *DB
)

type DB struct {
  Name      string
  Link      *sql.DB
  Addrs     string
  LastSql   string
}

func init() {
  DBList = make(map[string]*DB)
}

func Set(addrs, name string) (r bool) {
  dbLink, err := sql.Open("mysql", addrs)
  if err != nil {
    log.Fatal(err)
    return false
  }
  var d DB
  d = DB{name, dbLink, addrs, ""}
  DBList[name] = &d
  return true
}

func Active(name string) (r bool) {
  if val, ok := DBList[name]; ok {
    CurDB = val
    Conn = val.Link
    return true
  } else {
    return false
  }
}

func Query(sql string) (r Result){
  CurDB.LastSql = sql
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

func Exec(sql string) (r bool, lastId int64, affectedRows int64) {
  r = true
  CurDB.LastSql = sql
  res, err := Conn.Exec(sql)
  if err != nil {
    return false, 0, 0
  }
  affectedRows, err = res.RowsAffected()
  if err != nil {
    return r, 0, 0
  }
  lastId, err = res.LastInsertId()
  if err != nil {
    return r, affectedRows, 0
  }
  return
}

func Current() string {
  return CurDB.Name
}

func LogQuery(f string) {
  LogFile = f
}

func Info() string {
  return fmt.Sprintf("%s %s", CurDB.Name, CurDB.Addrs)
}

func Drop(tn string) bool {
  r, _, _ := Exec(fmt.Sprintf("DROP TABLE `%s`", tn))
  return r
}

func Insert(tn string, data AssocData) (r bool, lastId int64) {
  str := Format(data)
  sql := fmt.Sprintf("INSERT INTO `%s` SET %s", tn, str)
  r, lastId, _ = Exec(sql)
  return
}

func Format(data AssocData) string {
  arr := make([]string, len(data))
  for k, v := range data {
    var val string
    switch t := reflect.ValueOf(v).String(); t {
      case "string" : 
        switch v {
          case "true" :
            val = "1"
          case "false" :
            val = "0"
          case "NOW()" :
            val = "NOW()"
          case "NULL" :
            val = "NULL"
          default :
            val = fmt.Sprintf("'%s'", url.QueryEscape(v.(string)))
        }
      case "bool" :
        val = "0"
        if v == true {
          val = "1"
        }
      default :
        val = fmt.Sprintf("%v", v)
    }
    arr = append(arr, fmt.Sprintf("`%s` = %s", k, val))
  }
  return strings.Join(arr, ",")
}

func ListDB() {
  for i, j := range DBList {
    fmt.Println(i)
    fmt.Println(j)
  }
}

