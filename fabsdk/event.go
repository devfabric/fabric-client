package fabsdk

import (
	"fmt"

	pfab "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

type EventTxAdapter func(blockHeight uint64, txEv *Event) error //监听最新区块和区块内交易
type EventBlockAdapter func(blockInfo *Block) error             //监听最新区块

func (fab *FabricClient) RegisterBlockEvent(notifyBlockEv EventBlockAdapter, notifyTxEv EventTxAdapter) error {
	var (
		err error
	)

	fab.listenEvent.registration, fab.listenEvent.blockChain, err = fab.eventClient.RegisterBlockEvent()
	if err != nil {
		return err
	}

	go callback(fab.listenEvent.blockChain, fab.isExit, notifyBlockEv, notifyTxEv)
	return nil
}

func callback(blockChain <-chan *pfab.BlockEvent, isExit chan struct{}, notifyBlockEvent EventBlockAdapter, notifyTxEvent EventTxAdapter) {
	for {
		select {
		case cBlock := <-blockChain:
			if cBlock != nil {
				block, err := GetBlock(cBlock.Block)
				if err != nil {
					fmt.Println("GetBlock", err)
					continue
				}

				if notifyBlockEvent != nil {
					err = notifyBlockEvent(block)
					if err != nil {
						fmt.Println("notifyBlockEvent", err)
					}
				}
				for _, tx := range block.Data {
					if tx.ValidationCode != 0 {
						continue
					}
					txActions, ok := tx.Payload.(TransactionActions)
					if !ok {
						// log.Warn("is config tx skip")
						continue
					}
					for _, txAc := range txActions {
						if notifyTxEvent != nil {
							err = notifyTxEvent(block.Height, &txAc.Events)
							if err != nil {
								fmt.Println("notifyTxEvent", err)
							}
						}
					}
				}

				err = UpdateBlockHeight("./", block.Height)
				if err != nil {
					fmt.Println("UpdateBlockHeight", err)
				}
			}
		case <-isExit:
			return
		}
	}
}
