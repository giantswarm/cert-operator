package main

import (
	"log"

	"github.com/giantswarm/certctl/cli"
)

func main() {
	if err := cli.CLICmd.Execute(); err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}
}
