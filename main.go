package main

import (
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

	blockInfo, err := fabric.QueryLedger()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(blockInfo)

	// binfo, err := fabric.QueryBlock(2)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// fmt.Println(binfo)

	// binfo, err := fabric.QueryBlockByHash("3eba885e44edfe4293797ffeef568d777ab052793941f875d2a6d8000d51ca40")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// fmt.Println(binfo)

	payLoad, err := fabric.QueryChaincode(fabConfig.ChaincodeID, "User1", "check", nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(payLoad))

	payLoad, err = fabric.InvokeChaincodeWithEvent(fabConfig.ChaincodeID, "User1", "reg_org", [][]byte{[]byte("ss")})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(payLoad))

	// payLoad, err = fabric.InvokeChaincode(fabConfig.ChaincodeID, "User1", "reg_org", [][]byte{[]byte("ss")})
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// fmt.Println(string(payLoad))
	//90d3ea7c062cc5b9ad82a991ef15ac0de238a1b45ae5dd8e5eb8127178fa6bde

	// tx, err := fabric.QueryTransaction("90d3ea7c062cc5b9ad82a991ef15ac0de238a1b45ae5dd8e5eb8127178fa6bde")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// fmt.Println(tx)

	// channelInfo, err := fabric.QueryChannelConfig()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// fmt.Println(channelInfo)

	// _, _, err = fabric.RegisterUser("liuhy")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }

	events, err := fabric.GetEventFromBlock(28)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for i, v := range events {
		fmt.Println(i, v)
		fmt.Println(i, string(v.Payload))
	}
}
