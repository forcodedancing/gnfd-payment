package monitor

import (
	"encoding/csv"
	"fmt"
	paymentTypes "github.com/bnb-chain/greenfield/x/payment/types"
	abciTypes "github.com/cometbft/cometbft/abci/types"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

type GnfdBlockProcessor struct {
	client *GnfdCompositeClients
}

func NewGnfdBlockProcessor(client *GnfdCompositeClients) *GnfdBlockProcessor {
	return &GnfdBlockProcessor{
		client: client,
	}
}

func (p *GnfdBlockProcessor) Name() string {
	return "gnfd"
}

func (p *GnfdBlockProcessor) GetDatabaseBlockHeight() (uint64, error) {
	return 0, nil
}

func (p *GnfdBlockProcessor) GetBlockchainBlockHeight() (uint64, error) {
	return p.client.GetLatestBlockHeight()
}

func (p *GnfdBlockProcessor) Process(blockHeight uint64, writer *csv.Writer) error {
	//fmt.Println("handling block", blockHeight)
	results, err := p.client.GetBlockResults(int64(blockHeight))
	if err != nil {
		fmt.Printf("processor: %s, fail to get block results err: %s", p.Name(), err)
		return err
	}
	for _, result := range results.TxsResults {
		for _, event := range result.Events {
			switch event.Type {
			case "greenfield.payment.EventStreamRecordUpdate":
				p.handleUpdateStreamRecord(blockHeight, 1, event, writer)
			}
		}
	}
	for _, event := range results.EndBlockEvents {
		switch event.Type {
		case "greenfield.payment.EventStreamRecordUpdate":
			p.handleUpdateStreamRecord(blockHeight, 0, event, writer)
		}
	}
	return nil
}

func (p *GnfdBlockProcessor) handleUpdateStreamRecord(blockHeight uint64, isTxEvent int, event abciTypes.Event, writer *csv.Writer) {

	e, err := sdkTypes.ParseTypedEvent(event)
	if err != nil {
		fmt.Printf("processor: %s, fail to parse EventCreateGroup err: %s", p.Name(), err)
		panic(err)
	}
	streamRecordUpdate := e.(*paymentTypes.EventStreamRecordUpdate)

	if strings.ToLower(streamRecordUpdate.Account) == "0x417b1c42e90cf0933900b4404e5d83f3ae7b2e4e" {
		//fmt.Println("writing csv")
		row := []string{
			fmt.Sprintf("%d", blockHeight),
			fmt.Sprintf("%d", isTxEvent),
			streamRecordUpdate.Account,
			streamRecordUpdate.NetflowRate.String(),
			streamRecordUpdate.BufferBalance.String(),
			streamRecordUpdate.StaticBalance.String(),
			fmt.Sprintf("%d", streamRecordUpdate.SettleTimestamp),
			fmt.Sprintf("%d", streamRecordUpdate.CrudTimestamp),
		}
		err = writer.Write(row)
		if err != nil {
			fmt.Printf("processor: %s, fail to write row err: %s", p.Name(), err)
			panic(err)
		}
		writer.Flush()
	}
}
