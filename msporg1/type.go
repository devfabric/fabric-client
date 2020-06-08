package main

type PutBaseInfo struct {
	UserName   string                 `json:"username"`
	CardIDMap  map[uint8]string       `json:"cardidmap"`
	FieldKVMap map[string]interface{} `json:"fields,omitempty"` //可以为空
}

type UpdateBaseInfo struct {
	UserName   string                 `json:"username"`
	CardType   uint8                  `json:"cardtype"`
	CardID     string                 `json:"cardid"`
	CardIDMap  map[uint8]string       `json:"cardidmap,omitempty"` //可以为空
	FieldKVMap map[string]interface{} `json:"fields"`              //不可以为空
}

type PutInfo struct {
	UserName   string                 `json:"username"`
	CardType   uint8                  `json:"cardtype"`
	CardID     string                 `json:"cardid"`
	FieldKVMap map[string]interface{} `json:"fields"`
}

type UpdateInfo = PutInfo

const (
	EvPutBaseInfo    string = "EvPutBaseInfo"
	EvUpdateBaseInfo string = "EvUpdateBaseInfo"

	EvPutInfo    string = "EvPutInfo"
	EvUpdateInfo string = "EvUpdateInfo"
)

type GetInfo struct {
	UserName string `json:"username"`
	CardType uint8  `json:"cardtype"`
	CardID   string `json:"cardid"`
}

//插叙条件
type QueryCond struct {
	Field map[string]bool `json:"fields"`
	Limit int32           `json:"limit"`
	Skip  int32           `json:"skip"`
}

type SimpleCond struct {
	Limit int32 `json:"limit"`
	Skip  int32 `json:"skip"`
}
