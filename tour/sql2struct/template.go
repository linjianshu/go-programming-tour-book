package sql2struct

import (
	"fmt"
	"os"
	"src/go-programming-tour-book/tour/internal/word"
	"text/template"
)

//定义模版
const structTpl = `type {{ .TableName | ToCamelCase }} struct {
{{ range .Columns}}  {{ $length:= len .Comment}} {{ if gt $length 0}} //{{ .Comment }} {{ else}} // {{ .Name}} {{end}}
	{{ $typeLen :=len .Type }} {{if gt $typeLen 0}} {{.Name | ToCamelCase}}	{{.Type }} {{.Tag}} {{else }} {{.Name}} {{end}}
{{end}}}

func (model {{.TableName | ToCamelCase }}) TableName() string {
	return "{{.TableName}}"
}`

// StructTemplate 这个结构体 承载模版的对象
type StructTemplate struct {
	structTpl string
}

// StructColumn 用来存储转换后的go结构体中的所有字段信息
type StructColumn struct {
	Name    string
	Type    string
	Tag     string
	Comment string
}

// StructTemplateDB 用于存储待渲染的对象 也就是mysql中的表结构
type StructTemplateDB struct {
	TableName string
	Columns   []*StructColumn
}

// NewStructTemplate 新建一个StructTemplate指针对象 并将模版存入
func NewStructTemplate() *StructTemplate {
	return &StructTemplate{structTpl: structTpl}
}

// AssemblyColumns 将连接数据库获取到的数据库字段 转换成StructColumn结构体
//并且进行了数据库类型和结构体类型的转换 json tag的转换
//有了这个 就相当于有了待渲染的素材
func (t *StructTemplate) AssemblyColumns(tbColumns []*TableColumn) []*StructColumn {
	tplColumns := make([]*StructColumn, 0, len(tbColumns))
	for _, column := range tbColumns {
		tag := fmt.Sprintf("`"+"json:"+"\"%s\""+"`", column.ColumnName)
		tplColumns = append(tplColumns, &StructColumn{
			Name:    column.ColumnName,
			Type:    DBTypeToStructType[column.DataType],
			Tag:     tag,
			Comment: column.ColumnComment,
		})
	}
	return tplColumns
}

// Generate  定义一个模版解析方法 根据传入的表名和已经做好转换的StructColumn 对t内部的原始模版进行渲染
func (t *StructTemplate) Generate(tableName string, tplColumns []*StructColumn) error {
	tpl := template.Must(template.New("sql2struct").Funcs(template.FuncMap{
		//"ToCamelCase": word.UnderscoreToUpperCamelCase,
		"ToCamelCase": word.ToUpper,
	}).Parse(t.structTpl))
	tplDB := StructTemplateDB{
		TableName: tableName,
		Columns:   tplColumns,
	}

	file, err := os.OpenFile(fmt.Sprintf("./models/%s.log", tableName), os.O_WRONLY|os.O_CREATE, 0644)
	defer file.Close()
	if err != nil {
		return err
	}
	err = tpl.Execute(os.Stdout, tplDB)
	if err != nil {
		return err
	}
	err = tpl.Execute(file, tplDB)
	if err != nil {
		return err
	}
	return nil
}
