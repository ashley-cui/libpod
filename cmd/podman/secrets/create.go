package secrets

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/containers/common/pkg/completion"
	"github.com/containers/podman/v2/cmd/podman/registry"
	"github.com/containers/podman/v2/pkg/domain/entities"
	"github.com/spf13/cobra"
)

var (
	// podman secret _create_
	createCmd = &cobra.Command{
		Use:     "create SECRET",
		Short:   "Create a new secret",
		Long:    "asjkdghlaiuehaudfhgilkadlh",
		RunE:    create,
		Example: "podman secret create mysecret",
		Args:    cobra.ExactArgs(2),
	}
)

var (
	createOpts = entities.SecretCreateOptions{}
)

func init() {
	// Subscribe inspect sub command to manifest command
	registry.Commands = append(registry.Commands, registry.CliCommand{
		// _podman manifest inspect_ will support both ABIMode and TunnelMode
		Mode: []entities.EngineMode{entities.ABIMode},
		// The definition for this command
		Command: createCmd,
		// The parent command to proceed this command on the CLI
		Parent: secretCmd,
	})

	// This is where you would configure the cobra flags using inspectCmd.Flags()
	flags := createCmd.Flags()

	driverFlagName := "driver"
	flags.StringVar(&createOpts.Driver, driverFlagName, "file", "Specify secret driver")
	_ = createCmd.RegisterFlagCompletionFunc(driverFlagName, completion.AutocompleteNone)
}

// Business logic: cmd is inspectCmd, args is the positional arguments from os.Args
func create(cmd *cobra.Command, args []string) error {
	name := args[0]

	var err error
	path := args[1]

	var reader io.Reader
	if path == "-" {
		stat, err := os.Stdin.Stat()
		if err != nil {
			return err
		}
		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			return errors.New("you need to pass something into stdin if youre going to use -")

		}
		reader = os.Stdin
	} else {
		reader, err = os.Open(path)
		if err != nil {
			return err
		}
	}

	report, err := registry.ContainerEngine().SecretCreate(context.Background(), name, reader, createOpts)
	if err != nil {
		return err
	}
	fmt.Println(report.ID)
	return nil
}
