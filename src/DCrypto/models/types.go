package models

import (
	"sync"
)

type (
	AppConf struct {
		// center
		Center struct {
			Log      string `json:"log" yaml:"log"`             // 日志
			TcpPort  int    `json:"tcp_port" yaml:"tcp_port"`   // tcp上报服务监听端口
			HttpPort int    `json:"http_port" yaml:"http_port"` // http对外服务监听端口
		} `json:"center" yaml:"center"`

		// worker
		Worker struct {
			Log    string `json:"log" yaml:"log"`       // 日志
			Port   int    `json:"port" yaml:"port"`     // 服务监听端口
			Center string `json:"center" yaml:"center"` // 中心服务器地址
			Report int    `json:"report" yaml:"report"` // 上报时间间隔
		} `json:"worker" yaml:"worker"`
	}

	Frame struct {
		Index  int
		PriKey []byte
		Data   []byte
	}

	FrameRet struct {
		Index int
		Error string
		Data  []byte
	}

	Worker struct {
		IP         string `json:"ip"`
		Port       int    `json:"port"`
		UpdateTime int64  `json:"update_time"`
	}

	workerManager struct {
		sync.RWMutex

		workerList []*Worker
		workerHash map[string]int
		size       int
		capacity   int
	}
)
