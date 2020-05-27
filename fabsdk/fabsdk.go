package fabsdk

import (
	"os"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
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
