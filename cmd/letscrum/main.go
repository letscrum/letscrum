package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/letscrum/letscrum/models"
	"github.com/letscrum/letscrum/pkg/gredis"
	"github.com/letscrum/letscrum/pkg/logging"
	"github.com/letscrum/letscrum/pkg/settings"
	"github.com/letscrum/letscrum/routers"
)

func init() {
	settings.Setup("./config/config.yaml")
	models.Setup()
	godotenv.Load()
	logging.Setup()
	gredis.Setup()
}

// @title Golang Gin API
// @version 1.0
// @description An example of gin
// @termsOfService https://github.com/letscrum/letscrum
// @license.name MIT
// @license.url https://github.com/letscrum/letscrum/blob/master/LICENSE
func main() {
	gin.SetMode(settings.ServerSetting.RunMode)

	routersInit := routers.InitRouter()
	readTimeout := settings.ServerSetting.ReadTimeout
	writeTimeout := settings.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", settings.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	log.Printf("[info] start http server listening %s", endPoint)

	server.ListenAndServe()

	main.Execute()

	
	// If you want Graceful Restart, you need a Unix system and download github.com/fvbock/endless
	//endless.DefaultReadTimeOut = readTimeout
	//endless.DefaultWriteTimeOut = writeTimeout
	//endless.DefaultMaxHeaderBytes = maxHeaderBytes
	//server := endless.NewServer(endPoint, routersInit)
	//server.BeforeBegin = func(add string) {
	//	log.Printf("Actual pid is %d", syscall.Getpid())
	//}
	//
	//err := server.ListenAndServe()
	//if err != nil {
	//	log.Printf("Server err: %v", err)
	//}
}