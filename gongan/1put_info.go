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
		SeqNum: "BN1000",
		Data:   "xml-encode-bytes",
	}

	putItem2 := &PutInfo{
		SeqNum: "BN2000",
		Data:   "xml-encode-bytes",
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
