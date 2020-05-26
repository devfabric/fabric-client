package client

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	// cfg "github.com/jonluo94/cool/config"
	// "github.com/jonluo94/cool/log"
)

const (
	Admin = "Admin"
)

type FabricClient struct {
	ConnectionFile []byte
	Orgs           []string
	ChannelId      string
	GoPath         string

	userName string
	userOrg  string

	resmgmtClients []*resmgmt.Client
	sdk            *fabsdk.FabricSDK
	retry          resmgmt.RequestOption
}

func (f *FabricClient) Setup() error {
	sdk, err := fabsdk.New(config.FromRaw(f.ConnectionFile, "yaml"))
	if err != nil {
		return err
	}
	f.sdk = sdk

	resmgmtClients := make([]*resmgmt.Client, 0)
	for _, v := range f.Orgs {
		resmgmtClient, err := resmgmt.New(sdk.Context(fabsdk.WithUser(Admin), fabsdk.WithOrg(v)))
		if err != nil {
			return err
		}
		resmgmtClients = append(resmgmtClients, resmgmtClient)
	}
	f.resmgmtClients = resmgmtClients
	//重试
	f.retry = resmgmt.WithRetry(retry.DefaultResMgmtOpts)

	return nil
}

func (f *FabricClient) Close() {
	if f.sdk != nil {
		f.sdk.Close()
	}
}

func (f *FabricClient) SetUser(userName, userOrg string) {
	f.userName = userName
	f.userOrg = userOrg
}

func (f *FabricClient) GetKeyFile(id msp.SigningIdentity) (string, string) {
	priFile := hex.EncodeToString(id.PrivateKey().SKI()) + "_sk"
	pubFile := id.Identifier().ID + "@" + id.Identifier().MSPID + "-cert.pem"
	return priFile, pubFile
}

func (f *FabricClient) RegisterUser(userName string, userOrg string) (priFile string, pubFile string, err error) {
	//secret is userName+userOrg
	secret := userName + userOrg
	mspClient, err := mspclient.New(f.sdk.Context(), mspclient.WithOrg(userOrg))
	if err != nil {
		return "", "", err
	}
	//判断是否存在
	id, err := mspClient.GetSigningIdentity(userName)
	if err == nil {
		priFile, pubFile = f.GetKeyFile(id)
		return
	}
	//注册用户
	request := &mspclient.RegistrationRequest{Name: userName, Type: "client", Secret: secret}
	_, err = mspClient.Register(request)
	if err != nil && !strings.Contains(err.Error(), "is already registered") {
		return "", "", err
	}
	//登记保存证书到stores
	err = mspClient.Enroll(userName, mspclient.WithSecret(secret))
	if err != nil {
		return "", "", err
	}

	id, _ = mspClient.GetSigningIdentity(userName)
	priFile, pubFile = f.GetKeyFile(id)
	return priFile, pubFile, nil
}

func (f *FabricClient) CreateChannel(channelTx string) error {
	mspClient, err := mspclient.New(f.sdk.Context(), mspclient.WithOrg(f.Orgs[0]))
	if err != nil {
		return err
	}
	adminIdentity, err := mspClient.GetSigningIdentity(Admin)
	if err != nil {
		return err
	}
	req := resmgmt.SaveChannelRequest{
		ChannelID:         f.ChannelId,
		ChannelConfigPath: channelTx,
		SigningIdentities: []msp.SigningIdentity{adminIdentity},
	}
	_, err = f.resmgmtClients[0].SaveChannel(req, f.retry)
	if err != nil {
		return err
	}

	return nil
}

func (f *FabricClient) UpdateChannel(anchorsTx []string) error {

	for i, c := range f.resmgmtClients {

		mspClient, err := mspclient.New(f.sdk.Context(), mspclient.WithOrg(f.Orgs[i]))
		if err != nil {

			return err
		}
		adminIdentity, err := mspClient.GetSigningIdentity(Admin)
		if err != nil {

			return err
		}
		req := resmgmt.SaveChannelRequest{
			ChannelID:         f.ChannelId,
			ChannelConfigPath: anchorsTx[i],
			SigningIdentities: []msp.SigningIdentity{adminIdentity},
		}
		_, err = c.SaveChannel(req, f.retry)
		if err != nil {

			return err
		}

	}

	return nil
}

func (f *FabricClient) JoinChannel() error {

	for _, c := range f.resmgmtClients {
		err := c.JoinChannel(f.ChannelId, f.retry)
		if err != nil && !strings.Contains(err.Error(), "LedgerID already exists") {
			return err
		}

	}
	return nil

}

func (f *FabricClient) InstallChaincode(chaincodeId, chaincodePath, version string) error {
	ccPkg, err := gopackager.NewCCPackage(chaincodePath, f.GoPath)
	if err != nil {

		return err
	}

	req := resmgmt.InstallCCRequest{
		Name:    chaincodeId,
		Path:    chaincodePath,
		Version: version,
		Package: ccPkg,
	}

	for _, c := range f.resmgmtClients {
		_, err := c.InstallCC(req, f.retry)
		if err != nil {

			return err
		}

	}

	return nil
}

func (f *FabricClient) InstantiateChaincode(chaincodeId, chaincodePath, version string, policy string, args [][]byte) (string, error) {

	//"OR ('Org1MSP.member','Org2MSP.member')"
	ccPolicy, err := cauthdsl.FromString(policy)
	if err != nil {

		return "", err
	}
	resp, err := f.resmgmtClients[0].InstantiateCC(
		f.ChannelId,
		resmgmt.InstantiateCCRequest{
			Name:    chaincodeId,
			Path:    chaincodePath,
			Version: version,
			Args:    args,
			Policy:  ccPolicy,
		},
		f.retry,
	)

	return string(resp.TransactionID), nil
}

func (f *FabricClient) UpgradeChaincode(chaincodeId, chaincodePath, version string, policy string, args [][]byte) (string, error) {

	f.InstallChaincode(chaincodeId, chaincodePath, version)

	ccPolicy, err := cauthdsl.FromString(policy)
	if err != nil {

		return "", err
	}
	resp, err := f.resmgmtClients[0].UpgradeCC(
		f.ChannelId,
		resmgmt.UpgradeCCRequest{
			Name:    chaincodeId,
			Path:    chaincodePath,
			Version: version,
			Args:    args,
			Policy:  ccPolicy,
		},
		f.retry,
	)
	// logger.Infof("%s", resp.TransactionID)
	return string(resp.TransactionID), nil
}

func (f *FabricClient) QueryLedger() (*FabricBlockchainInfo, error) {
	ledger, err := ledger.New(f.sdk.ChannelContext(f.ChannelId, fabsdk.WithUser(f.userName), fabsdk.WithOrg(f.userOrg)))
	if err != nil {
		// logger.Error(err.Error())
		return nil, err
	}

	bci, err := ledger.QueryInfo()
	if err != nil {
		// logger.Error(err.Error())
		return nil, err
	}

	return parseFabricBlockchainInfo(bci), nil
}

func (f *FabricClient) QueryBlock(height uint64) (*FabricBlock, error) {
	ledger, err := ledger.New(f.sdk.ChannelContext(f.ChannelId, fabsdk.WithUser(f.userName), fabsdk.WithOrg(f.userOrg)))
	if err != nil {
		// logger.Error(err.Error())
		return nil, err
	}

	block, err := ledger.QueryBlock(height)
	bs, err := parseFabricBlock(blockParse(block))

	if err != nil {
		// logger.Error(err.Error())
		return nil, err
	}
	return bs, nil
}

func (f *FabricClient) QueryBlockByHash(hash []byte) (*FabricBlock, error) {

	ledger, err := ledger.New(f.sdk.ChannelContext(f.ChannelId, fabsdk.WithUser(f.userName), fabsdk.WithOrg(f.userOrg)))
	if err != nil {
		// logger.Error(err.Error())
		return nil, err
	}

	block, err := ledger.QueryBlockByHash(hash)
	bs, err := parseFabricBlock(blockParse(block))

	if err != nil {
		// logger.Error(err.Error())
		return nil, err
	}
	return bs, nil
}

func (f *FabricClient) QueryBlockByTxid(txid string) (*FabricBlock, error) {

	ledger, err := ledger.New(f.sdk.ChannelContext(f.ChannelId, fabsdk.WithUser(f.userName), fabsdk.WithOrg(f.userOrg)))
	if err != nil {
		// logger.Error(err.Error())
		return nil, err
	}

	block, err := ledger.QueryBlockByTxID(fab.TransactionID(txid))
	bs, err := parseFabricBlock(blockParse(block))

	if err != nil {
		// logger.Error(err.Error())
		return nil, err
	}
	return bs, nil
}

func (f *FabricClient) QueryTransaction(txid string) (*FabricTransaction, error) {

	ledger, err := ledger.New(f.sdk.ChannelContext(f.ChannelId, fabsdk.WithUser(f.userName), fabsdk.WithOrg(f.userOrg)))
	if err != nil {
		// logger.Error(err.Error())
		return nil, err
	}

	tx, err := ledger.QueryTransaction(fab.TransactionID(txid))
	if err != nil {
		// logger.Error(err.Error())
		return nil, err
	}

	return parseFabricTransaction(tx.GetValidationCode(), tx.GetTransactionEnvelope().Payload, tx.GetTransactionEnvelope().Signature)
}

func (f *FabricClient) QueryChannelConfig() (*FabricChannelConfig, error) {

	ledger, err := ledger.New(f.sdk.ChannelContext(f.ChannelId, fabsdk.WithUser(f.userName), fabsdk.WithOrg(f.userOrg)))
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

func (f *FabricClient) QueryChaincode(chaincodeId, fcn string, args [][]byte) ([]byte, error) {

	client, err := channel.New(f.sdk.ChannelContext(f.ChannelId, fabsdk.WithUser(f.userName), fabsdk.WithOrg(f.userOrg)))
	if err != nil {
		// logger.Error(err.Error())
		return nil, err
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

func (f *FabricClient) InvokeChaincodeWithEvent(chaincodeId, fcn string, args [][]byte) ([]byte, error) {
	eventId := fmt.Sprintf("event%d", time.Now().UnixNano())

	client, err := channel.New(f.sdk.ChannelContext(f.ChannelId, fabsdk.WithUser(f.userName), fabsdk.WithOrg(f.userOrg)))
	if err != nil {
		// logger.Error(err.Error())
		return nil, err
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
	case <-time.After(time.Second * 30):
		// logger.Info("不能根据指定的事件ID接收到相应的链码事件")
		return nil, fmt.Errorf("%s", "等到事件超时")
	}
	return []byte(resp.TransactionID), nil
}

func (f *FabricClient) InvokeChaincode(chaincodeId, fcn string, args [][]byte) ([]byte, error) {

	client, err := channel.New(f.sdk.ChannelContext(f.ChannelId, fabsdk.WithUser(f.userName), fabsdk.WithOrg(f.userOrg)))
	// fabsdk.WithTargets()
	if err != nil {
		// logger.Error(err.Error())
		return nil, err
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

func NewFabricClient(connectionFile []byte, channelId string, orgs []string) *FabricClient {
	fabric := &FabricClient{
		ConnectionFile: connectionFile,
		ChannelId:      channelId,
		Orgs:           orgs,
		GoPath:         os.Getenv("GOPATH"),
	}

	return fabric
}
