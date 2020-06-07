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

	fabric := fabsdk.NewFabricClient(fabConfig.ConfigFile, fabConfig.ChannelID, fabConfig.UserName, fabConfig.UserOrg)
	err = fabric.Setup(runDir)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	updateData := make([]*UpdateBaseInfo, 0)
	updateData1 := &UpdateBaseInfo{
		UserName: "user1",
		CardType: 1,
		CardID:   "user1-100001",
		CardIDMap: map[uint8]string{
			4: "user1-100004",
		},
		FieldKVMap: map[string]interface{}{
			"sex":   "女",
			"card4": "user1-100004",
		},
	}

	updateData2 := &UpdateBaseInfo{
		UserName: "user2",
		CardType: 1,
		CardID:   "user2-200001",
		FieldKVMap: map[string]interface{}{
			"name":  "user2",
			"sex":   "女",
			"card1": "user2-200001",
			"card2": "user2-200002",
			"card3": "user2-200003",
		},
	}

	updateData = append(updateData, updateData1)
	updateData = append(updateData, updateData2)

	arrayList, err := json.Marshal(updateData)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//printf
	{
		arrayList1, err := json.MarshalIndent(updateData, "", " ")
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

	payLoad, err := fabric.InvokeChaincodeWithEvent(fabConfig.ChaincodeID, "User1", "update_baseinfo", [][]byte{arrayList, []byte(EvUpdateBaseInfo)})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(payLoad))

}
