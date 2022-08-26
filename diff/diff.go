package diff

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
)

var (
	ErrNoCommonRoot = errors.New("no common root block")
)

type Client interface {
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	BlockNumber(ctx context.Context) (uint64, error)
}

func getNumbers(ctx context.Context, left, right Client) (uint64, uint64, error) {
	var (
		wg                      sync.WaitGroup
		leftNumber, rightNumber uint64
		leftErr, rightErr       error
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		leftNumber, leftErr = left.BlockNumber(ctx)
	}()
	go func() {
		defer wg.Done()
		rightNumber, rightErr = right.BlockNumber(ctx)
	}()
	wg.Wait()

	if leftErr != nil {
		return 0, 0, fmt.Errorf("left.BlockNumber: error: %w", leftErr)
	}
	if rightErr != nil {
		return 0, 0, fmt.Errorf("right.BlockNumber: error: %w", rightErr)
	}

	return leftNumber, rightNumber, nil
}

func getBlocks(ctx context.Context, left, right Client, blockNumber uint64) (*types.Header, *types.Header, error) {
	var (
		wg                    sync.WaitGroup
		leftBlock, rightBlock *types.Header
		leftErr, rightErr     error
	)
	bigBlockNumber := big.NewInt(int64(blockNumber))

	wg.Add(2)
	go func() {
		defer wg.Done()
		leftBlock, leftErr = left.HeaderByNumber(ctx, bigBlockNumber)
	}()
	go func() {
		defer wg.Done()
		rightBlock, rightErr = right.HeaderByNumber(ctx, bigBlockNumber)
	}()
	wg.Wait()

	if leftErr != nil {
		return nil, nil, fmt.Errorf("left.HeaderByNumber: error: %w", leftErr)
	}
	if rightErr != nil {
		return nil, nil, fmt.Errorf("right.HeaderByNumber: error: %w", rightErr)
	}

	return leftBlock, rightBlock, nil
}

func LastCommonBlock(ctx context.Context, left, right Client, offset uint64) (uint64, uint64, error) {
	leftLatestBlock, rightLatestBlock, err := getNumbers(ctx, left, right)
	if err != nil {
		return 0, 0, err
	}
	highestCommonBlock := min(leftLatestBlock, rightLatestBlock)
	log.Printf("highestCommonBlock = 0x%x (%d)", highestCommonBlock, highestCommonBlock)
	highestCommonBlock -= offset
	log.Printf("highestCommonBlock (safe value) = 0x%x (%d)", highestCommonBlock, highestCommonBlock)

	// fast path
	leftBlock, rightBlock, err := getBlocks(ctx, left, right, highestCommonBlock)
	if err != nil {
		return 0, 0, err
	}
	if leftBlock.Hash() == rightBlock.Hash() {
		return highestCommonBlock, highestCommonBlock, nil
	}

	res, err := search(highestCommonBlock+1, func(blockNumber uint64) (bool, error) {

		leftBlock, rightBlock, err := getBlocks(ctx, left, right, blockNumber)
		if err != nil {
			return false, err
		}

		result := leftBlock.Hash() != rightBlock.Hash()
		if result && blockNumber == 0 {
			return false, ErrNoCommonRoot
		}
		log.Printf("searchFunc(0x%x) = %v", blockNumber, result)
		return result, nil
	})
	return res - 1, highestCommonBlock, err
}

func min(x, y uint64) uint64 {
	if x < y {
		return x
	}
	return y
}

func search(n uint64, f func(uint64) (bool, error)) (uint64, error) {
	// Define f(-1) == false and f(n) == true.
	// Invariant: f(i-1) == false, f(j) == true.
	i, j := uint64(0), n
	for i < j {
		h := (i + j) >> 1
		// i ≤ h < j
		r, err := f(h)
		if err != nil {
			return 0, fmt.Errorf("search function error: %w", err)
		}
		if !r {
			i = h + 1 // preserves f(i-1) == false
		} else {
			j = h // preserves f(j) == true
		}
	}
	// i == j, f(i-1) == false, and f(j) (= f(i)) == true  =>  answer is i.
	return i, nil
}
