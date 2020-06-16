package main

//个人基本信息
type PersonInfo struct {
	FieldKVMap map[string]interface{} `json:"fields"`
}

//社保卡信息
type IndemCard struct {
	FieldKVMap map[string]interface{} `json:"fields"`
}

//银行卡
type BankCard struct {
	FieldKVMap map[string]interface{} `json:"fields"`
}

//个人基本权益信息
type AssetsInfo struct {
	FieldKVMap map[string]interface{} `json:"fields"`
}

type PutInfo struct {
	CardIDMap  map[uint8]string `json:"cardidmap"`
	PersonInfo *PersonInfo      `json:"personinfo,omitempty"`
	AssetsInfo *AssetsInfo      `json:"assetsinfo,omitempty"`
}

type PutBaseInfo struct {
	CardIDMap  map[uint8]string `json:"cardidmap"`
	PersonInfo *PersonInfo      `json:"personinfo,omitempty"`
}

type UpdateBaseInfo struct {
	CardType   uint8            `json:"cardtype"`
	CardID     string           `json:"cardid"`
	CardIDMap  map[uint8]string `json:"cardidmap,omitempty"` //可以为空
	PersonInfo PersonInfo       `json:"personinfo"`
}

type PutFavor struct {
	CardType   uint8       `json:"cardtype"`
	CardID     string      `json:"cardid"`
	IndemCard  *IndemCard  `json:"indemcard,omitempty"`  //可忽略
	BankCard   *BankCard   `json:"bankcard,omitempty"`   //可忽略
	AssetsInfo *AssetsInfo `json:"assetsinfo,omitempty"` //可忽略
}

type UpdateFavor = PutFavor

type GetInfo struct {
	CardType uint8  `json:"cardtype"`
	CardID   string `json:"cardid"`
}

///////////////////////////////////////////////////
// type PutBaseInfo struct {
// 	UserName   string                 `json:"username"`
// 	CardIDMap  map[uint8]string       `json:"cardidmap"`
// 	FieldKVMap map[string]interface{} `json:"fields,omitempty"` //可以为空
// }

// type UpdateBaseInfo struct {
// 	UserName   string                 `json:"username"`
// 	CardType   uint8                  `json:"cardtype"`
// 	CardID     string                 `json:"cardid"`
// 	CardIDMap  map[uint8]string       `json:"cardidmap,omitempty"` //可以为空
// 	FieldKVMap map[string]interface{} `json:"fields"`              //不可以为空
// }

// type UpdateInfo = PutInfo

const (
	EvRegisterOrg     string = "EvRegisterOrg"
	EvPutInfo         string = "EvPutInfo"
	EvPutBaseInfo     string = "EvPutBaseInfo"
	EvUpdateBaseInfo  string = "EvUpdateBaseInfo"
	EvPutFavorInfo    string = "EvPutFavorInfo"
	EvUpdateFavorInfo string = "EvUpdateFavorInfo"
	EvConfirmBaseInfo string = "EvConfirmBaseInfo"
	// confirm_baseinfo
	// EvUpdateInfo string = "EvUpdateInfo"
)

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
