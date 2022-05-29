package main

import (
	"flag"
	v1 "gitee.com/MoGD/go-study/k8s-client-go/api/v1"
	_ "gitee.com/MoGD/go-study/k8s-client-go/docs"
	"gitee.com/MoGD/go-study/k8s-client-go/global"
	"gitee.com/MoGD/go-study/k8s-client-go/initialize"
	"k8s.io/client-go/util/homedir"
	"net/http"
	"path/filepath"
	"time"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

// @title Swagger Example API
// @version 0.0.1
// @description This is a sample Server pets
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name x-token
// @BasePath /
func main() {
	// 传入的参数
	if home := homedir.HomeDir(); home != "" {
		global.KUBECONFIG = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		global.KUBECONFIG = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// create the k8sClient
	global.K8SCLIENT = initialize.InitK8sClient(global.KUBECONFIG)

	// 初始化路由
	router := initialize.InitRouters()
	// 注册 pod 路由组
	podGroup := router.Group("pod")
	{
		podGroup.GET("getNamespacePod/:namespace", v1.GetNamespacePod)
		podGroup.GET("getAllPod", v1.GetAllPod)
		podGroup.POST("getPod", v1.GetPod)
		podGroup.POST("createPod", v1.CreatePod)
	}

	// Custom HTTP configuration
	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	// server start to listen
	s.ListenAndServe()
}
