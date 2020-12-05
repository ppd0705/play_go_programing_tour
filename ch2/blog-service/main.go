package main

import (
	"block-service/global"
	"block-service/internal/model"
	"block-service/internal/routers"
	"block-service/pkg/logger"
	"block-service/pkg/setting"
	"block-service/pkg/tracer"
	"context"
	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func setupSetting() error {
	newSetting, err := setting.NewSetting()
	if err != nil {
		return err
	}
	err = newSetting.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}
	err = newSetting.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}
	err = newSetting.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}
	err = newSetting.ReadSection("JWT", &global.JWTSetting)
	if err != nil {
		return err
	}
	global.JWTSetting.Expire *= time.Second
	global.ServerSetting.ReadTimeOut *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second
	return nil
}

func setDBEngine() error {
	var err error
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}
	return nil
}

func setupLogger() error {
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename:  global.AppSetting.LogSavePath + "/" + global.AppSetting.LogFileName + global.AppSetting.LogFileExt,
		MaxSize:   600,
		MaxAge:    10,
		LocalTime: true,
	}, "", log.LstdFlags).WithCaller(2)
	return nil
}

func init() {
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v\n", err)
	}
	err = setDBEngine()
	if err != nil {
		log.Fatalf("init.setupDBengine err: %v\n", err)
	}
	err = setupLogger()
	if err != nil {
		log.Fatalf("init.setupLogger err: %v\n", err)
	}
	err = setupTracer()
	if err != nil {
		log.Fatalf("init.setupTracer err: %v\n", err)
	}
}

func setupTracer() error {
	jaegerTracer, _, err := tracer.NewJaegerTracer("blog-service", "127.0.0.1:6831")
	if err != nil {
		return err
	}
	global.Tracer = jaegerTracer
	return nil
}

// @title block_service
// @version 1.0
// @description hello blog
func main() {
	gin.SetMode(global.ServerSetting.RunMode)
	router := routers.NewRouter()
	s := &http.Server{
		Addr:           ":" + global.ServerSetting.HttpPort,
		Handler:        router,
		ReadTimeout:    global.ServerSetting.ReadTimeOut,
		WriteTimeout:   global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	global.Logger.Info("blog service")
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe err: %v\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}
