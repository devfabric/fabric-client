package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	client "github.com/devfabric/fabric-client/client"
	config "github.com/devfabric/fabric-client/config"
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

	connectConfig, err := ioutil.ReadFile(fabConfig.ConfigFile)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fabric := client.NewFabricClient(connectConfig, fabConfig.ChannelId, fabConfig.MulOrgs)
	defer fabric.Close()
	err = fabric.Setup()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//查询状态
	ledger, err := fabric.QueryLedger()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	log.Println((ledger))

	//查询账本
	queryFcn := "query"
	queryArgs := [][]byte{[]byte("a")}
	a, _ := fabric.QueryChaincode(fabConfig.ChannelId, queryFcn, queryArgs)
	log.Println("a的值: ", string(a))

}
