package main

import (
	"github.com/koind/telnet/internal"
	flag "github.com/spf13/pflag"
	"log"
)

var (
	address     string
	timeout     int64
	readTimeout int64
)

func init() {
	flag.StringVarP(&address, "address", "a", "", "resource address for the connection")
	flag.Int64VarP(&timeout, "timeout", "t", 0, "timeout to connect")
	flag.Int64VarP(&readTimeout, "read_timeout", "r", 0, "read timeout to connect")
}

func main() {
	flag.Parse()

	if address == "" {
		log.Fatal("Specify the address to connect")
	}

	if timeout == 0 {
		log.Fatal("Specify the timeout to connect")
	}

	if readTimeout == 0 {
		log.Fatal("Specify the read_timeout to connect")
	}

	options := internal.Options{
		Address:     address,
		Timeout:     timeout,
		ReadTimeout: readTimeout,
	}

	internal.Run(options)
}
