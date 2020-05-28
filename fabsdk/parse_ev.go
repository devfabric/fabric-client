package fabsdk

import (
	"github.com/gogo/protobuf/proto"

	// "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
	cm "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	pb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"

	"github.com/hyperledger/fabric/core/ledger/util"
	"github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/protos/utils"
	"github.com/pkg/errors"
)

type TxEvent struct {
	Number      uint64
	Status      string `json:"status"`
	Txid        string `json:"txid"`
	Timestamp   int64  `json:"timestamp"`
	ChaincodeID string
	EventName   string
	Payload     []byte
}

func fastParseBlock(block *cm.Block) ([]*TxEvent, error) {
	var (
		tranNo   int64 = -1
		err      error
		TxEvents = make([]*TxEvent, 0)
	)

	txsFilter := util.TxValidationFlags(block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])
	if len(txsFilter) == 0 {
		txsFilter = util.NewTxValidationFlags(len(block.Data.Data))
		block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER] = txsFilter
	}

	for _, txBytes := range block.Data.Data {
		txInfo := &TxEvent{
			Number: block.Header.Number,
		}

		tranNo++
		if txsFilter.IsInvalid(int(tranNo)) {
			txInfo.Status = "INVALID"
			continue
		} else {
			txInfo.Status = "VALID"
		}

		var env *common.Envelope
		if env, err = utils.GetEnvelopeFromBlock(txBytes); err != nil {
			return nil, err
		}

		var payload *common.Payload
		if payload, err = utils.GetPayload(env); err != nil {
			return nil, err
		}

		var chdr *common.ChannelHeader
		chdr, err = utils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
		if err != nil {
			return nil, err
		}
		txInfo.Txid = chdr.TxId
		txInfo.Timestamp = chdr.Timestamp.GetSeconds()

		if common.HeaderType(chdr.Type) != common.HeaderType_ENDORSER_TRANSACTION {
			continue
		}

		var tx *peer.Transaction
		if tx, err = utils.GetTransaction(payload.Data); err != nil {
			return nil, err
		}

		for _, action := range tx.Actions {
			cap, err := GetChaincodeActionPayload(action.Payload)
			if err != nil {
				return nil, errors.New(err.Error())
			}

			chaincodeAction, err := GetChaincodeAction(cap.Action.ProposalResponsePayload)
			if err != nil {
				return nil, errors.Wrap(err, "error unmarshaling GetChaincodeAction")
			}

			if len(chaincodeAction.Events) > 0 {
				chaincodeEvent := &pb.ChaincodeEvent{}
				err := proto.Unmarshal(chaincodeAction.Events, chaincodeEvent)
				if err != nil {
					return nil, errors.Wrap(err, "error unmarshaling ChaincodeEvent")
				}
				txInfo.ChaincodeID = chaincodeEvent.ChaincodeId
				txInfo.EventName = chaincodeEvent.EventName
				txInfo.Payload = chaincodeEvent.Payload
			}
		}
		TxEvents = append(TxEvents, txInfo)
	}
	return TxEvents, nil
}
