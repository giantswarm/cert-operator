package main

import (
	"fmt"

	"github.com/giantswarm/microkit/command"
	microserver "github.com/giantswarm/microkit/server"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/viper"

	"github.com/giantswarm/cert-operator/flag"
	"github.com/giantswarm/cert-operator/server"
	"github.com/giantswarm/cert-operator/service"
)

var (
	description string     = "The cert-operator handles certificates for Kubernetes clusters running on Giantnetes."
	f           *flag.Flag = flag.New()
	gitCommit   string     = "n/a"
	name        string     = "cert-operator"
	source      string     = "https://github.com/giantswarm/cert-operator"
)

func main() {
	var err error

	// Create a new logger which is used by all packages.
	var newLogger micrologger.Logger
	{
		newLogger, err = micrologger.New(micrologger.Config{})
		if err != nil {
			panic(fmt.Sprintf("%#v\n", err))
		}
	}

	// We define a server factory to create the custom server once all command
	// line flags are parsed and all microservice configuration is storted out.
	newServerFactory := func(v *viper.Viper) microserver.Server {
		// Create a new custom service which implements business logic.
		var newService *service.Service
		{
			serviceConfig := service.DefaultConfig()

			serviceConfig.Flag = f
			serviceConfig.Logger = newLogger
			serviceConfig.Viper = v

			serviceConfig.Description = description
			serviceConfig.GitCommit = gitCommit
			serviceConfig.ProjectName = name
			serviceConfig.Source = source

			newService, err = service.New(serviceConfig)
			if err != nil {
				panic(fmt.Sprintf("%#v\n", err))
			}
			go newService.Boot()
		}

		// Create a new custom server which bundles our endpoints.
		var newServer microserver.Server
		{
			c := server.Config{
				Logger:  newLogger,
				Service: newService,
				Viper:   v,

				ProjectName: name,
			}

			newServer, err = server.New(c)
			if err != nil {
				panic(fmt.Sprintf("%#v\n", err))
			}
		}

		return newServer
	}

	// Create a new microkit command which manages our custom microservice.
	var newCommand command.Command
	{
		c := command.Config{
			Logger:        newLogger,
			ServerFactory: newServerFactory,

			Description:    description,
			GitCommit:      gitCommit,
			Name:           name,
			Source:         source,
			VersionBundles: service.NewVersionBundles(),
		}

		newCommand, err = command.New(c)
		if err != nil {
			panic(fmt.Sprintf("%#v\n", err))
		}
	}

	daemonCommand := newCommand.DaemonCommand().CobraCommand()

	daemonCommand.PersistentFlags().String(f.Service.Kubernetes.Address, "http://127.0.0.1:6443", "Address used to connect to Kubernetes. When empty in-cluster config is created.")
	daemonCommand.PersistentFlags().Bool(f.Service.Kubernetes.InCluster, false, "Whether to use the in-cluster config to authenticate with Kubernetes.")
	daemonCommand.PersistentFlags().String(f.Service.Kubernetes.TLS.CAFile, "", "Certificate authority file path to use to authenticate with Kubernetes.")
	daemonCommand.PersistentFlags().String(f.Service.Kubernetes.TLS.CrtFile, "", "Certificate file path to use to authenticate with Kubernetes.")
	daemonCommand.PersistentFlags().String(f.Service.Kubernetes.TLS.KeyFile, "", "Key file path to use to authenticate with Kubernetes.")

	daemonCommand.PersistentFlags().Duration(f.Service.Resource.VaultCrt.ExpirationThreshold, 0, "Amount of time to renew certificates before their expiration date.")
	daemonCommand.PersistentFlags().String(f.Service.Resource.VaultCrt.Namespace, "", "Namespace used to manage Kubernetes secrets in.")

	daemonCommand.PersistentFlags().String(f.Service.Vault.Config.Address, "", "Address used to connect to Vault.")
	daemonCommand.PersistentFlags().String(f.Service.Vault.Config.Token, "", "Token used to authenticate against Vault.")
	daemonCommand.PersistentFlags().String(f.Service.Vault.Config.PKI.CA.TTL, "", "TTL used to generate a new Cluster CA.")
	daemonCommand.PersistentFlags().String(f.Service.Vault.Config.PKI.CommonName.Format, "", "Common name used to generate a new Cluster CA.")

	newCommand.CobraCommand().Execute()
}
