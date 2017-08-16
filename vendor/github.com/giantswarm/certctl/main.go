package main

import (
	"log"

	"github.com/giantswarm/certctl/cli"
	"github.com/giantswarm/microerror"
)

func main() {
	if err := cli.CLICmd.Execute(); err != nil {
		log.Fatalf("%#v\n", microerror.Mask(err))
	}
}
