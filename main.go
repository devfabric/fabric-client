package main

import (
	"fmt"
	"io/ioutil"
	"log"

	client "github.com/devfabric/fabric-client/client"
	config "github.com/devfabric/fabric-client/config"
)

func mmain() {
	fabConfig, err := config.LoadHPCacheConfig("./")
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

}
