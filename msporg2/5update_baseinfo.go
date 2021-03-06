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

	updateBaseList := make([]*UpdateBaseInfo, 0)
	updateItem1 := &UpdateBaseInfo{
		CardType: 1,
		CardID:   "身份证S3000",
		CardIDMap: map[uint8]string{
			4: "国家统一号码G3000",
		},
		PersonInfo: PersonInfo{
			FieldKVMap: map[string]interface{}{
				"职业":     "科学家",
				"国家统一号码": "国家统一号码G1000",
			},
		},
	}

	updateItem2 := &UpdateBaseInfo{
		CardType: 1,
		CardID:   "身份证S4000",
		CardIDMap: map[uint8]string{
			4: "国家统一号码G4000",
		},
		PersonInfo: PersonInfo{
			FieldKVMap: map[string]interface{}{
				"职业":     "医生",
				"国家统一号码": "国家统一号码G4000",
			},
		},
	}

	updateBaseList = append(updateBaseList, updateItem1)
	updateBaseList = append(updateBaseList, updateItem2)

	arrayList, err := json.Marshal(updateBaseList)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//printf
	{
		arrayList1, err := json.MarshalIndent(updateBaseList, "", " ")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(string(arrayList1))
	}

	payLoad, err := fabric.InvokeChaincodeWithEvent(fabConfig.ChaincodeID, "User1", "update_baseinfo", [][]byte{arrayList, []byte(EvUpdateBaseInfo)})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(payLoad))

}
