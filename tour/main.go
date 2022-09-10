package main

import (
	"log"
	"src/go-programming-tour-book/tour/cmd"
	sql2struct2 "src/go-programming-tour-book/tour/sql2struct"
)

func main() {
	//这里要调用的不是 github上的cmd 而是cmd文件夹里的execute方法
	err := cmd.Execute()
	if err != nil {
		log.Fatalf("cmd.Execute err: %v", err)
	}
}

var username string
var password string
var host string
var charset string
var dbType string
var dbName string
var tableName string

//oracle特殊
var serviceName string

//测试用
func main1() {
	//这里要调用的不是 github上的cmd 而是cmd文件夹里的execute方法
	username = "JACMES"
	password = "hfutIE2019"
	host = "localhost:1521"
	dbType = "oci8"
	dbName = "JACMES"
	tableName = "AREA"
	serviceName = "LocalCummins"

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
}
