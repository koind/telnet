package main

import (
	"github.com/koind/telnet/internal"
	flag "github.com/spf13/pflag"
	"log"
)

var (
	address string
	timeout int64
)

func init() {
	flag.StringVarP(&address, "address", "a", "", "resource address for the connection")
	flag.Int64VarP(&timeout, "timeout", "t", 0, "timeout to connect")
}

func main() {
	flag.Parse()

	if address == "" {
		log.Fatal("Specify the address to connect")
	}

	if timeout == 0 {
		log.Fatal("Specify the timeout to connect")
	}

	options := internal.Options{
		Address: address,
		Timeout: timeout,
	}

	internal.Run(options)
}
