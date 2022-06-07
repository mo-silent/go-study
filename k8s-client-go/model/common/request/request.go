package request

type PodReques struct {
	PodName       string            `json:"podName"`
	ContainerName string            `json:"containerName"`
	Image         string            `json:"image"`
	Labels        map[string]string `json:"labels"`
	ContainerPort int32             `json:"containerPort"`
}
