package main

import (
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
	fmt.Println("runDir=", runDir)

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

	blockInfo, err := fabric.QueryLedger()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(blockInfo)

	err = fabric.RegisterBlockEvent(doEventBlockProcess, doEventTxProcess)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("exit waitting....")
	select {}
}

func doEventBlockProcess(blockInfo *fabsdk.Block) error {
	fmt.Println(blockInfo.Height)
	fmt.Println(blockInfo.DataHash)
	fmt.Println(blockInfo.PreviousHash)

	// events, err := fabric.GetEventFromBlock(28)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// for i, v := range events {
	// 	fmt.Println(i, v)
	// 	fmt.Println(i, string(v.Payload))
	// }
	return nil
}

func doEventTxProcess(blockHeight uint64, txEv *fabsdk.Event) error {
	fmt.Println(blockHeight)
	fmt.Println(txEv.ChaincodeId, txEv.EventName, txEv.TxId)
	fmt.Println(string(txEv.Payload))
	return nil
}
