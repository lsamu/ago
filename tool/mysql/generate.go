package mysql

import (
    "bytes"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "github.com/jinzhu/gorm"
    "github.com/urfave/cli/v2"
    "go.uber.org/zap"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "text/template"
    "time"
    "unicode"
)

// https://github.com/arrayhua/go_grpc_gorm_micro
// Generate 生成
func Generate(c *cli.Context) (err error) {
    return
}

var DB *gorm.DB

func InitGormMysql(username, password, host, port, database string, debug bool) {
    dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        username, password, host, port, database,
    )
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
    TempateBasePath string // 模板路径
    OutBasePath     string // 生成目录
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
    var templateStruct TemplateStruct
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
    // 获取 basePath 文件夹下所有tpl文件
    tplFileList, err := s.getAllTplFile(s.TempateBasePath, nil)
    if err != nil {
        return err
    }
    fmt.Println(tplFileList)
    // 遍历所有表，获取表结构
    for _, value := range tableNames {
        columns, err := s.GetColumns(dbName, value.TableName)
        if err != nil {
            continue
        }
        fmt.Printf("%+v", columns)
        templateStruct.TableName = value.TableName
        templateStruct.ModelName = s.Case2Camel(value.TableName)
        templateStruct.RouterName = "req.Path"
        templateStruct.Fields = columns
        // 生成源码
        err = s.createTemp(tplFileList, templateStruct)
        if err != nil {
            log.Println(" 生成源码失败" + err.Error())
            continue
        }
    }
    return
}

type tplData struct {
    template     *template.Template
    locationPath string
    autoCodePath string
}

// createTemp  createTemp
func (s *Schema) createTemp(tplFileList []string, templateStruct TemplateStruct) (err error) {
    basePath := s.TempateBasePath
    dataList := make([]tplData, 0, len(tplFileList))
    fileList := make([]string, 0, len(tplFileList))
    // 定义一个函数add
    // 这个函数要么只有一个返回值，要么有俩返回值且第二个返回值必须是error类型
    add := func(params int) (int, error) {
        return 100 + params, nil
    }
    // 定义首字母小写
    case2CamelAndLcfirst := s.Case2CamelAndLcfirst
    // 定义首字母大写
    case2CamelAndUcfirst := s.Case2CamelAndUcfirst

    // 根据文件路径生成 tplData 结构体，待填充数据
    for _, value := range tplFileList {
        dataList = append(dataList, tplData{locationPath: value})
    }
    // 生成 *Template, 填充 template 字段
    for index, value := range dataList {
        textByte, err := ioutil.ReadFile(value.locationPath)
        if err != nil {
            return err
        }
        dataList[index].template, err = template.New("").Funcs(template.FuncMap{"add": add, "case2CamelAndLcfirst": case2CamelAndLcfirst, "case2CamelAndUcfirst": case2CamelAndUcfirst}).Parse(string(textByte))
        if err != nil {
            return err
        }
    }

    // 生成文件路径，填充 autoCodePath 字段
    for index, value := range dataList {
        trimBase := strings.TrimPrefix(value.locationPath, basePath)
        if lastSeparator := strings.LastIndex(trimBase, "/"); lastSeparator != -1 {
            origFileName := strings.TrimSuffix(trimBase[lastSeparator+1:], ".tpl")
            firstDot := strings.Index(origFileName, ".")
            if firstDot != -1 {
                dataList[index].autoCodePath = s.OutBasePath + trimBase[:lastSeparator] + "/" + templateStruct.TableName + origFileName
            }
        }
    }

    // 生成文件
    for _, value := range dataList {
        fileList = append(fileList, value.autoCodePath)
        dir := filepath.Dir(value.autoCodePath)
        err := os.MkdirAll(dir, os.ModePerm)
        f, err := os.OpenFile(value.autoCodePath, os.O_CREATE|os.O_WRONLY, 0755)
        if err != nil {
            return err
        }
        if err = value.template.Execute(f, templateStruct); err != nil {
            return err
        }
        _ = f.Close()
    }
    return
}

// getAllTplFile getAllTplFile
func (s *Schema) getAllTplFile(pathName string, fileList []string) ([]string, error) {
    files, err := ioutil.ReadDir(pathName)
    for _, fi := range files {
        if fi.IsDir() {
            fileList, err = s.getAllTplFile(pathName+"/"+fi.Name(), fileList)
            if err != nil {
                return nil, err
            }
        } else {
            if strings.HasSuffix(fi.Name(), ".tpl") {
                fileList = append(fileList, pathName+"/"+fi.Name())
            }
        }
    }
    return fileList, err
}

// Camel2Case 驼峰式写法转为下划线写法
func (s *Schema) Camel2Case(name string) string {
    buffer := NewBuffer()
    for i, r := range name {
        if unicode.IsUpper(r) {
            if i != 0 {
                buffer.Append('_')
            }
            buffer.Append(unicode.ToLower(r))
        } else {
            buffer.Append(r)
        }
    }
    return buffer.String()
}

// Case2Camel 下划线写法转为驼峰写法
func (s *Schema) Case2Camel(name string) string {
    name = strings.Replace(name, "_", " ", -1)
    name = strings.Title(name)
    return strings.Replace(name, " ", "", -1)
}

// Ucfirst 首字母大写
func (s *Schema) Ucfirst(str string) string {
    for i, v := range str {
        return string(unicode.ToUpper(v)) + str[i+1:]
    }
    return ""
}

// Case2CamelAndUcfirst 下划线写法转为驼峰写法并且首字母大写
func (s *Schema) Case2CamelAndUcfirst(name string) string {
    return s.Ucfirst(s.Case2Camel(name))
}

// Case2CamelAndLcfirst 下划线写法转为驼峰写法并且首字母小写
func (s *Schema) Case2CamelAndLcfirst(name string) string {
    return s.Lcfirst(s.Case2Camel(name))
}

// Lcfirst 首字母小写
func (s *Schema) Lcfirst(str string) string {
    for i, v := range str {
        return string(unicode.ToLower(v)) + str[i+1:]
    }
    return ""
}

// Buffer 内嵌bytes.Buffer，支持连写
type Buffer struct {
    *bytes.Buffer
}

func NewBuffer() *Buffer {
    return &Buffer{Buffer: new(bytes.Buffer)}
}

func (b *Buffer) Append(i interface{}) *Buffer {
    switch val := i.(type) {
    case int:
        b.append(strconv.Itoa(val))
    case int64:
        b.append(strconv.FormatInt(val, 10))
    case uint:
        b.append(strconv.FormatUint(uint64(val), 10))
    case uint64:
        b.append(strconv.FormatUint(val, 10))
    case string:
        b.append(val)
    case []byte:
        b.Write(val)
    case rune:
        b.WriteRune(val)
    }
    return b
}

func (b *Buffer) append(s string) *Buffer {
    defer func() {
        if err := recover(); err != nil {
            log.Println("*****内存不够了！******")
        }
    }()
    b.WriteString(s)
    return b
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
