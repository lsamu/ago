package clickhouse

import (
    "database/sql"
    "fmt"
    "github.com/ClickHouse/clickhouse-go"
    "log"
    "time"
)

//https://github.com/ClickHouse/clickhouse-go

func InitClickHouse() {

}

func GetClickHouse() {

}

func NewClickHouse() *ClickHouseOrm {
    orm := &ClickHouseOrm{}
    return orm
}

type (
    ClickHouseOrm struct {
        db *sql.DB
    }
)

func (o *ClickHouseOrm) Open() {
    connect, err := sql.Open("clickhouse", "tcp://127.0.0.1:9000?debug=true")
    if err != nil {
        log.Fatal(err)
    }
    if err := connect.Ping(); err != nil {
        if exception, ok := err.(*clickhouse.Exception); ok {
            fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
        } else {
            fmt.Println(err)
        }
        return
    }
    o.db = connect
}

func (o *ClickHouseOrm) Exec() (err error){
    _, err = o.db.Exec(`
		CREATE TABLE IF NOT EXISTS example (
			country_code FixedString(2),
			os_id        UInt8,
			browser_id   UInt8,
			categories   Array(Int16),
			action_day   Date,
			action_time  DateTime
		) engine=Memory
	`)

    if err != nil {
        log.Fatal(err)
    }

    return err
}

func (o *ClickHouseOrm) Query() (rows *sql.Rows, err error){
    rows, err = o.db.Query("SELECT country_code, os_id, browser_id, categories, action_day, action_time FROM example")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    for rows.Next() {
        var (
            country               string
            os, browser           uint8
            categories            []int16
            actionDay, actionTime time.Time
        )
        if err := rows.Scan(&country, &os, &browser, &categories, &actionDay, &actionTime); err != nil {
            log.Fatal(err)
        }
        log.Printf("country: %s, os: %d, browser: %d, categories: %v, action_day: %s, action_time: %s", country, os, browser, categories, actionDay, actionTime)
    }

    if err := rows.Err(); err != nil {
        log.Fatal(err)
    }
    return
}

func (o *ClickHouseOrm) First() {

}


func (o *ClickHouseOrm) Find() {

}


