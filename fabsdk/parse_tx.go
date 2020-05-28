package fabsdk

import (
	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	ledgerUtil "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/core/ledger/util"
	cb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	fabriccmn "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"

	// "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
	pb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/rwsetutil"

	"github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
)

const (
	HeaderType_MESSAGE = iota
	HeaderType_CONFIG
	HeaderType_CONFIG_UPDATE
	HeaderType_ENDORSER_TRANSACTION
	HeaderType_ORDERER_TRANSACTION
	HeaderType_DELIVER_SEEK_INFO
	HeaderType_CHAINCODE_PACKAGE
	HeaderType_PEER_ADMIN_OPERATION
	HeaderType_TOKEN_TRANSACTION
)

type TransactionActions []Action

type Action struct {
	ChaincodeID   string
	ChainCodeArgs [][]byte

	Endorsements [][]byte
	Response     ChaincodeAction
	TxRwSet      string
	Events       Event
	NsRwSets     []NsRwSet
}

type Event struct {
	ChaincodeId string
	TxId        string
	EventName   string
	Payload     []byte
}
type ChaincodeAction struct {
	Status  int32
	Message string
	Payload []byte
}

type NsRwSet struct {
	NameSpace string
	Reads     []KVRead
	Writes    []KVWrite
}
type KVWrite struct {
	Key      string
	IsDelete bool
	Value    []byte
}
type KVRead struct {
	Key string
	Version
}
type Version struct {
	BlockNum uint64
	TxNum    uint64
}

type Transaction struct {
	TxID            string
	TransactionType int32
	ValidationCode  int32 //交易确认码
	ChannelID       string
	Timestamp       *timestamp.Timestamp

	Payload interface{}
}

type Block struct {
	Height       uint64 `json:"height,omitempty"`
	PreviousHash []byte `json:"current_block_hash,omitempty"`
	DataHash     []byte `json:"previous_block_hash,omitempty"`

	Data []Transaction
}

func GetBlock(fabrBlock *fabriccmn.Block) (*Block, error) {
	ret := Block{
		Height:       fabrBlock.Header.Number,
		PreviousHash: fabrBlock.Header.PreviousHash,
		DataHash:     fabrBlock.Header.DataHash,
	}

	validationCode := []int32{}
	txValidationFlags := ledgerUtil.TxValidationFlags(fabrBlock.Metadata.Metadata[fabriccmn.BlockMetadataIndex_TRANSACTIONS_FILTER])
	for i := 0; i < len(txValidationFlags); i++ {
		validationCode = append(validationCode, int32(txValidationFlags[i]))
	}
	validationCodeLen := len(validationCode)

	datas := []Transaction{}
	for i, v := range fabrBlock.Data.Data {
		data, err := getBlockData(v)
		if err != nil {
			return nil, err
		}
		if i < validationCodeLen {
			data.ValidationCode = validationCode[i]
		}
		datas = append(datas, *data)
	}

	ret.Data = datas

	return &ret, nil
}

func getBlockData(dataBytes []byte) (*Transaction, error) {
	env := &cb.Envelope{}
	err := proto.Unmarshal(dataBytes, env)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshaling Envelope")
	}

	payload := &cb.Payload{}
	err = proto.Unmarshal(env.Payload, payload)
	if err != nil {
		return nil, errors.Wrap(err, "no payload in envelope")
	}

	chdr, err := UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	envelope, err := getPayload(int(chdr.Type), payload.Data)
	if err != nil {
		return nil, err
	}

	ret := Transaction{
		// ValidationCode:  validationCode,
		TxID:            chdr.TxId,
		TransactionType: chdr.Type,
		ChannelID:       chdr.ChannelId,
		Timestamp:       chdr.Timestamp,
		Payload:         envelope,
	}

	return &ret, nil
}

func UnmarshalChannelHeader(bytes []byte) (*cb.ChannelHeader, error) {
	chdr := &cb.ChannelHeader{}
	err := proto.Unmarshal(bytes, chdr)
	return chdr, errors.Wrap(err, "error unmarshaling ChannelHeader")
}

func getPayload(headerType int, data []byte) (interface{}, error) {
	if headerType == HeaderType_CONFIG {
		return nil, nil
	} else if headerType == HeaderType_CONFIG_UPDATE {
		return nil, nil
	} else if headerType == HeaderType_ENDORSER_TRANSACTION {
		return getTransaction(data)
	}
	return nil, errors.New("headerType err")
}

func getTransaction(txBytes []byte) (TransactionActions, error) {
	ret := TransactionActions{}

	tx := &peer.Transaction{}
	err := proto.Unmarshal(txBytes, tx)
	if err != nil {
		return nil, errors.Wrap(err, "Bad envelope:error unmarshaling Transaction")
	}

	for _, action := range tx.Actions {
		cap, err := GetChaincodeActionPayload(action.Payload)
		if err != nil {
			return nil, errors.New(err.Error())
		}

		cis, err := GetChaincodeInvocationSpec(cap.ChaincodeProposalPayload)
		if err != nil {
			return nil, errors.Wrap(err, "error unmarshaling ChaincodeInvocationSpec")
		}

		endorsements := [][]byte{}
		for _, endorsement := range cap.Action.Endorsements {
			endorsements = append(endorsements, endorsement.Endorser)
		}

		chaincodeAction, err := GetChaincodeAction(cap.Action.ProposalResponsePayload)
		if err != nil {
			return nil, errors.Wrap(err, "error unmarshaling GetChaincodeAction")
		}
		ccac := ChaincodeAction{
			Status:  chaincodeAction.Response.Status,
			Message: chaincodeAction.Response.Message,
			Payload: chaincodeAction.Response.Payload,
		}

		events, err := GetEvents(chaincodeAction.Events)
		if err != nil {
			return nil, errors.Wrap(err, "error unmarshaling ChaincodeEvent")
		}

		mytxRWSet, err := GetRWSet(chaincodeAction.Results)
		if err != nil {
			return nil, errors.Wrap(err, "error unmarshaling ChaincodeEvent")
		}

		ac := Action{
			ChaincodeID:   cis.ChaincodeSpec.ChaincodeId.Name,
			ChainCodeArgs: cis.ChaincodeSpec.Input.Args,
			// Endorsements:  endorsements,
			Response: ccac,
			Events:   *events,
			NsRwSets: mytxRWSet,
		}
		ret = append(ret, ac)
	}

	return ret, nil
}

// GetChaincodeActionPayload Get ChaincodeActionPayload from bytes
func GetChaincodeActionPayload(capBytes []byte) (*peer.ChaincodeActionPayload, error) {
	cap := &peer.ChaincodeActionPayload{}
	err := proto.Unmarshal(capBytes, cap)
	return cap, errors.Wrap(err, "error unmarshaling ChaincodeActionPayload")
}

// GetChaincodeInvocationSpec Get GetChaincodeInvocationSpec from bytes
func GetChaincodeInvocationSpec(chaincodeProposalPayload []byte) (*pb.ChaincodeInvocationSpec, error) {
	cpp := &pb.ChaincodeProposalPayload{}
	err := proto.Unmarshal(chaincodeProposalPayload, cpp)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshaling ChaincodeProposalPayload")
	}

	cis := &pb.ChaincodeInvocationSpec{}
	err = proto.Unmarshal(cpp.Input, cis)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshaling ChaincodeInvocationSpec")
	}

	return cis, nil
}

// GetChaincodeAction Get GetChaincodeAction from bytes
func GetChaincodeAction(ProposalResponsePayload []byte) (*pb.ChaincodeAction, error) {
	prp := &pb.ProposalResponsePayload{}
	err := proto.Unmarshal(ProposalResponsePayload, prp)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshaling ProposalResponsePayload")
	}
	chaincodeAction := &pb.ChaincodeAction{}
	err = proto.Unmarshal(prp.Extension, chaincodeAction)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshaling ChaincodeAction")
	}
	return chaincodeAction, nil
}

// GetEvents Get GetEvents from bytes
func GetEvents(eventByte []byte) (*Event, error) {
	events := Event{}
	if len(eventByte) > 0 {
		chaincodeEvent := &pb.ChaincodeEvent{}
		err := proto.Unmarshal(eventByte, chaincodeEvent)
		if err != nil {
			return nil, errors.Wrap(err, "error unmarshaling ChaincodeEvent")
		}
		events.EventName = chaincodeEvent.EventName
		events.Payload = chaincodeEvent.Payload
		events.TxId = chaincodeEvent.TxId
		events.ChaincodeId = chaincodeEvent.ChaincodeId
	}
	return &events, nil
}

// GetRWSet Get GetRWSet from bytes
func GetRWSet(setByte []byte) ([]NsRwSet, error) {
	mytxRWSet := []NsRwSet{}
	if len(setByte) > 0 {
		txRWSet := &rwsetutil.TxRwSet{}
		if err := txRWSet.FromProtoBytes(setByte); err != nil {
			return nil, errors.Wrap(err, "get txrwset error")
		}
		for _, nsRWSet := range txRWSet.NsRwSets {
			if nsRWSet.KvRwSet == nil {
				continue
			}
			reads := []KVRead{}
			for _, r := range nsRWSet.KvRwSet.Reads {
				kvRead := KVRead{
					Key:     r.Key,
					Version: Version{},
				}
				if r.Version != nil {
					kvRead.Version.BlockNum = r.Version.BlockNum
					kvRead.Version.TxNum = r.Version.TxNum
				}
				reads = append(reads, kvRead)
			}

			writes := []KVWrite{}
			for _, w := range nsRWSet.KvRwSet.Writes {
				kvwrite := KVWrite{
					Key:      w.Key,
					IsDelete: w.IsDelete,
					Value:    w.Value,
				}
				writes = append(writes, kvwrite)
			}
			nsRwSet := NsRwSet{
				NameSpace: nsRWSet.NameSpace,
				Reads:     reads,
				Writes:    writes,
			}
			mytxRWSet = append(mytxRWSet, nsRwSet)
		}
	}

	return mytxRWSet, nil
}
