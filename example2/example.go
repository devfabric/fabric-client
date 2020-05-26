package main

import (
	"io/ioutil"
	"log"
	"time"

	client "github.com/devfabric/fabric-client/client"
)

func main() {
	orgs := []string{"Org1", "Org2"}
	channelId := "mychannel"
	connectConfig, _ := ioutil.ReadFile("./first-network.yaml")
	chaincodeId := "mycc"

	/*操作fabric start*/
	fabric := client.NewFabricClient(connectConfig, channelId, orgs)
	defer fabric.Close()
	fabric.Setup()

	//查询状态
	ledger, _ := fabric.QueryLedger()
	log.Println((ledger))

	//查询账本
	queryFcn := "query"
	queryArgs := [][]byte{[]byte("a")}
	a, _ := fabric.QueryChaincode(chaincodeId, queryFcn, queryArgs)
	log.Println("a的值: ", string(a))
	//invoke 账本
	invokeFcn := "invoke"
	invokeArgs := [][]byte{[]byte("a"), []byte("b"), []byte("10")}
	txid, _ := fabric.InvokeChaincode(chaincodeId, invokeFcn, invokeArgs)
	log.Println(string(txid))
	time.Sleep(10 * time.Second)
	//查询账本
	a, _ = fabric.QueryChaincode(chaincodeId, queryFcn, queryArgs)
	log.Println("a的值: ", string(a))

	//invoke账本
	txid, _ = fabric.InvokeChaincode(chaincodeId, invokeFcn, invokeArgs)
	log.Println(string(txid))
	time.Sleep(10 * time.Second)
	//查询账本
	a, _ = fabric.QueryChaincode(chaincodeId, queryFcn, queryArgs)
	log.Println("a的值: ", string(a))

}
