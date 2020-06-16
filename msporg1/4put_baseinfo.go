package main

import (
	"encoding/json"
	"fmt"

	"os"
	"path/filepath"
	"strings"

	"github.com/devfabric/fabric-client/config"
	"github.com/devfabric/fabric-client/fabsdk"
)

func GetCurrentDirectory() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}
	return strings.Replace(dir, "\\", "/", -1), nil
}

func main() {
	runDir, err := GetCurrentDirectory()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//设置环境变量，防止应用未设置
	workDirForFabSDK := os.Getenv("WORKDIR")
	if workDirForFabSDK == "" {
		os.Setenv("WORKDIR", runDir)
	}

	fabConfig, err := config.LoadFabircConfig(runDir)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fabric := fabsdk.NewFabricClient(fabConfig.ConfigFile, fabConfig.ChannelID, fabConfig.UserName, fabConfig.UserOrg)
	err = fabric.Setup(runDir)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	putBaseList := make([]*PutBaseInfo, 0)
	putItem1 := &PutBaseInfo{
		CardIDMap: map[uint8]string{
			1: "身份证S3000",
			2: "护照H3000",
			3: "军官证J3000",
		},
		PersonInfo: &PersonInfo{
			FieldKVMap: map[string]interface{}{
				"姓名":    "明明",
				"性别":    "男",
				"身份证":   "身份证S3000",
				"护照":    "护照H3000",
				"军官证":   "军官证J3000",
				"证件有效期": "2025-10-1",
				"出生日期":  "1982.10.1",
				"职业":    "飞行员",
				"户口所在地": "北京市朝阳区XXX",
				"工作地址":  "北京市海淀区XXXX",
				"联系电话":  "177XXXXXX",
			},
		},
	}

	putItem2 := &PutBaseInfo{
		CardIDMap: map[uint8]string{
			1: "身份证S4000",
			2: "护照H4000",
			3: "军官证J4000",
		},
		PersonInfo: &PersonInfo{
			FieldKVMap: map[string]interface{}{
				"姓名":    "晨晨",
				"性别":    "女",
				"身份证":   "身份证S4000",
				"护照":    "护照H4000",
				"军官证":   "军官证J4000",
				"证件有效期": "2090-10-1",
				"出生日期":  "1985.10.1",
				"职业":    "工程师",
				"户口所在地": "北京市朝阳区XXX",
				"工作地址":  "北京市海淀区XXXX",
				"联系电话":  "15699XXXXXX",
			},
		},
	}

	putBaseList = append(putBaseList, putItem1)
	putBaseList = append(putBaseList, putItem2)

	arrayList, err := json.Marshal(putBaseList)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//printf
	{
		dataBytes, err := json.MarshalIndent(putBaseList, "", " ")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(string(dataBytes))
	}

	payLoad, err := fabric.InvokeChaincodeWithEvent(fabConfig.ChaincodeID, "User1", "put_baseinfo", [][]byte{arrayList, []byte(EvPutInfo)})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(payLoad))

}
