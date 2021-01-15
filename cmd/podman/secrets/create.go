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
	createCmd = &cobra.Command{
		Use:   "create SECRET",
		Short: "Create a new secret",
		Long:  "Create a secret. Default driver is file (unencrypted).",
		RunE:  create,
		Args:  cobra.ExactArgs(2),
		Example: `podman secret create mysecret /path/to/secret
		printf "secretdata" | podman secret create mysecret -`,
	}
)

var (
	createOpts = entities.SecretCreateOptions{}
)

func init() {
	registry.Commands = append(registry.Commands, registry.CliCommand{
		Mode:    []entities.EngineMode{entities.ABIMode},
		Command: createCmd,
		Parent:  secretCmd,
	})

	flags := createCmd.Flags()

	driverFlagName := "driver"
	flags.StringVar(&createOpts.Driver, driverFlagName, "file", "Specify secret driver")
	_ = createCmd.RegisterFlagCompletionFunc(driverFlagName, completion.AutocompleteNone)
}

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
			return errors.New("if `-` is used, data must be passed into stdin")

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
