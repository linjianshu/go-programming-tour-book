package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	sql2struct2 "src/go-programming-tour-book/tour/sql2struct"
	"src/go-programming-tour-book/tour/zap"
)

var username string
var password string
var host string
var charset string
var dbType string
var dbName string
var tableName string

//oracle特殊
var serviceName string

var sqlCmd = &cobra.Command{
	Use:   "sql",
	Short: "sql转换和处理",
	Long:  "sql转换和处理",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var sql2structCmd = &cobra.Command{
	Use:   "struct",
	Short: "sql转换",
	Long:  "sql转换",
	Run: func(cmd *cobra.Command, args []string) {
		dbinfo := &sql2struct2.DBInfo{
			DBType:      dbType,
			Host:        host,
			UserName:    username,
			Password:    password,
			Charset:     charset,
			ServiceName: serviceName,
		}
		dbModel := sql2struct2.NewDBModel(dbinfo)
		err := dbModel.Connect()
		if err != nil {
			log.Fatalf("dbModel,Connect err: %v ", err)
		}
		columns, err := dbModel.GetColumns(dbName, tableName)
		if err != nil {
			log.Fatalf("dbModel.GetColumns err: %v", err)
		}

		template := sql2struct2.NewStructTemplate()
		templateColumns := template.AssemblyColumns(columns)
		err = template.Generate(tableName, templateColumns)
		if err != nil {
			log.Fatalf("template.Generate err: %v", err)
		}

		if z == 1 {
			zap.WriteLog(fmt.Sprintf("go run main.go sql struct --username %s --password %s --db %s --host %s --serviceName %s --table %s --type %s -z 1", username, password, dbName, host, serviceName, tableName, dbType))
			zap.WriteLog(fmt.Sprintf("输出结果 : %s ", "请查看models文件夹!"))
		}
	},
}

func init() {
	sqlCmd.AddCommand(sql2structCmd)
	sql2structCmd.Flags().StringVarP(&username, "username", "", "", "请输入数据库的账号")
	sql2structCmd.Flags().StringVarP(&password, "password", "", "", "请输入数据库的密码")
	sql2structCmd.Flags().StringVarP(&host, "host", "", "127.0.0.1:3306", "请输入数据库的HOST")
	sql2structCmd.Flags().StringVarP(&charset, "charset", "", "utf8mb4", "请输入数据库的编码")
	sql2structCmd.Flags().StringVarP(&dbType, "type", "", "mysql", "请输入数据库的类型")
	//oracle特殊
	sql2structCmd.Flags().StringVarP(&serviceName, "serviceName", "", "LocalCummins", "请输入数据库的服务名")
	sql2structCmd.Flags().StringVarP(&dbName, "db", "", "", "请输入数据库的名称")
	sql2structCmd.Flags().StringVarP(&tableName, "table", "", "", "请输入表名称")
}
