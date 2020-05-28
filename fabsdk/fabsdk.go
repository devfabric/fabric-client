package fabsdk

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	contextApi "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	pfab "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/events/deliverclient/seek"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

type ListenBlock struct {
	registration fab.Registration
	blockChain   <-chan *pfab.BlockEvent
}

type FabricClient struct {
	ConnectionFile string
	ChannelID      string

	DefaultName string
	DefaultOrg  string

	GoPath         string
	sdk            *fabsdk.FabricSDK
	resmgmtClients []*resmgmt.Client
	retry          resmgmt.RequestOption

	ledgerClient *ledger.Client //获取最新区块信息客户端

	eventClient *event.Client //事件监听客户端
	listenEvent ListenBlock
	isExit      chan struct{}
}

func NewFabricClient(connectionFile string, channelId string, name string, orgs string) *FabricClient {
	return &FabricClient{
		ConnectionFile: connectionFile,
		ChannelID:      channelId,
		DefaultName:    name,
		DefaultOrg:     orgs,
		GoPath:         os.Getenv("GOPATH"),
	}
}

func (fab *FabricClient) Setup(rootDir string) error {
	var (
		err                      error
		org1ChannelClientContext contextApi.ChannelProvider
	)
	fab.sdk, err = fabsdk.New(config.FromFile(fab.ConnectionFile))
	if err != nil {
		return err
	}

	resmgmtClients := make([]*resmgmt.Client, 0)
	resmgmtClient, err := resmgmt.New(fab.sdk.Context(fabsdk.WithUser(Admin), fabsdk.WithOrg(fab.DefaultOrg)))
	if err != nil {
		return err
	}
	resmgmtClients = append(resmgmtClients, resmgmtClient)
	//重试
	fab.retry = resmgmt.WithRetry(retry.DefaultResMgmtOpts)

	//block
	{
		fab.ledgerClient, err = ledger.New(fab.sdk.ChannelContext(fab.ChannelID, fabsdk.WithUser(fab.DefaultName), fabsdk.WithOrg(fab.DefaultOrg)))
		if err != nil {
			return err
		}
	}

	//event
	{
		blockNum, err := GetBlockHeight(rootDir)
		if err != nil {
			return err
		}

		org1ChannelClientContext = fab.sdk.ChannelContext(fab.ChannelID, fabsdk.WithUser(fab.DefaultName), fabsdk.WithOrg(fab.DefaultOrg))
		evnetOpts, err := newEvnetOpts("from", blockNum)
		if err != nil {
			return err
		}
		fab.eventClient, err = event.New(org1ChannelClientContext, evnetOpts...)
		if err != nil {
			return err
		}
		fab.listenEvent = ListenBlock{}

	}
	fab.isExit = make(chan struct{})
	return nil
}

func (fab *FabricClient) Close() {
	if fab.isExit != nil {
		close(fab.isExit)
	}

	if fab.listenEvent.registration != nil {
		fab.eventClient.Unregister(fab.listenEvent.registration)
	}

	if fab.sdk != nil {
		fab.sdk.Close()
	}
}

//写块高
const BlockInfoFile = "./configs/blockInfo"

type BlockInfo struct {
	Height uint64 `json:"height"`
}

func GetBlockHeight(dir string) (uint64, error) {
	blockInfoFile := filepath.Join(dir, BlockInfoFile)
	_, err := os.Stat(blockInfoFile)
	if err != nil {
		var initHeight uint64
		blockInfo := BlockInfo{Height: initHeight}
		bytes, err := json.Marshal(&blockInfo)
		if err != nil {
			return 0, err
		}
		err = ioutil.WriteFile(blockInfoFile, bytes, os.ModePerm)
		if err != nil {
			return 0, err
		}
		return initHeight, nil
	}

	bytes, err := ioutil.ReadFile(blockInfoFile)
	if err != nil {
		return 0, err
	}

	blockInfo := BlockInfo{}
	err = json.Unmarshal(bytes, &blockInfo)
	if err != nil {
		return 0, err
	}

	return blockInfo.Height, nil
}

func UpdateBlockHeight(dir string, height uint64) error {
	blockInfo := BlockInfo{Height: height}
	bytes, err := json.Marshal(&blockInfo)
	if err != nil {
		return err
	}

	blockInfoFile := filepath.Join(dir, BlockInfoFile)
	err = ioutil.WriteFile(blockInfoFile, bytes, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func newEvnetOpts(seekType string, blockNum uint64) ([]event.ClientOption, error) {
	if seekType != "oldest" && seekType != "newest" && seekType != "from" {
		return nil, errors.New("seek type error,must be one of [oldest,newest,from]")
	}

	evnetOpts := []event.ClientOption{}
	evnetOpts = append(evnetOpts, event.WithBlockEvents())
	evnetOpts = append(evnetOpts, event.WithSeekType(seek.Type(seekType)))

	if seekType == "from" {
		evnetOpts = append(evnetOpts, event.WithBlockNum(blockNum))
	}
	return evnetOpts, nil
}
