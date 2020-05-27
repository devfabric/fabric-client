package fabsdk

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	// cm "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	// "github.com/hyperledger/fabric-protos-go/common"
	// cm "github.com/hyperledger/fabric/protos/common"
	// "github.com/hyperledger/fabric/protos/common"
	// cfg "github.com/jonluo94/cool/config"
)

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
