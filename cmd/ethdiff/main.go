package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

var version = "undefined"

var (
	reqTimeout      = flag.Duration("req-timeout", 5*time.Second, "timeout for single request")
	totalTimeout    = flag.Duration("total-timeout", 20*time.Second, "whole operation timeout")
	showVersion     = flag.Bool("version", false, "show program version and exit")
	confirmations   = flag.Int("confirmations", 10, "required number of confirmations")
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

	return 0
}

func main() {
	log.Default().SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	os.Exit(run())
}
