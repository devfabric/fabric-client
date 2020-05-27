package fabsdk

import (
	"github.com/cloudflare/cfssl/log"
	pfab "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

type EventAdapter func(cbName string, payload []byte) error

func (fab *FabricClient) RegisterBlockEvent(doProcess EventAdapter) error {
	var (
		err        error
		blockChain <-chan *pfab.BlockEvent
	)

	fab.registration, blockChain, err = fab.eventClient.RegisterBlockEvent()
	if err != nil {
		return err
	}

	go callback(blockChain, doProcess)
	return nil
}

func callback(blockChain <-chan *pfab.BlockEvent, doProcess EventAdapter) {
	for {
		select {
		case cBlock := <-blockChain:
			if cBlock != nil {
				// enBlock := blockParse(cBlock.Block)
				block, err := GetBlock(cBlock.Block)
				if err != nil {
					log.Error(err)
					return
				}

				for _, tx := range block.Data {
					if tx.ValidationCode != 0 {

						continue
					}
					acs, ok := tx.Payload.(TransactionActions)
					if !ok {
						// log.Warn("is config tx skip")
						continue
					}
					_ = acs
					// for _, ac := range acs {
					// 	_, ok := types.EventMap[ac.Events.EventName]
					// 	if !ok {
					// 		continue
					// 	}

					// 	cce := &fab.CCEvent{
					// 		EventName:   ac.Events.EventName,
					// 		Payload:     ac.Events.Payload,
					// 		TxID:        ac.Events.TxId,
					// 		BlockNumber: block.Height,
					// 	}

					// 	E.EventCh <- cce
					// }

				}
			}

		}
	}
}
