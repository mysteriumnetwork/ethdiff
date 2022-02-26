package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mysteriumnetwork/ethdiff/diff"
)

var version = "undefined"

var (
	totalTimeout = flag.Duration("total-timeout", 1*time.Minute, "whole operation timeout")
	showVersion  = flag.Bool("version", false, "show program version and exit")
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

	left, right := flag.Arg(0), flag.Arg(1)

	leftClientFuture, rightClientFuture := asyncClientConnect(ctx, left), asyncClientConnect(ctx, right)
	leftClientResult, rightClientResult := <-leftClientFuture, <-rightClientFuture
	if leftClientResult.Err != nil {
		log.Fatalf("asyncClientConnect(%q) error: %v", left, leftClientResult.Err)
	}
	if rightClientResult.Err != nil {
		log.Fatalf("asyncClientConnect(%q) error: %v", right, rightClientResult.Err)
	}

	leftClient, rightClient := leftClientResult.Client, rightClientResult.Client

	lastCommonBlock, err := diff.LastCommonBlock(ctx, leftClient, rightClient)
	if err != nil {
		log.Fatalf("LastCommonBlock(%v, %v) error: %v", flag.Arg(0), flag.Arg(1), err)
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
	Err    error
}

func asyncClientConnect(ctx context.Context, address string) <-chan *clientConnectResponse {
	r := make(chan *clientConnectResponse, 1)

	go func() {
		defer close(r)
		client, err := ethclient.DialContext(ctx, address)
		r <- &clientConnectResponse{
			Client: client,
			Err:    err,
		}
	}()

	return r
}
