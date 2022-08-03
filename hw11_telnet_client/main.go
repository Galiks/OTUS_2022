package main

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/spf13/pflag"
)

var (
	timeout    time.Duration
	host, port string
)

func main() {
	pflag.DurationVarP(&timeout, "timeout", "t", time.Second*10, "connection timeout")
	pflag.Parse()
	if len(pflag.Args()) < 2 {
		log.Fatal("invalid count of parameters")
	}
	host = pflag.Arg(0)
	port = pflag.Arg(1)
	tc := NewTelnetClient(net.JoinHostPort(host, port), timeout, os.Stdin, os.Stdout)
	if err := Run(tc); err != nil {
		log.Fatal(err)
	}
	// Place your code here,
	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?
}
