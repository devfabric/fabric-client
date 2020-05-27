package fabsdk

import (
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	pfab "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	// cm "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	// "github.com/hyperledger/fabric-protos-go/common"
	// cm "github.com/hyperledger/fabric/protos/common"
	// "github.com/hyperledger/fabric/protos/common"
	// cfg "github.com/jonluo94/cool/config"
)

type FabricClient struct {
	ConnectionFile string
	ChannelID      string

	DefaultName string
	DefaultOrg  string

	GoPath         string
	sdk            *fabsdk.FabricSDK
	resmgmtClients []*resmgmt.Client
	retry          resmgmt.RequestOption
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

func (fab *FabricClient) Setup() error {
	var (
		err error
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

	return nil
}

func (fab *FabricClient) QueryLedger() (*FabricBlockchainInfo, error) {
	ledger, err := ledger.New(fab.sdk.ChannelContext(fab.ChannelID, fabsdk.WithUser(fab.DefaultName), fabsdk.WithOrg(fab.DefaultOrg)))
	if err != nil {
		return nil, err
	}

	bci, err := ledger.QueryInfo()
	if err != nil {
		return nil, err
	}
	return parseFabricBlockchainInfo(bci), nil
}

func (fab *FabricClient) QueryBlock(height uint64) (*FabricBlock, error) {
	ledger, err := ledger.New(fab.sdk.ChannelContext(fab.ChannelID, fabsdk.WithUser(fab.DefaultName), fabsdk.WithOrg(fab.DefaultOrg)))
	if err != nil {
		return nil, err
	}

	block, err := ledger.QueryBlock(height)
	if err != nil {
		return nil, err
	}

	bs, err := parseFabricBlock(blockParse(block))
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func (fab *FabricClient) QueryBlockByHash(hash string) (*FabricBlock, error) {
	ledger, err := ledger.New(fab.sdk.ChannelContext(fab.ChannelID, fabsdk.WithUser(fab.DefaultName), fabsdk.WithOrg(fab.DefaultOrg)))
	if err != nil {
		return nil, err
	}
	hashbys, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}

	block, err := ledger.QueryBlockByHash(hashbys)
	if err != nil {
		return nil, err
	}

	bs, err := parseFabricBlock(blockParse(block))
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func (fab *FabricClient) QueryBlockByTxid(txid string) (*FabricBlock, error) {
	ledger, err := ledger.New(fab.sdk.ChannelContext(fab.ChannelID, fabsdk.WithUser(fab.DefaultName), fabsdk.WithOrg(fab.DefaultOrg)))
	if err != nil {
		return nil, err
	}

	block, err := ledger.QueryBlockByTxID(pfab.TransactionID(txid))
	if err != nil {
		return nil, err
	}

	bs, err := parseFabricBlock(blockParse(block))
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func (fab *FabricClient) QueryTransaction(txid string) (*FabricTransaction, error) {
	ledger, err := ledger.New(fab.sdk.ChannelContext(fab.ChannelID, fabsdk.WithUser(fab.DefaultName), fabsdk.WithOrg(fab.DefaultOrg)))
	if err != nil {
		// logger.Error(err.Error())
		return nil, err
	}

	tx, err := ledger.QueryTransaction(pfab.TransactionID(txid))
	if err != nil {
		// logger.Error(err.Error())
		return nil, err
	}

	return parseFabricTransaction(tx.GetValidationCode(), tx.GetTransactionEnvelope().Payload, tx.GetTransactionEnvelope().Signature)
}

func (fab *FabricClient) QueryChannelConfig() (*FabricChannelConfig, error) {
	ledger, err := ledger.New(fab.sdk.ChannelContext(fab.ChannelID, fabsdk.WithUser(fab.DefaultName), fabsdk.WithOrg(fab.DefaultOrg)))
	if err != nil {
		// logger.Error(err.Error())
		return nil, err
	}

	cfg, err := ledger.QueryConfig()
	if err != nil {
		// logger.Error(err.Error())
		return nil, err
	}
	return parseChannelConfig(cfg), nil
}

func (fab *FabricClient) QueryChaincode(chaincodeId, uname string, fcn string, args [][]byte) ([]byte, error) {
	var (
		client *channel.Client
		err    error
	)
	if uname != "" {
		client, err = channel.New(fab.sdk.ChannelContext(fab.ChannelID, fabsdk.WithUser(uname), fabsdk.WithOrg(fab.DefaultOrg)))
		if err != nil {
			return nil, err
		}
	} else {
		client, err = channel.New(fab.sdk.ChannelContext(fab.ChannelID, fabsdk.WithUser(fab.DefaultName), fabsdk.WithOrg(fab.DefaultOrg)))
		if err != nil {
			return nil, err
		}
	}

	resp, err := client.Query(channel.Request{
		ChaincodeID: chaincodeId,
		Fcn:         fcn,
		Args:        args,
	})
	if err != nil {
		// logger.Error(err.Error())
		return nil, err
	}
	// logger.Infof(string(resp.Payload))
	return resp.Payload, nil
}

func (fab *FabricClient) InvokeChaincodeWithEvent(chaincodeId, uname string, fcn string, args [][]byte) ([]byte, error) {
	eventId := fmt.Sprintf("event%d", time.Now().UnixNano())

	var (
		client *channel.Client
		err    error
	)

	if uname != "" {
		client, err = channel.New(fab.sdk.ChannelContext(fab.ChannelID, fabsdk.WithUser(uname), fabsdk.WithOrg(fab.DefaultOrg)))
		if err != nil {
			// logger.Error(err.Error())
			return nil, err
		}
	} else {
		client, err = channel.New(fab.sdk.ChannelContext(fab.ChannelID, fabsdk.WithUser(fab.DefaultName), fabsdk.WithOrg(fab.DefaultOrg)))
		if err != nil {
			// logger.Error(err.Error())
			return nil, err
		}
	}

	// 注册事件
	reg, notifier, err := client.RegisterChaincodeEvent(chaincodeId, eventId)
	if err != nil {
		// logger.Errorf("注册链码事件失败: %s", err)
		return nil, err
	}
	defer client.UnregisterChaincodeEvent(reg)

	req := channel.Request{
		ChaincodeID: chaincodeId,
		Fcn:         fcn,
		Args:        append(args, []byte(eventId)),
	}
	resp, err := client.Execute(req)
	if err != nil {
		// logger.Error(err.Error())
		return nil, err
	}

	select {
	case ccEvent := <-notifier:
		// logger.Infof("接收到链码事件: %v\n", ccEvent)
		return []byte(ccEvent.TxID), nil
	case <-time.After(time.Second * 3):
		// logger.Info("不能根据指定的事件ID接收到相应的链码事件")
		return nil, fmt.Errorf("%s", "等到事件超时")
	}
	return []byte(resp.TransactionID), nil
}

func (fab *FabricClient) InvokeChaincode(chaincodeId, uname string, fcn string, args [][]byte) ([]byte, error) {
	var (
		client *channel.Client
		err    error
	)

	if uname != "" {
		client, err = channel.New(fab.sdk.ChannelContext(fab.ChannelID, fabsdk.WithUser(uname), fabsdk.WithOrg(fab.DefaultOrg)))
		// fabsdk.WithTargets()
		if err != nil {
			// logger.Error(err.Error())
			return nil, err
		}
	} else {
		client, err = channel.New(fab.sdk.ChannelContext(fab.ChannelID, fabsdk.WithUser(fab.DefaultName), fabsdk.WithOrg(fab.DefaultOrg)))
		// fabsdk.WithTargets()
		if err != nil {
			// logger.Error(err.Error())
			return nil, err
		}
	}

	req := channel.Request{
		ChaincodeID: chaincodeId,
		Fcn:         fcn,
		Args:        args,
	}
	resp, err := client.Execute(req)
	if err != nil {
		// logger.Error(err.Error())
		return nil, err
	}
	return []byte(resp.TransactionID), nil
}
