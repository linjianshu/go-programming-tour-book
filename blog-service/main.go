package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"src/go-programming-tour-book/blog-service/global"
	"src/go-programming-tour-book/blog-service/internal/model"
	"src/go-programming-tour-book/blog-service/internal/routers"
	"src/go-programming-tour-book/blog-service/pkg/setting"
	"time"
)

func main() {
	//test

	//r := gin.Default()
	//r.GET("/ping", func(c *gin.Context) {
	//	c.JSON(200,gin.H{"message":"pong"})
	//})
	//
	//r.Run(":8081")

	gin.SetMode(global.ServerSetting.RunMode)
	router := routers.NewRouter()
	s := &http.Server{
		Addr:              ":" + global.ServerSetting.HttpPort,
		Handler:           router,
		ReadTimeout:       global.ServerSetting.ReadTimeout,
		ReadHeaderTimeout: global.ServerSetting.WriteTimeout,
		MaxHeaderBytes:    1 << 20,
	}
	s.ListenAndServe()
}

func init() {
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}

	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init.setupDBEngine err: %v", err)
	}
}

func setupSetting() error {
	setting, err := setting.NewSetting()
	if err != nil {
		return err
	}
	err = setting.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}
	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second
	return nil
}

func setupDBEngine() error {
	var err error
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}
	return nil
}
