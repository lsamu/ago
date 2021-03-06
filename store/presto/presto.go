package presto

import (
    "database/sql"
    "errors"
    "fmt"
    "github.com/blockloop/scan"
    _ "github.com/prestodb/presto-go-client/presto"
)

type Presto struct {
    db    *sql.DB
    rows  *sql.Rows
    Error error
}

//Open 打开
/**
格式：http[s]://user[:pass]@host[:port][?parameters]
例子：http://root@bigdata.node-1:8081?catalog=hive&schema=test
*/
func Open(dsn string) (p *Presto, err error) {
    //dsn = "http://admin@master:8523?catalog=hive&schema=hub"
    db, err := sql.Open("presto", dsn)
    if err != nil {
        return nil, errors.New("open db error")
    }
    p = &Presto{
        db: db,
    }
    return p, nil
}

func Format(userName, password, host, port, catalog, schema string) string  {
    dsn := fmt.Sprintf("http://%s:%s@%s:%s?catalog=%s&schema=%s",
        userName,
        password,
        host,
        port,
        catalog,
        schema,
    )
    return dsn
}

//Close 关闭
func (p *Presto) Close() {
    err := p.db.Close()
    if err != nil {
        p.Error = err
    }
}

//Exec 执行sql
func (p *Presto) Exec(sql string, args ...interface{}) *Presto {
    result, err := p.db.Exec(sql, args)
    if err != nil {
        p.Error = err
        return p
    }
    fmt.Println(result.LastInsertId())
    return p
}

//Query 查询
func (p *Presto) Query(sql string, args ...interface{}) *Presto {
    rows, err := p.db.Query(sql, args...)
    if err != nil {
        p.Error = err
        return p
    }
    p.rows = rows
    return p
}

//Find 赋值
func (p *Presto) Find(out interface{}) *Presto {
    err := scan.Rows(out, p.rows)
    p.Error = err
    return p
}

//First 单行  首行首列
func (p *Presto) First(out interface{}) *Presto {
    err := scan.Row(out, p.rows)
    p.Error = err
    return p
}

//Map 字典
func (p *Presto) Map() (list []map[string]interface{}, err error) {
    columns, _ := p.rows.Columns()
    columnLength := len(columns)
    cache := make([]interface{}, columnLength)
    for index, _ := range cache {
        var a interface{}
        cache[index] = &a
    }
    list = []map[string]interface{}{}
    for p.rows.Next() {
        _ = p.rows.Scan(cache...)
        item := make(map[string]interface{})
        for i, data := range cache {
            item[columns[i]] = *data.(*interface{}) //取实际类型
        }
        list = append(list, item)
    }
    _ = p.rows.Close()
    return list, nil
}