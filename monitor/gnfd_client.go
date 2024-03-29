package monitor

import (
	"context"
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	bfttypes "github.com/cometbft/cometbft/types"
	"github.com/fcd/gnfd-payment/util"
	"sync"

	sdkclient "github.com/bnb-chain/greenfield-go-sdk/client"
)

type GnfdClient struct {
	sdkclient.IClient
	Height int64
}

type GnfdCompositeClients struct {
	clients []*GnfdClient
}

func NewGnfdCompositClients(rpcAddrs []string, chainId string, useWebsocket bool) *GnfdCompositeClients {
	clients := make([]*GnfdClient, 0)
	for i := 0; i < len(rpcAddrs); i++ {
		sdkClient, err := sdkclient.New(chainId, rpcAddrs[i], sdkclient.Option{DefaultAccount: nil, UseWebSocketConn: useWebsocket})
		if err != nil {
			util.Logger.Errorf("rpc node %s is not available", rpcAddrs[i])
			continue
		}
		clients = append(clients, &GnfdClient{
			IClient: sdkClient,
		})
		if len(clients) == 0 {
			panic("no Greenfield client available")
		}
	}
	return &GnfdCompositeClients{
		clients: clients,
	}
}

func getClientBlockHeight(clientChan chan *GnfdClient, wg *sync.WaitGroup, client *GnfdClient) {
	defer wg.Done()
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), RPCTimeout)
	defer cancel()
	status, err := client.GetStatus(ctxWithTimeout)
	if err != nil {
		return
	}
	latestHeight := status.SyncInfo.LatestBlockHeight
	client.Height = latestHeight
	clientChan <- client
}

func (c *GnfdCompositeClients) GetBlock(height int64) (*bfttypes.Block, error) {
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
	defer cancel()
	block, err := c.pickClient().GetBlockByHeight(ctx, height)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (c *GnfdCompositeClients) GetBlockResults(height int64) (*ctypes.ResultBlockResults, error) {
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
	defer cancel()
	blockResults, err := c.pickClient().GetBlockResultByHeight(ctx, height)
	if err != nil {
		return nil, err
	}
	return blockResults, nil
}

func (c *GnfdCompositeClients) GetLatestBlockHeight() (uint64, error) {
	return uint64(c.pickClient().Height), nil
}

func (c *GnfdCompositeClients) pickClient() *GnfdClient {
	wg := new(sync.WaitGroup)
	wg.Add(len(c.clients))
	clientCh := make(chan *GnfdClient)
	waitCh := make(chan struct{})
	go func() {
		for _, c := range c.clients {
			go getClientBlockHeight(clientCh, wg, c)
		}
		wg.Wait()
		close(waitCh)
	}()
	var maxHeight int64
	maxHeightClient := c.clients[0]
	for {
		select {
		case c := <-clientCh:
			if c.Height > maxHeight {
				maxHeight = c.Height
				maxHeightClient = c
			}
		case <-waitCh:
			return maxHeightClient
		}
	}
}
