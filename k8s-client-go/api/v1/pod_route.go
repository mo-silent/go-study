package v1

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"gitee.com/MoGD/go-study/k8s-client-go/global"
	"gitee.com/MoGD/go-study/k8s-client-go/model/common/response"
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// @Tags GetAllPod
// @Summary 获取所有 Pod 信息
// @Produce application/json
// @Success 200 {object} response.CommonResponse
// @Router /pod/getAllPod [get]
func GetAllPod(c *gin.Context) {
	// list pod
	pods, err := global.K8SCLIENT.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	//fmt.Sprintf("There are %d pods in the cluster\n", len(pods.Items))
	c.JSON(http.StatusOK, response.CommonResponse{
		Message: pods.Items,
	})
}

// @Tags GetNamespacePod
// @Summary 获取单个命名空间中所有的  Pod 信息
// @Produce application/json
// @Param   namespace    path  string  false "命名空间" default(default)
// @Success 200 {object} response.CommonResponse
// @Router /pod/getNamespacePod/{namespace} [get]
func GetNamespacePod(c *gin.Context) {
	// get namespace
	namespcae := c.Param("namespace")
	if namespcae == "{namespace}" || namespcae == "" ||
		strings.TrimSpace(namespcae) == "" {
		namespcae = "default"
	}
	// list pod
	pods, err := global.K8SCLIENT.CoreV1().Pods(namespcae).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	//fmt.Sprintf("There are %d pods in the cluster\n", len(pods.Items))
	c.JSON(http.StatusOK, response.CommonResponse{
		Message: pods.Items,
	})
}

// @Tags GetPod
// @Summary 获取单个 Pod 信息
// @Produce application/json
// @Param   namespace  query  string  false "命名空间" default(default)
// @Param   name    query  string  true "pod名称"
// @Success 200 {object} response.CommonResponse
// @Router /pod/getPod [post]
func GetPod(c *gin.Context) {
	// get namespace
	namespcae := c.DefaultQuery("namespace", "default")
	name := c.Query("name")
	fmt.Println(namespcae, name)
	if namespcae == "" || strings.TrimSpace(namespcae) == "" {
		namespcae = "default"
	}
	// list one pod
	pods, err := global.K8SCLIENT.CoreV1().Pods(namespcae).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	//fmt.Sprintf("There are %d pods in the cluster\n", len(pods.Items))
	c.JSON(http.StatusOK, response.CommonResponse{
		Message: pods,
	})
}

// @Tags CreatePod
// @Summary 获取单个 Pod 信息
// @Produce application/json
// @Success 200 {object} response.CommonResponse
// @Router /pod/createPod [post]
func CreatePod(c *gin.Context) {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-pod",
			Labels: map[string]string{
				"app": "demo",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "demo-pod",
					Image: "nginx:1.12",
					Ports: []corev1.ContainerPort{
						{
							Name:          "http",
							Protocol:      corev1.ProtocolTCP,
							ContainerPort: 80,
						},
					},
				},
			},
			RestartPolicy: "Always",
		},
	}
	// 获取 pod 接口
	podClient := global.K8SCLIENT.CoreV1().Pods("default")
	// 创建 Pod
	_, err := podClient.Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		c.JSON(http.StatusOK, response.CommonResponse{
			Message: "create pod fail!",
		})
		panic(err.Error())
	}
	// 循环获取 pod 状态，检查为 Running 状态后，返回 pod 信息
	for {
		podStatus, _ := podClient.Get(context.TODO(), "demo-pod", metav1.GetOptions{})
		if podStatus.Status.Phase == "Running" {
			c.JSON(http.StatusOK, response.CommonResponse{
				Message: podStatus,
			})
			break
		}
	}

}

// @Tags DeletePod
// @Summary 删除单个 Pod
// @Produce application/json
// @Param   namespace  query  string  false "命名空间" default(default)
// @Param   name    query  string  true "pod名称"
// @Success 200 {object} response.CommonResponse
// @Router /pod/deletePod [post]
func DeletePod(c *gin.Context) {
	// 获取命名空间和 pod 名称
	namespcae := c.DefaultQuery("namespace", "default")
	name := c.Query("name")
	// 获取 pod 接口
	podClient := global.K8SCLIENT.CoreV1().Pods(namespcae)
	// 删除 Pod
	err := podClient.Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		c.JSON(http.StatusOK, response.CommonResponse{
			Message: "delete pod fail!",
		})
		panic(err.Error())
	}

	c.JSON(http.StatusOK, response.CommonResponse{
		Message: fmt.Sprintf("delete pod %v success", name),
	})

}
