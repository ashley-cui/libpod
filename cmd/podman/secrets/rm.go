package secrets

import (
	"context"
	"fmt"

	"github.com/containers/podman/v2/cmd/podman/registry"
	"github.com/containers/podman/v2/pkg/domain/entities"
	"github.com/spf13/cobra"
)

var (
	// podman secret _rm_
	rmCmd = &cobra.Command{
		Use:     "rm SECRET",
		Short:   "Remove a secret",
		Long:    "asjkdghlaiuehaudfhgilkadlh",
		RunE:    rm,
		Example: "podman secret rm mysecret",
		Args:    cobra.ExactArgs(1),
	}
)

func init() {
	// Subscribe inspect sub command to manifest command
	registry.Commands = append(registry.Commands, registry.CliCommand{
		// _podman manifest inspect_ will support both ABIMode and TunnelMode
		Mode: []entities.EngineMode{entities.ABIMode},
		// The definition for this command
		Command: rmCmd,
		// The parent command to proceed this command on the CLI
		Parent: secretCmd,
	})

	// This is where you would configure the cobra flags using inspectCmd.Flags()
}

// Business logic: cmd is inspectCmd, args is the positional arguments from os.Args
func rm(cmd *cobra.Command, args []string) error {
	fmt.Println("reeeeeeemooved!")
	_, err := registry.ContainerEngine().SecretRm(context.Background(), args[0])
	if err != nil {
		return err
	}
	fmt.Println(args[0])
	return nil
}
