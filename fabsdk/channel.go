package fabsdk

import (
	"strings"

	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	// cfg "github.com/jonluo94/cool/config"
	// "github.com/jonluo94/cool/log"
)

const (
	Admin = "Admin"
)

func (f *FabricClient) CreateChannel(channelTx string) error {
	mspClient, err := mspclient.New(f.sdk.Context(), mspclient.WithOrg(f.DefaultOrg))
	if err != nil {
		return err
	}
	adminIdentity, err := mspClient.GetSigningIdentity(Admin)
	if err != nil {
		return err
	}
	req := resmgmt.SaveChannelRequest{
		ChannelID:         f.ChannelID,
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

		mspClient, err := mspclient.New(f.sdk.Context(), mspclient.WithOrg(f.DefaultOrg))
		if err != nil {

			return err
		}
		adminIdentity, err := mspClient.GetSigningIdentity(Admin)
		if err != nil {

			return err
		}
		req := resmgmt.SaveChannelRequest{
			ChannelID:         f.ChannelID,
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
		err := c.JoinChannel(f.ChannelID, f.retry)
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
		f.ChannelID,
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
		f.ChannelID,
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
