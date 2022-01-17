package mysql

import (
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "github.com/jinzhu/gorm"
    "github.com/urfave/cli/v2"
    "go.uber.org/zap"
    "time"
)

// https://github.com/arrayhua/go_grpc_gorm_micro
// Generate 生成
func Generate(c *cli.Context) (err error) {
    return
}

var DB *gorm.DB

func InitGormMysql() {
    dsn := ""
    db, err := gorm.Open("mysql", dsn)
    if err != nil {
        fmt.Println("mysql连接失败", zap.Any("err", err))
        panic(err)
    }
    // See "Important settings" section.
    db.DB().SetConnMaxLifetime(time.Minute * 3)
    db.DB().SetMaxOpenConns(100)
    db.DB().SetMaxIdleConns(100)
    db.LogMode(true)
    fmt.Println("mysql连接成功")
    DB = db
}

type Schema struct {
}

// GetTables 获取所有表
func (s *Schema) GetTables(dbName string) (tables []TableSchema, err error) {
    err = DB.Raw("select table_name as table_name from information_schema.tables where TABLE_SCHEMA = ?", dbName).Scan(&tables).Error
    return tables, err
}

// GetColumns 获取所有列
func (s *Schema) GetColumns(dbName string, tableName string) (columns []ColumnSchema, err error) {
    err = DB.Raw("SELECT COLUMN_NAME column_name,DATA_TYPE data_type,COLUMN_KEY column_key,EXTRA extra,CASE DATA_TYPE WHEN 'longtext' THEN c.CHARACTER_MAXIMUM_LENGTH WHEN 'varchar' THEN c.CHARACTER_MAXIMUM_LENGTH WHEN 'double' THEN CONCAT_WS( ',', c.NUMERIC_PRECISION, c.NUMERIC_SCALE ) WHEN 'decimal' THEN CONCAT_WS( ',', c.NUMERIC_PRECISION, c.NUMERIC_SCALE ) WHEN 'int' THEN c.NUMERIC_PRECISION WHEN 'bigint' THEN c.NUMERIC_PRECISION ELSE '' END AS data_type_long,COLUMN_COMMENT column_comment FROM INFORMATION_SCHEMA.COLUMNS c WHERE TABLE_NAME = ? AND TABLE_SCHEMA = ?", tableName, dbName).Scan(&columns).Error
    return columns, err
}

// Generate Generate
func (s *Schema) Generate(dbName, tableName string) (err error) {
    var tableNames []TableSchema
    if tableName == "" { // 获取所有表
        tableNames, err = s.GetTables(dbName)
        if err != nil {
            return err
        }
    } else {
        var custableName TableSchema
        custableName.TableName = tableName
        tableNames = append(tableNames, custableName)
    }
    // 遍历所有表，获取表结构
    for _, value := range tableNames {
        fmt.Println(value)
        // 组装数据
        err, columns := s.GetColumns(dbName, value.TableName)
        if err != nil {
            continue
        }
        fmt.Println(columns)
    }
    return
}

// createTemp  createTemp
func (s *Schema) createTemp(tplFileList []string, templateStruct TemplateStruct) {

}

// TableSchema TableSchema
type TableSchema struct {
    TableName string `json:"tableName"`
}

// ColumnSchema ColumnSchema
type ColumnSchema struct {
    ColumnName    string `json:"columnName" gorm:"column:column_name"`
    DataType      string `json:"dataType" gorm:"column:data_type"`
    COLUMNKEY     string `json:"columnKey" gorm:"column:column_key"`
    EXTRA         string `json:"extra" gorm:"column:extra"`
    DataTypeLong  string `json:"dataTypeLong" gorm:"column:data_type_long"`
    ColumnComment string `json:"columnComment" gorm:"column:column_comment"`
}

// TemplateStruct TemplateStruct
type TemplateStruct struct {
    ModelName  string         `json:"structName"` // SysApis
    TableName  string         `json:"tableName"`  // sys_apis
    RouterName string         `json:"routerName"` // sysApis
    Fields     []ColumnSchema `json:"fields"`
}
