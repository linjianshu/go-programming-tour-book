package sql2struct

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/wendal/go-oci8"
	"golang.org/x/text/encoding/simplifiedchinese"
	"strings"
)

var DBTypeToStructType = map[string]string{
	"int":       "int32",
	"tinyint":   "int8",
	"smallint":  "int",
	"mediumint": "int64",
	"bigint":    "int64",
	"bit":       "int",
	"bool":      "bool",
	"enum":      "string",
	"set":       "string",
	"varchar":   "string",

	//oracle 特殊
	"CHAR":      "string",
	"DATE":      "time.Time",
	"FLOAT":     "float",
	"NUMBER":    "int64",
	"NVARCHAR2": "string",
	"RAW":       "[]byte",
	"TIMESTAMP": "time.Time",
	"VARCHAR2":  "string",
}

type DBModel struct {
	DBEngine *sql.DB
	DBInfo   *DBInfo
}

type DBInfo struct {
	DBType   string
	Host     string
	UserName string
	Password string
	Charset  string
	//oracle特殊
	ServiceName string
}

type TableColumn struct {
	ColumnName    string
	DataType      string
	IsNullable    string
	ColumnKey     string
	ColumnType    string
	ColumnComment string
}

// NewDBModel 使用该函数传入DBInfo指针 操纵指针给DBModel赋值 返回一个DBModel指针
func NewDBModel(info *DBInfo) *DBModel {
	return &DBModel{DBInfo: info}
}

// Connect 使用该方法,通过DBModel自带的DBInfo结构体的参数,拼接成dsn(datasourcename),然后通过拼接后的字符串
//调用mysql驱动 返回的对象装进自身的DBEngine中
func (m *DBModel) Connect() error {
	var s string
	var dsn string
	var err error
	switch t := m.DBInfo.DBType; t {
	case "mysql":
		s = "%s:%s@tcp(%s)/information_schema?" + "charset=%s&parseTime=True&loc=Local"
		dsn = fmt.Sprintf(s, m.DBInfo.UserName, m.DBInfo.Password, m.DBInfo.Host, m.DBInfo.Charset)
	case "oci8":
		s = "%s/%s@%s/%s"
		dsn = fmt.Sprintf(s, m.DBInfo.UserName, m.DBInfo.Password, m.DBInfo.Host, m.DBInfo.ServiceName)
	}
	m.DBEngine, err = sql.Open(m.DBInfo.DBType, dsn)
	if err != nil {
		return err
	}
	return nil
}

// GetColumns  拼接sql语句 在mysql的information_schema表中查询表的信息 查到后遍历 遍历后塞到返回的TableColumn指针对象中
func (m *DBModel) GetColumns(dbName, tableName string) ([]*TableColumn, error) {
	var query string
	var rows *sql.Rows
	var err error
	switch t := m.DBInfo.DBType; t {
	case "mysql":
		query = `SELECT
	column_name,
		data_type,
		column_key,
		is_nullable,
		column_type,
		column_comment
	FROM
	COLUMNS
	WHERE
	table_schema = ?
	AND
	table_name = ?
	`
		rows, err = m.DBEngine.Query(query, dbName, tableName)
		if err != nil {
			return nil, err
		}
		if rows == nil {
			return nil, errors.New("没有数据")
		}

	case "oci8":

		query = `SELECT 
       T1.COLUMN_NAME column_name,
       T1.DATA_TYPE || '(' || T1.DATA_LENGTH || ')' data_type,
       T1.NULLABLE column_key,
       T1.NULLABLE is_nullable,
       T1.DATA_TYPE || '(' || T1.DATA_LENGTH || ')' column_type,
       T2.COMMENTS column_comment
  	FROM USER_TAB_COLS T1, USER_COL_COMMENTS T2
 	WHERE T1.TABLE_NAME = T2.TABLE_NAME
   	AND T1.COLUMN_NAME = T2.COLUMN_NAME 
  `
		query = fmt.Sprintf(query + fmt.Sprintf("and T1.TABLE_NAME = '%s'", tableName))
		rows, err = m.DBEngine.Query(query)
		if err != nil {
			return nil, err
		}
		if rows == nil {
			return nil, errors.New("没有数据")
		}
	}

	defer rows.Close()

	var columns []*TableColumn
	for rows.Next() {
		var column TableColumn
		err = rows.Scan(&column.ColumnName, &column.DataType, &column.ColumnKey, &column.IsNullable, &column.ColumnType, &column.ColumnComment)
		if err != nil {
			return nil, err
		}
		columns = append(columns, &column)
	}
	//oracle DATA_TYPE去空格
	for _, v := range columns {
		if strings.Contains(v.DataType, "(") {
			index := strings.Index(v.DataType, "(")
			v.DataType = v.DataType[:index]
			v.ColumnComment, _ = simplifiedchinese.GBK.NewDecoder().String(v.ColumnComment)
		}
	}
	return columns, nil
}
