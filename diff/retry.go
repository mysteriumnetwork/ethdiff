package diff

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

type RetryingClient struct {
	client  Client
	retries uint
}

func NewRetryingClient(client Client, retries uint) *RetryingClient {
	if retries < 1 {
		retries = 1
	}

	return &RetryingClient{
		client:  client,
		retries: retries,
	}
}

func (c *RetryingClient) BlockByNumber(ctx context.Context, number *big.Int) (block *types.Block, err error) {
	for i := uint(0); i < c.retries; i++ {
		block, err = c.client.BlockByNumber(ctx, number)
		if err == nil {
			return
		}
		log.Printf("attempt #%d failed, error = %v", i+1, err)
	}
	return
}

func (c *RetryingClient) BlockNumber(ctx context.Context) (number uint64, err error) {
	for i := uint(0); i < c.retries; i++ {
		number, err = c.client.BlockNumber(ctx)
		if err == nil {
			return number, err
		}
		log.Printf("attempt #%d failed, error = %v", i+1, err)
	}
	return
}
