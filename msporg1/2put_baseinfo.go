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

	puData := make([]*PutBaseInfo, 0)
	putData1 := &PutBaseInfo{
		UserName: "user1",
		CardIDMap: map[uint8]string{
			1: "user1-100001",
			2: "user1-100002",
			3: "user1-100003",
		},
		FieldKVMap: map[string]interface{}{
			"name":  "user1",
			"sex":   "男",
			"card1": "user1-100001",
			"card2": "user1-100002",
			"card3": "user1-100003",
		},
	}

	putData2 := &PutBaseInfo{
		UserName: "user2",
		CardIDMap: map[uint8]string{
			1: "user2-200001",
			2: "user2-200002",
			3: "user2-200003",
		},
		FieldKVMap: map[string]interface{}{
			"name":  "user2",
			"sex":   "女",
			"card1": "user2-200001",
			"card2": "user2-200002",
			"card3": "user2-200003",
		},
	}

	puData = append(puData, putData1)
	puData = append(puData, putData2)

	arrayList, err := json.Marshal(puData)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//printf
	{
		arrayList1, err := json.MarshalIndent(puData, "", " ")
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

	payLoad, err := fabric.InvokeChaincodeWithEvent(fabConfig.ChaincodeID, "User1", "put_baseinfo", [][]byte{arrayList, []byte(EvPutBaseInfo)})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(payLoad))

}
