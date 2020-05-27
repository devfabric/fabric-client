package fabsdk

import (
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
		case event := <-blockChain:
			_ = event
			// processChan <- event
			//UpdateBlockHeight
		}
	}
}
