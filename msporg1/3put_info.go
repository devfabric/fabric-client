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

	putInfoList := make([]*PutInfo, 0)
	putItem1 := &PutInfo{
		CardIDMap: map[uint8]string{
			1: "身份证S1000",
			2: "护照H1000",
			3: "军官证J1000",
		},
		PersonInfo: &PersonInfo{
			FieldKVMap: map[string]interface{}{
				"姓名":    "张三",
				"性别":    "男",
				"身份证":   "身份证S1000",
				"护照":    "护照H1000",
				"军官证":   "军官证J1000",
				"证件有效期": "2025-10-1",
				"出生日期":  "1982.10.1",
				"职业":    "工程师",
				"户口所在地": "北京市朝阳区XXX",
				"工作地址":  "北京市海淀区XXXX",
				"联系电话":  "185XXXXXX",
			},
		},
		AssetsInfo: &AssetsInfo{
			FieldKVMap: map[string]interface{}{
				"社保身份": 0,
				"残疾级别": 1,
			},
		},
	}

	putItem2 := &PutInfo{
		CardIDMap: map[uint8]string{
			1: "身份证S2000",
			2: "护照H2000",
			3: "军官证J2000",
		},
		PersonInfo: &PersonInfo{
			FieldKVMap: map[string]interface{}{
				"姓名":    "李四",
				"性别":    "女",
				"身份证":   "身份证S2000",
				"护照":    "护照H2000",
				"军官证":   "军官证J2000",
				"证件有效期": "2090-10-1",
				"出生日期":  "1985.10.1",
				"职业":    "工程师",
				"户口所在地": "北京市朝阳区XXX",
				"工作地址":  "北京市海淀区XXXX",
				"联系电话":  "156XXXXXX",
			},
		},
		AssetsInfo: &AssetsInfo{
			FieldKVMap: map[string]interface{}{
				"社保身份": 0,
				"残疾级别": 0,
			},
		},
	}

	putInfoList = append(putInfoList, putItem1)
	putInfoList = append(putInfoList, putItem2)

	arrayList, err := json.Marshal(putInfoList)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//printf
	{
		dataBytes, err := json.MarshalIndent(putInfoList, "", " ")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(string(dataBytes))
	}

	// // fmt.Println("success:", string(arrayList))
	// {
	// 	var puDataTest []*PutDataReq
	// 	err := json.Unmarshal(arrayList, &puDataTest)
	// 	if err != nil {
	// 		fmt.Println(err.Error())
	// 		return
	// 	}
	// 	// fmt.Println("success1:", puDataTest)
	// 	for i := range puDataTest {
	// 		fmt.Println(puDataTest[i])
	// 	}
	// }

	payLoad, err := fabric.InvokeChaincodeWithEvent(fabConfig.ChaincodeID, "User1", "put_info", [][]byte{arrayList, []byte(EvPutInfo)})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(payLoad))

}
