package main

type PutInfo struct {
	PackID string `json:"packid"` //批次号
	CardID string `json:"cardid"` //身份证
	Value  []byte `json:"value"`  //xml文件内容
	Dtime  int64  `json:"dtime"`  //上链时间
	SysID  int    `json:"sysid"`  //链外系统id
}

type GetInfo struct {
	PackID string `json:"packid"`
	CardID string `json:"cardid,omitempty"`
}

const (
	EvPutInfo string = "EvPutInfo"
)
