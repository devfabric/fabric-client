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

	fabConfig, err := config.LoadHPCacheConfig(runDir)
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

	updateInfoList := make([]*UpdateInfo, 0)
	updateData1 := &UpdateInfo{
		UserName: "user1",
		CardType: 1,
		CardID:   "user1-100001",
		FieldKVMap: map[string]interface{}{
			"bank1": "user1-bank1-000001-update",
			"bank2": "user1-bank2-000001-update",
		},
	}

	updateData2 := &UpdateInfo{
		UserName: "user2",
		CardType: 1,
		CardID:   "user2-200001",
		FieldKVMap: map[string]interface{}{
			"bank1": "user2-bank1-000001-update",
			"bank2": "user2-bank2-000001-update",
		},
	}

	updateInfoList = append(updateInfoList, updateData1)
	updateInfoList = append(updateInfoList, updateData2)

	arrayList, err := json.Marshal(updateInfoList)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//printf
	{
		arrayList1, err := json.MarshalIndent(updateInfoList, "", " ")
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

	payLoad, err := fabric.InvokeChaincodeWithEvent(fabConfig.ChaincodeID, "User1", "update_info", [][]byte{arrayList, []byte(EvPutInfo)})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(payLoad))

}