package fabsdk

import (
	"encoding/hex"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	pfab "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

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
