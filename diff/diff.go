package diff

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

type clientConnectResponse struct {
	Client *ethclient.Client
	Err error
}

func asyncClientConnect(ctx context.Context, address string) <-chan *clientConnectResponse {
	r := make(chan *clientConnectResponse, 1)

	go func() {
		defer close(r)
		client, err := ethclient.DialContext(ctx, address)
		r <- &clientConnectResponse{
			Client: client,
			Err: err,
		}
	}()

	return r
}

type blockNumberResponse struct {
	BlockNumber uint64
	Err error
}

func asyncBlockNumber(ctx context.Context, client *ethclient.Client) <- chan *blockNumberResponse {
	r := make(chan *blockNumberResponse, 1)

	go func() {
		defer close(r)
		number, err := client.BlockNumber(ctx)
		r <- &blockNumberResponse{
			BlockNumber: number,
			Err: err,
		}
	}()

	return r
}

func LastCommonBlock(ctx context.Context, reqTimeout time.Duration, left, right string) (uint64, error) {
	ctx1, cl := context.WithTimeout(ctx, reqTimeout)
	defer cl()

	leftClientFuture, rightClientFuture := asyncClientConnect(ctx1, left), asyncClientConnect(ctx1, right)
	leftClientResult, rightClientResult := <-leftClientFuture, <-rightClientFuture
	if leftClientResult.Err != nil {
		return 0, fmt.Errorf("asyncClientConnect(%q) error: %w", left, leftClientResult.Err)
	}
	if rightClientResult.Err != nil {
		return 0, fmt.Errorf("asyncClientConnect(%q) error: %w", right, rightClientResult.Err)
	}

	leftClient, rightClient := leftClientResult.Client, rightClientResult.Client

	ctx1, cl = context.WithTimeout(ctx, reqTimeout)
	defer cl()

	leftBlockNumberFuture, rightBlockNumberFuture := asyncBlockNumber(ctx1, leftClient), asyncBlockNumber(ctx1, rightClient)
	leftLatestBlock, rightLatestBlock := <-leftBlockNumberFuture, <-rightBlockNumberFuture
	if leftLatestBlock.Err != nil {
		return 0, fmt.Errorf("asyncBlockNumber(%q) error: %w", left, leftLatestBlock.Err)
	}
	if rightLatestBlock.Err != nil {
		return 0, fmt.Errorf("asyncBlockNumber(%q) error: %w", right, rightLatestBlock.Err)
	}

	highestCommonBlock := max(leftLatestBlock.BlockNumber, rightLatestBlock.BlockNumber)

	return highestCommonBlock, nil
}

func max(x, y uint64) uint64 {
	if x > y {
		return x
	}
	return y
}
