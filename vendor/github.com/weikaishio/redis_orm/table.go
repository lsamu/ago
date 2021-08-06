package redis_orm

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

/*
SET @table_schema='employees';
SELECT
    table_name,
    table_type,
    engine,
    table_rows,
    avg_row_length,
    data_length,
    index_length,
    table_collation,
    create_time
FROM
    information_schema.tables
WHERE
    table_schema = @table_schema
ORDER BY table_name;
*/
type SchemaTablesTb struct {
	Id            int64  `redis_orm:"pk autoincr comment 'ID'"`
	TableName     string `redis_orm:"unique comment '唯一'"`
	TableComment  string `redis_orm:"dft '' comment '表注释'"` //暂时没用上
	PrimaryKey    string `redis_orm:"comment '主键字段'"`
	AutoIncrement string `redis_orm:"comment '自增字段'"`
	IsSync2DB     bool   `redis_orm:"comment '是否同步到数据库'"`
	Created       string `redis_orm:"comment '创建时间字段'"`
	Updated       string `redis_orm:"comment '更新时间字段'"`
	//Version       int32  `redis_orm:"comment '版本'"`
	CreatedAt     int64  `redis_orm:"created_at comment '创建时间'"`
	UpdatedAt     int64  `redis_orm:"updated_at comment '修改时间'"`
}

func SchemaTablesFromTable(table *Table) *SchemaTablesTb {
	return &SchemaTablesTb{
		Id:            table.TableId,
		TableName:     table.Name,
		TableComment:  table.Comment,
		PrimaryKey:    table.PrimaryKey,
		AutoIncrement: table.AutoIncrement,
		IsSync2DB:     table.IsSync2DB,
		Created:       table.Created,
		Updated:       table.Updated,
		//Version:       table.Version,
	}
}

type Table struct {
	TableId int64
	Name    string
	Comment string
	//Version int32
	//Type          reflect.Type
	ColumnsSeq    []string
	ColumnsMap    map[string]*Column
	IndexesMap    map[string]*Index
	PrimaryKey    string
	AutoIncrement string
	IsSync2DB     bool
	Created       string
	Updated       string
	mutex         sync.RWMutex
}

func TableFromSchemaTables(table *SchemaTablesTb) *Table {
	return &Table{
		TableId:       table.Id,
		Name:          table.TableName,
		Comment:       table.TableComment,
		PrimaryKey:    table.PrimaryKey,
		AutoIncrement: table.AutoIncrement,
		IsSync2DB:     table.IsSync2DB,
		ColumnsMap:    make(map[string]*Column),
		IndexesMap:    make(map[string]*Index),
		Created:       table.Created,
		Updated:       table.Updated,
		//Version:       table.Version,
	}
}

func NewEmptyTable() *Table {
	return &Table{Name: "",
		ColumnsSeq: make([]string, 0),
		ColumnsMap: make(map[string]*Column),
		IndexesMap: make(map[string]*Index),
	}
}
func (table *Table) GetAutoIncrKey() string {
	if table.AutoIncrement != "" {
		return fmt.Sprintf("%s%s", KeyAutoIncrPrefix, strings.ToLower(table.AutoIncrement))
	} else {
		return ""
	}
}
func (table *Table) GetTableKey() string {
	return fmt.Sprintf("%s%s", KeyTbPrefix, strings.ToLower(table.Name))
}
func (table *Table) AddIndex(typ string, indexColumn, columnName, comment string, isUnique bool, seq byte) {
	var indexType IndexType
	switch typ {
	case reflect.String.String():
		indexType = IndexType_IdScore

	case reflect.Int.String(), reflect.Int8.String(), reflect.Int16.String(), reflect.Int32.String(), reflect.Int64.String(), reflect.Uint.String(), reflect.Uint8.String():
		fallthrough
	case reflect.Uint16.String(), reflect.Uint32.String(), reflect.Uint64.String(), reflect.Float32.String(), reflect.Float64.String():
		indexType = IndexType_IdMember

	case reflect.Uintptr.String(), reflect.Ptr.String():
		fallthrough
	case reflect.Complex64.String(), reflect.Complex128.String(), reflect.Array.String(), reflect.Chan.String(), reflect.Interface.String(), reflect.Map.String():
		fallthrough
	case reflect.Slice.String(), reflect.Struct.String(), reflect.Bool.String(), reflect.UnsafePointer.String():
		fallthrough
	default:
		indexType = IndexType_UnSupport
	}

	if indexType == IndexType_UnSupport {
		return
	}
	if isUnique {
		indexType = IndexType_IdScore
	}
	index := &Index{
		NameKey:     table.GetIndexKey(columnName),
		Seq:         seq,
		IndexColumn: strings.Split(indexColumn, "&"),
		Comment:     comment,
		Type:        indexType,
		IsUnique:    isUnique,
	}
	table.mutex.Lock()
	table.IndexesMap[strings.ToLower(indexColumn)] = index
	table.mutex.Unlock()
}
func (table *Table) GetIndexKey(col string) string {
	return fmt.Sprintf("%s%s_%s", KeyIndexPrefix, strings.ToLower(table.Name), strings.ToLower(col))
}
func (table *Table) AddColumn(col *Column) {
	if col.IsCombinedIndex {
		return
	}
	table.ColumnsSeq = append(table.ColumnsSeq, col.Name)
	colName := col.Name
	table.ColumnsMap[colName] = col

	if col.IsPrimaryKey {
		table.PrimaryKey = col.Name
	}
	if col.IsAutoIncrement {
		table.AutoIncrement = col.Name
	}
	if col.IsCreated {
		table.Created = col.Name
	}
	if col.IsUpdated {
		table.Updated = col.Name
	}
}
