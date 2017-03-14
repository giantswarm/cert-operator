// Package daemon implements the daemon command for any microservice.
package daemon

import (
	"os"
	"os/signal"
	"sync"

	"github.com/spf13/cobra"

	"github.com/giantswarm/microkit/command/daemon/flag"
	microerror "github.com/giantswarm/microkit/error"
	microflag "github.com/giantswarm/microkit/flag"
	"github.com/giantswarm/microkit/logger"
	"github.com/giantswarm/microkit/server"
)

var (
	Flag = flag.New()
)

// Config represents the configuration used to create a new daemon command.
type Config struct {
	// Dependencies.
	Logger        logger.Logger
	ServerFactory func() server.Server
}

// DefaultConfig provides a default configuration to create a new daemon command
// by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger:        nil,
		ServerFactory: nil,
	}
}

// New creates a new daemon command.
func New(config Config) (Command, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "logger must not be empty")
	}
	if config.ServerFactory == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "server factory must not be empty")
	}

	newCommand := &command{
		// Internals.
		cobraCommand:  nil,
		logger:        config.Logger,
		serverFactory: config.ServerFactory,
	}

	newCommand.cobraCommand = &cobra.Command{
		Use:   "daemon",
		Short: "Execute the daemon of the microservice.",
		Long:  "Execute the daemon of the microservice.",
		Run:   newCommand.Execute,
	}

	newCommand.cobraCommand.PersistentFlags().StringSlice(Flag.Config.Dirs, []string{"."}, "List of config file directories.")
	newCommand.cobraCommand.PersistentFlags().StringSlice(Flag.Config.Files, []string{"config"}, "List of the config file names. All viper supported extensions can be used.")

	newCommand.cobraCommand.PersistentFlags().String(Flag.Server.Listen.Address, "http://127.0.0.1:8000", "Address used to make the server listen to.")
	newCommand.cobraCommand.PersistentFlags().String(Flag.Server.TLS.CaFile, "", "File path of the TLS root CA file, if any.")
	newCommand.cobraCommand.PersistentFlags().String(Flag.Server.TLS.CrtFile, "", "File path of the TLS public key file, if any.")
	newCommand.cobraCommand.PersistentFlags().String(Flag.Server.TLS.KeyFile, "", "File path of the TLS private key file, if any.")

	return newCommand, nil
}

type command struct {
	// Internals.
	cobraCommand  *cobra.Command
	logger        logger.Logger
	serverFactory func() server.Server
}

func (c *command) CobraCommand() *cobra.Command {
	return c.cobraCommand
}

func (c *command) Execute(cmd *cobra.Command, args []string) {
	v := c.serverFactory().Config().Viper

	// Merge the given command line flags with the given environment variables and
	// the given config files, if any. The merged flags will be applied to the
	// given viper.
	err := microflag.Merge(v, cmd.Flags(), v.GetStringSlice(Flag.Config.Dirs), v.GetStringSlice(Flag.Config.Files))
	if err != nil {
		panic(err)
	}

	var newServer server.Server
	{
		serverConfig := c.serverFactory().Config()
		serverConfig.ListenAddress = v.GetString(Flag.Server.Listen.Address)
		serverConfig.TLSCAFile = v.GetString(Flag.Server.TLS.CaFile)
		serverConfig.TLSCrtFile = v.GetString(Flag.Server.TLS.CrtFile)
		serverConfig.TLSKeyFile = v.GetString(Flag.Server.TLS.KeyFile)
		newServer, err = server.New(serverConfig)
		if err != nil {
			panic(err)
		}
		go newServer.Boot()
	}

	// Listen to OS signals.
	listener := make(chan os.Signal, 2)
	signal.Notify(listener, os.Interrupt, os.Kill)

	<-listener

	go func() {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			newServer.Shutdown()
		}()

		os.Exit(0)
	}()

	<-listener

	os.Exit(0)
}
