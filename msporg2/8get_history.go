package main

import (
	"encoding/json"
	"fmt"

	"fabric-client/config"
	"fabric-client/fabsdk"
	"os"
	"path/filepath"
	"strings"
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
	fmt.Println("runDir=", runDir)

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

	getInfo := &GetInfo{
		UserName: "user1",
		CardType: 1,
		CardID:   "user1-100001",
	}

	getInfoBys, err := json.Marshal(getInfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//printf
	{
		arrayList1, err := json.MarshalIndent(getInfo, "", " ")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(string(arrayList1))
	}

	payLoad, err := fabric.QueryChaincode(fabConfig.ChaincodeID, "User1", "get_history", [][]byte{getInfoBys})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(payLoad))

}
