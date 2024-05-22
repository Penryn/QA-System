package main

import (
	"QA-System/internal/middleware"
	mongodb "QA-System/internal/pkg/database/mongodb"
	mysql "QA-System/internal/pkg/database/mysql"
	"QA-System/internal/pkg/log"
	"QA-System/internal/pkg/queue/asynq"
	"QA-System/internal/pkg/session"
	"QA-System/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化日志系统
	log.ZapInit()
	// 初始化数据库
    mysql.MysqlInit()
	mongodb.MongodbInit()
	// 初始化gin
	r := gin.Default()
	r.Use(middlewares.ErrHandler())
	r.NoMethod(middlewares.HandleNotFound)
	r.NoRoute(middlewares.HandleNotFound)
	r.Static("/static", "./public/static")
	r.Static("/xlsx", "./public/xlsx")
	session.Init(r)
	router.Init(r)
	go asynq.AsynqInit()
	err := r.Run()
	if err != nil {
		log.Logger.Fatal("Failed to start the server:" + err.Error())
	}
}