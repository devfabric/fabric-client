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

	// connectConfig, err := ioutil.ReadFile(fabConfig.ConfigFile)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }

	fabric := fabsdk.NewFabricClient(fabConfig.ConfigFile, fabConfig.ChannelID, fabConfig.UserName, fabConfig.UserOrg)
	err = fabric.Setup(runDir)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	putFavorList := make([]*PutFavor, 0)
	putItem1 := &PutFavor{
		CardType: 1,
		CardID:   "身份证S3000",
		IndemCard: &IndemCard{
			FieldKVMap: map[string]interface{}{
				"社保账号": "身份证S3000",
				"民生卡号": "M3000",
				"卡状态":  1,
			},
		},
		BankCard: &BankCard{
			FieldKVMap: map[string]interface{}{
				"银行卡号": "Y3000",
				"发卡银行": "工商银行",
			},
		},
		AssetsInfo: &AssetsInfo{
			FieldKVMap: map[string]interface{}{
				"社保身份": 0,
			},
		},
	}

	putItem2 := &PutFavor{
		CardType: 1,
		CardID:   "身份证S4000",
		IndemCard: &IndemCard{
			FieldKVMap: map[string]interface{}{
				"社保账号": "身份证S4000",
				"民生卡号": "M4000",
				"卡状态":  1,
			},
		},
		BankCard: &BankCard{
			FieldKVMap: map[string]interface{}{
				"银行卡号": "Y4000",
				"发卡银行": "北京银行",
			},
		},
		AssetsInfo: &AssetsInfo{
			FieldKVMap: map[string]interface{}{
				"社保身份": 0,
			},
		},
	}

	putFavorList = append(putFavorList, putItem1)
	putFavorList = append(putFavorList, putItem2)

	arrayList, err := json.Marshal(putFavorList)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//printf
	{
		arrayList1, err := json.MarshalIndent(putFavorList, "", " ")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(string(arrayList1))
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

	payLoad, err := fabric.InvokeChaincodeWithEvent(fabConfig.ChaincodeID, "User1", "put_favorinfo", [][]byte{arrayList, []byte(EvPutFavorInfo)})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(payLoad))

}
