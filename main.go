package main

import (
	"Garyen-go/pkg/setting"
	"Garyen-go/router"
	"fmt"
	"net/http"
)

func main() {
	// 启动http服务
	r := router.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HttpPort),
		Handler:        r,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	_ = s.ListenAndServe()
}
