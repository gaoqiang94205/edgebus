package common

import "time"

const (
	//测试请求
	Ping = iota
	Application
	Other
)

//请求参数解析相关
const (
	Url    = "url"
	Method = "method"
	Edge   = "edge"
)

//边缘端发送请求
//请求结构
type SideRequest struct {
	Target string
	Url    string
	Method string
	Body   string
	Type   int
}

//返回结果
type SideResponse struct {
	content string
	url     string
}

//心跳请求
type HeartRequest struct {
	client string
	time   time.Time
}

type UpgradeRequest struct {
	Name string `json:"name" binding:"required"`
}

//向边缘端注入云端信息
type CloudInfoRequest struct {
	Addr string
	Path string
}

//服务端发送请求
//
