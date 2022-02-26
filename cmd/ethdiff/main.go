package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mysteriumnetwork/ethdiff/diff"
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

	lastCommonBlock, err := diff.LastCommonBlock(ctx, *reqTimeout, flag.Arg(0), flag.Arg(1))
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
