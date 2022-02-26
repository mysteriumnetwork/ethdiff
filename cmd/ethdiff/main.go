package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

var version = "undefined"

var (
	reqTimeout      = flag.Duration("req-timeout", 5*time.Second, "timeout for single request")
	totalTimeout    = flag.Duration("total-timeout", 20*time.Second, "whole operation timeout")
	showVersion     = flag.Bool("version", false, "show program version and exit")
)

func run() int {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage:\n\n%s [options...] <left RPC address> <right RPC address>\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		return 0
	}

	if flag.NArg() != 2 {
		flag.Usage()
		return 2
	}

	ctx, cl := context.WithTimeout(context.Background(), *totalTimeout)
	defer cl()

	lastCommonBlock, err := LastCommonBlock(ctx, *reqTimeout, flag.Arg(0), flag.Arg(1))
	if err != nil {
		log.Fatalf("LastCommonBlock(%v, %v, %v) error: %v", *reqTimeout, flag.Arg(0), flag.Arg(1), err)
	}

	fmt.Printf("0x%x\n", lastCommonBlock)

	return 0
}

func main() {
	log.Default().SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	os.Exit(run())
}

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
